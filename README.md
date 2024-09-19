# GitLabCommand - glc - README

This is design to be a simple go script to perform some action on gitlab as describe in the help.

```
$ go run glc.go --help
Usage of glc:
  -bulkmr string
        CSV file with bulk merge request. CSV Header: pid,path,source,target,title
  -debug
        Enable verbose logging
  -default string
        Set the default branch, 2nd parameter is -gid or -pid
  -dumpcfg
        Dump the configurations
  -findgrp string
        Find the group. Example: -findgrp Common
  -forcetag
        depends on -tagcsv. This will force tag. If tag exist, it will delete and recreate
  -gid string
        The group Id to perform an action
  -ls string
        List all the repo in a Group
  -lspb
        List all the protected branch settings in a group, 2nd parameter is -gid or -pid
  -pid string
        The project id to perform an action
  -setpb
        Set branch(es) to be protected by the default configuration, 2nd parameter is -gid or -pid
  -setptag
        Set tag(s) to be protected by the default configuration, 2nd parameter is -gid or -pid
  -tagcsv string
        CSV file with tagging information. The tag will suffix with yymmdd-hash. CSV Header: pid,path,prefix,branch,message
  -tagnosuffix
        depends on -tagcsv and this will disable the suffix of yymmdd-hash
```

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
	go run glc.go -findgrp common       # Search and List out gitlab group with "common" string
	go run glc.go -ls 332               # List all projects in group id 332
	go run glc.go -setpb -gid 332       # Set protected branch based on glc.yaml default_protected_branches setting in all project in group 332
	go run glc.go -tagcsv uat-tag.csv   # tag all the repo listed in the uat-tag.csv file.
	```
