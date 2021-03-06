package main

import "github.com/poncheska/iot-mousetrap/pkg/api"


// @title Smart Mousetrap
// @version 1.0
// @description Server for the IOT project Smart Mousetrap

// @host smart-mousetrap.herokuapp.com
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	api.Start()
}
