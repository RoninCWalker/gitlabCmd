package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	url2 "net/url"
	"os"
	"runtime"
	"strconv"
	"time"

	"gopkg.in/yaml.v3"
)

//*****************************************************************************
// Structs, Enum and Method for the Configurations
//*****************************************************************************

type config struct {
	GitlabToken              string                  `yaml:"gitlab_token"`
	GitlabUrl                string                  `yaml:"gitlab_url"`
	DefaultProtectedBranches []protectedBranchConfig `yaml:"default_protected_branches"`
}

type protectedBranchConfig struct {
	Name  string         `yaml:"name"`
	Push  accessDescEnum `yaml:"push"`
	Merge accessDescEnum `yaml:"merge"`
}

func (in *config) readConfig() *config {
	data, err := ioutil.ReadFile("./glc.yaml")
	handleError(err)
	err = yaml.Unmarshal(data, in)
	handleError(err)
	return in
}
func (in *config) ToString() string {
	result, err := yaml.Marshal(in)
	handleError(err)
	return string(result)
}

//*****************************************************************************
// Structs for REST
//*****************************************************************************
type gitlabGroup struct {
	Id     int    `json:"id"`
	WebUrl string `json:"web_url"`
	Name   string `json:"name"`
	Path   string `json:"full_path"`
}

type gitlabProject struct {
	gitlabGroup
	HttpUrlToRepo     string `json:"http_url_to_repo"`
	DefaultBranch     string `json:"default_branch"`
	PathWithNamespace string `json:"path_with_namespace"`
}

type gitlabProtectedBranch struct {
	Name              string               `json:"name"`
	PushAccessLevels  []gitlabBranchAccess `json:"push_access_levels"`
	MergeAccessLevels []gitlabBranchAccess `json:"merge_access_levels"`
}

type gitlabBranchAccess struct {
	AccessLevel int    `json:"access_level"`
	Description string `json:"access_level_description"`
}

type gitlabBranch struct {
	Name               string `json:"name"`
	Merged             bool   `json:"merged"`
	Protected          bool   `json:"protected"`
	Default            bool   `json:"default"`
	DevelopersCanPush  bool   `json:"developers_can_push"`
	DevelopersCanMerge bool   `json:"developers_can_merge"`
	CanPush            bool   `json:"can_push"`
	WebURL             string `json:"web_url"`
	Commit             struct {
		AuthorEmail    string   `json:"author_email"`
		AuthorName     string   `json:"author_name"`
		AuthoredDate   string   `json:"authored_date"`
		CommittedDate  string   `json:"committed_date"`
		CommitterEmail string   `json:"committer_email"`
		CommitterName  string   `json:"committer_name"`
		ID             string   `json:"id"`
		ShortID        string   `json:"short_id"`
		Title          string   `json:"title"`
		Message        string   `json:"message"`
		ParentIds      []string `json:"parent_ids"`
	} `json:"commit"`
}

