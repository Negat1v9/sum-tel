package userhttp

import (
	"context"
	"net/http"
	"time"

	userservice "github.com/Negat1v9/sum-tel/services/core/internal/user/service"
	"github.com/Negat1v9/sum-tel/services/core/pkg/utils"
	"github.com/Negat1v9/sum-tel/shared/config"
	"github.com/Negat1v9/sum-tel/shared/logger"
)

const (
	responseDataNameUsers = "user"
)

type UserHandler struct {
	log *logger.Logger
	cfg *config.CoreConfig
	s   *userservice.UserService
}

func NewUserHandler(log *logger.Logger, cfg *config.CoreConfig, s *userservice.UserService) *UserHandler {
	return &UserHandler{
		log: log,
		cfg: cfg,
		s:   s,
	}
}

// login or register user with telegram auth data
func (h *UserHandler) LoginOrRegister(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*time.Duration(h.cfg.WebConfig.ReadTimeout))
	defer cancel()

	telegramInitData := r.URL.Query().Get("telegram_init_data")
	if telegramInitData == "" {
		utils.WriteErrResponse(w, utils.NewError(http.StatusUnprocessableEntity, "telegram_init_data is required", nil))
		return
	}
	userLoginResponse, err := h.s.LoginOrRegister(ctx, telegramInitData)
	if err != nil {
		utils.LogResponseErr(r, h.log, err)
		utils.WriteErrResponse(w, err)
		return
	}

	utils.WriteJsonResponse(w, responseDataNameUsers, http.StatusOK, userLoginResponse)
}

// get user info by username or id
func (h *UserHandler) User(w http.ResponseWriter, r *http.Request) {
	// TODO implement me
	panic("implement me")
}
