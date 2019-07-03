package main

import (
	handlers "github.com/andresvillavicenciowizeline/proxy-app/api/handlers"
	server "github.com/andresvillavicenciowizeline/proxy-app/api/server"
	utils "github.com/andresvillavicenciowizeline/proxy-app/api/utils"
)

func main() {
	utils.LoadEnv()
	app := server.SetUp()
	handlers.HandlerRedirection(app)
	server.RunServer(app)
}
