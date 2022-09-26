package gorillamux

import (
	"net/http"
	"sample/api/core"

	"github.com/gorilla/mux"
)

type GorillaMuxAdapter struct {
	MainRequest core.MainRequest
	router      *mux.Router
}

func New(router *mux.Router) *GorillaMuxAdapter {
	return &GorillaMuxAdapter{
		router: router,
	}
}

func (h *GorillaMuxAdapter) MainFnAdapter(args core.MainRequestArgs) (core.MainResponseArgs, error) {
	req, err := h.MainRequest.MainArgsToHTTPRequest(args)
	return proxyAdapter(h, req, err)
}

func proxyAdapter(h *GorillaMuxAdapter, req *http.Request, err error) (core.MainResponseArgs, error) {
	if err != nil {
		return core.ErrorResponse(500), core.NewLoggedError("Could not convert proxy event to request: %v", err)
	}

	w := core.NewMainResponseWriter()
	h.router.ServeHTTP(http.ResponseWriter(w), req)

	return w.GetMainResponse()
}
