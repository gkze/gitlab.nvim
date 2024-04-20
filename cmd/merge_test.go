package main

import (
	"errors"
	"net/http"
	"testing"

	"github.com/xanzy/go-gitlab"
)

func acceptMergeRequest(pid interface{}, mergeRequest int, opt *gitlab.AcceptMergeRequestOptions, options ...gitlab.RequestOptionFunc) (*gitlab.MergeRequest, *gitlab.Response, error) {
	return &gitlab.MergeRequest{}, makeResponse(http.StatusOK), nil
}

func acceptMergeRequestErr(pid interface{}, mergeRequest int, opt *gitlab.AcceptMergeRequestOptions, options ...gitlab.RequestOptionFunc) (*gitlab.MergeRequest, *gitlab.Response, error) {
	return nil, nil, errors.New("Some error from Gitlab")
}

func acceptAndMergeNon200(pid interface{}, mergeRequest int, opt *gitlab.AcceptMergeRequestOptions, options ...gitlab.RequestOptionFunc) (*gitlab.MergeRequest, *gitlab.Response, error) {
	return nil, makeResponse(http.StatusSeeOther), nil
}

func TestAcceptAndMergeHandler(t *testing.T) {
	t.Run("Accepts and merges a merge request", func(t *testing.T) {
		request := makeRequest(t, http.MethodPost, "/mr/merge", AcceptMergeRequestRequest{})
		server, _ := createRouterAndApi(fakeClient{acceptMergeRequest: acceptMergeRequest})
		data := serveRequest(t, server, request, SuccessResponse{})
		assert(t, data.Message, "MR merged successfully")
		assert(t, data.Status, http.StatusOK)
	})

	t.Run("Disallows non-POST methods", func(t *testing.T) {
		request := makeRequest(t, http.MethodGet, "/mr/merge", AcceptMergeRequestRequest{})
		server, _ := createRouterAndApi(fakeClient{acceptMergeRequest: acceptMergeRequest})
		data := serveRequest(t, server, request, ErrorResponse{})
		checkBadMethod(t, *data, http.MethodPost)
	})

	t.Run("Handles errors from Gitlab client", func(t *testing.T) {
		request := makeRequest(t, http.MethodPost, "/mr/merge", AcceptMergeRequestRequest{})
		server, _ := createRouterAndApi(fakeClient{acceptMergeRequest: acceptMergeRequestErr})
		data := serveRequest(t, server, request, ErrorResponse{})
		checkErrorFromGitlab(t, *data, "Could not merge MR")
	})

	t.Run("Handles non-200s from Gitlab client", func(t *testing.T) {
		request := makeRequest(t, http.MethodPost, "/mr/merge", AcceptMergeRequestRequest{})
		server, _ := createRouterAndApi(fakeClient{acceptMergeRequest: acceptAndMergeNon200})
		data := serveRequest(t, server, request, ErrorResponse{})
		checkNon200(t, *data, "Could not merge MR", "/mr/merge")
	})
}
