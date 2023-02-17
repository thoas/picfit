package picfit

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/url"
	"os"
	filepathpkg "path/filepath"
	"strings"
	"time"

	"github.com/cstockton/go-conv"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/thoas/picfit/util"
	"github.com/ulule/gostorages"

	"github.com/thoas/picfit/config"
	"github.com/thoas/picfit/engine"
	"github.com/thoas/picfit/failure"
	"github.com/thoas/picfit/hash"
	"github.com/thoas/picfit/image"
	"github.com/thoas/picfit/logger"
	"github.com/thoas/picfit/payload"
	"github.com/thoas/picfit/store"
)

type Processor struct {
	config             *config.Config
	destinationStorage gostorages.Storage
	engine             *engine.Engine
	logger             logger.Logger
	sourceStorage      gostorages.Storage
	store              store.Store
}

// Upload uploads a file to its storage
func (p *Processor) Upload(payload *payload.Multipart) (*image.ImageFile, error) {
	var fh io.ReadCloser

	fh, err := payload.Data.Open()
	if err != nil {
		return nil, err
	}

	dataBytes := bytes.Buffer{}

	_, err = dataBytes.ReadFrom(fh)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to read data from uploaded file")
	}
	if err := fh.Close(); err != nil {
		return nil, err
	}

	err = p.sourceStorage.Save(payload.Data.Filename, gostorages.NewContentFile(dataBytes.Bytes()))
	if err != nil {
		return nil, errors.Wrapf(err, "unable to save data on storage as: %s", payload.Data.Filename)
	}

	return &image.ImageFile{
		Filepath: payload.Data.Filename,
		Storage:  p.sourceStorage,
	}, nil
}

// Store stores an image file with the defined filepath
func (p *Processor) Store(ctx context.Context, filepath string, i *image.ImageFile) error {
	starttime := time.Now()
	if err := i.Save(); err != nil {
		return err
	}

	endtime := time.Now()
	p.logger.Info("Save file to storage",
		logger.Duration("duration", endtime.Sub(starttime)),
		logger.String("file", i.Filepath))

	starttime = time.Now()
	if err := p.store.Set(ctx, i.Key, i.Filepath); err != nil {
		return err
	}
	endtime = time.Now()
	defaultMetrics.histogram.WithLabelValues(
		"store",
		strings.ToLower(filepathpkg.Ext(filepath)),
	).Observe(endtime.Sub(starttime).Seconds())

	p.logger.Info("Save key to store",
		logger.String("key", i.Key),
		logger.Duration("duration", endtime.Sub(starttime)),
		logger.String("filepath", i.Filepath))

	// Write children info only when we actually want to be able to delete things.
	if p.config.Options.EnableCascadeDelete {
		parentKey := hash.Tokey(filepath)

		parentKey = fmt.Sprintf("%s:children", parentKey)

		if err := p.store.AppendSlice(ctx, parentKey, i.Key); err != nil {
			return err
		}

		p.logger.Info("Put key into set in store",
			logger.String("set", parentKey),
			logger.String("value", filepath),
			logger.String("key", i.Key))
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
		err = p.destinationStorage.Delete(dstfile)
		if err != nil {
			return errors.Wrapf(err, "unable to delete %s on storage", dstfile)
		}
	}

	err = p.store.Delete(ctx, key)
	if err != nil {
		return errors.Wrapf(err, "unable to delete key %s", key)
	}

	p.logger.Info("Deleting child",
		logger.String("key", key))

	return nil
}

