package app

import (
	"sync"

	db "github.com/CamiloC-pvt/re-challenge/app/db"
	server "github.com/CamiloC-pvt/re-challenge/app/server"
)

type ReChallenge struct {
	postgresConnection *db.PostgresConnection
	server             server.IServer
}

var (
	app  *ReChallenge
	once sync.Once
)

func InitReChallenge(port int) *ReChallenge {
	once.Do(func() {
		app = &ReChallenge{}
		app.postgresConnection = db.NewPostgresConnection()
		app.server = server.NewGinServer(port)
	})

	return app
}

// StartReChallenge initialize the connection to the DB, define the routes of the server and start it
func (a *ReChallenge) StartReChallenge() {
	a.postgresConnection.Connect()
	a.server.DefineRoutes()
	a.server.Start()
}
