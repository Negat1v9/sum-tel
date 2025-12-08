package newshttp

import (
	"context"
	"net/http"
	"time"

	"github.com/Negat1v9/sum-tel/services/core/internal/middleware"
	newsservice "github.com/Negat1v9/sum-tel/services/core/internal/news/service"
	"github.com/Negat1v9/sum-tel/services/core/pkg/utils"
	"github.com/Negat1v9/sum-tel/shared/config"
	"github.com/Negat1v9/sum-tel/shared/logger"
)

const (
	responseDataNameNews = "news"
)

type NewsHandler struct {
	log *logger.Logger
	cfg *config.CoreConfig
	s   *newsservice.NewsService
}

func NewNewsHandler(log *logger.Logger, cfg *config.CoreConfig, s *newsservice.NewsService) *NewsHandler {
	return &NewsHandler{
		log: log,
		cfg: cfg,
		s:   s,
	}
}

func (h *NewsHandler) News(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*time.Duration(h.cfg.WebConfig.ReadTimeout))
	defer cancel()

	userID := r.Context().Value(middleware.CtxUserIDKey).(int)

	v := r.URL.Query()

	newsList, err := h.s.News(ctx, userID, utils.GetLimitParam(v, 50), utils.GetOffset(v))
	if err != nil {
		utils.LogResponseErr(r, h.log, err)
		utils.WriteErrResponse(w, err)
		return
	}

	utils.WriteJsonResponse(w, responseDataNameNews, 200, newsList)
}

// retrun list of news sources from parser service by newsID
func (h *NewsHandler) NewsSources(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("not implemented"))
}
