package order

import (
	"fmt"
	"html/template"
	"net/http"
	"sort"
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	db "github.com/CamiloC-pvt/re-challenge/app/db"
	order_business "github.com/CamiloC-pvt/re-challenge/app/order/business"
	order_domain "github.com/CamiloC-pvt/re-challenge/app/order/domain"
	order_infraestructure "github.com/CamiloC-pvt/re-challenge/app/order/infraestructure"
	pack_business "github.com/CamiloC-pvt/re-challenge/app/pack/business"
	pack_domain "github.com/CamiloC-pvt/re-challenge/app/pack/domain"
	pack_infraestructure "github.com/CamiloC-pvt/re-challenge/app/pack/infraestructure"
)

const GIN_WEB_CONTROLLER_ROUTE = "/"

var (
	orderGinWebControllerInstance *OrderGinWebController
	orderGinWebControllerOnce     sync.Once
)

type OrderGinWebController struct {
	orderBusiness order_domain.IOrderBusiness
	packBusiness  pack_domain.IPackBusiness
}

func NewOrderGinWebController(router *gin.RouterGroup) order_domain.IOrderController {
	orderGinWebControllerOnce.Do(func() {
		postgresConnection := db.NewPostgresConnection()

		orderPostgresRepo := order_infraestructure.NewOrderPostgresRepo(postgresConnection)
		packPostgresRepo := pack_infraestructure.NewPackPostgresRepo(postgresConnection)

		orderGinWebControllerInstance = &OrderGinWebController{}
		orderGinWebControllerInstance.orderBusiness = order_business.NewOrderBusiness(orderPostgresRepo, packPostgresRepo)
		orderGinWebControllerInstance.packBusiness = pack_business.NewPackBusiness(orderPostgresRepo, packPostgresRepo)

		orderGinWebControllerInstance.SetRoutes(router)
	})

	return orderGinWebControllerInstance
}

func (c *OrderGinWebController) SetRoutes(router *gin.RouterGroup) {
	router.GET(fmt.Sprintf("%s/", GIN_WEB_CONTROLLER_ROUTE), c.webAppHandler())
	router.OPTIONS(fmt.Sprintf("%s/", GIN_WEB_CONTROLLER_ROUTE), c.webAppHandler())
}

func (c *OrderGinWebController) webAppHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		type OrderFe struct {
			Created        time.Time
			Deleted        bool
			ID             int32
			Modified       time.Time
			Packs          []order_domain.OrderPack
			Size           int32
			TotalDelivered int32
		}

		type WebAppData struct {
			AvailablePacks []pack_domain.Pack
			Orders         []OrderFe
		}

		// Orders
		orders, err := c.orderBusiness.GetAll()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Packs
		availablePacks, err := c.packBusiness.GetAll()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		sort.Slice(availablePacks, func(i, j int) bool {
			return availablePacks[i].Size < availablePacks[j].Size
		})

		ordersTemplate, err := template.ParseFiles("./app/order/web/index.html")
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		feOrders := []OrderFe{}
		for _, dbOrder := range orders {
			orderDelivered := int32(0)

			for _, orderPackage := range dbOrder.Packs {
				orderDelivered += orderPackage.Amount * orderPackage.Size
			}

			feOrders = append(feOrders, OrderFe{
				Created:        dbOrder.Created,
				Deleted:        dbOrder.Deleted,
				ID:             dbOrder.ID,
				Modified:       dbOrder.Modified,
				Packs:          dbOrder.Packs,
				Size:           dbOrder.Size,
				TotalDelivered: orderDelivered,
			})
		}

		ordersTemplate.Execute(ctx.Writer, WebAppData{
			AvailablePacks: availablePacks,
			Orders:         feOrders,
		})
	}
}
