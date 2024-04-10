package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/xanzy/go-gitlab"
)

/* The data coming from the client is the same,
but the Gitlab endpoints + resources we handle are different */

type PostDraftCommentRequest struct {
	Note           string     `json:"note"`
	StartCommitSHA string     `json:"start_commit_sha"`
	NewLine        *int       `json:"new_line,omitempty"`
	OldLine        *int       `json:"old_line,omitempty"`
	LineRange      *LineRange `json:"line_range,omitempty"`
	FileName       string     `json:"file_name"`
}
type DeleteDraftCommentRequest struct{}
type EditDraftCommentRequest struct{}

type DraftNoteResponse struct {
	SuccessResponse
	DraftNote *gitlab.DraftNote
}

/* commentHandler creates, edits, and deletes draft discussions (comments, multi-line comments) */
func (a *api) draftCommentHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case http.MethodPost:
		a.postDraftComment(w, r)
	case http.MethodPatch:
		a.editDraftComment(w, r)
	case http.MethodDelete:
		a.deleteDraftComment(w, r)
	default:
		w.Header().Set("Access-Control-Allow-Methods", fmt.Sprintf("%s, %s, %s", http.MethodDelete, http.MethodPost, http.MethodPatch))
		handleError(w, InvalidRequestError{}, "Expected DELETE, POST or PATCH", http.StatusMethodNotAllowed)
	}
}

/* postComment creates a draft comment */
func (a *api) postDraftComment(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		handleError(w, err, "Could not read request body", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	var postDraftCommentRequest PostDraftCommentRequest
	err = json.Unmarshal(body, &postDraftCommentRequest)
	if err != nil {
		handleError(w, err, "Could not unmarshal data from request body", http.StatusBadRequest)
		return
	}

	opt := gitlab.CreateDraftNoteOptions{
		Note:     &postDraftCommentRequest.Note,
		CommitID: &postDraftCommentRequest.StartCommitSHA,
		// InReplyToDiscussionID *string          `url:"in_reply_to_discussion_id,omitempty" json:"in_reply_to_discussion_id,omitempty"`
		// Position              *PositionOptions `url:"position,omitempty" json:"position,omitempty"`
	}

	draftNote, res, err := a.client.CreateDraftNote(a.projectInfo.ProjectId, a.projectInfo.MergeId, &opt)

	if err != nil {
		handleError(w, err, "Could not create draft note", http.StatusInternalServerError)
		return
	}

	if res.StatusCode >= 300 {
		handleError(w, GenericError{endpoint: "/mr/draft/comment"}, "Could not create draft note", res.StatusCode)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := DraftNoteResponse{
		SuccessResponse: SuccessResponse{
			Message: "Draft note created successfully",
			Status:  http.StatusOK,
		},
		DraftNote: draftNote,
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		handleError(w, err, "Could not encode response", http.StatusInternalServerError)
	}

}

/* deleteComment deletes a draft comment */
func (a *api) deleteDraftComment(w http.ResponseWriter, r *http.Request) {}

/* deleteComment edits a draft comment */
func (a *api) editDraftComment(w http.ResponseWriter, r *http.Request) {}
