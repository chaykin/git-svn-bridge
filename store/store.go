package store

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"git-svn-bridge/conf"
	"git-svn-bridge/repo"
	"git-svn-bridge/usr"
	"github.com/peterbourgon/diskv/v3"
	"strings"
)

var store *diskv.Diskv

func StoreRepo(repo repo.Repo) {
	repoStoreItem := repoStoreItem{Name: repo.GetName(), SvnUrl: repo.GetSvnUrl()}
	storeItem(getRepoKey(repoStoreItem.Name), repoStoreItem)
}

func HasRepo(name string) bool {
	return getStore().Has(getRepoKey(name))
}

func GetRepo(name string) repo.Repo {
	var storeItem repoStoreItem
	getItem(getRepoKey(name), &storeItem)

	return repo.CreateRepo(storeItem.Name, storeItem.SvnUrl)
}

func StoreUser(user usr.User) {
	userStoreItem := userStoreItem{user.GetSvnUserName(), user.GetSvnPasswordEncrypted(), user.GetGitUserName(), user.GetEmail()}
	storeItem(getUserKey(user.GetRepo(), user.GetEmail()), userStoreItem)
}

func GetUser(repo repo.Repo, email string) usr.User {
	var storeItem userStoreItem
	getItem(getUserKey(repo, email), &storeItem)

	return usr.CreateEncryptedUser(repo, storeItem.SvnUserName, storeItem.SvnPass, storeItem.GitUserName, storeItem.Email)
}

func storeItem(key string, item interface{}) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	err := enc.Encode(item)
	if err != nil {
		panic(fmt.Errorf("could not encode item for store: %w", err))
	}

	err = getStore().Write(key, buf.Bytes())
	if err != nil {
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

	err = dec.Decode(item)
	if err != nil {
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
	return strings.Join(pathKey.Path, "/") + pathKey.FileName[:len(pathKey.FileName)]
}