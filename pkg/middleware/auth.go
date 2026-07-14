package middleware

import (
	baseMw "github.com/AgileExecutives/serverbase/modules/base/middleware"
	"github.com/AgileExecutives/serverbase/pkg/core"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware returns authentication middleware constructed from a
// `core.ModuleContext` by delegating to the base module's middleware provider.
func AuthMiddleware(ctx core.ModuleContext) gin.HandlerFunc {
	p := baseMw.NewAuthMiddleware(ctx, ctx.Logger)
	prov := baseMw.NewAuthMiddlewareProvider(p)
	return prov.Handler()
}

// AuthOptions mirrors legacy auth options
type AuthOptions struct {
	SingleTenant bool
}

// AuthMiddlewareWithOptions constructs middleware using explicit options.
func AuthMiddlewareWithOptions(ctx core.ModuleContext, opts AuthOptions) gin.HandlerFunc {
	return baseMw.BuildAuthHandler(ctx.DB, ctx.Logger, baseMw.Options{SingleTenant: opts.SingleTenant})
}
