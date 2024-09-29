package server

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/gin-gonic/gin"

	order "github.com/CamiloC-pvt/re-challenge/app/order/interface"
	pack "github.com/CamiloC-pvt/re-challenge/app/pack/interface"
)

type GinServer struct {
	port   int
	server *gin.Engine
}

var (
	ginServerInstance *GinServer
	ginServerOnce     sync.Once
)

func NewGinServer(port int) IServer {
	ginServerOnce.Do(func() {
		ginServerInstance = &GinServer{}
		ginServerInstance.port = port
		ginServerInstance.server = gin.New()
	})

	return ginServerInstance
}

// DefineRoutes define the routes for the Gin server
func (g *GinServer) DefineRoutes() {
	webAppGroupID := g.server.Group("")
	v1ApiGroupID := g.server.Group("api/v1")

	g.server.NoRoute(g.notFound())

	// API
	order.NewOrderGinApiController(v1ApiGroupID)
	pack.NewPackController(v1ApiGroupID)

	// Web App
	order.NewOrderGinWebController(webAppGroupID)
}

func (g *GinServer) Start() {
	err := g.server.Run(fmt.Sprintf("%s:%d", "", g.port))
	if err != nil {
		log.Panicf("There was an error starting the GinServer: %s", err.Error())
		os.Exit(1)
	}
}

// Private
func (g *GinServer) notFound() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		htmlFile, err := os.Open("./app/web/notFound.html")
		if err != nil {
			fmt.Println(err.Error())
		}
		defer htmlFile.Close()

		byteValue, _ := io.ReadAll(htmlFile)
		response := string(byteValue)

		ctx.Writer.WriteHeader(http.StatusNotFound)
		ctx.Writer.Write([]byte(response))
	}
}
