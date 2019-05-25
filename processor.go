package picfit

import (
	"bytes"
	"fmt"
	"io"
	"net/url"
	"strings"

	conv "github.com/cstockton/go-conv"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
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
	logger             logger.Logger
	SourceStorage      gostorages.Storage
	DestinationStorage gostorages.Storage
	store              store.Store
	engine             *engine.Engine
}

// Upload uploads a file to its storage
func (p *Processor) Upload(payload *payload.Multipart) (*image.ImageFile, error) {
	var fh io.ReadCloser

	fh, err := payload.Data.Open()
	if err != nil {
		return nil, err
	}
	defer fh.Close()

	dataBytes := bytes.Buffer{}

	_, err = dataBytes.ReadFrom(fh)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to read data from uploaded file")
	}

	err = p.SourceStorage.Save(payload.Data.Filename, gostorages.NewContentFile(dataBytes.Bytes()))
	if err != nil {
		return nil, errors.Wrapf(err, "unable to save data on storage as: %s", payload.Data.Filename)
	}

	return &image.ImageFile{
		Filepath: payload.Data.Filename,
		Storage:  p.SourceStorage,
	}, nil
}

// Store stores an image file with the defined filepath
func (p *Processor) Store(filepath string, i *image.ImageFile) error {
	err := i.Save()
	if err != nil {
		return err
	}

	p.logger.Info("Save file to storage",
		logger.String("file", i.Filepath))

	err = p.store.Set(i.Key, i.Filepath)
	if err != nil {
		return err
	}

	p.logger.Info("Save key to store",
		logger.String("key", i.Key),
		logger.String("filepath", i.Filepath))

	// Write children info only when we actually want to be able to delete things.
	if p.config.Options.EnableCascadeDelete {
		parentKey := hash.Tokey(filepath)

		parentKey = fmt.Sprintf("%s:children", parentKey)

		err = p.store.AppendSlice(parentKey, i.Key)
		if err != nil {
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
func (p *Processor) DeleteChild(key string) error {
	// Now, every child is a hash which points to a key/value pair in
	// Store which in turn points to a file in dst storage.
	dstfileRaw, err := p.store.Get(key)
	if err != nil {
		return errors.Wrapf(err, "unable to retrieve key %s", key)
	}

	if dstfileRaw != nil {
		dstfile, err := conv.String(dstfileRaw)
		if err != nil {
			return errors.Wrapf(err, "unable to cast %v to string", dstfileRaw)
		}

		// And try to delete it all.
		err = p.DestinationStorage.Delete(dstfile)
		if err != nil {
			return errors.Wrapf(err, "unable to delete %s on storage", dstfile)
		}
	}

	err = p.store.Delete(key)
	if err != nil {
		return errors.Wrapf(err, "unable to delete key %s", key)
	}

	p.logger.Info("Deleting child",
		logger.String("key", key))

	return nil
}

// Delete removes a file from store and storage
func (p *Processor) Delete(filepath string) error {
	p.logger.Info("Deleting file on source storage",
		logger.String("file", filepath))

	if !p.SourceStorage.Exists(filepath) {
		p.logger.Info("File does not exist anymore on source storage",
			logger.String("file", filepath))

		return errors.Wrapf(failure.ErrFileNotExists, "unable to delete, file does not exist: %s", filepath)
	}

	err := p.SourceStorage.Delete(filepath)
	if err != nil {
		return errors.Wrapf(err, "unable to delete %s on source storage", filepath)
	}

	parentKey := hash.Tokey(filepath)

	childrenKey := fmt.Sprintf("%s:children", parentKey)

	exists, err := p.store.Exists(childrenKey)
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
	children, err := p.store.GetSlice(childrenKey)
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

		err = p.DeleteChild(key)
		if err != nil {
			return errors.Wrapf(err, "unable to delete child %s", key)
		}
	}

	// Delete them right away, we don't care about them anymore.
	p.logger.Info("Delete set %s",
		logger.String("set", childrenKey))

	err = p.store.Delete(childrenKey)
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
	)

	modifiedSince := c.Request.Header.Get("If-Modified-Since")
	if modifiedSince != "" && force == "" {
		exists, err := p.store.Exists(storeKey)
		if err != nil {
			return nil, err
		}

		if exists {
			p.logger.Info("Key already exists on store, file not modified",
				logger.String("key", storeKey),
				logger.String("modified-since", modifiedSince))

			return nil, failure.ErrFileNotModified
		}
	}

	if force == "" {
		// try to retrieve image from the k/v rtore
		filepathRaw, err := p.store.Get(storeKey)
		if err != nil {
			return nil, err
		}

		if filepathRaw != nil {
			filepath, err := conv.String(filepathRaw)
			if err != nil {
				return nil, err
			}

			p.logger.Info("Key found in store",
				logger.String("key", storeKey),
				logger.String("filepath", filepath))

			return p.fileFromStorage(storeKey, filepath, options.Load)
		}

		// Image not found from the Store, we need to process it
		// URL available in Query String
		p.logger.Info("Key not found in store",
			logger.String("key", storeKey))
	} else {
		p.logger.Info("Force activated, key will be re-processed",
			logger.String("key", storeKey))
	}

	return p.processImage(c, storeKey, options.Async)
}

func (p *Processor) fileFromStorage(key string, filepath string, load bool) (*image.ImageFile, error) {
	var (
		file = &image.ImageFile{
			Key:      key,
			Storage:  p.DestinationStorage,
			Filepath: filepath,
			Headers:  map[string]string{},
		}
		err error
	)

	if load {
		file, err = image.FromStorage(p.DestinationStorage, filepath)
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
	)

	file := &image.ImageFile{
		Key:     storeKey,
		Storage: p.DestinationStorage,
		Headers: map[string]string{},
	}

	qs := c.MustGet("parameters").(map[string]interface{})

	u, exists := c.Get("url")
	if exists {
		file, err = image.FromURL(u.(*url.URL), p.config.Options.DefaultUserAgent)
	} else {
		// URL provided we use http protocol to retrieve it
		filepath = qs["path"].(string)
		if !p.SourceStorage.Exists(filepath) {
			return nil, errors.Wrapf(failure.ErrFileNotExists, "unable to process image, file does exist: %s", filepath)
		}

		file, err = image.FromStorage(p.SourceStorage, filepath)
	}
	if err != nil {
		return nil, errors.Wrap(err, "unable to process image")
	}

	parameters, err := p.newParameters(file, qs)
	if err != nil {
		return nil, errors.Wrap(err, "unable to process image")
	}

	file, err = p.engine.Transform(parameters.Output, parameters.Operations)
	if err != nil {
		return nil, errors.Wrap(err, "unable to process image")
	}

	filename := p.ShardFilename(storeKey)
	file.Filepath = fmt.Sprintf("%s.%s", filename, file.Format())
	file.Storage = p.DestinationStorage
	file.Key = storeKey
	file.Headers["ETag"] = storeKey

	if async == true {
		go p.Store(filepath, file)
	} else {
		err = p.Store(filepath, file)
		if err != nil {
			return nil, errors.Wrapf(err, "unable to store processed image: %s", filepath)
		}
	}

	return file, nil
}

// ShardFilename shards a filename based on config
func (p Processor) ShardFilename(filename string) string {
	cfg := p.config

	results := hash.Shard(filename, cfg.Shard.Width, cfg.Shard.Depth, cfg.Shard.RestOnly)

	return strings.Join(results, "/")
}

func (p Processor) GetKey(key string) (interface{}, error) {
	return p.store.Get(key)
}

func (p Processor) KeyExists(key string) (bool, error) {
	return p.store.Exists(key)
}

func (p Processor) FileExists(name string) bool {
	return p.SourceStorage.Exists(name)
}

func (p Processor) OpenFile(name string) (gostorages.File, error) {
	return p.SourceStorage.Open(name)
}
