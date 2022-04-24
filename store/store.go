package store

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"git-svn-bridge/conf"
	"github.com/peterbourgon/diskv/v3"
	"strings"
)

var store *diskv.Diskv

func storeItem(key string, item interface{}) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	if err := enc.Encode(item); err != nil {
		panic(fmt.Errorf("could not encode item for store: %w", err))
	}

	if err := getStore().Write(key, buf.Bytes()); err != nil {
		panic(fmt.Errorf("could not write item to store: %w", err))
	}
}

func getItem(key string, item interface{}) {
	itemBytes, err := getStore().Read(key)
	if err != nil {
		panic(fmt.Errorf("could not read item '%s' from store: %w", key, err))
	}

	buf := bytes.NewBuffer(itemBytes)
	dec := gob.NewDecoder(buf)

	if err := dec.Decode(item); err != nil {
		panic(fmt.Errorf("could not decode item '%s': %w", key, err))
	}
}

func getStore() *diskv.Diskv {
	if store == nil {
		config := conf.GetConfig()
		store = diskv.New(diskv.Options{
			BasePath:          config.DbRoot,
			AdvancedTransform: advancedTransform,
			InverseTransform:  inverseTransform,
			CacheSizeMax:      config.DbCacheSize,
		})
	}

	return store
}

func advancedTransform(key string) *diskv.PathKey {
	path := strings.Split(key, "/")
	last := len(path) - 1
	return &diskv.PathKey{
		Path:     path[:last],
		FileName: path[last],
	}
}

func inverseTransform(pathKey *diskv.PathKey) (key string) {
	parentPath := strings.Join(pathKey.Path, "/")
	if parentPath != "" {
		parentPath += "/"
	}
	return parentPath + pathKey.FileName
}
