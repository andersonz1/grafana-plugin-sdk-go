package mock_test

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/andersonz1/grafana-plugin-sdk-go/experimental/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRoundTripper_RoundTrip(t *testing.T) {
	tests := []struct {
		name    string
		rt      *mock.RoundTripper
		req     *http.Request
		want    *http.Response
		wantErr error
		test    func(t *testing.T, res *http.Response)
	}{
		{
			name: "default mock client should return valid result",
			test: func(t *testing.T, res *http.Response) {
				t.Helper()
				require.NotNil(t, res)
				assert.Equal(t, http.StatusOK, res.StatusCode)
				assert.Equal(t, "200 OK", res.Status)
				assert.Equal(t, io.NopCloser(bytes.NewBufferString("{}")), res.Body)
			},
		},
		{
			name: "should return body if present",
			rt:   &mock.RoundTripper{Body: `{ "message" : "ok" }`},
			test: func(t *testing.T, res *http.Response) {
				t.Helper()
				require.NotNil(t, res)
				assert.Equal(t, http.StatusOK, res.StatusCode)
				assert.Equal(t, "200 OK", res.Status)
				assert.Equal(t, io.NopCloser(bytes.NewBufferString(`{ "message" : "ok" }`)), res.Body)
			},
		},
		{
			name: "should return body if present and respect status code",
			rt:   &mock.RoundTripper{Body: `{ "message" : "error" }`, StatusCode: 500, Status: "err"},
			test: func(t *testing.T, res *http.Response) {
				t.Helper()
				require.NotNil(t, res)
				assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
				assert.Equal(t, "err", res.Status)
				assert.Equal(t, io.NopCloser(bytes.NewBufferString(`{ "message" : "error" }`)), res.Body)
			},
		},
		{
			name: "should return file content if present",
			req:  exampleRequest(http.MethodPost, "http://foo/ok"),
			rt:   &mock.RoundTripper{FileName: "testdata/ok.json"},
			test: func(t *testing.T, res *http.Response) {
				t.Helper()
				require.NotNil(t, res)
				assert.Equal(t, http.StatusOK, res.StatusCode)
				assert.Equal(t, "200 OK", res.Status)
				b, _ := os.ReadFile("testdata/ok.json")
				rb, _ := io.ReadAll(res.Body)
				assert.Equal(t, b, rb)
			},
		},
		{
			name: "should return file content if present and respect status code",
			rt:   &mock.RoundTripper{FileName: "testdata/error.json", StatusCode: 500, Status: "err"},
			test: func(t *testing.T, res *http.Response) {
				t.Helper()
				require.NotNil(t, res)
				assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
				assert.Equal(t, "err", res.Status)
				b, _ := os.ReadFile("testdata/error.json")
				rb, _ := io.ReadAll(res.Body)
				assert.Equal(t, b, rb)
			},
		},
		{
			name:    "should return matched response from HAR file",
			rt:      &mock.RoundTripper{HARFileName: "testdata/example.har"},
			req:     exampleRequest(http.MethodGet, "https://grafana.com/api/plugins/two"),
			wantErr: errors.New("no matched request found in HAR file"),
			test: func(t *testing.T, res *http.Response) {
				t.Helper()
				require.NotNil(t, res)
				assert.Equal(t, http.StatusOK, res.StatusCode)
				assert.Equal(t, "OK", res.Status)
				assert.Equal(t, io.NopCloser(bytes.NewBufferString(`plugin-two`)), res.Body)
			},
		},
		{
			name:    "should throw error when no response matched from HAR file",
			rt:      &mock.RoundTripper{HARFileName: "testdata/example.har"},
			req:     exampleRequest(http.MethodGet, "https://grafana.com/api/plugins"),
			wantErr: errors.New("no matched request found in HAR file"),
		},
		{
			name: "should throw error when authentication is expected",
			rt:   &mock.RoundTripper{BasicAuthEnabled: true, BasicAuthUser: "foo", BasicAuthPassword: "bar"},
			req:  exampleRequest(http.MethodGet, "https://grafana.com/api/plugins"),
			test: func(t *testing.T, res *http.Response) {
				t.Helper()
				require.Equal(t, http.StatusUnauthorized, res.StatusCode)
			},
		},
		{
			name: "should not throw error when authentication is expected and details present",
			rt:   &mock.RoundTripper{BasicAuthEnabled: true, BasicAuthUser: "foo", BasicAuthPassword: "bar"},
			req:  exampleRequest(http.MethodGet, "https://foo:bar@grafana.com/api/plugins"),
			test: func(t *testing.T, res *http.Response) {
				t.Helper()
				require.Equal(t, http.StatusOK, res.StatusCode)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rt := &mock.RoundTripper{}
			if tt.rt != nil {
				rt = tt.rt
			}
			got, err := rt.RoundTrip(tt.req)
			if got != nil {
				defer got.Body.Close()
			}
			if tt.wantErr != nil {
				require.NotNil(t, err)
				assert.Equal(t, tt.wantErr, err)
				return
			}
			require.Nil(t, err)
			require.NotNil(t, got)
			if tt.want != nil {
				assert.Equal(t, tt.want, got)
			}
			if tt.test != nil {
				tt.test(t, got)
			}
		})
	}
}

func exampleRequest(method string, u string) *http.Request {
	req, _ := http.NewRequest(method, u, nil)
	return req
}