type gitlabMergeRequest struct {
	ID           int       `json:"id"`
	Iid          int       `json:"iid"`
	ProjectID    int       `json:"project_id"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	State        string    `json:"state"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	TargetBranch string    `json:"target_branch"`
	SourceBranch string    `json:"source_branch"`
	Upvotes      int       `json:"upvotes"`
	Downvotes    int       `json:"downvotes"`
	Author       struct {
		ID        int         `json:"id"`
		Name      string      `json:"name"`
		Username  string      `json:"username"`
		State     string      `json:"state"`
		AvatarURL interface{} `json:"avatar_url"`
		WebURL    string      `json:"web_url"`
	} `json:"author"`
	Assignee struct {
		ID        int         `json:"id"`
		Name      string      `json:"name"`
		Username  string      `json:"username"`
		State     string      `json:"state"`
		AvatarURL interface{} `json:"avatar_url"`
		WebURL    string      `json:"web_url"`
	} `json:"assignee"`
	SourceProjectID int      `json:"source_project_id"`
	TargetProjectID int      `json:"target_project_id"`
	Labels          []string `json:"labels"`
	Draft           bool     `json:"draft"`
	WorkInProgress  bool     `json:"work_in_progress"`
	Milestone       struct {
		ID          int       `json:"id"`
		Iid         int       `json:"iid"`
		ProjectID   int       `json:"project_id"`
		Title       string    `json:"title"`
		Description string    `json:"description"`
		State       string    `json:"state"`
		CreatedAt   time.Time `json:"created_at"`
		UpdatedAt   time.Time `json:"updated_at"`
		DueDate     string    `json:"due_date"`
		StartDate   string    `json:"start_date"`
		WebURL      string    `json:"web_url"`
	} `json:"milestone"`
	MergeWhenPipelineSucceeds bool        `json:"merge_when_pipeline_succeeds"`
	MergeStatus               string      `json:"merge_status"`
	MergeError                interface{} `json:"merge_error"`
	Sha                       string      `json:"sha"`
	MergeCommitSha            interface{} `json:"merge_commit_sha"`
	SquashCommitSha           interface{} `json:"squash_commit_sha"`
	UserNotesCount            int         `json:"user_notes_count"`
	DiscussionLocked          interface{} `json:"discussion_locked"`
	ShouldRemoveSourceBranch  bool        `json:"should_remove_source_branch"`
	ForceRemoveSourceBranch   bool        `json:"force_remove_source_branch"`
	AllowCollaboration        bool        `json:"allow_collaboration"`
	AllowMaintainerToPush     bool        `json:"allow_maintainer_to_push"`
	WebURL                    string      `json:"web_url"`
	References                struct {
		Short    string `json:"short"`
		Relative string `json:"relative"`
		Full     string `json:"full"`
	} `json:"references"`
	TimeStats struct {
		TimeEstimate        int         `json:"time_estimate"`
		TotalTimeSpent      int         `json:"total_time_spent"`
		HumanTimeEstimate   interface{} `json:"human_time_estimate"`
		HumanTotalTimeSpent interface{} `json:"human_total_time_spent"`
	} `json:"time_stats"`
	Squash       bool   `json:"squash"`
	Subscribed   bool   `json:"subscribed"`
	ChangesCount string `json:"changes_count"`
	MergedBy     struct {
		ID        int    `json:"id"`
		Name      string `json:"name"`
		Username  string `json:"username"`
		State     string `json:"state"`
		AvatarURL string `json:"avatar_url"`
		WebURL    string `json:"web_url"`
	} `json:"merged_by"`
	MergedAt                    time.Time   `json:"merged_at"`
	ClosedBy                    interface{} `json:"closed_by"`
	ClosedAt                    interface{} `json:"closed_at"`
	LatestBuildStartedAt        time.Time   `json:"latest_build_started_at"`
	LatestBuildFinishedAt       time.Time   `json:"latest_build_finished_at"`
	FirstDeployedToProductionAt interface{} `json:"first_deployed_to_production_at"`
	Pipeline                    struct {
		ID     int    `json:"id"`
		Sha    string `json:"sha"`
		Ref    string `json:"ref"`
		Status string `json:"status"`
		WebURL string `json:"web_url"`
	} `json:"pipeline"`
	DiffRefs struct {
		BaseSha  string `json:"base_sha"`
		HeadSha  string `json:"head_sha"`
		StartSha string `json:"start_sha"`
	} `json:"diff_refs"`
	DivergedCommitsCount int `json:"diverged_commits_count"`
	TaskCompletionStatus struct {
		Count          int `json:"count"`
		CompletedCount int `json:"completed_count"`
	} `json:"task_completion_status"`
}

type gitlabError struct {
	Error string `json:"error"`
}

// Access Level Types
type accessLevelEnum int
type accessDescEnum string

const (
	NO_ONE_LEVEL                 accessLevelEnum = 0
	DEVELOPERS_MAINTAINERS_LEVEL                 = 30
	MAINTAINERS_LEVEL                            = 40
	NO_ONE_DESC                  accessDescEnum  = "No One"
	DEVELOPERS_MAINTAINERS_DESC                  = "Developers + Maintainers"
	MAINTAINERS_DESC                             = "Maintainers"
)

var accessDescMap = map[accessDescEnum]accessLevelEnum{
	NO_ONE_DESC:                 NO_ONE_LEVEL,
	DEVELOPERS_MAINTAINERS_DESC: DEVELOPERS_MAINTAINERS_LEVEL,
	MAINTAINERS_DESC:            MAINTAINERS_LEVEL,
}

