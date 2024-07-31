package main

import (
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func Test_GetData(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mockURL := "http://localhost/mockguy"

	mockResponse := []map[string]interface{}{
		{
			"userId": 1,
			"id":     1,
			"title":  "TITLE 1",
			"body":   "BODY #1",
		},
		{
			"userId": 2,
			"id":     2,
			"title":  "TITLE 2",
			"body":   "BODY #2",
		},
		{
			"userId": 3,
			"id":     3,
			"title":  "TITLE 3",
			"body":   "BODY #3",
		},
	}

	httpmock.RegisterResponder("GET", mockURL, func(r *http.Request) (*http.Response, error) {
		resp, err := httpmock.NewJsonResponse(http.StatusOK, mockResponse)
		if err != nil {
			return httpmock.NewStringResponse(500, ""), nil
		}
		return resp, nil
	})

	resp, err := GetData(mockURL)
	assert.Nil(t, err)
	assert.Equal(t, []*Source{
		{ID: 1, Title: "TITLE 1", Body: "BODY #1"},
		{ID: 2, Title: "TITLE 2", Body: "BODY #2"},
		{ID: 3, Title: "TITLE 3", Body: "BODY #3"},
	}, resp)
}
