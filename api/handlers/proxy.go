package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/andresvillavicenciowizeline/proxy-app/api/middleware"
	"github.com/kataras/iris"
)

// HandlerRedirection should redirect traffic
func HandlerRedirection(app *iris.Application) {
	app.Get("/ping", middleware.ProxyMiddleware, proxyHandler)
}

var domains []string

func proxyHandler(c iris.Context) {

	for _, item := range middleware.FinalQueue {
		fmt.Printf("Domain: %s Weight: %d\n", item.Domain, item.Weight)
		domains = append(domains, item.Domain)
	}
	response, err := json.Marshal(domains)

	if err != nil {
		c.JSON(iris.Map{"status": 400, "result": "parse error"})
		return
	}
	c.JSON(iris.Map{"result": string(response)})
}
