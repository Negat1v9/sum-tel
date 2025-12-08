package newshttp

import (
	"net/http"

	"github.com/Negat1v9/sum-tel/services/core/internal/middleware"
)

func NewNewsRouter(h *NewsHandler, mw *middleware.MiddleWareManager) http.Handler {
	mux := http.NewServeMux()

	private := http.NewServeMux()

	private.HandleFunc("GET /latest", h.News)
	private.HandleFunc("GET /sources/{id}", h.NewsSources)

	// all private routes with auth middleware
	mux.Handle("/", mw.AuthUser(private))

	return mux
}
