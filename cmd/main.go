package main

import (
	"github.com/nickhansel/nucleus/api/analytics"
	email3 "github.com/nickhansel/nucleus/api/analytics/email"
	email2 "github.com/nickhansel/nucleus/api/campaigns/email"
	"github.com/nickhansel/nucleus/api/campaigns/text"
	"github.com/nickhansel/nucleus/config"
	"github.com/nickhansel/nucleus/cron/email"
	textCron "github.com/nickhansel/nucleus/cron/text"
	"github.com/nickhansel/nucleus/sendinblue"

	"github.com/gin-gonic/gin"
	// "github.com/nickhansel/nucleus/api"
	"github.com/nickhansel/nucleus/api/auth"
	campaign "github.com/nickhansel/nucleus/api/campaigns"
	"github.com/nickhansel/nucleus/api/customers"
	"github.com/nickhansel/nucleus/api/middleware"
	org "github.com/nickhansel/nucleus/api/organization"
	"github.com/nickhansel/nucleus/api/transactions"

	// "github.com/nickhansel/nucleus/cron"
	fbAcc "github.com/nickhansel/nucleus/fb/account"
	fb "github.com/nickhansel/nucleus/fb/ads"
	fbAud "github.com/nickhansel/nucleus/fb/audiences"
	"github.com/nickhansel/nucleus/twilio"
)

func main() {

	r := gin.Default()

	config.Connect()

	email.GetEmailCampaignAnalytics()
	email.ScheduleGetEmailBounces()
	textCron.ScheduleGetTextBounces()

	//segmentQL.FindCustomersWhoPurchasedItem(4, 19, "", "", 14, 0)

	// cron.ScheduleTask("2023-01-22 11:27:10")
	// 2023-01-13 20:04:27.299298 -0600 CST m=+36.150158126
	// pass middleware.JWT() to the r.Use function to use the middleware
	r.GET("/login", auth.LoginUser)

	r.POST("/organization", middleware.JwtAuthMiddleware(), org.CreateOrg)
	r.GET("/organization/:orgId", middleware.JwtAuthMiddleware(), middleware.CheckOrgMiddleware(), org.GetOrg)
	r.PUT("/organization/:orgId", middleware.JwtAuthMiddleware(), middleware.CheckOrgMiddleware(), org.UpdateOrg)

	r.GET("/customers/:orgId", middleware.JwtAuthMiddleware(), middleware.CheckOrgMiddleware(), customers.GetCustomers)
	// r.POST("/work", customers.CreateCustomerGroup)
	r.POST("/customers/:orgId", middleware.JwtAuthMiddleware(), middleware.CheckOrgMiddleware(), customers.CreateCustomerGroup)
	r.POST("/ql/:orgId/", middleware.JwtAuthMiddleware(), middleware.CheckOrgMiddleware(), customers.CreateCustomerGroupSegmentQL)
	r.GET("/customers/:orgId/groups/:groupId", middleware.JwtAuthMiddleware(), middleware.CheckOrgMiddleware(), customers.GetCustomerGroup)

	r.GET("/purchases/:orgId", middleware.JwtAuthMiddleware(), middleware.CheckOrgMiddleware(), transactions.GetPurchases)

	r.POST("/organization/:orgId/verify_email", middleware.JwtAuthMiddleware(), middleware.CheckOrgMiddleware(), sendinblue.SendVerifyEmail)
	r.POST("/organization/:orgId/verify_code", middleware.JwtAuthMiddleware(), middleware.CheckOrgMiddleware(), sendinblue.Verify)

	r.POST("/twilio/:orgId", middleware.JwtAuthMiddleware(), middleware.CheckOrgMiddleware(), twilio.RegisterOrgTwilioNumber)
	r.POST("/twilio/:orgId/send", middleware.JwtAuthMiddleware(), middleware.CheckOrgMiddleware(), twilio.SendTextAPI)
	//r.GET("/twilio/:orgId/messages", middleware.JwtAuthMiddleware(), middleware.CheckOrgMiddleware(), twilio.GetMessages)

	r.POST("/fb/:orgId", middleware.JwtAuthMiddleware(), middleware.CheckOrgMiddleware(), fb.CreateCampaign)
	r.POST("/fb/:orgId/adset", middleware.JwtAuthMiddleware(), middleware.CheckOrgMiddleware(), fb.CreateAdSet)
	r.GET("/fb/:orgId/pagename", middleware.JwtAuthMiddleware(), middleware.CheckOrgMiddleware(), fbAcc.GetPageID)

	r.GET("/fb/:orgId/url", middleware.JwtAuthMiddleware(), middleware.CheckOrgMiddleware(), fb.CreateAd)
	// r.POST("/fb/:orgId/create_audience/:customer_group_id", middleware.JwtAuthMiddleware(), middleware.CheckOrgMiddleware(), fbAud.CreateCustomAudience)
	r.POST("/fb/:orgId/audiences/:customer_group_id", middleware.JwtAuthMiddleware(), middleware.CheckOrgMiddleware(), fbAud.CreateCustomAudience)
	r.PUT("/fb/:orgId/audiences/:customer_group_id", middleware.JwtAuthMiddleware(), middleware.CheckOrgMiddleware(), fbAud.UpdateCustomAudience)

	r.POST("/campaigns/:orgId/text", middleware.JwtAuthMiddleware(), middleware.CheckOrgMiddleware(), text.CreateTextCampaign)
	r.POST("/campaigns/:orgId/email", middleware.JwtAuthMiddleware(), middleware.CheckOrgMiddleware(), email2.CreateEmailCampaign)
	r.GET("/campaigns/:orgId", middleware.JwtAuthMiddleware(), middleware.CheckOrgMiddleware(), campaign.GetCampaign)

	r.GET("/metrics/:orgId/email/:email_campaign_id", middleware.JwtAuthMiddleware(), middleware.CheckOrgMiddleware(), email3.GetEmailAnalytics)
	r.GET("/metrics/:orgId/totals", middleware.JwtAuthMiddleware(), middleware.CheckOrgMiddleware(), analytics.GetTotalRevenue)
	// r.POST("/aws", aws.UploadImage)

	// use api.getCustomers to handle the request

	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})

	err := r.Run()
	if err != nil {
		return
	} // listen and serve on 0.0.0.0:8080
}
