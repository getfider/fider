# https://medium.com/@kcmueller/delete-local-git-branches-that-were-deleted-on-remote-repository-b596b71b530c
git fetch -p
git branch -vv | grep ' gone]' | awk '{print $1}' | xargs git branch -D