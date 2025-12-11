package userservice

import (
	"context"
	"fmt"
	"time"

	"github.com/Negat1v9/sum-tel/services/core/internal/model"
	"github.com/Negat1v9/sum-tel/services/core/internal/store"
	"github.com/Negat1v9/sum-tel/services/core/pkg/utils"
	initdata "github.com/telegram-mini-apps/init-data-golang"
)

var (
	maxInitDataAge = time.Hour * 12
)

type UserService struct {
	store *store.Storage

	tgToken   string
	jwtSecret []byte
}

func NewUserService(store *store.Storage, tgToken string, jwtSecret []byte) *UserService {
	return &UserService{
		store:     store,
		tgToken:   tgToken,
		jwtSecret: jwtSecret,
	}
}

func (s *UserService) LoginOrRegister(ctx context.Context, telegramInitData string) (*model.UserLoginResponse, error) {
	mn := "UserService.LoginOrRegister"
	err := initdata.Validate(telegramInitData, s.tgToken, maxInitDataAge)
	if err != nil {
		return nil, fmt.Errorf("%s.Validate: %w", mn, err)
	}

	initData, err := initdata.Parse(telegramInitData)
	if err != nil {
		return nil, fmt.Errorf("%s.Parse: %w", mn, err)
	}

	user, err := s.store.UserRepo().GetByTelegramID(ctx, initData.User.ID)
	if err != nil {
		return nil, fmt.Errorf("%s.GetByTelegramID: %w", mn, err)
	}

	token, err := utils.GenerateJwtToken(&utils.Claims{
		UserID: user.ID,
		Role:   user.Role,
	}, s.jwtSecret)
	if err != nil {
		return nil, fmt.Errorf("%s.GenerateJwtToken: %w", mn, err)
	}

	return &model.UserLoginResponse{
		Token: token,
	}, nil
}
