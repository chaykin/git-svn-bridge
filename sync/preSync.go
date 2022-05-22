package sync

import (
	"git-svn-bridge/log"
	"git-svn-bridge/store"
	"git-svn-bridge/vcs/gitutils"
	"strings"
)

type PreSyncManager struct {
	man *Manager
}

func NewPreSync(repoName string) *PreSyncManager {
	man := New(repoName)
	return &PreSyncManager{man: man}
}

func (pSyncMan *PreSyncManager) PreSync(ref, oldSha, newSha string) {
	mergeBase := strings.TrimSpace(gitutils.GetMergeBase(pSyncMan.man.getCentralRepoPath(), oldSha, newSha))
	if mergeBase != oldSha {
		log.StdErrFatalf("Non-fast-forward commits are not allowed")
	}

	branchName := gitutils.GetBranchName(ref)
	relation := store.GetRelation(pSyncMan.man.repo, branchName)
	if relation != nil && relation.GetParent() == branchName {
		log.StdErrFatalf("Commits to branch %s is not allowed due to conflicts. Resolve conflicts manually in branch %s and commit (it will be merged automatically)",
			branchName, relation.GetChild())
	}
}
