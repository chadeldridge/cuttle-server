package api

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/chadeldridge/cuttle/core"
	"github.com/stretchr/testify/require"
)

func testHandler(
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

func TestRoutesAddRoutes(t *testing.T) {
	require := require.New(t)
	mux := http.NewServeMux()
	logger := core.NewLogger(nil, "cuttle: ", 0, false)

	addRoutes(mux, &HTTPServer{logger: logger})
	require.NotNil(mux, "addRoutes() returned nil")

	w := httptest.NewRecorder()
	err := encode(w, http.StatusOK, struct{ Message string }{Message: "you did it"})
	require.NoError(err, "encode() returned an error: %s", err)

	exp := `{"Message":"you did it"}` + "\n"
	got := w.Body.String()
	require.Equal(exp, got, "encode() returned wrong body")
}

func TestRoutesEncode(t *testing.T) {
	require := require.New(t)

	t.Run("valid", func(t *testing.T) {
		w := httptest.NewRecorder()
		err := encode(w, http.StatusOK, struct{ Message string }{Message: "you did it"})
		require.NoError(err, "encode() returned an error: %s", err)

		exp := `{"Message":"you did it"}` + "\n"
		got := w.Body.String()
		require.Equal(exp, got, "encode() returned wrong body")
	})

	t.Run("invalid", func(t *testing.T) {
		w := httptest.NewRecorder()
		err := encode(w, http.StatusOK, make(chan struct{}))
		require.Error(err, "encode() did not return an error")
		require.Equal("encoder: json: unsupported type: chan struct {}", err.Error(), "encode() returned wrong error")
	})
}

func TestRoutesDecode(t *testing.T) {
	require := require.New(t)
	body := strings.NewReader(`{"Message":"you did it"}` + "\n")
	req := http.Request{Body: io.NopCloser(body)}

	t.Run("valid", func(t *testing.T) {
		data, err := decode[struct{ Message string }](&req)
		require.NoError(err, "decode() returned an error: %s", err)
		require.Equal(struct{ Message string }{Message: "you did it"}, data, "decode() returned wrong data")
	})

	t.Run("invalid", func(t *testing.T) {
		data, err := decode[struct{ Data int }](&req)
		require.Error(err, "decode() did not return an error")
		require.Equal(struct{ Data int }{}, data, "decode() returned wrong data")
	})
}

func TestRoutesHandleTest(t *testing.T) {
	require := require.New(t)
	logger := core.NewLogger(nil, "cuttle: ", 0, false)

	resp := testHandler(t, handleTest(logger), "GET", "/v1/test", nil, http.StatusOK)
	exp := struct{ Message string }{Message: "you did it"}
	got, err := decode[struct{ Message string }](&http.Request{Body: resp.Result().Body})
	require.NoError(err, "decode() returned an error: %s", err)
	require.Equal(exp, got, "handler returned wrong body")
}
