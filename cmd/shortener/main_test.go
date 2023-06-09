package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/KorsakovPV/shortener/cmd/shortener/storage"
	"github.com/KorsakovPV/shortener/internal/apiserver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testRequest(t *testing.T, ts *httptest.Server, method, path string, jsonBody io.Reader) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, jsonBody)
	require.NoError(t, err)

	resp, err := ts.Client().Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp, string(respBody)
}

func TestRouter(t *testing.T) {
	ts := httptest.NewServer(apiserver.Router())
	defer ts.Close()

	successBodyRegex := `^http://127.0.0.1:8080/[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`

	testCases := []struct {
		url          string
		body         string
		method       string
		expectedCode int
		expectedBody string
	}{
		{url: `/`, method: http.MethodGet, expectedCode: http.StatusBadRequest, expectedBody: ""},
		{url: `/`, method: http.MethodPost, expectedCode: http.StatusCreated, expectedBody: successBodyRegex},
		{url: `/`, method: http.MethodPut, expectedCode: http.StatusBadRequest, expectedBody: ""},
		{url: `/`, method: http.MethodPatch, expectedCode: http.StatusBadRequest, expectedBody: ""},
		{url: `/`, method: http.MethodDelete, expectedCode: http.StatusBadRequest, expectedBody: ""},
		{url: `/api/shorten`, body: `{"url": "https://practicum.yandex.ru"}`, method: http.MethodPost, expectedCode: http.StatusCreated, expectedBody: ""},
	}

	for _, tc := range testCases {
		t.Run(tc.method, func(t *testing.T) {
			jsonBody := []byte(tc.body)
			bodyReader := bytes.NewReader(jsonBody)
			_ = storage.InitStorage()
			resp, get := testRequest(t, ts, tc.method, tc.url, bodyReader)
			defer resp.Body.Close()

			assert.Equal(t, tc.expectedCode, resp.StatusCode, "Код ответа не совпадает с ожидаемым")
			if tc.expectedBody != "" {
				match, _ := regexp.MatchString(tc.expectedBody, get)
				fmt.Println(get, match)
			}
		})
	}
}
