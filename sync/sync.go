package sync

import (
	"errors"
	"fmt"
	"git-svn-bridge/conf"
	"git-svn-bridge/repo"
	"git-svn-bridge/store"
	"git-svn-bridge/usr"
	"git-svn-bridge/vcs/gitsvn"
	"git-svn-bridge/vcs/gitutils"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"strings"
)

type Manager struct {
	repo *repo.Repo

	bridgeRepo  *git.Repository
	centralRepo *git.Repository

	isSvnFetched bool
}

func New(repoName string) (*Manager, error) {
	if !store.HasRepo(repoName) {
		return nil, errors.New("There is no repoInfo with name " + repoName)
	}

	repoInfo := store.GetRepo(repoName)

	bridgeRepoPath := repoInfo.GetBridgeRepoPath()
	centralRepoPath := repoInfo.GetGitRepoPath()

	bridgeRepo, err := git.PlainOpen(bridgeRepoPath)
	if err != nil {
		return nil, fmt.Errorf("could not open bridge repoInfo '%s': %w", repoName, err)
	}

	centralRepo, err := git.PlainOpen(centralRepoPath)
	if err != nil {
		return nil, fmt.Errorf("could not open central repoInfo '%s': %w", repoName, err)
	}

	return &Manager{repo: &repoInfo, bridgeRepo: bridgeRepo, centralRepo: centralRepo, isSvnFetched: false}, nil
}

func (man *Manager) SyncAllRefs() error {
	err := man.fetchSvnChanges()
	if err != nil {
		return fmt.Errorf("could not fetch SVN changes for repo '%s': %w", man.repo.GetName(), err)
	}

	allRefs, err := man.getAllRefs()
	if err != nil {
		return fmt.Errorf("could not list all refs for repo '%s': %w", man.repo.GetName(), err)
	}

	return man.SyncRefs(allRefs)
}

func (man *Manager) SyncRefs(refs []string) error {
	for _, ref := range refs {
		err := man.syncRef(ref)
		if err != nil {
			return fmt.Errorf("could not sync ref %s for repo '%s': %w", ref, man.repo.GetName(), err)
		}
	}

	return nil
}

func (man *Manager) syncRef(ref string) error {
	err := man.fetchSvnChanges()
	if err != nil {
		return fmt.Errorf("could not fetch SVN changes for repo '%s'(%s): %w", man.repo.GetName(), ref, err)
	}

	err = man.fetchGitChanges(ref)
	if err != nil {
		return fmt.Errorf("could not fetch GIT changes for repo '%s'(%s): %w", man.repo.GetName(), ref, err)
	}

	err = man.pushRefToSvn(ref)
	if err != nil {
		return fmt.Errorf("could not push changes to SVN for repo '%s'(%s): %w", man.repo.GetName(), ref, err)
	}

	err = man.pushRefToGit(ref)
	if err != nil {
		return fmt.Errorf("could not push changes to GIT for repo '%s'(%s): %w", man.repo.GetName(), ref, err)
	}

	return nil
}

func (man *Manager) fetchSvnChanges() error {
	if man.isSvnFetched {
		return nil
	}
	man.isSvnFetched = true

	gitSvnExecutor := gitsvn.CreateExecutor(man.getSysUser())
	return gitSvnExecutor.Fetch(man.getBridgeRepoPath())
}

func (man *Manager) fetchGitChanges(ref string) error {
	err := man.checkoutToRef(ref)
	if err != nil {
		return fmt.Errorf("could not checkout bridge repo '%s' to ref %s: %w", man.repo.GetName(), ref, err)
	}

	centralRepoRefExists, err := gitutils.IsRefExists(man.centralRepo, ref)
	if err != nil {
		return err
	}

	if centralRepoRefExists {
		branchName := gitutils.GetBranchName(ref)
		return gitutils.PullAndRebase(man.getBridgeRepoPath(), gitutils.GitCentralRepoName, branchName)
	}

	return nil
}

