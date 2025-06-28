package user

import (
	"net/http"

	"github.com/EdisonTantra/lemonPajak/internal/core/cons"
	"github.com/EdisonTantra/lemonPajak/internal/core/domain"
	lemonPort "github.com/EdisonTantra/lemonPajak/internal/core/port"
	"github.com/EdisonTantra/lemonPajak/pkg/lib/logat"
	"github.com/EdisonTantra/lemonPajak/pkg/lib/tracer"
	"github.com/gin-gonic/gin"
)

var _ lemonPort.AppsHandler = (*Handler)(nil)

type Handler struct {
	userSvc lemonPort.UserService
}

type HandlerOpts struct {
	SvcUser lemonPort.UserService
}

func New(
	opts *HandlerOpts,
) *Handler {
	return &Handler{
		userSvc: opts.SvcUser,
	}
}

func (h *Handler) Mount(router *gin.Engine) {
	router.GET("/v1/user/:external_id", h.UserDetail())
}

func (h *Handler) UserDetail() gin.HandlerFunc {
	return func(gCtx *gin.Context) {
		ctx := gCtx.Request.Context()
		tr := tracer.StartTrace(ctx, "HandlerHTTP-UserDetail")
		defer tr.Finish()
		ctx = tr.Context()
		extID := gCtx.Param("external_id")

		logat.GetLogger().Info(ctx, "start user detail usecase", cons.EventLogNameUserDetail, map[string]string{
			"external_id": extID,
		})
		defer logat.GetLogger().Info(ctx, "end user detail usecase", cons.EventLogNameUserDetail, nil)

		if extID == "" {
			logat.GetLogger().Error(ctx, "error ext id empty", cons.EventLogNameUserDetail, nil)
			return
		}

		resp, err := h.userSvc.GetExternalUser(ctx, &domain.RequestGetUser{
			ExternalID: extID,
		})
		if err != nil {
			logat.GetLogger().Error(ctx, "error resp get", cons.EventLogNameUserDetail, err)
			return
		}

		gCtx.JSON(http.StatusOK, ResponseUserDetail{
			ID:          resp.ID,
			ExternalID:  resp.ExternalID,
			Username:    resp.Username,
			Email:       resp.Email,
			PhoneNumber: resp.PhoneNumber,
			FirstName:   resp.Profile.FirstName,
			LastName:    resp.Profile.LastName,
			Gender:      resp.Profile.Gender,
			Age:         resp.Profile.Age,
			Description: resp.Profile.Description,
			CreatedAt:   resp.CreatedAt,
			UpdatedAt:   resp.UpdatedAt,
		})
	}
}
