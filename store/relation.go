package store

import (
	"fmt"
	"git-svn-bridge/rel"
	"git-svn-bridge/repo"
)

type relStoreItem struct {
	Parent string
	Child  string
}

func StoreRelation(repo *repo.Repo, rel *rel.Relation) {
	relStoreItem := relStoreItem{Parent: rel.GetParent(), Child: rel.GetChild()}
	storeItem(getRelKey(repo, rel.GetParent()), relStoreItem)
	storeItem(getRelKey(repo, rel.GetChild()), relStoreItem)
}

func GetRelation(repo *repo.Repo, name string) *rel.Relation {
	key := getRelKey(repo, name)
	if getStore().Has(key) {
		var relStoreItem relStoreItem
		getItem(key, &relStoreItem)

		return rel.New(relStoreItem.Parent, relStoreItem.Child)
	}
	return nil
}

func RemoveRelation(repo *repo.Repo, name string) {
	relation := GetRelation(repo, name)
	if relation != nil {
		for _, key := range []string{getRelKey(repo, relation.GetChild()), getRelKey(repo, relation.GetParent())} {
			if getStore().Has(key) {
				if err := getStore().Erase(key); err != nil {
					panic(fmt.Errorf("could not erase key '%s': %w", key, err))
				}
			}
		}
	}
}

func getRelKey(repo *repo.Repo, name string) string {
	return fmt.Sprintf("%s/rel/%s", repo.GetName(), name)
}
