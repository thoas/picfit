package picfit

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/url"
	"os"
	filepathpkg "path/filepath"
	"strings"
	"time"

	"github.com/cstockton/go-conv"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/thoas/picfit/storage"
	"github.com/thoas/picfit/util"
	"github.com/ulule/gostorages"

	"github.com/thoas/picfit/config"
	"github.com/thoas/picfit/engine"
	"github.com/thoas/picfit/failure"
	"github.com/thoas/picfit/hash"
	"github.com/thoas/picfit/image"
	"github.com/thoas/picfit/payload"
	"github.com/thoas/picfit/store"
)

type Processor struct {
	Logger *slog.Logger

	config             *config.Config
	destinationStorage *storage.Storage
	engine             *engine.Engine
	sourceStorage      *storage.Storage
	store              store.Store
}

// Upload uploads a file to its storage
func (p *Processor) Upload(ctx context.Context, payload *payload.Multipart) (*image.ImageFile, error) {
	var fh io.ReadCloser

	fh, err := payload.Data.Open()
	if err != nil {
		return nil, err
	}

	filepath := payload.Data.Filename
	if err := p.sourceStorage.Save(ctx, fh, filepath); err != nil {
		return nil, errors.Wrapf(err, "unable to save data on storage as: %s", filepath)
	}
	if err := fh.Close(); err != nil {
		return nil, errors.WithStack(err)
	}
	return &image.ImageFile{
		Filepath: payload.Data.Filename,
		Storage:  p.sourceStorage,
	}, nil
}

// Store stores an image file with the defined filepath
func (p *Processor) Store(ctx context.Context, log *slog.Logger, filepath string, i *image.ImageFile) error {
	starttime := time.Now()
	if err := i.Save(ctx); err != nil {
		return err
	}

	endtime := time.Now()
	log.InfoContext(ctx, "Save file to storage",
		slog.String("duration", endtime.Sub(starttime).String()),
	)

	starttime = time.Now()
	if err := p.store.Set(ctx, i.Key, i.Filepath); err != nil {
		return errors.WithStack(err)
	}
	endtime = time.Now()
	defaultMetrics.histogram.WithLabelValues(
		"store",
		strings.ToLower(filepathpkg.Ext(filepath)),
	).Observe(endtime.Sub(starttime).Seconds())

	log.InfoContext(ctx, "Save key to store",
		slog.String("duration", endtime.Sub(starttime).String()),
	)

	// Write children info only when we actually want to be able to delete things.
	if p.config.Options.EnableCascadeDelete {
		parentKey := hash.Tokey(filepath)

		parentKey = fmt.Sprintf("%s:children", parentKey)

		if err := p.store.AppendSlice(ctx, parentKey, i.Key); err != nil {
			return errors.WithStack(err)
		}

		log.InfoContext(ctx, "Put key into set in store",
			slog.String("set", parentKey),
			slog.String("value", filepath),
		)
	}

	return nil
}

// DeleteChild remove a child from store and storage
func (p *Processor) DeleteChild(ctx context.Context, key string) error {
	// Now, every child is a hash which points to a key/value pair in
	// Store which in turn points to a file in dst storage.
	dstfileRaw, err := p.store.Get(ctx, key)
	if err != nil {
		return errors.Wrapf(err, "unable to retrieve key %s", key)
	}

	if dstfileRaw != nil {
		dstfile, err := conv.String(dstfileRaw)
		if err != nil {
			return errors.Wrapf(err, "unable to cast %v to string", dstfileRaw)
		}

		// And try to delete it all.
		err = p.destinationStorage.Delete(ctx, dstfile)
		if err != nil {
			return errors.Wrapf(err, "unable to delete %s on storage", dstfile)
		}
	}

	if err := p.store.Delete(ctx, key); err != nil {
		return errors.Wrapf(err, "unable to delete key %s", key)
	}

	p.Logger.InfoContext(ctx, "Deleting child",
		slog.String("key", key))

	return nil
}

