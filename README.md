# GitLabCommand - glc - README

## How to use ?
1. ensure glc.yaml is configured with the valid token and url

   ```yaml
	gitlab_token: <your valid token>
   	gitlab_url: https://gitlab.com
   	default_protected_branches:
	  - name: develop*
	    push: No One
	    merge: Developers + Maintainers
	  - name: release*
	    push: No One
	    merge: Maintainers
	  - name: master
	    push: No One
	    merge: Maintainers
	```

2. Ensure golang is installed

   ```
   brew install go
   ```

3. Run glc.go to list out all the options
   
   ```
   go run glc.go -h
   ```

4. Example use cases:

	```
	go run glc.go -findgrp common   	# Search and List out gitlab group with "common" string
	go run glc.go -ls 332               # List all projects in group id 332
	go run glc.go -setpb -gid 332       # Set protected branch based on glc.yaml default_protected_branches setting in all project in group 332
	go run glc.go -tagcsv uat-tag.csv   # tag all the repo listed in the uat-tag.csv file.
	```
