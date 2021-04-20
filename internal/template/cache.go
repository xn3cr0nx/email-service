package template

import (
	"io/fs"
	"os"
	"path/filepath"

	"github.com/xn3cr0nx/email-service/pkg/logger"
)

type TemplateCache struct {
	Dir   string
	Cache map[string][]byte
}

var cache *TemplateCache

func NewTemplateCache(templateDir *string) (*TemplateCache, error) {
	if cache != nil {
		return cache, nil
	}
	cache = &TemplateCache{Dir: *templateDir}
	cache.Cache = make(map[string][]byte)
	addToCache := func(path string, f fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if filepath.Ext(path) != ".html" || f.IsDir() {
			return nil
		}

		if f.Type().IsRegular() {
			b, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			cache.Cache[path] = b
		}
		return nil
	}
	if err := filepath.WalkDir(*templateDir, addToCache); err != nil {
		logger.Error("Cache", err, logger.Params{"template dir": templateDir})
		return nil, err
	}

	logger.Info("Cache", "Initialized", logger.Params{"templates": len(cache.Cache), "base_dir": cache.Dir})
	return cache, nil
}

func (c *TemplateCache) Get(path string) []byte {
	return c.Cache[path]
}