// Delete removes a file from store and storage
func (p *Processor) Delete(ctx context.Context, filepath string) error {
	p.logger.Info("Deleting file on source storage",
		logger.String("file", filepath))

	if !p.sourceStorage.Exists(filepath) {
		p.logger.Info("File does not exist anymore on source storage",
			logger.String("file", filepath))

		return errors.Wrapf(failure.ErrFileNotExists, "unable to delete, file does not exist: %s", filepath)
	}

	err := p.sourceStorage.Delete(filepath)
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
		p.logger.Info("Children key does not exist in set",
			logger.String("key", childrenKey),
			logger.String("set", parentKey))

		return nil
	}

	// Get the list of items to cleanup.
	children, err := p.store.GetSlice(ctx, childrenKey)
	if err != nil {
		return errors.Wrapf(err, "unable to retrieve children set %s", childrenKey)
	}

	if children == nil {
		p.logger.Info("No children to delete in set",
			logger.String("set", parentKey))

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
	p.logger.Info("Delete set %s",
		logger.String("set", childrenKey))

	err = p.store.Delete(ctx, childrenKey)
	if err != nil {
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
		log      = p.logger.With(logger.String("key", storeKey))
	)

	modifiedSince := c.Request.Header.Get("If-Modified-Since")
	if modifiedSince != "" && force == "" {
		exists, err := p.store.Exists(ctx, storeKey)
		if err != nil {
			return nil, err
		}

		if exists {
			log.Info("Key already exists on store, file not modified",
				logger.String("modified-since", modifiedSince))

			return nil, failure.ErrFileNotModified
		}
	}

	if force == "" {
		// try to retrieve image from the k/v rtore
		filepathRaw, err := p.store.Get(ctx, storeKey)
		if err != nil {
			return nil, err
		}

		if filepathRaw != nil {
			filepath, err := conv.String(filepathRaw)
			if err != nil {
				return nil, err
			}

			log.Info("Key found in store",
				logger.String("filepath", filepath))

			starttime := time.Now()
			img, err := p.fileFromStorage(storeKey, filepath, options.Load)
			// no such file, just reprocess (maybe file cache was purged)
			if err != nil {
				if os.IsNotExist(err) {
					return p.processImage(c, storeKey, options.Async)
				}

				return nil, err
			}

			filesize := util.ByteCountDecimal(int64(len(img.Content())))
			endtime := time.Now()
			log.Info("Image retrieved from storage",
				logger.Duration("duration", endtime.Sub(starttime)),
				logger.String("size", filesize),
				logger.String("image", img.Filepath))

			defaultMetrics.histogram.WithLabelValues(
				"load",
				strings.ToLower(filepathpkg.Ext(filepath)),
			).Observe(endtime.Sub(starttime).Seconds())

			return img, nil
		}

		// Image not found from the Store, we need to process it
		// URL available in Query String
		log.Info("Key not found in store")
	} else {
		log.Info("Force activated, key will be re-processed")
	}

	return p.processImage(c, storeKey, options.Async)
}

func (p *Processor) fileFromStorage(key string, filepath string, load bool) (*image.ImageFile, error) {
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
		file, err = image.FromStorage(p.destinationStorage, filepath)
		if err != nil {
			return nil, err
		}
	}

	file.Headers["ETag"] = key
	return file, nil
}

func (p *Processor) processImage(c *gin.Context, storeKey string, async bool) (*image.ImageFile, error) {
	var (
		filepath string
		err      error
		ctx      = c.Request.Context()
		log      = p.logger.With(logger.String("key", storeKey))
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
		file, err = image.FromURL(u.(*url.URL), p.config.Options.DefaultUserAgent)
	} else {
		// URL provided we use http protocol to retrieve it
		filepath = qs["path"].(string)
		if !p.sourceStorage.Exists(filepath) {
			return nil, errors.Wrapf(failure.ErrFileNotExists, "unable to process image, file does exist: %s", filepath)
		}

		file, err = image.FromStorage(p.sourceStorage, filepath)
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
	log.Info("Retrieved image to process from storage",
		logger.Duration("duration", endtime.Sub(starttime)),
		logger.String("image", file.Filepath),
		logger.String("size", filesize))

	parameters, err := p.NewParameters(file, qs)
	if err != nil {
		return nil, errors.Wrap(err, "unable to process image")
	}

	starttime = time.Now()
	file, err = p.engine.Transform(parameters.output, parameters.operations)
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

	log.Info("Image processed",
		logger.String("image", file.Filepath),
		logger.Duration("duration", endtime.Sub(starttime)),
		logger.String("size", filesize))

	if async == true {
		go func() {
			if err := p.Store(context.Background(), filepath, file); err != nil {
				p.logger.Error("async store", logger.Error(err))
			}
		}()
	} else {
		if err := p.Store(ctx, filepath, file); err != nil {
			return nil, errors.Wrapf(err, "unable to store processed image: %s", filepath)
		}
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

func (p *Processor) FileExists(name string) bool {
	return p.sourceStorage.Exists(name)
}

func (p *Processor) OpenFile(name string) (gostorages.File, error) {
	return p.sourceStorage.Open(name)
}
