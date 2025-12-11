package userhttp

import (
	"net/http"

	"github.com/Negat1v9/sum-tel/services/core/internal/middleware"
)

func NewUsersRouter(h *UserHandler, mw *middleware.MiddleWareManager) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /login", h.LoginOrRegister)

	mux.Handle("GET /me", mw.AuthUser(http.HandlerFunc(h.User)))

	return mux
}
