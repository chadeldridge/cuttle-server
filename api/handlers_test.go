package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chadeldridge/cuttle/core"
	"github.com/chadeldridge/cuttle/router"
	"github.com/chadeldridge/cuttle/test_helpers"
	"github.com/stretchr/testify/require"
)

func TestRoutesAddRoutes(t *testing.T) {
	require := require.New(t)
	mux := http.NewServeMux()
	// logger := core.NewLogger(nil, "cuttle: ", 0, false)

	// Add routes here if neccasary.
	// addAPIRoutes(mux, &HTTPServer{logger: logger})
	require.NotNil(mux, "addRoutes() returned nil")

	w := httptest.NewRecorder()
	err := router.RenderJSON(w, http.StatusOK, struct{ Message string }{Message: "you did it"})
	require.NoError(err, "encode() returned an error: %s", err)

	exp := `{"Message":"you did it"}` + "\n"
	got := w.Body.String()
	require.Equal(exp, got, "encode() returned wrong body")
}

func TestRoutesHandleTest(t *testing.T) {
	require := require.New(t)
	logger := core.NewLogger(nil, "cuttle: ", 0, false)

	resp := test_helpers.TestHandler(t, handleTest(logger), "GET", "/v1/test", nil, http.StatusOK)
	exp := struct{ Message string }{Message: "you did it"}
	got, err := router.ReadJSON[struct{ Message string }](&http.Request{Body: resp.Result().Body})
	require.NoError(err, "decode() returned an error: %s", err)
	require.Equal(exp, got, "handler returned wrong body")
}
