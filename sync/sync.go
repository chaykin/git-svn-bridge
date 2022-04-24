package sync

import (
	"fmt"
	"git-svn-bridge/conf"
	"git-svn-bridge/log"
	"git-svn-bridge/repo"
	"git-svn-bridge/store"
	"git-svn-bridge/usr"
	"git-svn-bridge/vcs/gitsvn"
	"git-svn-bridge/vcs/gitutils"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"strings"
)

const maxRetryCount = 3

type Manager struct {
	repo *repo.Repo

	bridgeRepo  *git.Repository
	centralRepo *git.Repository

	isSvnFetched bool
	retryCount   int
}

func New(repoName string) *Manager {
	if !store.HasRepo(repoName) {
		panic(fmt.Errorf("there is no repo with name %s", repoName))
	}

	repoInfo := store.GetRepo(repoName)

	bridgeRepoPath := repoInfo.GetBridgeRepoPath()
	centralRepoPath := repoInfo.GetGitRepoPath()

	bridgeRepo, err := git.PlainOpen(bridgeRepoPath)
	if err != nil {
		panic(fmt.Errorf("could not open bridge repo '%s': %w", bridgeRepoPath, err))
	}

	centralRepo, err := git.PlainOpen(centralRepoPath)
	if err != nil {
		panic(fmt.Errorf("could not open central repo '%s': %w", centralRepoPath, err))
	}

	return &Manager{repo: &repoInfo, bridgeRepo: bridgeRepo, centralRepo: centralRepo, isSvnFetched: false, retryCount: 0}
}

func (man *Manager) SyncAllRefs() {
	man.fetchSvnChanges()

	allRefs := man.getAllRefs()
	man.SyncRefs(allRefs)
}

func (man *Manager) SyncRefs(refs []string) {
	for _, ref := range refs {
		man.syncRef(ref)
	}
}

func (man *Manager) syncRef(ref string) {
	man.fetchSvnChanges()
	man.fetchGitChanges(ref)
	man.pushRefToSvn(ref)
	man.pushRefToGit(ref)

	store.RemoveRelation(man.repo, gitutils.GetBranchName(ref))
	man.retryCount = 0
}

func (man *Manager) fetchSvnChanges() {
	if man.isSvnFetched {
		return
	}
	man.isSvnFetched = true

	gitSvnExecutor := gitsvn.CreateExecutor(man.getSysUser())
	gitSvnExecutor.Fetch(man.getBridgeRepoPath())
}

func (man *Manager) fetchGitChanges(ref string) {
	man.checkoutToRef(ref)

	centralRepoRefExists := gitutils.IsRefExists(man.centralRepo, ref)
	if centralRepoRefExists {
		branchName := gitutils.GetBranchName(ref)

		defer func() {
			if r := recover(); r != nil {
				gitutils.AbortRebase(man.getBridgeRepoPath())

				//There are a conflicts! Handle with it...
				conflictBranch := (&conflictHandler{man: man}).handleConflict(ref)
				log.StdErrFatalf("Could not sync changes with SVN. You must resolve conflicts manually in branch %s", conflictBranch)
			}
		}()

		gitutils.PullAndRebase(man.getBridgeRepoPath(), gitutils.GitCentralRepoName, branchName)
	}
}

func (man *Manager) pushRefToSvn(ref string) {
	// store the SVN URL and author of the last commit for `git svn dcommit`
	gitUserName := gitutils.GetGitAuthor(man.repo, man.getBridgeRepoPath())

	man.checkoutToOriginRef(getParentRef(man.repo, ref))
	man.mergeWithConflicts(ref)
	man.tryCommit(gitUserName, ref)
}

func (man *Manager) mergeWithConflicts(ref string) {
	branchName := gitutils.GetBranchName(ref)

	// create a merged log message for the merge commit to SVN. Note that we squash log messages together for
	// merge commits. Experiment with the fomat, e.g. '%ai | [%an] %s' etc
	commitMessage := gitutils.BuildCommitMessage(man.getBridgeRepoPath(), branchName)

	defer func() {
		if r := recover(); r != nil {
			gitutils.AbortMerge(man.getBridgeRepoPath())

			//There are a conflicts! Handle with it...
			conflictBranch := (&conflictHandler{man: man}).handleConflict(ref)
			log.StdErrFatalf("Could not sync changes with SVN. You must resolve conflicts manually in branch %s", conflictBranch)
		}
	}()

	// Merge changes from master to the SVN-tracking branch and commit to SVN.
	// Note that we always record the merge with --no-ff
	gitutils.MergeNoFF(man.getBridgeRepoPath(), commitMessage, branchName)
}