//*****************************************************************************
// Common Functions
//*****************************************************************************
func handleError(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func trace() string {
	pc := make([]uintptr, 15)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	return fmt.Sprintf("%s:%d %s", frame.File, frame.Line, frame.Function)
}

type httpMethodEnum string

const (
	GET  httpMethodEnum = "GET"
	POST httpMethodEnum = "POST"
	PUT  httpMethodEnum = "PUT"
)

func doHttpRequest(method httpMethodEnum, url string, body io.Reader) (*http.Response, error) {
	log.Printf("%s - param: %s,%s,%v", trace(), method, url, body)
	client := &http.Client{}
	req, err := http.NewRequest(string(method), url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("PRIVATE-TOKEN", conf.GitlabToken)
	res, err := client.Do(req)
	return res, err
}

//*****************************************************************************
// functions for Flags
//*****************************************************************************

// Print out the list of group found in the search (string) string.
func findGroup(search string) {
	log.Printf("%s - param: %s", trace(), search)
	path := "/api/v4/groups?search="
	url := conf.GitlabUrl + path + search
	res, err := doHttpRequest(GET, url, nil)
	handleError(err)
	body, err := io.ReadAll(res.Body)
	handleError(err)
	defer res.Body.Close()
	var groups []gitlabGroup
	json.Unmarshal(body, &groups)
	for _, g := range groups {
		fmt.Printf("%4d, %s, %s, %s\n", g.Id, g.Path, g.Name, g.WebUrl)
	}
}

// Set default branch on a project based on Project Id (pid) or all
// project in a Group Id (gid)
func setDefaultBranch(defaultBranch string, pid string, gid string) {
	log.Printf("%s - param:  %s, %s, %s", trace(), defaultBranch, pid, gid)
	var projects []gitlabProject = getProjectList(pid, gid)
	for _, p := range projects {
		url := fmt.Sprintf("%s/api/v4/projects/%d?default_branch=%s",
			conf.GitlabUrl,
			p.Id,
			defaultBranch)
		res, err := doHttpRequest(PUT, url, nil)
		handleError(err)
		fmt.Printf("%d, %s, %s, %s\n", p.Id, p.HttpUrlToRepo, defaultBranch, res.Status)
	}
}

// Set proctected branch based on the default_protected_branches
// based on Project Id (pid) or Group Id (gid)
func setProtectedBranch(pid string, gid string) {
	log.Printf("%s - param: %s,%s", trace(), pid, gid)
	var projects []gitlabProject = getProjectList(pid, gid)
	for _, p := range projects {
		for _, c := range conf.DefaultProtectedBranches {
			url := conf.GitlabUrl +
				fmt.Sprintf("/api/v4/projects/%d/protected_branches?name=%s&push_access_level=%d&merge_access_level=%d",
					p.Id, url2.QueryEscape(c.Name), accessDescMap[c.Push], accessDescMap[c.Merge])
			res, err := doHttpRequest(POST, url, nil)
			handleError(err)
			fmt.Printf("%d, %s, %s, %s\n", p.Id, p.Name, c.Name, res.Status)
		}
	}
}

// Returns list of Project Id in String from PID and GID
func getProjectList(pid string, gid string) []gitlabProject {
	log.Printf("%s - param: %s, %s", trace(), pid, gid)
	if len(pid)+len(gid) == 0 {
		flag.PrintDefaults()
		log.Fatalln("-repo or -grp is not defined")
	}
	var projects []gitlabProject
	if len(gid) > 0 {
		projects = listProjects(gid, false)
	}
	if len(pid) > 0 {
		prj := getProject(pid)
		projects = append(projects, prj)
	}
	return projects
}

// List protected branch in STDOUT based on Project Id (pid) or Group Id (god)
func listProtectedBranch(pid string, gid string) {
	log.Printf("%s - param: %s, %s", trace(), pid, gid)
	var projects []gitlabProject = getProjectList(pid, gid)

	for _, p := range projects {
		url := conf.GitlabUrl + "/api/v4/projects/" + strconv.Itoa(p.Id) + "/protected_branches"
		res, err := doHttpRequest(GET, url, nil)
		handleError(err)
		body, err := io.ReadAll(res.Body)
		handleError(err)
		defer res.Body.Close()
		var protectedBranches []gitlabProtectedBranch
		err = json.Unmarshal(body, &protectedBranches)
		handleError(err)
		for _, b := range protectedBranches {
			for _, m := range b.MergeAccessLevels {
				fmt.Printf("%4d, %s, merge: %s : %d: %s\n", p.Id, p.PathWithNamespace, b.Name, m.AccessLevel, m.Description)
			}
			for _, a := range b.PushAccessLevels {
				fmt.Printf("%4d, %s, push: %s - %d: %s\n", p.Id, p.PathWithNamespace, b.Name, a.AccessLevel, a.Description)
			}
		}
	}
}

// Get project details by Project Id (pid)
func getProject(pid string) gitlabProject {
	log.Printf("%s - param: %s", trace(), pid)
	url := conf.GitlabUrl + "/api/v4/projects/" + pid
	res, err := doHttpRequest(GET, url, nil)

	handleError(err)
	body, err := io.ReadAll(res.Body)
	handleError(err)
	defer res.Body.Close()
	var result gitlabProject
	json.Unmarshal(body, &result)
	return result
}

// Return list  ofall the projects in the Groupd Id (gid).
// If print = true, print the result to stdout
func listProjects(gid string, print bool) []gitlabProject {
	log.Printf("%s - param: %s, %t", trace(), gid, print)
	path := "/api/v4/groups/" + gid + "/projects"
	url := conf.GitlabUrl + path
	prj, err := doHttpRequest(GET, url, nil)
	handleError(err)
	body, err := io.ReadAll(prj.Body)
	handleError(err)
	prj.Body.Close()
	var prjs []gitlabProject
	json.Unmarshal(body, &prjs)
	if print && len(prjs) > 0 {
		for _, p := range prjs {
			fmt.Printf("%d, %s, %s\n", p.Id, p.PathWithNamespace, p.HttpUrlToRepo)

		}
	}
	return prjs
}

// Return the 8 character commit hash from Project Id (pid) and the branch (branch)
func getProjectBranchCommitHash(pid string, branch string) (string, error) {
	log.Printf("%s - param: %s, %s", trace(), pid, branch)
	url := conf.GitlabUrl + "/api/v4/projects/" + pid + "/repository/branches/" + url2.QueryEscape(branch)
	res, err := doHttpRequest(GET, url, nil)
	handleError(err)
	if res.StatusCode >= 400 {
		return "", fmt.Errorf("%s - %s", branch, res.Status)
	}
	body, err := io.ReadAll(res.Body)
	handleError(err)
	defer res.Body.Close()
	var result gitlabBranch
	json.Unmarshal(body, &result)
	return result.Commit.ShortID, nil
}

// CSV columns
type tagCSVRecord struct {
	Pid     string
	Path    string
	Prefix  string
	Branch  string
	Message string
}

// Tag the repo based on the data in the CSV File (csvfile) with suffix yymmdd-hash.
// If nosuffix is true, then the tag will not have the suffix.
func tagCSV(csvfile string, nosuffix bool) {
	log.Printf("%s - param: %s", trace(), csvfile)
	var csvEntries []tagCSVRecord
	var tag string

	// open file
	f, err := os.Open(csvfile)
	handleError(err)
	defer f.Close()

	// parse CSV
	lines, err := csv.NewReader(f).ReadAll()
	handleError(err)

	for _, line := range lines {
		data := tagCSVRecord{
			Pid:     line[0],
			Path:    line[1],
			Prefix:  line[2],
			Branch:  line[3],
			Message: line[4],
		}
		csvEntries = append(csvEntries, data)
	}
	today := time.Now()
	todayStr := fmt.Sprintf("%02d%02d%02d", today.Year(), today.Month(), today.Day())
	todayStr = todayStr[2:]
	for _, project := range csvEntries[1:] {
		projectHash, err := getProjectBranchCommitHash(project.Pid, project.Branch)
		if err != nil {
			fmt.Printf("%s, %s, %s : %s\n", project.Pid, project.Path, project.Branch, err.Error())
		} else {
			if nosuffix {
				tag = project.Prefix
			} else {
				tag = fmt.Sprintf("%s-%s-%s", project.Prefix, todayStr, projectHash)
			}
			url := fmt.Sprintf("%s/api/v4/projects/%s/repository/tags?tag_name=%s&ref=%s&message=%s",
				conf.GitlabUrl, project.Pid, tag, url2.QueryEscape(project.Branch), url2.QueryEscape(project.Message))
			res, err := doHttpRequest(POST, url, nil)
			handleError(err)
			fmt.Printf("%s, %s/%s/-/tags/%s, %s : %s\n", project.Pid, conf.GitlabUrl, project.Path, tag, project.Branch, res.Status)
		}

	}
}

type mergeRequestCSVRecord struct {
	Pid    string
	Path   string
	Source string
	Target string
	Title  string
}

func bulkMergeRequest(csvfile string) {
	log.Printf("%s - param: %s", trace(), csvfile)
	var csvEntries []mergeRequestCSVRecord

	// open file
	f, err := os.Open(csvfile)
	handleError(err)
	defer f.Close()

	// parse CSV
	lines, err := csv.NewReader(f).ReadAll()
	handleError(err)

	for _, line := range lines {
		data := mergeRequestCSVRecord{
			Pid:    line[0],
			Path:   line[1],
			Source: line[2],
			Target: line[3],
			Title:  line[4],
		}
		csvEntries = append(csvEntries, data)
	}

	for _, project := range csvEntries[1:] {
		url := fmt.Sprintf("%s/api/v4/projects/%s/merge_requests?source_branch=%s&target_branch=%s&title=%s",
			conf.GitlabUrl, project.Pid, url2.QueryEscape(project.Source), url2.QueryEscape(project.Target), url2.QueryEscape(project.Title))
		res, err := doHttpRequest(POST, url, nil)
		handleError(err)
		body, err := io.ReadAll(res.Body)
		handleError(err)
		defer res.Body.Close()
		var result gitlabMergeRequest
		var errorResponse gitlabError
		json.Unmarshal(body, &result)
		json.Unmarshal(body, &errorResponse)

		// set the status message
		statusMessage := errorResponse.Error
		if len(statusMessage) == 0 {
			statusMessage = res.Status
		}

		// set the path
		path := project.Path
		if len(result.WebURL) > 0 {
			path = result.WebURL
		}
		fmt.Printf("%s, %s, %s\n", project.Pid, path, statusMessage)

	}

}

var (
	conf        config
	debug       *bool
	findgrp     *string
	dumpcfg     *bool
	default_    *string
	ls          *string
	lspbranch   *bool
	setpbranch  *bool
	grpid       *string
	prjid       *string
	tagcsv      *string
	tagnosuffix *bool
	bulkmr      *string // Bulk Merge Request
)

func init() {
	conf.readConfig()
	log.SetOutput(io.Discard)
}

func main() {
	debug = flag.Bool("debug", false, "Enable verbose logging")
	findgrp = flag.String("findgrp", "", "Find the group. Example: -findgrp Common")
	dumpcfg = flag.Bool("dumpcfg", false, "Dump the configurations")
	default_ = flag.String("default", "", "Set the default branch, 2nd parameter is -gid or -pid")
	ls = flag.String("ls", "", "List all the repo in a Group")
	lspbranch = flag.Bool("lspb", false, "List all the protected branch settings in a group, 2nd parameter is -gid or -pid")
	setpbranch = flag.Bool("setpb", false, "Set branch(es) to be protected by the default configuration, 2nd parameter is -gid or -pid")
	grpid = flag.String("gid", "", "The group Id to perform an action")
	prjid = flag.String("pid", "", "The project id to perform an action")
	tagcsv = flag.String("tagcsv", "", "CSV file with tagging information. The tag will suffix with yymmdd-hash. CSV Header: pid,path,prefix,branch,message")
	tagnosuffix = flag.Bool("tagnosuffix", false, "depends on -tagcsv and this will disable the suffix of yymmdd-hash")
	bulkmr = flag.String("bulkmr", "", "CSV file with bulk merge request. CSV Header: pid,path,source,target,title")
	flag.Parse()

	if *debug {
		log.SetOutput(os.Stdout)
	}
	if *dumpcfg {
		fmt.Println(conf.ToString())
	}
	if len(*findgrp) > 0 {
		findGroup(*findgrp)
	}
	if len(*default_) > 0 {
		setDefaultBranch(*default_, *prjid, *grpid)
	}
	if len(*ls) > 0 {
		listProjects(*ls, true)

	}
	if *lspbranch {
		listProtectedBranch(*prjid, *grpid)
	}
	if *setpbranch {
		setProtectedBranch(*prjid, *grpid)
	}
	if len(*tagcsv) > 0 {
		tagCSV(*tagcsv, *tagnosuffix)
	}
	if len(*tagcsv) == 0 && *tagnosuffix {
		log.Fatalln("-tagcsv <csvfile> not defined")
	}
	if len(*bulkmr) > 0 {
		bulkMergeRequest(*bulkmr)
	}

}
