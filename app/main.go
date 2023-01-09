package main

import (
	"net/http"

	"gorm.io/gorm"

	"github.com/nickhansel/nucleus/pkg/config"

	"github.com/gin-gonic/gin"
	"github.com/nickhansel/nucleus/pkg/api"
	"github.com/nickhansel/nucleus/pkg/model"
	"gorm.io/gen"
)

func main() {

	r := gin.Default()

	config.Connect()

	r.GET("/customers/:orgId", api.GetCustomers)
	r.GET("/purchases", func(c *gin.Context) {
		purchases := []model.Purchase{}
		// add the purchased items and the variation related to the purchased item to the purchase struct
		err := config.DB.Preload("PurchasedItems").Preload("PurchasedItems.Variation").Find(&purchases).Error

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"purchases": purchases,
		})
	})

	// use api.getCustomers to handle the request

	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})

	r.Run() // listen and serve on 0.0.0.0:8080
}

// function to generate the strcut for the db table
func generateTable(db *gorm.DB) {
	// generate struct for the db table using the gorm gen package
	g := gen.NewGenerator(gen.Config{
		OutPath: "../query",
		Mode:    gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface, // generate mode
	})

	// gormdb, _ := gorm.Open(mysql.Open("root:@(127.0.0.1:3306)/demo?charset=utf8mb4&parseTime=True&loc=Local"))
	g.UseDB(db) // reuse your gorm db

	g.GenerateAllTable()

	g.Execute() // generate

}
