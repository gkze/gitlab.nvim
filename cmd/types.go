package main

import (
	"bytes"
	"fmt"
	"io"

	"github.com/xanzy/go-gitlab"
)

type ErrorResponse struct {
	Message string `json:"message"`
	Details string `json:"details"`
	Status  int    `json:"status"`
}

type SuccessResponse struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

type GenericError struct {
	endpoint string
}

func (e GenericError) Error() string {
	return fmt.Sprintf("An error occurred on the %s endpoint", e.endpoint)
}

type InvalidRequestError struct{}

func (e InvalidRequestError) Error() string {
	return "Invalid request type"
}

/* The ClientInterface interface implements all the methods that our handlers need */
type ClientInterface interface {
	CreateMergeRequest(pid interface{}, opt *gitlab.CreateMergeRequestOptions, options ...gitlab.RequestOptionFunc) (*gitlab.MergeRequest, *gitlab.Response, error)
	ListProjectMergeRequests(pid interface{}, opt *gitlab.ListProjectMergeRequestsOptions, options ...gitlab.RequestOptionFunc) ([]*gitlab.MergeRequest, *gitlab.Response, error)
	GetMergeRequest(pid interface{}, mergeRequestIID int, opt *gitlab.GetMergeRequestsOptions, options ...gitlab.RequestOptionFunc) (*gitlab.MergeRequest, *gitlab.Response, error)
	AcceptMergeRequest(pid interface{}, mergeRequestIID int, opt *gitlab.AcceptMergeRequestOptions, options ...gitlab.RequestOptionFunc) (*gitlab.MergeRequest, *gitlab.Response, error)
	UpdateMergeRequest(pid interface{}, mergeRequestIID int, opt *gitlab.UpdateMergeRequestOptions, options ...gitlab.RequestOptionFunc) (*gitlab.MergeRequest, *gitlab.Response, error)
	UploadFile(pid interface{}, content io.Reader, filename string, options ...gitlab.RequestOptionFunc) (*gitlab.ProjectFile, *gitlab.Response, error)
	GetMergeRequestDiffVersions(pid interface{}, mergeRequestIID int, opt *gitlab.GetMergeRequestDiffVersionsOptions, options ...gitlab.RequestOptionFunc) ([]*gitlab.MergeRequestDiffVersion, *gitlab.Response, error)
	ApproveMergeRequest(pid interface{}, mergeRequestIID int, opt *gitlab.ApproveMergeRequestOptions, options ...gitlab.RequestOptionFunc) (*gitlab.MergeRequestApprovals, *gitlab.Response, error)
	UnapproveMergeRequest(pid interface{}, mergeRequestIID int, options ...gitlab.RequestOptionFunc) (*gitlab.Response, error)
	ListMergeRequestDiscussions(pid interface{}, mergeRequestIID int, opt *gitlab.ListMergeRequestDiscussionsOptions, options ...gitlab.RequestOptionFunc) ([]*gitlab.Discussion, *gitlab.Response, error)
	ResolveMergeRequestDiscussion(pid interface{}, mergeRequestIID int, discussion string, opt *gitlab.ResolveMergeRequestDiscussionOptions, options ...gitlab.RequestOptionFunc) (*gitlab.Discussion, *gitlab.Response, error)
	CreateMergeRequestDiscussion(pid interface{}, mergeRequestIID int, opt *gitlab.CreateMergeRequestDiscussionOptions, options ...gitlab.RequestOptionFunc) (*gitlab.Discussion, *gitlab.Response, error)
	UpdateMergeRequestDiscussionNote(pid interface{}, mergeRequestIID int, discussion string, note int, opt *gitlab.UpdateMergeRequestDiscussionNoteOptions, options ...gitlab.RequestOptionFunc) (*gitlab.Note, *gitlab.Response, error)
	DeleteMergeRequestDiscussionNote(pid interface{}, mergeRequestIID int, discussion string, note int, options ...gitlab.RequestOptionFunc) (*gitlab.Response, error)
	CreateDraftNote(pid interface{}, mergeRequestIID int, opt *gitlab.CreateDraftNoteOptions, options ...gitlab.RequestOptionFunc) (*gitlab.DraftNote, *gitlab.Response, error)
	ListDraftNotes(pid interface{}, mergeRequest int, opt *gitlab.ListDraftNotesOptions, options ...gitlab.RequestOptionFunc) ([]*gitlab.DraftNote, *gitlab.Response, error)
	DeleteDraftNote(pid interface{}, mergeRequest int, note int, options ...gitlab.RequestOptionFunc) (*gitlab.Response, error)
	UpdateDraftNote(pid interface{}, mergeRequest int, note int, opt *gitlab.UpdateDraftNoteOptions, options ...gitlab.RequestOptionFunc) (*gitlab.DraftNote, *gitlab.Response, error)
	AddMergeRequestDiscussionNote(pid interface{}, mergeRequestIID int, discussion string, opt *gitlab.AddMergeRequestDiscussionNoteOptions, options ...gitlab.RequestOptionFunc) (*gitlab.Note, *gitlab.Response, error)
	ListAllProjectMembers(pid interface{}, opt *gitlab.ListProjectMembersOptions, options ...gitlab.RequestOptionFunc) ([]*gitlab.ProjectMember, *gitlab.Response, error)
	RetryPipelineBuild(pid interface{}, pipeline int, options ...gitlab.RequestOptionFunc) (*gitlab.Pipeline, *gitlab.Response, error)
	ListPipelineJobs(pid interface{}, pipelineID int, opts *gitlab.ListJobsOptions, options ...gitlab.RequestOptionFunc) ([]*gitlab.Job, *gitlab.Response, error)
	GetLatestPipeline(pid interface{}, opt *gitlab.GetLatestPipelineOptions, options ...gitlab.RequestOptionFunc) (*gitlab.Pipeline, *gitlab.Response, error)
	GetTraceFile(pid interface{}, jobID int, options ...gitlab.RequestOptionFunc) (*bytes.Reader, *gitlab.Response, error)
	ListLabels(pid interface{}, opt *gitlab.ListLabelsOptions, options ...gitlab.RequestOptionFunc) ([]*gitlab.Label, *gitlab.Response, error)
	ListMergeRequestAwardEmojiOnNote(pid interface{}, mergeRequestIID int, noteID int, opt *gitlab.ListAwardEmojiOptions, options ...gitlab.RequestOptionFunc) ([]*gitlab.AwardEmoji, *gitlab.Response, error)
	CreateMergeRequestAwardEmojiOnNote(pid interface{}, mergeRequestIID int, noteID int, opt *gitlab.CreateAwardEmojiOptions, options ...gitlab.RequestOptionFunc) (*gitlab.AwardEmoji, *gitlab.Response, error)
	DeleteMergeRequestAwardEmojiOnNote(pid interface{}, mergeRequestIID, noteID, awardID int, options ...gitlab.RequestOptionFunc) (*gitlab.Response, error)
	CurrentUser(options ...gitlab.RequestOptionFunc) (*gitlab.User, *gitlab.Response, error)
}
