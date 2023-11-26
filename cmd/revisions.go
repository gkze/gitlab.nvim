package main

import (
	"encoding/json"
	"net/http"

	"github.com/xanzy/go-gitlab"
)

type RevisionsResponse struct {
	SuccessResponse
	Revisions []*gitlab.MergeRequestDiffVersion
}

func RevisionsHandler(w http.ResponseWriter, r *http.Request, c HandlerClient, d *ProjectInfo) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodGet {
		w.Header().Set("Access-Control-Allow-Methods", http.MethodGet)
		HandleError(w, InvalidRequestError{}, "Expected GET", http.StatusMethodNotAllowed)
		return
	}

	versionInfo, res, err := c.GetMergeRequestDiffVersions(d.ProjectId, d.MergeId, &gitlab.GetMergeRequestDiffVersionsOptions{})
	if err != nil {
		HandleError(w, err, "Could not get diff version info", http.StatusInternalServerError)
		return
	}

	if res.StatusCode >= 300 {
		HandleError(w, GenericError{endpoint: "/mr/revisions"}, "Could not get diff version info", res.StatusCode)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := RevisionsResponse{
		SuccessResponse: SuccessResponse{
			Message: "Revisions fetched successfully",
			Status:  http.StatusOK,
		},
		Revisions: versionInfo,
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		HandleError(w, err, "Could not encode response", http.StatusInternalServerError)
	}

}
