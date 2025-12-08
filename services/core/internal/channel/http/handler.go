package channelhttp

import (
	"context"
	"net/http"
	"time"

	channelservice "github.com/Negat1v9/sum-tel/services/core/internal/channel/service"
	"github.com/Negat1v9/sum-tel/services/core/internal/middleware"
	"github.com/Negat1v9/sum-tel/services/core/pkg/utils"
	"github.com/Negat1v9/sum-tel/shared/config"
	"github.com/Negat1v9/sum-tel/shared/logger"
)

const (
	responseDataNameChannel      = "channel"
	responseDataNameSubscription = "subscription"
)

type ChannelHandler struct {
	log *logger.Logger
	cfg *config.CoreConfig
	s   *channelservice.ChannelService
}

func NewChannelHandler(log *logger.Logger, cfg *config.CoreConfig, s *channelservice.ChannelService) *ChannelHandler {
	return &ChannelHandler{
		log: log,
		cfg: cfg,
		s:   s,
	}
}

func (h *ChannelHandler) CreateChannel(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*time.Duration(h.cfg.WebConfig.ReadTimeout))
	defer cancel()

	username := r.URL.Query().Get("username")
	if username == "" {
		utils.WriteErrResponse(w, utils.NewError(http.StatusUnprocessableEntity, "username is required", nil))
		return
	}
	userID := r.Context().Value(middleware.CtxUserIDKey).(int)

	createdChannel, err := h.s.CreateChannel(ctx, userID, username)
	if err != nil {
		utils.LogResponseErr(r, h.log, err)
		utils.WriteErrResponse(w, err)
		return
	}

	utils.WriteJsonResponse(w, responseDataNameChannel, http.StatusCreated, createdChannel)
}

func (h *ChannelHandler) GetChannel(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*time.Duration(h.cfg.WebConfig.ReadTimeout))
	defer cancel()

	username := r.URL.Query().Get("username")
	if username == "" {
		utils.WriteErrResponse(w, utils.NewError(http.StatusUnprocessableEntity, "username is required", nil))
		return
	}

	channel, err := h.s.GetChannelByUsername(ctx, username)
	if err != nil {
		utils.LogResponseErr(r, h.log, err)
		utils.WriteErrResponse(w, err)
		return
	}

	utils.WriteJsonResponse(w, responseDataNameChannel, http.StatusOK, channel)
}

func (h *ChannelHandler) SubscribeChannel(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*time.Duration(h.cfg.WebConfig.ReadTimeout))
	defer cancel()

	channelIDParam := r.URL.Query().Get("channel_id")
	if channelIDParam == "" {
		utils.WriteErrResponse(w, utils.NewError(http.StatusUnprocessableEntity, "channel_id is required", nil))
		return
	}

	userID := r.Context().Value(middleware.CtxUserIDKey).(int)

	subscription, err := h.s.SubscribeChannel(ctx, userID, channelIDParam)
	if err != nil {
		utils.LogResponseErr(r, h.log, err)
		utils.WriteErrResponse(w, err)
		return
	}

	utils.WriteJsonResponse(w, responseDataNameSubscription, http.StatusCreated, subscription)
}

func (h *ChannelHandler) GetSubscriptions(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*time.Duration(h.cfg.WebConfig.ReadTimeout))
	defer cancel()

	userID := r.Context().Value(middleware.CtxUserIDKey).(int)

	v := r.URL.Query()

	subscriptions, err := h.s.UsersSubscriptions(ctx, userID, utils.GetLimitParam(v, 30), utils.GetOffset(v))
	if err != nil {
		utils.LogResponseErr(r, h.log, err)
		utils.WriteErrResponse(w, err)
		return
	}

	utils.WriteJsonResponse(w, "", http.StatusOK, subscriptions)
}
