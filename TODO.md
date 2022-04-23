Sync, вроде, работает в простом варианте. 
сделать механизм обработки конфликтов

А еще обработать: Transaction is out of date: File '/Prj/trunk/main.txt' is out of date": git svn --authors-file="$HOME/git/git-svn-bridge-authors" fetch




Committing to svn://localhost/trunk ...
M       main.txt

ERROR from SVN:
Transaction is out of date: File '/trunk/main.txt' is out of date
W: 927df16f68d929337aa79a02a0c2bf2efd2892d8 and refs/remotes/origin/trunk differ, using rebase:
:100644 100644 387d108502b3ec1a8b868c60831071977c83d258 77e9aa9c69c1a23ff30b577dac1be62b392d9ed9 M      main.txt
Successfully rebased and updated detached HEAD.
ERROR: Not all changes have been committed into SVN, however the committed
ones (if any) seem to be successfully integrated into the working tree.
Please see the above messages for details.


Exit code: 1


---------------------------------------
Автослияние main.txt
КОНФЛИКТ (содержимое): Конфликт слияния в main.txt
Не удалось провести автоматическое слияние; исправьте конфликты и сделайте коммит результата.

Error: could not sync ref refs/heads/master for repo 'Prj': could not push changes to SVN for repo 'Prj'(refs/heads/master): could not merge (--no-ff) for repo 'repos/bridge/Prj': exit status 1
Usage:
git-svn-bridge sync <ref1>..<refN>  [flags]

Flags:
-h, --help          help for sync
-r, --repo string   repository name

could not sync ref refs/heads/master for repo 'Prj': could not push changes to SVN for repo 'Prj'(refs/heads/master): could not merge (--no-ff) for repo 'repos/bridge/Prj': exit status 1
