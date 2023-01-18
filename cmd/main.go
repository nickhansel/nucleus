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
	fbAcc "github.com/nickhansel/nucleus/fb/account"
	fb "github.com/nickhansel/nucleus/fb/ads"
	fbAud "github.com/nickhansel/nucleus/fb/audiences"
	"github.com/nickhansel/nucleus/sendgrid"
	"github.com/nickhansel/nucleus/twilio"
	"gorm.io/gen"
)

func main() {

	r := gin.Default()

	config.Connect()

	// pass middleware.JWT() to the r.Use function to use the middleware
	r.GET("/login", auth.LoginUser)

	r.POST("/organization", middleware.JwtAuthMiddleware(), org.CreateOrg)
	r.GET("/organization/:orgId", middleware.JwtAuthMiddleware(), middleware.CheckOrgMiddleware(), org.GetOrg)
	r.PUT("/organization/:orgId", middleware.JwtAuthMiddleware(), middleware.CheckOrgMiddleware(), org.UpdateOrg)

	r.GET("/customers/:orgId", middleware.JwtAuthMiddleware(), middleware.CheckOrgMiddleware(), customers.GetCustomers)
	// r.POST("/work", customers.CreateCustomerGroup)
	r.POST("/customers/:orgId", middleware.JwtAuthMiddleware(), middleware.CheckOrgMiddleware(), customers.CreateCustomerGroup)
	r.GET("/customers/:orgId/groups/:groupId", middleware.JwtAuthMiddleware(), middleware.CheckOrgMiddleware(), customers.GetCustomerGroup)

	r.GET("/purchases/:orgId", middleware.JwtAuthMiddleware(), middleware.CheckOrgMiddleware(), transactions.GetPurchases)

	r.POST("/sendgrid/:orgId", middleware.JwtAuthMiddleware(), middleware.CheckOrgMiddleware(), sendgrid.VerifySendgridEmail)
	r.POST("/sendgrid/:orgId/resend", middleware.JwtAuthMiddleware(), middleware.CheckOrgMiddleware(), sendgrid.ResendVerificationEmail)

	r.POST("/twilio/:orgId", middleware.JwtAuthMiddleware(), middleware.CheckOrgMiddleware(), twilio.RegisterOrgTwilioNumber)
	r.POST("/twilio/:orgId/send", middleware.JwtAuthMiddleware(), middleware.CheckOrgMiddleware(), twilio.SendText)

	r.POST("/fb/:orgId", middleware.JwtAuthMiddleware(), middleware.CheckOrgMiddleware(), fb.CreateCampaign)
	r.POST("/fb/:orgId/adset", middleware.JwtAuthMiddleware(), middleware.CheckOrgMiddleware(), fb.CreateAdSet)
	r.GET("/fb/:orgId/pagename", middleware.JwtAuthMiddleware(), middleware.CheckOrgMiddleware(), fbAcc.GetPageID)

	r.GET("/fb/:orgId/url", middleware.JwtAuthMiddleware(), middleware.CheckOrgMiddleware(), fb.CreateAd)
	// r.POST("/fb/:orgId/create_audience/:customer_group_id", middleware.JwtAuthMiddleware(), middleware.CheckOrgMiddleware(), fbAud.CreateCustomAudience)
	r.POST("/fb/:orgId/audiences/:customer_group_id", middleware.JwtAuthMiddleware(), middleware.CheckOrgMiddleware(), fbAud.CreateCustomAudience)

	// r.POST("/aws", aws.UploadImage)

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
