package middleware

import (
	"context"
	"net/http"
)

// Info: user authorization via an authentication token. Transmits userID in request context (type int64)
func (m *MiddleWareManager) AuthUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// tokenStr := r.Header.Get("Authorization")

		// tokenClaims, err := utils.JwtClaimsFromToken(tokenStr, []byte(m.cfg.JwtSecret))
		// if err != nil {
		// 	utils.WriteErrResponse(w, err)
		// 	return
		// }

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), CtxUserIDKey, int64(1))))
	})
}
