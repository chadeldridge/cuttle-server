package test_helpers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHandler(
	t *testing.T,
	h http.Handler,
	method, path string,
	body io.Reader,
	expCode int,
) *httptest.ResponseRecorder {
	require := require.New(t)
	req, err := http.NewRequest(method, path, body)
	require.NoError(err, "http.NewRequest() returned an error: %s", err)

	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)
	require.Equal(expCode, rr.Code, "handler returned wrong status code")

	return rr
}
