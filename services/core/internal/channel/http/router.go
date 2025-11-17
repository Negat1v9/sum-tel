package channelhttp

import (
	"net/http"

	"github.com/Negat1v9/sum-tel/services/core/internal/middleware"
)

func NewChannelRouter(h *ChannelHandler, mw *middleware.MiddleWareManager) http.Handler {
	mux := http.NewServeMux()

	private := http.NewServeMux()

	private.HandleFunc("POST /new", h.CreateChannel)
	private.HandleFunc("GET /info", h.GetChannel)
	private.HandleFunc("GET /subscriptions", h.GetSubscriptions)
	private.HandleFunc("POST /subscribe", h.SubscribeChannel)
	// all private routes with auth middleware
	mux.Handle("/", mw.AuthUser(private))

	return mux
}