// Delete removes a file from store and storage
func (p *Processor) Delete(ctx context.Context, filepath string) error {
	p.Logger.InfoContext(ctx, "Deleting file on source storage",
		slog.String("file", filepath))

	if !p.FileExists(ctx, filepath) {
		p.Logger.InfoContext(ctx, "File does not exist anymore on source storage",
			slog.String("file", filepath))

		return errors.Wrapf(failure.ErrFileNotExists, "unable to delete, file does not exist: %s", filepath)
	}

	err := p.sourceStorage.Delete(ctx, filepath)
	if err != nil {
		return errors.Wrapf(err, "unable to delete %s on source storage", filepath)
	}

	parentKey := hash.Tokey(filepath)

	childrenKey := fmt.Sprintf("%s:children", parentKey)

	exists, err := p.store.Exists(ctx, childrenKey)
	if err != nil {
		return errors.Wrapf(err, "unable to verify if %s exists", childrenKey)
	}

	if !exists {
		p.Logger.InfoContext(ctx, "Children key does not exist in set",
			slog.String("key", childrenKey),
			slog.String("set", parentKey))

		return nil
	}

	// Get the list of items to cleanup.
	children, err := p.store.GetSlice(ctx, childrenKey)
	if err != nil {
		return errors.Wrapf(err, "unable to retrieve children set %s", childrenKey)
	}

	if children == nil {
		p.Logger.InfoContext(ctx, "No children to delete in set",
			slog.String("set", parentKey))

		return nil
	}

	for _, s := range children {
		key, err := conv.String(s)
		if err != nil {
			return err
		}

		err = p.DeleteChild(ctx, key)
		if err != nil {
			return errors.Wrapf(err, "unable to delete child %s", key)
		}
	}

	// Delete them right away, we don't care about them anymore.
	p.Logger.InfoContext(ctx, "Delete set %s",
		slog.String("set", childrenKey))

	if err := p.store.Delete(ctx, childrenKey); err != nil {
		return errors.Wrapf(err, "unable to delete key %s", childrenKey)
	}

	return nil
}

// ProcessContext processes a gin.Context generates and retrieves an ImageFile
func (p *Processor) ProcessContext(c *gin.Context, opts ...Option) (*image.ImageFile, error) {
	var (
		storeKey = c.MustGet("key").(string)
		force    = c.Query("force")
		options  = newOptions(opts...)
		ctx      = c.Request.Context()
		log      = p.Logger.With(slog.String("key", storeKey))
	)

	modifiedSince := c.Request.Header.Get("If-Modified-Since")
	if modifiedSince != "" && force == "" {
		exists, err := p.store.Exists(ctx, storeKey)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		if exists {
			log.InfoContext(ctx, "Key already exists on store, file not modified",
				slog.String("modified-since", modifiedSince))

			return nil, failure.ErrFileNotModified
		}
	}

	if force == "" {
		// try to retrieve image from the k/v rtore
		filepathRaw, err := p.store.Get(ctx, storeKey)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		if filepathRaw != nil {
			filepath, err := conv.String(filepathRaw)
			if err != nil {
				return nil, errors.WithStack(err)
			}

			log.InfoContext(ctx, "Key found in store",
				slog.String("filepath", filepath))

			starttime := time.Now()
			img, err := p.fileFromStorage(ctx, storeKey, filepath, options.Load)
			// no such file, just reprocess (maybe file cache was purged)
			if err != nil {
				if os.IsNotExist(err) {
					return p.processImage(c, storeKey)
				}

				return nil, errors.WithStack(err)
			}

			filesize := util.ByteCountDecimal(int64(len(img.Content())))
			endtime := time.Now()
			log.InfoContext(ctx, "Image retrieved from storage",
				slog.String("duration", endtime.Sub(starttime).String()),
				slog.String("size", filesize),
				slog.String("image", img.Filepath))

			defaultMetrics.histogram.WithLabelValues(
				"load",
				strings.ToLower(filepathpkg.Ext(filepath)),
			).Observe(endtime.Sub(starttime).Seconds())

			return img, nil
		}

		// Image not found from the Store, we need to process it
		// URL available in Query String
		log.InfoContext(ctx, "Key not found in store")
	} else {
		log.InfoContext(ctx, "Force activated, key will be re-processed")
	}

	return p.processImage(c, storeKey)
}

