package order

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"

	db "github.com/CamiloC-pvt/re-challenge/app/db"
	order_business "github.com/CamiloC-pvt/re-challenge/app/order/business"
	order_domain "github.com/CamiloC-pvt/re-challenge/app/order/domain"
	order_infraestructure "github.com/CamiloC-pvt/re-challenge/app/order/infraestructure"
	pack_infraestructure "github.com/CamiloC-pvt/re-challenge/app/pack/infraestructure"
)

const GIN_API_CONTROLLER_ROUTE = "/order"

var (
	orderGinApiControllerInstance *OrderGinApiController
	orderGinApiControllerOnce     sync.Once
)

type OrderGinApiController struct {
	orderBusiness order_domain.IOrderBusiness
}

func NewOrderGinApiController(router *gin.RouterGroup) order_domain.IOrderController {
	orderGinApiControllerOnce.Do(func() {
		postgresConnection := db.NewPostgresConnection()

		orderPostgresRepo := order_infraestructure.NewOrderPostgresRepo(postgresConnection)
		packPostgresRepo := pack_infraestructure.NewPackPostgresRepo(postgresConnection)

		orderGinApiControllerInstance = &OrderGinApiController{}
		orderGinApiControllerInstance.orderBusiness = order_business.NewOrderBusiness(orderPostgresRepo, packPostgresRepo)

		orderGinApiControllerInstance.SetRoutes(router)
	})

	return orderGinApiControllerInstance
}

func (c *OrderGinApiController) SetRoutes(router *gin.RouterGroup) {
	router.DELETE(fmt.Sprintf("%s/cancel", GIN_API_CONTROLLER_ROUTE), c.Cancel())
	router.OPTIONS(fmt.Sprintf("%s/cancel", GIN_API_CONTROLLER_ROUTE), c.Cancel())

	router.GET(fmt.Sprintf("%s/", GIN_API_CONTROLLER_ROUTE), c.GetAll())
	router.OPTIONS(fmt.Sprintf("%s/", GIN_API_CONTROLLER_ROUTE), c.GetAll())

	router.POST(fmt.Sprintf("%s/create", GIN_API_CONTROLLER_ROUTE), c.Create())
	router.OPTIONS(fmt.Sprintf("%s/create", GIN_API_CONTROLLER_ROUTE), c.Create())
}

func (c *OrderGinApiController) Cancel() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		strOrderID := ctx.Query("order_id")

		orderID, err := strconv.ParseInt(strOrderID, 10, 32)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("wrong order ID: %s", strOrderID)})
			return
		}

		err = c.orderBusiness.Cancel(int32(orderID))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		} else {
			ctx.JSON(http.StatusOK, gin.H{"response": "Order Canceled"})
		}
	}
}

func (c *OrderGinApiController) Create() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		strOrderSize := ctx.Query("order_size")

		orderSize, err := strconv.ParseInt(strOrderSize, 10, 32)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("wrong order size: %s", strOrderSize)})
			return
		}

		order, err := c.orderBusiness.Create(int32(orderSize))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		} else {
			ctx.JSON(http.StatusOK, order)
		}
	}
}

func (c *OrderGinApiController) GetAll() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		orders, err := c.orderBusiness.GetAll()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		} else {
			ctx.JSON(http.StatusOK, orders)
		}
	}
}
