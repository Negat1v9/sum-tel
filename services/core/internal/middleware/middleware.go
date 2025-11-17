package middleware

import (
	"net/http"

	"github.com/Negat1v9/sum-tel/shared/config"
)

const (
	// CtxUserIDKey - used to store the user id in the request context
	CtxUserIDKey int = 0
)

type Middleware func(http.Handler) http.Handler

// Info: used to build all the necessary middleware as a constructor
type MiddleWareManager struct {
	cfg *config.CoreConfig
}

func New(cfg *config.CoreConfig) *MiddleWareManager {
	return &MiddleWareManager{
		cfg: cfg,
	}
}

// Info: create middleware for all requests for CORS, logging and same
func (mw *MiddleWareManager) BasicMW() Middleware {
	return createStack(cors)
}

func createStack(xs ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := 0; i < len(xs); i++ {
			x := xs[i]
			next = x(next)
		}
		return next
	}
}
