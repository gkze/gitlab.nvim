package main

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"os"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/xanzy/go-gitlab"
)

type DebugSettings struct {
	GoRequest  bool `json:"go_request"`
	GoResponse bool `json:"go_response"`
}

type ConnectionOptions struct {
	Insecure bool `json:"insecure"`
}

type ProjectInfo struct {
	ProjectId string
	MergeId   int
}

/* The Client struct embeds all the methods from Gitlab for the different services */
type Client struct {
	*gitlab.MergeRequestsService
	*gitlab.MergeRequestApprovalsService
	*gitlab.DiscussionsService
	*gitlab.ProjectsService
	*gitlab.ProjectMembersService
	*gitlab.JobsService
	*gitlab.PipelinesService
	*gitlab.LabelsService
	*gitlab.AwardEmojiService
	*gitlab.UsersService
	*gitlab.DraftNotesService
}

/* initGitlabClient parses and validates the project settings and initializes the Gitlab client. */
func initGitlabClient() (error, *Client) {

	if len(os.Args) < 7 {
		return errors.New("Must provide gitlab url, port, auth token, debug settings, log path, and connection settings"), nil
	}

	gitlabInstance := os.Args[1]
	if gitlabInstance == "" {
		return errors.New("GitLab instance URL cannot be empty"), nil
	}

	authToken := os.Args[3]
	if authToken == "" {
		return errors.New("Auth token cannot be empty"), nil
	}

	/* Parse debug settings and initialize logger handlers */
	debugSettings := os.Args[4]
	var debugObject DebugSettings
	err := json.Unmarshal([]byte(debugSettings), &debugObject)
	if err != nil {
		return fmt.Errorf("Could not parse debug settings: %w, %s", err, debugSettings), nil
	}

	/* Parse connection options */
	connectionSettings := os.Args[6]
	var connectionObject ConnectionOptions
	err = json.Unmarshal([]byte(connectionSettings), &connectionObject)
	if err != nil {
		return fmt.Errorf("Could not parse connection settings: %w, %s", err, connectionSettings), nil
	}

	var apiCustUrl = fmt.Sprintf(gitlabInstance + "/api/v4")

	gitlabOptions := []gitlab.ClientOptionFunc{
		gitlab.WithBaseURL(apiCustUrl),
	}

	if debugObject.GoRequest {
		gitlabOptions = append(gitlabOptions, gitlab.WithRequestLogHook(requestLogger))
	}

	if debugObject.GoResponse {
		gitlabOptions = append(gitlabOptions, gitlab.WithResponseLogHook(responseLogger))
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: connectionObject.Insecure,
		},
	}

	retryClient := retryablehttp.NewClient()
	retryClient.HTTPClient.Transport = tr
	gitlabOptions = append(gitlabOptions, gitlab.WithHTTPClient(retryClient.HTTPClient))

	client, err := gitlab.NewClient(authToken, gitlabOptions...)

	if err != nil {
		return fmt.Errorf("Failed to create client: %v", err), nil
	}

	return nil, &Client{
		MergeRequestsService:         client.MergeRequests,
		MergeRequestApprovalsService: client.MergeRequestApprovals,
		DiscussionsService:           client.Discussions,
		ProjectsService:              client.Projects,
		ProjectMembersService:        client.ProjectMembers,
		JobsService:                  client.Jobs,
		PipelinesService:             client.Pipelines,
		LabelsService:                client.Labels,
		AwardEmojiService:            client.AwardEmoji,
		UsersService:                 client.Users,
		DraftNotesService:            client.DraftNotes,
	}
}

/* initProjectSettings fetch the project ID using the client */
func initProjectSettings(c *Client, gitInfo GitProjectInfo) (error, *ProjectInfo) {

	opt := gitlab.GetProjectOptions{}
	project, _, err := c.GetProject(gitInfo.projectPath(), &opt)

	if err != nil {
		return fmt.Errorf(fmt.Sprintf("Error getting project at %s", gitInfo.RemoteUrl), err), nil
	}
	if project == nil {
		return fmt.Errorf(fmt.Sprintf("Could not find project at %s", gitInfo.RemoteUrl), err), nil
	}

	if project == nil {
		return fmt.Errorf("No projects you are a member of contained remote URL %s", gitInfo.RemoteUrl), nil
	}

	projectId := fmt.Sprint(project.ID)

	return nil, &ProjectInfo{
		ProjectId: projectId,
	}

}

/* handleError is a utililty handler that returns errors to the client along with their statuses and messages */
func handleError(w http.ResponseWriter, err error, message string, status int) {
	w.WriteHeader(status)
	response := ErrorResponse{
		Message: message,
		Details: err.Error(),
		Status:  status,
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		handleError(w, err, "Could not encode error response", http.StatusInternalServerError)
	}
}

var requestLogger retryablehttp.RequestLogHook = func(l retryablehttp.Logger, r *http.Request, i int) {
	file := openLogFile()
	defer file.Close()

	token := r.Header.Get("Private-Token")
	r.Header.Set("Private-Token", "REDACTED")
	res, err := httputil.DumpRequest(r, true)
	if err != nil {
		log.Fatalf("Error dumping request: %v", err)
		os.Exit(1)
	}
	r.Header.Set("Private-Token", token)

	_, err = file.Write([]byte("\n-- REQUEST --\n")) //nolint:all
	_, err = file.Write(res)                         //nolint:all
	_, err = file.Write([]byte("\n"))                //nolint:all
}

var responseLogger retryablehttp.ResponseLogHook = func(l retryablehttp.Logger, response *http.Response) {
	file := openLogFile()
	defer file.Close()

	res, err := httputil.DumpResponse(response, true)
	if err != nil {
		log.Fatalf("Error dumping response: %v", err)
		os.Exit(1)
	}

	_, err = file.Write([]byte("\n-- RESPONSE --\n")) //nolint:all
	_, err = file.Write(res)                          //nolint:all
	_, err = file.Write([]byte("\n"))                 //nolint:all
}

func openLogFile() *os.File {
	logFile := os.Args[5]
	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		if os.IsNotExist(err) {
			log.Printf("Log file %s does not exist", logFile)
		} else if os.IsPermission(err) {
			log.Printf("Permission denied for log file %s", logFile)
		} else {
			log.Printf("Error opening log file %s: %v", logFile, err)
		}

		os.Exit(1)
	}

	return file
}
