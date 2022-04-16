package repo

import (
	"git-svn-bridge/conf"
	"path/filepath"
)

type Repo struct {
	name   string
	svnUrl string
}

func CreateRepo(name, snvUrl string) Repo {
	return Repo{name, snvUrl}
}

func (r *Repo) GetName() string {
	return r.name
}

func (r *Repo) GetSvnUrl() string {
	return r.svnUrl
}

func (r *Repo) GetGitRepoPath() string {
	return r.getRepoPath("git")
}

func (r *Repo) GetBridgeRepoPath() string {
	return r.getRepoPath("bridge")
}

func (r *Repo) getRepoPath(repoType string) string {
	return filepath.Join(conf.GetConfig().ReposRoot, repoType, r.name)
}
