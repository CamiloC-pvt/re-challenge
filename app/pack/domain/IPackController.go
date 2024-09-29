package pack

import "github.com/gin-gonic/gin"

type IPackController interface {
	SetRoutes(router *gin.RouterGroup)
}
