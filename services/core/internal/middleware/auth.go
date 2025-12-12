package middleware

import (
	"context"
	"net/http"

	"github.com/Negat1v9/sum-tel/services/core/internal/model"
	"github.com/Negat1v9/sum-tel/services/core/pkg/utils"
)

// Info: user authorization via an authentication token. Transmits userID in request context (type int64)
func (m *MiddleWareManager) AuthUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenStr := r.Header.Get("Authorization")

		tokenClaims, err := utils.JwtClaimsFromToken(tokenStr, []byte(m.cfg.JwtSecret))
		if err != nil {
			utils.WriteErrResponse(w, err)
			return
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), CtxUserIDKey, tokenClaims.UserID)))
	})
}

func (m *MiddleWareManager) AuthAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenStr := r.Header.Get("Authorization")

		tokenClaims, err := utils.JwtClaimsFromToken(tokenStr, []byte(m.cfg.JwtSecret))
		if err != nil {
			utils.WriteErrResponse(w, err)
			return
		}

		// not admin - return error
		if tokenClaims.Role != model.RoleAdmin {
			utils.WriteErrResponse(w, utils.NewError(http.StatusForbidden, "access denied", nil))
			return
		}

		// Pass the user ID in the context for further processing

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), CtxUserIDKey, tokenClaims.UserID)))
	})
}
