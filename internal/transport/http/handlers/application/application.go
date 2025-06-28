package application

import (
	lemonPort "github.com/EdisonTantra/lemonPajak/internal/core/port"
	"github.com/gin-gonic/gin"
)

var _ lemonPort.AppsHandler = (*Handler)(nil)

type Handler struct {
	svcApp lemonPort.ApplicationService
}

type HandlerOpts struct {
	SvcApp lemonPort.ApplicationService
}

func New(opts *HandlerOpts) *Handler {
	return &Handler{
		svcApp: opts.SvcApp,
	}
}

func (h *Handler) Mount(router *gin.Engine) {
	router.POST("/v1/e-faktur/validation", h.EFakturValidation())
}

func (h *Handler) EFakturValidation() gin.HandlerFunc {
	return func(gCtx *gin.Context) {
		//TODO validation for multipart headers

		//into usecase
		// response dto json
	}
}
