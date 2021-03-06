package main

import (
	"fmt"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/labstack/echo"
	"github.com/magrandera/SPaaS/common"
	"github.com/magrandera/SPaaS/config"
	"github.com/magrandera/SPaaS/server/controller"
	"github.com/magrandera/SPaaS/server/routing"
)

func initialize(e *echo.Echo) {
	config.New(filepath.Join(common.HomeDir(), ".spaas"), ".spaas.json")
	if err := config.Save(); err != nil {
		fmt.Println(err.Error())
	}
	config.Cfg.Config.WatchConfig()
	config.Cfg.Config.OnConfigChange(func(_ fsnotify.Event) {
		fmt.Println("Config file changed")
	})
	routing.GlobalMiddleware(e)
	routing.SetupRoutes(e)
	controller.InitDocker()
	routing.InitReverseProxy()
}

func main() {
	e := echo.New()
	initialize(e)
	e.Logger.Fatal(e.Start(":8080"))
}
