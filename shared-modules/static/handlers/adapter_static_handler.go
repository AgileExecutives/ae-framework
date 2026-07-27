package handlers

import (
    "github.com/AgileExecutives/serverbase/pkg/core"
    "github.com/gin-gonic/gin"
)

// StaticHandlerAdapter wraps the shared static handlers and provides a
// ModuleContext-aware constructor similar to the serverbase helper.
type StaticHandlerAdapter struct{
    delegate *StaticHandlers
}

func NewStaticHandlerWithCtx(ctx core.ModuleContext) *StaticHandlerAdapter {
    var logger core.Logger = ctx.Logger
    if logger == nil {
        logger = core.NewLogger()
    }
    repo := NewFSStaticRepo("./statics/json")
    return &StaticHandlerAdapter{delegate: NewStaticHandlers(logger, repo)}
}

func (h *StaticHandlerAdapter) ServeStaticJSON(c *gin.Context) { h.delegate.ServeStaticJSON(c) }
func (h *StaticHandlerAdapter) ListStaticJSON(c *gin.Context)  { h.delegate.ListStaticJSON(c) }

// Legacy convenience functions
var legacyStatic = NewStaticHandlerWithCtx(core.ModuleContext{})

func ServeStaticJSON(c *gin.Context) { legacyStatic.ServeStaticJSON(c) }
func ListStaticJSON(c *gin.Context)  { legacyStatic.ListStaticJSON(c) }
