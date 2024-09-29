package order

import "github.com/gin-gonic/gin"

type IOrderController interface {
	SetRoutes(router *gin.RouterGroup)
}
