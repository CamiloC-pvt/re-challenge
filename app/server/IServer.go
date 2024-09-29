package server

type IServer interface {
	DefineRoutes()
	Start()
}
