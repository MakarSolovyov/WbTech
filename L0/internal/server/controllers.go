package controllers

import (
	"encoding/json"
	"l0/internal/cache"
	"l0/internal/model"
	kafka "l0/internal/pkg"
	"strings"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	InitController *kafka.InitController
	LRU            *cache.LRUCache
}

func (ic *Controller) GenerateOrders(c *gin.Context) {

	var orderCount = 10

	msg := kafka.DoRequest(ic.InitController.Producer, ic.InitController.Worker, orderCount,
		"post_order", "post_order_response")

	if !strings.Contains(msg, "Сгенерировано заказов: ") {
		c.JSON(200, gin.H{
			"error":   true,
			"message": msg,
		})

		return
	}

	c.JSON(200, gin.H{
		"error":   false,
		"message": msg,
	})
}

func (ic *Controller) GetOrderById(c *gin.Context) {
	uid := c.Param("order_uid")

	var order model.Order

	value, found := ic.LRU.Get(uid)
	if found {
		order = value.(model.Order)
	} else {
		msg := kafka.DoRequest(ic.InitController.Producer, ic.InitController.Worker,
			uid, "get_order_by_id", "get_order_by_id_response")

		err := json.Unmarshal([]byte(msg), &order)
		if err != nil {
			//log.Println(err)

			c.JSON(200, gin.H{
				"error":   true,
				"message": msg,
			})

			return
		}

		ic.LRU.Add(order.Order_uid, order)
	}

	c.JSON(200, gin.H{
		"error":   false,
		"message": order,
	})
}
