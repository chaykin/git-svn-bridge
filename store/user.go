package store

import (
	"fmt"
	"git-svn-bridge/repo"
	"git-svn-bridge/usr"
)

type userStoreItem struct {
	SvnUserName string
	SvnPass     string
	GitUserName string
	Email       string
}

func StoreUser(user usr.User) {
	userStoreItem := userStoreItem{user.GetSvnUserName(), user.GetSvnPasswordEncrypted(), user.GetGitUserName(), user.GetEmail()}
	storeItem(getUserKey(user.GetRepo(), user.GetGitUserName()), userStoreItem)
}

func GetUser(repo *repo.Repo, gitUserName string) usr.User {
	var storeItem userStoreItem
	getItem(getUserKey(repo, gitUserName), &storeItem)

	return usr.CreateEncryptedUser(repo, storeItem.SvnUserName, storeItem.SvnPass, storeItem.GitUserName, storeItem.Email)
}

func GetAllUsers(repo *repo.Repo) []usr.User {
	prefix := getUserKey(repo, "")
	userKeysChan := getStore().KeysPrefix(prefix, nil)

	var users []usr.User
	for userKey := range userKeysChan {
		var storeItem userStoreItem
		getItem(userKey, &storeItem)
		user := usr.CreateEncryptedUser(repo, storeItem.SvnUserName, storeItem.SvnPass, storeItem.GitUserName, storeItem.Email)

		users = append(users, user)
	}

	return users
}

func getUserKey(repo *repo.Repo, gitUserName string) string {
	return fmt.Sprintf("%s/users/%s", repo.GetName(), gitUserName)
}
