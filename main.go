package main

import (
	"fmt"

	"github.com/fc-testzone/apiserver/auth"
	"github.com/fc-testzone/apiserver/content"
	"github.com/fc-testzone/apiserver/net"
	"github.com/fc-testzone/apiserver/utils"

	"go.uber.org/dig"
)

func main() {
	container := dig.New()

	container.Provide(utils.NewLog)
	container.Provide(utils.NewConfigs)
	container.Provide(auth.NewAuthorizer)
	container.Provide(content.NewContent)
	container.Provide(net.NewWebServer)
	container.Provide(NewApp)

	err := container.Invoke(func(app *App) {
		app.Start()
	})

	if err != nil {
		fmt.Printf("Fatal error: %s\n", err.Error())
	}
}
