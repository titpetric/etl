select name,commit_sha from github_tags where created_at IS NULL or created_at=''