func (p *Processor) fileFromStorage(ctx context.Context, key string, filepath string, load bool) (*image.ImageFile, error) {
	var (
		file = &image.ImageFile{
			Key:      key,
			Storage:  p.destinationStorage,
			Filepath: filepath,
			Headers:  map[string]string{},
		}
		err error
	)

	if load {
		file, err = image.FromStorage(ctx, p.destinationStorage, filepath)
		if err != nil {
			return nil, err
		}
	}

	file.Headers["ETag"] = key
	return file, nil
}

func (p *Processor) processImage(c *gin.Context, storeKey string) (*image.ImageFile, error) {
	var (
		filepath string
		err      error
		ctx      = c.Request.Context()
		log      = p.Logger.With(slog.String("key", storeKey))
	)

	file := &image.ImageFile{
		Key:     storeKey,
		Storage: p.destinationStorage,
		Headers: map[string]string{},
	}

	qs := c.MustGet("parameters").(map[string]interface{})
	starttime := time.Now()
	u, exists := c.Get("url")
	if exists {
		file, err = image.FromURL(ctx, u.(*url.URL), p.config.Options.DefaultUserAgent)
	} else {
		// URL provided we use http protocol to retrieve it
		filepath = qs["path"].(string)
		if !p.FileExists(ctx, filepath) {
			return nil, errors.Wrapf(failure.ErrFileNotExists, "unable to process image, file does exist: %s", filepath)
		}

		file, err = image.FromStorage(ctx, p.sourceStorage, filepath)
	}
	if err != nil {
		return nil, errors.Wrap(err, "unable to process image")
	}
	endtime := time.Now()

	defaultMetrics.histogram.WithLabelValues(
		"load",
		strings.ToLower(filepathpkg.Ext(filepath)),
	).Observe(endtime.Sub(starttime).Seconds())

	filesize := util.ByteCountDecimal(int64(len(file.Content())))

	log = log.With(
		slog.String("image", file.Filepath),
		slog.String("size", filesize),
	)

	log.InfoContext(ctx, "Retrieved image to process from storage",
		slog.String("duration", endtime.Sub(starttime).String()))

	parameters, err := p.NewParameters(ctx, file, qs)
	if err != nil {
		return nil, errors.Wrap(err, "unable to process image")
	}

	starttime = time.Now()
	file, err = p.engine.Transform(ctx, parameters.output, parameters.operations)
	if err != nil {
		return nil, errors.Wrap(err, "unable to process image")
	}
	endtime = time.Now()
	defaultMetrics.histogram.WithLabelValues(
		"transform",
		strings.ToLower(filepathpkg.Ext(filepath)),
	).Observe(endtime.Sub(starttime).Seconds())

	filesize = util.ByteCountDecimal(int64(len(file.Content())))
	filename := p.ShardFilename(storeKey)
	file.Filepath = fmt.Sprintf("%s.%s", filename, file.Format())
	file.Storage = p.destinationStorage
	file.Key = storeKey
	file.Headers["ETag"] = storeKey

	log = log.With(
		slog.String("image", file.Filepath),
		slog.String("size", filesize),
	)

	log.InfoContext(ctx, "Image processed",
		slog.String("duration", endtime.Sub(starttime).String()))

	if err := p.Store(ctx, log, filepath, file); err != nil {
		return nil, errors.Wrapf(err, "unable to store processed image: %s", filepath)
	}

	return file, nil
}

// ShardFilename shards a filename based on config
func (p *Processor) ShardFilename(filename string) string {
	cfg := p.config

	results := hash.Shard(filename, cfg.Shard.Width, cfg.Shard.Depth, cfg.Shard.RestOnly)

	return strings.Join(results, "/")
}

func (p *Processor) GetKey(ctx context.Context, key string) (interface{}, error) {
	return p.store.Get(ctx, key)
}

func (p *Processor) KeyExists(ctx context.Context, key string) (bool, error) {
	return p.store.Exists(ctx, key)
}

func (p *Processor) FileExists(ctx context.Context, name string) bool {
	_, err := p.sourceStorage.Stat(ctx, name)
	return !errors.Is(err, gostorages.ErrNotExist)
}

func (p *Processor) OpenFile(ctx context.Context, name string) (io.ReadCloser, error) {
	return p.sourceStorage.Open(ctx, name)
}