func (man *Manager) pushRefToSvn(ref string) error {
	// store the SVN URL and author of the last commit for `git svn dcommit`
	gitUserName, err := gitutils.GetGitAuthor(man.repo, man.getBridgeRepoPath())
	if err != nil {
		return err
	}

	err = man.checkoutToOriginRef(ref)
	if err != nil {
		return fmt.Errorf("could not checkout detached head for repo '%s'(%s): %w", man.repo.GetName(), ref, err)
	}

	branchName := gitutils.GetBranchName(ref)

	// create a merged log message for the merge commit to SVN. Note that we squash log messages together for
	// merge commits. Experiment with the fomat, e.g. '%ai | [%an] %s' etc
	commitMessage, err := gitutils.BuildCommitMessage(man.getBridgeRepoPath(), branchName)
	if err != nil {
		return err
	}

	// Merge changes from master to the SVN-tracking branch and commit to SVN.
	// Note that we always record the merge with --no-ff
	err = gitutils.MergeNoFF(man.getBridgeRepoPath(), commitMessage, branchName)
	if err != nil {
		return err
	}

	//commit changes to SVN
	commitUser := store.GetUser(man.repo, gitUserName)
	gitSvnExecutor := gitsvn.CreateExecutor(commitUser)
	err = gitSvnExecutor.Commit(man.getBridgeRepoPath())
	if err != nil {
		panic(err)
		//TODO надо различать разные ошибки...

		/*
			# if "Transaction is out of date: File '/Prj/trunk/main.txt' is out of date": git svn --authors-file="$HOME/git/git-svn-bridge-authors" fetch
			# if  КОНФЛИКТ (содержимое): Конфликт слияния в trunk/main.txt: git rebase --abort git merge --abort
			# git checkout svn/git-svn
			# git branch bridge-branch-1 sha-1
			# git checkout bridge-branch-1
			# git merge svn/git-svn
			#local CENTRAL_REPO_PATH="`git remote -v show | awk 'NR > 1 { exit }; { print $2 };'`"
			#pushd "$CENTRAL_REPO_PATH"
			#git fetch svn-bridge bridge-branch-1:bridge-branch-1
			#popd
		*/
	}

	return nil
}

func (man *Manager) pushRefToGit(ref string) error {
	err := man.checkoutToRef(ref)
	if err != nil {
		return err
	}

	branchName := gitutils.GetBranchName(ref)

	// --- merge changes from the SVN-tracking branch back to master ---
	err = gitutils.Merge(man.getBridgeRepoPath(), branchName)
	if err != nil {
		return err
	}

	// fetch changes to central repo master from SVN bridge master (note that cannot just
	//`git push git-central-repo master` as that would trigger the central repo update hook and deadlock)
	return gitutils.Fetch(man.getCentralRepoPath(), "bridge", branchName)
}

// Checkout bridge repo detached head
func (man *Manager) checkoutToOriginRef(ref string) error {
	worktree, err := man.bridgeRepo.Worktree()
	if err != nil {
		return err
	}

	branchName := gitutils.GetBranchName(ref)
	if branchName == "master" {
		branchName = "trunk"
	}

	return worktree.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewRemoteReferenceName("origin", branchName),
	})
}

// Checkout bridge repo to local branch, corresponding the <ref>. Branch will be created if needed
func (man *Manager) checkoutToRef(ref string) error {
	refExists, err := gitutils.IsRefExists(man.bridgeRepo, ref)
	if err != nil {
		return err
	}

	worktree, err := man.bridgeRepo.Worktree()
	if err != nil {
		return err
	}

	if !refExists {
		//checkout to origin branch first
		branch := plumbing.NewRemoteReferenceName("origin", gitutils.GetBranchName(ref))
		err = worktree.Checkout(&git.CheckoutOptions{Branch: branch})
		if err != nil {
			return err
		}
	}

	return worktree.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName(gitutils.GetBranchName(ref)),
		Create: !refExists,
	})
}

func (man *Manager) getAllRefs() ([]string, error) {
	references, err := man.bridgeRepo.References()
	if err != nil {
		return nil, err
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

	return allRefs, err
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
