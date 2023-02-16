package main

import (
	"github.com/nickhansel/nucleus/api/analytics"
	email3 "github.com/nickhansel/nucleus/api/analytics/email"
	email2 "github.com/nickhansel/nucleus/api/campaigns/email"
	"github.com/nickhansel/nucleus/api/campaigns/text"
	"github.com/nickhansel/nucleus/api/customers/groups"
	"github.com/nickhansel/nucleus/config"
	"github.com/nickhansel/nucleus/segmentQL"
	"github.com/nickhansel/nucleus/sendinblue"
	"github.com/nickhansel/nucleus/shopify"

	"github.com/gin-gonic/gin"
	"github.com/nickhansel/nucleus/api/auth"
	campaign "github.com/nickhansel/nucleus/api/campaigns"
	"github.com/nickhansel/nucleus/api/customers"
	"github.com/nickhansel/nucleus/api/middleware"
	org "github.com/nickhansel/nucleus/api/organization"
	"github.com/nickhansel/nucleus/api/transactions"

	apiFlows "github.com/nickhansel/nucleus/api/flows"
	fbAcc "github.com/nickhansel/nucleus/fb/account"
	fb "github.com/nickhansel/nucleus/fb/ads"
	fbAud "github.com/nickhansel/nucleus/fb/audiences"
	"github.com/nickhansel/nucleus/flows"
	"github.com/nickhansel/nucleus/twilio"
)

func main() {

	r := gin.Default()

	config.Connect()

	//email.GetEmailCampaignAnalytics()
	//email.ScheduleGetEmailBounces()
	//textCron.ScheduleGetTextBounces()

	groups.GetTopCustomers(19)

	segmentQL.Parse(838221504934608897, 838194565431033857, "2023-01-08T04:03:54.895Z", "2023-02-16T04:03:54.895Z", 10, 1000)

	r.GET("/login", auth.LoginUser)
	r.POST("/signup", auth.SignUp)

	r.POST("/shopify/:orgId/oauth", middleware.JwtAuthMiddleware(), middleware.CheckOrgMiddleware(), shopify.Oauth)
	r.GET("/shopify/:orgId/load", middleware.JwtAuthMiddleware(), middleware.CheckOrgMiddleware(), shopify.GetInitialData)

	r.POST("/organization", middleware.JwtAuthMiddleware(), org.CreateOrg)
	r.GET("/organization/:orgId", middleware.JwtAuthMiddleware(), middleware.CheckOrgMiddleware(), org.GetOrg)
	r.PUT("/organization/:orgId", middleware.JwtAuthMiddleware(), middleware.CheckOrgMiddleware(), org.UpdateOrg)

	r.GET("/customers/:orgId", middleware.JwtAuthMiddleware(), middleware.CheckOrgMiddleware(), customers.GetCustomers)
	// r.POST("/work", customers.CreateCustomerGroup)
	r.POST("/customers/:orgId", middleware.JwtAuthMiddleware(), middleware.CheckOrgMiddleware(), groups.CreateCustomerGroup)
	r.POST("/ql/:orgId/", middleware.JwtAuthMiddleware(), middleware.CheckOrgMiddleware(), groups.CreateCustomerGroupSegmentQL)
	r.GET("/ql/:orgId/segment", middleware.JwtAuthMiddleware(), middleware.CheckOrgMiddleware(), groups.SegmentCustomers)
	r.GET("/customers/:orgId/groups/:groupId", middleware.JwtAuthMiddleware(), middleware.CheckOrgMiddleware(), groups.GetCustomerGroup)

	r.GET("/customers/:orgId/groups", middleware.JwtAuthMiddleware(), middleware.CheckOrgMiddleware(), groups.ListCustomerGroups)

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
	r.GET("/campaigns/:orgId/all", middleware.JwtAuthMiddleware(), middleware.CheckOrgMiddleware(), campaign.GetAllCampaigns)

	r.GET("/metrics/:orgId/email/:email_campaign_id", middleware.JwtAuthMiddleware(), middleware.CheckOrgMiddleware(), email3.GetEmailAnalytics)
	r.GET("/metrics/:orgId/totals", middleware.JwtAuthMiddleware(), middleware.CheckOrgMiddleware(), analytics.GetTotalRevenue)
	r.GET("/metrics/:orgId/item", middleware.JwtAuthMiddleware(), middleware.CheckOrgMiddleware(), analytics.GetReveneuByItem)

	r.POST("/flows/:orgId", middleware.JwtAuthMiddleware(), middleware.CheckOrgMiddleware(), flows.CreateFlow)
	r.POST("/flows/:orgId/sms", middleware.JwtAuthMiddleware(), middleware.CheckOrgMiddleware(), apiFlows.ScheduleTextFlows)
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
