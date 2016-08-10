package application

import (
	"fmt"
	"net/url"
	"strings"

	"golang.org/x/net/context"

	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/thoas/gokvstores"
	"github.com/thoas/picfit/config"
	"github.com/thoas/picfit/engine"
	"github.com/thoas/picfit/errs"
	"github.com/thoas/picfit/hash"
	"github.com/thoas/picfit/image"
	"github.com/thoas/picfit/kvstore"
	"github.com/thoas/picfit/logger"
	"github.com/thoas/picfit/storage"
)

// Load creates a net/context from a file config path
func Load(path string) (context.Context, error) {
	cfg, err := config.Load(path)

	if err != nil {
		return nil, err
	}

	ctx := config.NewContext(context.Background(), *cfg)

	sourceStorage, destinationStorage, err := storage.NewStoragesFromConfig(cfg)

	if err != nil {
		return nil, err
	}

	ctx = storage.NewSourceContext(ctx, sourceStorage)
	ctx = storage.NewDestinationContext(ctx, destinationStorage)

	keystore, err := kvstore.NewKVStoreFromConfig(cfg)

	if err != nil {
		return nil, err
	}

	ctx = kvstore.NewContext(ctx, keystore)

	e := &engine.GoImageEngine{
		DefaultFormat:  cfg.Options.DefaultFormat,
		Format:         cfg.Options.Format,
		DefaultQuality: cfg.Options.Quality,
	}

	ctx = engine.NewContext(ctx, e)
	ctx = logger.NewContext(ctx, *logrus.New())

	return ctx, nil
}

// Store stores an image file with the defined filepath
func Store(ctx context.Context, filepath string, i *image.ImageFile) error {
	l := logger.FromContext(ctx)

	cfg := config.FromContext(ctx)

	k := kvstore.FromContext(ctx)
	con := k.Connection()
	defer con.Close()

	err := i.Save()

	if err != nil {
		l.Fatal(err)
		return err
	}

	l.Infof("Save thumbnail %s to storage", i.Filepath)

	prefix := cfg.KVStore.Prefix

	storeKey := i.Key

	key := i.Key

	if prefix != "" {
		storeKey = prefix + storeKey
	}

	err = con.Set(storeKey, i.Filepath)

	if err != nil {
		l.Fatal(err)

		return err
	}

	l.Infof("Save key %s => %s to kvstore", storeKey, i.Filepath)

	// Write children info only when we actually want to be able to delete things.
	if cfg.Options.EnableDelete {
		err = con.SetAdd(filepath+":children", key)

		if err != nil {
			l.Fatal(err)
			return err
		}

		l.Infof("Put key into set %s:children => %s in kvstore", filepath, key)
	}

	return nil
}

// ImageCleanup removes a file from kvstore and storage
func ImageCleanup(ctx context.Context, filepath string) error {
	k := kvstore.FromContext(ctx)
	con := k.Connection()
	defer con.Close()

	l := logger.FromContext(ctx)

	l.Infof("Deleting source storage file: %s", filepath)

	err := storage.SourceFromContext(ctx).Delete(filepath)

	if err != nil {
		return err
	}

	childrenPath := filepath + ":children"

	// Get the list of items to cleanup.
	children := con.SetMembers(childrenPath)

	// Delete them right away, we don't care about them anymore.
	l.Infof("Deleting children set: %s", childrenPath)

	err = con.Delete(childrenPath)

	if err != nil {
		return err
	}

	// No children? Okay..
	if children == nil {
		return nil
	}

	store := storage.DestinationFromContext(ctx)

	for _, s := range children {
		key, err := gokvstores.String(s)

		if err != nil {
			return err
		}

		// Now, every child is a hash which points to a key/value pair in
		// KVStore which in turn points to a file in dst storage.
		dstfile, err := gokvstores.String(con.Get(key))

		if err != nil {
			return err
		}

		// And try to delete it all. Ignore errors.
		err = store.Delete(dstfile)

		if err != nil {
			return err
		}

		err = con.Delete(key)

		if err != nil {
			return err
		}

		l.Infof("Deleting child %s and its KV store entry %s", dstfile, key)
	}

	return nil
}

// ImageFileFromContext generates an ImageFile from gin context
func ImageFileFromContext(c *gin.Context, async bool, load bool) (*image.ImageFile, error) {
	key := c.MustGet("key").(string)

	k := kvstore.FromContext(c)
	con := k.Connection()

	cfg := config.FromContext(c)

	l := logger.FromContext(c)

	defer con.Close()

	destStorage := storage.DestinationFromContext(c)

	var file = &image.ImageFile{
		Key:     key,
		Storage: destStorage,
		Headers: map[string]string{},
	}
	var err error

	prefix := cfg.KVStore.Prefix

	storeKey := key

	if prefix != "" {
		storeKey = prefix + key
	}

	// Image from the KVStore found
	stored, err := gokvstores.String(con.Get(storeKey))

	file.Filepath = stored

	if stored != "" {
		l.Infof("Key %s found in kvstore: %s", storeKey, stored)

		if load {
			file, err = image.FromStorage(destStorage, stored)

			if err != nil {
				return nil, err
			}
		}
	} else {
		l.Infof("Key %s not found in kvstore", storeKey)

		u, exists := c.Get("url")

		parameters := c.MustGet("parameters").(map[string]string)

		// Image not found from the KVStore, we need to process it
		// URL available in Query String
		if exists {
			file, err = image.FromURL(u.(*url.URL))
		} else {
			// URL provided we use http protocol to retrieve it
			s := storage.SourceFromContext(c)

			filepath := parameters["path"]

			if !s.Exists(filepath) {
				return nil, errs.ErrFileNotExists
			}

			file, err = image.FromStorage(s, filepath)
		}

		if err != nil {
			return nil, err
		}

		op := c.MustGet("op").(*engine.Operation)

		file, err = engine.FromContext(c).Transform(file, op, parameters)

		if err != nil {
			return nil, err
		}

		filename := ShardFilename(c, key)

		file.Filepath = fmt.Sprintf("%s.%s", filename, file.Format())
	}

	file.Key = key
	file.Storage = destStorage

	file.Headers["ETag"] = key

	if stored == "" {
		filepath := c.Param("path")

		if async == true {
			go Store(c, filepath, file)
		} else {
			err = Store(c, filepath, file)
		}
	}

	return file, err
}

// ShardFilename shards a filename based on config
func ShardFilename(ctx context.Context, filename string) string {
	cfg := config.FromContext(ctx)

	results := hash.Shard(filename, cfg.Shard.Width, cfg.Shard.Depth, true)

	return strings.Join(results, "/")
}
