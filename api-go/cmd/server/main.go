package main

import (
	"log"
	"net/http"

	"github.com/DanielRCor/talsory_interseguro/api-go/internal/client"
	"github.com/DanielRCor/talsory_interseguro/api-go/internal/config"
	"github.com/DanielRCor/talsory_interseguro/api-go/internal/httpapi"
)

func main() {
	cfg := config.Load()
	statisticsClient := client.NewStatisticsClient(cfg.NodeAPIURL, &http.Client{
		Timeout: cfg.HTTPClientTimeout,
	})

	app := httpapi.NewApp(httpapi.AppDependencies{
		Config:           cfg,
		StatisticsClient: statisticsClient,
	})

	log.Fatal(app.Listen(":" + cfg.Port))
}
