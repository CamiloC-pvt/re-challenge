package pack

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"

	db "github.com/CamiloC-pvt/re-challenge/app/db"
	pack_business "github.com/CamiloC-pvt/re-challenge/app/pack/business"
	pack_domain "github.com/CamiloC-pvt/re-challenge/app/pack/domain"
	pack_infraestructure "github.com/CamiloC-pvt/re-challenge/app/pack/infraestructure"
)

const CONTROLLER_ROUTE = "/pack"

var (
	packGinControllerInstance *PackGinController
	packGinControllerOnce     sync.Once
)

type PackGinController struct {
	packBusiness pack_domain.IPackBusiness
}

func NewPackController(router *gin.RouterGroup) pack_domain.IPackController {
	packGinControllerOnce.Do(func() {
		postgresConnection := db.NewPostgresConnection()

		packPostgresRepo := pack_infraestructure.NewPackPostgresRepo(postgresConnection)

		packGinControllerInstance = &PackGinController{}
		packGinControllerInstance.packBusiness = pack_business.NewPackBusiness(packPostgresRepo)

		packGinControllerInstance.SetRoutes(router)
	})

	return packGinControllerInstance
}

func (c *PackGinController) SetRoutes(router *gin.RouterGroup) {
	router.DELETE(fmt.Sprintf("%s/delete", CONTROLLER_ROUTE), c.Delete())
	router.OPTIONS(fmt.Sprintf("%s/delete", CONTROLLER_ROUTE), c.Delete())

	router.GET(fmt.Sprintf("%s/", CONTROLLER_ROUTE), c.GetAll())
	router.OPTIONS(fmt.Sprintf("%s/", CONTROLLER_ROUTE), c.GetAll())

	router.POST(fmt.Sprintf("%s/create", CONTROLLER_ROUTE), c.Create())
	router.OPTIONS(fmt.Sprintf("%s/create", CONTROLLER_ROUTE), c.Create())
}

func (c *PackGinController) Create() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		strPackID := ctx.Query("size")

		packSize, err := strconv.ParseInt(strPackID, 10, 32)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("wrong pack size: %s", strPackID)})
			return
		}

		newPackID, err := c.packBusiness.Create(int32(packSize))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		} else {
			ctx.JSON(http.StatusOK, gin.H{"pack_id": newPackID})
		}
	}
}

func (c *PackGinController) Delete() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		strPackID := ctx.Query("pack_id")

		packID, err := strconv.ParseInt(strPackID, 10, 32)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("wrong pack ID: %s", strPackID)})
			return
		}

		err = c.packBusiness.Delete(int32(packID))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		} else {
			ctx.JSON(http.StatusOK, gin.H{"response": "Done"})
		}
	}
}

func (c *PackGinController) GetAll() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		packs, err := c.packBusiness.GetAll()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		} else {
			ctx.JSON(http.StatusOK, packs)
		}
	}
}
