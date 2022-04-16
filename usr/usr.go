package usr

import (
	"git-svn-bridge/crypt"
	"git-svn-bridge/repo"
)

type User struct {
	repo            repo.Repo
	svnUserName     string
	svnPass         string
	gitUserName     string
	gitUserFullName string
	email           string
}

func (u *User) GetRepo() repo.Repo {
	return u.repo
}

func (u *User) GetSvnUserName() string {
	return u.svnUserName
}

func (u *User) GetSvnPassword() string {
	return crypt.Decrypt(u.svnPass)
}

func (u *User) GetSvnPasswordEncrypted() string {
	return u.svnPass
}

func (u *User) GetGitUserName() string {
	return u.gitUserName
}

func (u *User) GetGitUserFullName() string {
	return u.gitUserFullName
}

func (u *User) GetEmail() string {
	return u.email
}

func CreateUser(repo repo.Repo, svnUserName, svnPass, gitUserName, gitUserFullName, email string) User {
	return User{repo, svnUserName, crypt.Encrypt(svnPass), gitUserName, gitUserFullName, email}
}

func CreateEncryptedUser(repo repo.Repo, svnUserName, svnPass, gitUserName, gitUserFullName, email string) User {
	return User{repo, svnUserName, svnPass, gitUserName, gitUserFullName, email}
}
