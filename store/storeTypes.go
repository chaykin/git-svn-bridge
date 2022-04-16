package store

import (
	"fmt"
	"git-svn-bridge/repo"
)

type repoStoreItem struct {
	Name   string
	SvnUrl string
}

type userStoreItem struct {
	SvnUserName     string
	SvnPass         string
	GitUserName     string
	GitUserFullName string
	Email           string
}

func getRepoKey(name string) string {
	return fmt.Sprintf("%s/repo", name)
}

func getUserKey(repo *repo.Repo, gitUserName string) string {
	return fmt.Sprintf("%s/users/%s", repo.GetName(), gitUserName)
}
