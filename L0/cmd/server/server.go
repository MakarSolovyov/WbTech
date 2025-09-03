package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"l0/internal/cache"
	kafka "l0/internal/pkg"
	controllers "l0/internal/server"

	"github.com/gin-gonic/gin"
)

func main() {

	var brokers []string = []string{"localhost:9092"}

	worker, err := kafka.CreateConsumer(brokers)
	if err != nil {
		panic(err)
	}
	defer worker.Close()
	producer, err := kafka.CreateProducer(brokers)
	if err != nil {
		panic(err)
	}
	defer producer.Close()

	LRUcache := cache.CreateLRU(100)

	cache.FillRecentData(LRUcache)
	initLen := len(LRUcache.Items)
	fmt.Printf("Загружен кэш из БД, количество загруженных заказов: %d\n", initLen)

	kafkaController := &kafka.InitController{
		Producer: producer,
		Worker:   worker,
	}

	controllerInstance := &controllers.Controller{
		InitController: kafkaController,
		LRU:            LRUcache,
	}

	gin.SetMode(gin.ReleaseMode)
	// gin.SetMode(gin.DebugMode)

	router := gin.New()

	staticDir := getFrontendDir()

	router.StaticFS("/static", http.Dir(staticDir))

	router.NoRoute(getMainPage)
	router.GET("/orders", getMainPage)
	router.GET("/", redirectOnMainPage)

	router.POST("/generate", controllerInstance.GenerateOrders)

	orderRoutes := router.Group("/orders")
	{
		orderRoutes.GET("/:order_uid", controllerInstance.GetOrderById)
	}

	if err := router.Run(":8080"); err != nil {
		log.Fatal(err.Error())
	}
}

func getMainPage(ctx *gin.Context) {
	ctx.File(filepath.Join(getFrontendDir(), "server.html"))
}

func redirectOnMainPage(ctx *gin.Context) {
	ctx.Redirect(http.StatusMovedPermanently, "/orders")
}

func getFrontendDir() string {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	frontendDir := filepath.Join(dir, "../../frontend")
	return frontendDir
}
