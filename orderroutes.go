package main

import (
	"bukidlink/db"
	"net/http"

	"github.com/gin-gonic/gin"
)

func updateOrderStatusHandler(c *gin.Context) {
	var id string = c.Query("id")
	var status string = c.Query("status")

	if err := db.UpdateOrderStatus(id, status); err != nil {
		retInternalServErr(err, c)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "order status updated"})
}

func getUsersOrdersHandler(c *gin.Context) {
	userId := c.Param("user_id")

	orders, err := db.QueryUsersOrders(userId)
	if err != nil {
		retInternalServErr(err, c)
		return
	}

	c.JSON(http.StatusOK, orders)
}

func postOrderHandler(c *gin.Context) {
	var order db.Order
	if err := c.BindJSON(&order); err != nil {
		retBadReqErr(err, c)
		return
	}

	orderId, err := db.InsertOrder(order)
	if err != nil {
		retInternalServErr(err, c)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "order created", "order_id": orderId})
}

func deleteOrderHandler(c *gin.Context) {
	orderId := c.Query("order_id")

	if err := db.DeleteOrder(orderId); err != nil {
		retInternalServErr(err, c)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "order deleted"})
}
