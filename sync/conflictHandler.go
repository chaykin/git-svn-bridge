package sync

import (
	"fmt"
	"git-svn-bridge/rel"
	"git-svn-bridge/repo"
	"git-svn-bridge/store"
	"git-svn-bridge/vcs/gitutils"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"strings"
	"time"
)

const BridgeConflictBranchPrefix = "bridge-conflict-"

type conflictHandler struct {
	man *Manager
}

// Creates special branch for manual conflicts resolution
func (h *conflictHandler) handleConflict(ref string) string {
	branchName := gitutils.GetBranchName(ref)
	conflictBranchName := h.createConflictBranchName()
	relation := rel.New(branchName, conflictBranchName)
	store.StoreRelation(h.man.repo, relation)

	h.man.checkoutToOriginRef(ref)

	worktree, err := h.man.bridgeRepo.Worktree()
	if err != nil {
		panic(fmt.Errorf("could not get worktree for repo '%s': %w", h.man.getBridgeRepoPath(), err))
	}

	err = worktree.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName(conflictBranchName),
		Create: true,
	})
	if err != nil {
		panic(fmt.Errorf("could not checkout '%s' for repo '%s': %w", conflictBranchName, h.man.getBridgeRepoPath(), err))
	}

	gitutils.Fetch(h.man.getCentralRepoPath(), "bridge", conflictBranchName)

	return conflictBranchName
}

func getParentRef(repo *repo.Repo, ref string) string {
	branchName := gitutils.GetBranchName(ref)
	relation := store.GetRelation(repo, branchName)
	if relation == nil {
		return ref
	}

	return strings.Replace(ref, "/"+branchName, "/"+relation.GetParent(), 1)
}

func (h *conflictHandler) createConflictBranchName() string {
	return BridgeConflictBranchPrefix + time.Now().Format("2006-01-02T15_04_05.000")
}
