#!/bin/bash

# Please run this script with a test user account.
#
# NOTE THAT THE SCRIPT CHANGES SVN CONFIGURATION AS FOLLOWS:
#
#  'store-plaintext-passwords = yes'
#
# IN ~/.subversion/servers
#
# AND GIT CONFIGURATION AS FOLLOWS:
#
#  git config --global user.name "Git-SVN Bridge (GIT SIDE)"
#  git config --global user.email "git-svn-bridge@company.com"

set -u

# assure that .subversion is created
svn info /tmp

set -e
set -x

mkdir {bin,git,svn,test}

# === Create SVN repo ===

STORE_PASSWD='store-plaintext-passwords = yes'
grep -qx "$STORE_PASSWD" ~/.subversion/servers \
	|| echo "$STORE_PASSWD" >> ~/.subversion/servers

cd ~/svn
svnadmin create svn-repo
svn co file://`pwd`/svn-repo svn-checkout

cd svn-checkout
mkdir -p trunk/src
echo 'int main() { return 0; }' > trunk/src/main.cpp
svn add trunk
svn ci -m "First commit."

# svnserve

SVNSERVE_PIDFILE="$HOME/svn/svnserve.pid"
SVNSERVE_LOGFILE="$HOME/svn/svnserve.log"
SVNSERVE_CONFFILE="$HOME/svn/svnserve.conf"
SVNSERVE_USERSFILE="$HOME/svn/svnserve.users"

>> $SVNSERVE_LOGFILE

cat > "$SVNSERVE_CONFFILE" << EOT
[general]
realm = git-SVN test
anon-access = none
password-db = $SVNSERVE_USERSFILE
EOT

cat > "$SVNSERVE_USERSFILE" << EOT
[users]
git-svn-bridge = git-svn-bridge
alice = alice
bob = bob
carol = carol
dave = dave
EOT

TAB="`printf '\t'`"

cat > ~/svn/Makefile << EOT
svnserve-start:
${TAB}svnserve -d \\
${TAB}${TAB}--pid-file "$SVNSERVE_PIDFILE" \\
${TAB}${TAB}--log-file "$SVNSERVE_LOGFILE" \\
${TAB}${TAB}--config-file "$SVNSERVE_CONFFILE" \\
${TAB}${TAB}-r ~/svn/svn-repo
svnserve-stop:
${TAB}kill \`cat "$SVNSERVE_PIDFILE"\`
EOT

make -f ~/svn/Makefile svnserve-start

# === SVN repo created ===

# === Create Git and Git-Bridge repos ===

git config --global user.name "Git-SVN Bridge (GIT SIDE)"
git config --global user.email "git-svn-bridge@company.com"

cd ~/git
git init --bare git-central-repo-trunk.git
cd git-central-repo-trunk.git
git remote add svn-bridge ../git-svn-bridge-trunk

SVN_REPO_URL="svn://localhost/trunk"
cd ~/git
git svn init --prefix=svn/ $SVN_REPO_URL git-svn-bridge-trunk
cd git-svn-bridge-trunk
AUTHORS='/tmp/git-svn-bridge-authors'
echo 'git-svn-bridge = Git SVN Bridge <git-svn-bridge@company.com>' > $AUTHORS
echo -e "\n>>> USE 'git-svn-bridge' AS PASSWORD <<<\n"
git svn fetch --authors-file="$AUTHORS" --log-window-size 10000

git branch -a -v

git remote add git-central-repo ../git-central-repo-trunk.git
git push --all git-central-repo

cd ~/git
git clone git-central-repo-trunk.git git-central-repo-clone
cd git-central-repo-clone
git log

cd ~/git/git-central-repo-trunk.git
cat > hooks/update << 'EOT'
#!/bin/bash
set -u
refname=$1
shaold=$2
shanew=$3

# we are only interested in commits to master
[[ "$refname" = "refs/heads/master" ]] || exit 0
# don't allow non-fast-forward commits
if [[ $(git merge-base "$shanew" "$shaold") != "$shaold" ]]; then
    echo "Non-fast-forward commits to master are not allowed"
    exit 1
fi
EOT

chmod +x hooks/update

# === Git and Git-Bridge repos created ===

# === Test manual synchronization
# Commit to Git repo
cd ~/git/git-central-repo-clone
echo 'int test() { return 0; }' > ./src/test.cpp
git add src/test.cpp
git commit -m "Second commit. (git side)"
git push

# Fetch updates from Git to Git-Bridge
cd ~/git/git-svn-bridge-trunk
git checkout master
git pull --rebase git-central-repo master

# Merge new commits to remote branch, related to SVN
git checkout svn/git-svn
git merge --no-ff --no-log -m "Second commit (git side)" master

# Commit to SVN
git svn dcommit --username git-svn-bridge

# Fetch updates from SVN to Git-Bridge
git checkout master
git merge svn/git-svn

# Fetch updates from Git-Bridge to Git
# (note that cannot just `git push git-central-repo master` as that would trigger the central repo update hook and deadlock)
cd ~/git/git-central-repo-trunk.git
git fetch svn-bridge master:master

# Check commits in SVN
cd ~/svn/svn-checkout
svn up
svn log

# Check commits in GIT
cd ~/git/git-central-repo-clone/
git pull
git log
