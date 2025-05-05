BRANCH := $(shell git rev-parse --abbrev-ref HEAD)

branch:
	git checkout $(ARGS) > /dev/null 2>&1 || git checkout -b $(ARGS)

stash:
	git stash save --keep-index --include-untracked

unstash:
	git stash apply
  
push:
	git push origin $(BRANCH)

push!:
	git push --force origin $(BRANCH)

pull:
	git pull origin $(BRANCH)

pull!:
	git pull origin $(BRANCH) --rebase

trunk:
	git checkout master

uncommit:
	git reset --soft HEAD^

unmerge:
	git merge --abort
  
unrebase:
	git rebase --abort


merge-to:
	@$(eval current_branch := $(BRANCH))
	git checkout $(ARGS)
	git merge $(current_branch) --no-edit
	git push origin $(ARGS)
	git checkout $(current_branch)

remote:
	git remote -v