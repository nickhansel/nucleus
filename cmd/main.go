package main

import (
	"gorm.io/gorm"

	"github.com/nickhansel/nucleus/config"

	"github.com/gin-gonic/gin"
	// "github.com/nickhansel/nucleus/api"
	"github.com/nickhansel/nucleus/api/auth"
	"github.com/nickhansel/nucleus/api/customers"
	"github.com/nickhansel/nucleus/api/middleware"
	org "github.com/nickhansel/nucleus/api/organization"
	"github.com/nickhansel/nucleus/api/transactions"
	"gorm.io/gen"
)

func main() {

	r := gin.Default()

	config.Connect()

	// pass middleware.JWT() to the r.Use function to use the middleware
	r.GET("/login", auth.LoginUser)
	r.GET("/organization/:orgId", middleware.JwtAuthMiddleware(), middleware.CheckOrgMiddleware(), org.GetOrg)
	r.POST("/organization/:orgId", middleware.JwtAuthMiddleware(), org.CreateOrg)
	r.GET("/customers/:orgId", middleware.JwtAuthMiddleware(), middleware.CheckOrgMiddleware(), customers.GetCustomers)
	r.GET("/purchases/:orgId", middleware.JwtAuthMiddleware(), middleware.CheckOrgMiddleware(), transactions.GetPurchases)

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
