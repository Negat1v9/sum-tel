package server

import (
	"net/http"

	channelhttp "github.com/Negat1v9/sum-tel/services/core/internal/channel/http"
	channelservice "github.com/Negat1v9/sum-tel/services/core/internal/channel/service"
	"github.com/Negat1v9/sum-tel/services/core/internal/middleware"
	newshttp "github.com/Negat1v9/sum-tel/services/core/internal/news/http"
	newsservice "github.com/Negat1v9/sum-tel/services/core/internal/news/service"
	userhttp "github.com/Negat1v9/sum-tel/services/core/internal/user/http"
	userservice "github.com/Negat1v9/sum-tel/services/core/internal/user/service"
	"github.com/Negat1v9/sum-tel/services/core/pkg/metrics"
)

func (s *Server) MapHandlers(chService *channelservice.ChannelService, newsService *newsservice.NewsService, userService *userservice.UserService) {
	router := http.NewServeMux()

	// run metrics prometheus
	metrics, err := metrics.NewMetric(":9090", "core_api")
	if err != nil {
		s.log.Errorf("Failed to initialize metrics: %v", err)
		return
	}

	mw := middleware.New(s.cfg, metrics)

	// initialize handlers
	channelHandler := channelhttp.NewChannelHandler(s.log, s.cfg, chService)
	newsHandler := newshttp.NewNewsHandler(s.log, s.cfg, newsService)
	userHandler := userhttp.NewUserHandler(s.log, s.cfg, userService)

	// map handlers to routes
	channelRouter := channelhttp.NewChannelRouter(channelHandler, mw)
	router.Handle("/channels/", http.StripPrefix("/channels", channelRouter))

	newsRouter := newshttp.NewNewsRouter(newsHandler, mw)
	router.Handle("/news/", http.StripPrefix("/news", newsRouter))

	userRouter := userhttp.NewUsersRouter(userHandler, mw)
	router.Handle("/users/", http.StripPrefix("/users", userRouter))

	apiV1Routes := http.NewServeMux()

	apiV1Routes.Handle("/api/", http.StripPrefix("/api", router))

	basicMW := middleware.CreateStack(middleware.Cors, mw.MetricsMiddleware)

	s.server.Handler = basicMW(apiV1Routes)
}
