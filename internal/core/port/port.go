package port

import (
	"github.com/gin-gonic/gin"
)

type AppsRepository interface {
	Close() error
	RegisterStore()
	GetUserStore() UserStore
	GetAppStore() ApplicationStore
}

type AppsHandler interface {
	Mount(router *gin.Engine)
}
