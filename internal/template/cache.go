package template

import (
	"os"
	"path/filepath"

	"github.com/xn3cr0nx/email-service/pkg/logger"
)

type (
	TemplateCache map[string][]byte
)

var cache TemplateCache

func NewTemplateCache(dir *string) (TemplateCache, error) {
	if len(cache) != 0 {
		return cache, nil
	}
	cache = make(TemplateCache)
	addToCache := func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if filepath.Ext(path) != ".html" || f.IsDir() {
			return nil
		}

		if f.Mode().IsRegular() {
			b, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			cache[path] = b
		}
		return nil
	}
	if err := filepath.Walk(*dir, addToCache); err != nil {
		logger.Error("Cache", err, logger.Params{"dir": dir})
		return cache, err
	}

	return cache, nil
}
