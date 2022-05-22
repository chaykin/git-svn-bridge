package gitHook

import (
	"fmt"
	"git-svn-bridge/repo"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type Factory struct {
	repo *repo.Repo
}

func New(repo *repo.Repo) *Factory {
	return &Factory{repo: repo}
}

func (f *Factory) CreateUpdateHook() {
	app := f.getApp()

	hook := f.getHookTemplate("update.hook")
	hook = strings.ReplaceAll(hook, "${WORKDIR}", filepath.Dir(app))
	hook = strings.ReplaceAll(hook, "${APP}", app)
	hook = strings.ReplaceAll(hook, "${REPO}", f.repo.GetName())

	f.writeHook("update", hook)
}

func (f *Factory) getApp() string {
	app, err := os.Executable()
	if err != nil {
		panic(fmt.Errorf("could not get executable: %w", err))
	}

	return app
}

func (f *Factory) getHookTemplate(name string) string {
	hookFileName := filepath.Join("gitHookTemplates", name)
	b, err := ioutil.ReadFile(hookFileName)
	if err != nil {
		panic(fmt.Errorf("could not read hook file %s: %w", hookFileName, err))
	}

	return string(b)
}

func (f *Factory) writeHook(name, content string) {
	hooksDirName := filepath.Join(f.repo.GetGitRepoPath(), "hooks")
	if err := os.MkdirAll(hooksDirName, 0774); err != nil {
		panic(fmt.Errorf("could not create hooks directory: %w", err))
	}

	hookFileName := filepath.Join(hooksDirName, name)
	if err := ioutil.WriteFile(hookFileName, []byte(content), 0554); err != nil {
		panic(fmt.Errorf("could not write hook file %s: %w", hookFileName, err))
	}
}
