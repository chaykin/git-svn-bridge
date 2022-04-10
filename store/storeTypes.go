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
	SvnUserName string
	SvnPass     string
	GitUserName string
	Email       string
}

func getRepoKey(name string) string {
	return fmt.Sprintf("%s/repo", name)
}

func getUserKey(repo repo.Repo, email string) string {
	return fmt.Sprintf("%s/users/%s", repo.GetName(), email)
}