//commit changes to SVN
func (man *Manager) tryCommit(gitUserName, ref string) {
	commitUser := store.GetUser(man.repo, gitUserName)
	gitSvnExecutor := gitsvn.CreateExecutor(commitUser)

	defer func() {
		if r := recover(); r != nil {
			man.retryCount++
			if man.retryCount >= maxRetryCount {
				panic(fmt.Errorf("too many tries to push to SVN ror repo '%s'(%s): %s", man.repo.GetName(), ref, r))
			}

			// Guess somebody commit to SVN just now. Try to fetch changes again
			man.isSvnFetched = false
			man.syncRef(ref)
		}
	}()

	gitSvnExecutor.Commit(man.getBridgeRepoPath())
}

func (man *Manager) pushRefToGit(ref string) {
	ref = getParentRef(man.repo, ref)
	man.checkoutToRef(ref)

	branchName := gitutils.GetBranchName(ref)

	// --- merge changes from the SVN-tracking branch back to master ---
	gitutils.Merge(man.getBridgeRepoPath(), branchName)

	// fetch changes to central repo master from SVN bridge master (note that cannot just
	//`git push git-central-repo master` as that would trigger the central repo update hook and deadlock)
	gitutils.Fetch(man.getCentralRepoPath(), "bridge", branchName)
}

// Checkout bridge repo detached head
func (man *Manager) checkoutToOriginRef(ref string) {
	worktree, err := man.bridgeRepo.Worktree()
	if err != nil {
		panic(fmt.Errorf("could not get worktree for repo '%s': %w", man.getBridgeRepoPath(), err))
	}

	branchName := gitutils.GetBranchName(ref)
	if branchName == "master" {
		branchName = "trunk"
	}

	err = worktree.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewRemoteReferenceName("origin", branchName),
	})

	if err != nil {
		panic(fmt.Errorf("could not checkout '%s' for repo '%s': %w", branchName, man.getBridgeRepoPath(), err))
	}
}

// Checkout bridge repo to local branch, corresponding the <ref>. Branch will be created if needed
func (man *Manager) checkoutToRef(ref string) {
	refExists := gitutils.IsRefExists(man.bridgeRepo, ref)

	worktree, err := man.bridgeRepo.Worktree()
	if err != nil {
		panic(fmt.Errorf("could not get worktree for repo '%s': %w", man.getBridgeRepoPath(), err))
	}

	branchName := gitutils.GetBranchName(ref)
	if !refExists {
		//checkout to origin branch first
		branch := plumbing.NewRemoteReferenceName("origin", branchName)
		err = worktree.Checkout(&git.CheckoutOptions{Branch: branch})
		if err != nil {
			panic(fmt.Errorf("could not checkout origin '%s' for repo '%s': %w", branchName, man.getBridgeRepoPath(), err))
		}
	}

	err = worktree.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName(gitutils.GetBranchName(ref)),
		Create: !refExists,
	})

	if err != nil {
		panic(fmt.Errorf("could not checkout '%s' for repo '%s': %w", branchName, man.getBridgeRepoPath(), err))
	}
}

func (man *Manager) getAllRefs() []string {
	references, err := man.bridgeRepo.References()
	if err != nil {
		panic(fmt.Errorf("could not get refs for repo %s: %w", man.getBridgeRepoPath(), err))
	}

	var allRefs []string
	err = references.ForEach(func(ref *plumbing.Reference) error {
		refName := ref.Name().String()
		if strings.HasPrefix(refName, "refs/remotes/origin/") {
			refName = strings.Replace(refName, "refs/remotes/origin/", "refs/heads/", 1)
			refName = strings.Replace(refName, "/trunk", "/master", 1)
			allRefs = append(allRefs, refName)
		}
		return nil
	})

	if err != nil {
		panic(fmt.Errorf("an error occurred while getting refs for repo %s: %w", man.getBridgeRepoPath(), err))
	}

	return allRefs
}

func (man *Manager) getSysUser() usr.User {
	return store.GetUser(man.repo, conf.GetConfig().SystemGitUserName)
}

func (man *Manager) getBridgeRepoPath() string {
	return man.repo.GetBridgeRepoPath()
}

func (man *Manager) getCentralRepoPath() string {
	return man.repo.GetGitRepoPath()
}
