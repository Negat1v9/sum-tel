package server

import (
	"net/http"

	channelhttp "github.com/Negat1v9/sum-tel/services/core/internal/channel/http"
	channelservice "github.com/Negat1v9/sum-tel/services/core/internal/channel/service"
	"github.com/Negat1v9/sum-tel/services/core/internal/middleware"
)

func (s *Server) MapHandlers(chService *channelservice.ChannelService) {
	router := http.NewServeMux()

	mw := middleware.New(s.cfg)

	// initialize handlers
	channelHandler := channelhttp.NewChannelHandler(s.log, s.cfg, chService)

	// map handlers to routes
	channelRouter := channelhttp.NewChannelRouter(channelHandler, mw)
	router.Handle("/channels/", http.StripPrefix("/channels", channelRouter))

	apiV1Routes := http.NewServeMux()

	apiV1Routes.Handle("/api/", http.StripPrefix("/api", router))

	basicMW := mw.BasicMW()

	s.server.Handler = basicMW(apiV1Routes)
}
