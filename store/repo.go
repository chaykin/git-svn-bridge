package store

import (
	"fmt"
	"git-svn-bridge/repo"
)

type repoStoreItem struct {
	Name   string
	SvnUrl string
}

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

func getRepoKey(name string) string {
	return fmt.Sprintf("%s/repo", name)
}
