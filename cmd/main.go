package main

import (
	"github.com/gin-gonic/gin"
	"github.com/nickhansel/nucleus/api/analytics"
	email3 "github.com/nickhansel/nucleus/api/analytics/email"
	"github.com/nickhansel/nucleus/api/auth"
	campaign "github.com/nickhansel/nucleus/api/campaigns"
	email2 "github.com/nickhansel/nucleus/api/campaigns/email"
	"github.com/nickhansel/nucleus/api/campaigns/text"
	"github.com/nickhansel/nucleus/api/customers"
	"github.com/nickhansel/nucleus/api/customers/groups"
	"github.com/nickhansel/nucleus/api/generation"
	"github.com/nickhansel/nucleus/api/items"
	"github.com/nickhansel/nucleus/api/middleware"
	org "github.com/nickhansel/nucleus/api/organization"
	"github.com/nickhansel/nucleus/api/templates"
	"github.com/nickhansel/nucleus/api/transactions"
	"github.com/nickhansel/nucleus/config"
	"github.com/nickhansel/nucleus/cron/email"
	"github.com/nickhansel/nucleus/sendinblue"
	"github.com/nickhansel/nucleus/shopify"

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

	email.GetEmailCampaignAnalytics()
	email.ScheduleGetEmailBounces()
	//textCron.ScheduleGetTextBounces()

	//fbCron.ScheduleGetFBMetrics()

	//groups.GetTopCustomers(19)

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "ok",
		})
	})

	r.GET("/login", auth.LoginUser)
	r.POST("/signup", auth.SignUp)

	r.POST("/organization/:orgId/invite", middleware.JwtAuthMiddleware(), middleware.CheckOrgMiddleware(), auth.SendInviteEmail)
	r.POST("/organization/:orgId/invite/accept", auth.AcceptInvite)

	r.GET("/generate/text", middleware.JwtAuthMiddleware(), generation.GenerateText)

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
	r.GET("/customers/:orgId/customer/:customerId", middleware.JwtAuthMiddleware(), middleware.CheckOrgMiddleware(), customers.GetCustomerById)

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
	r.PUT("/campaigns/:orgId/text", middleware.JwtAuthMiddleware(), middleware.CheckOrgMiddleware(), campaign.UpdateSMSCampaign)
	r.POST("/campaigns/:orgId/email", middleware.JwtAuthMiddleware(), middleware.CheckOrgMiddleware(), email2.CreateEmailCampaign)
	r.PUT("/campaigns/:orgId/email", middleware.JwtAuthMiddleware(), middleware.CheckOrgMiddleware(), campaign.UpdateEmailCampaign)
	r.GET("/campaigns/:orgId", middleware.JwtAuthMiddleware(), middleware.CheckOrgMiddleware(), campaign.GetCampaign)
	r.PUT("/campaigns/:orgId", middleware.JwtAuthMiddleware(), middleware.CheckOrgMiddleware(), campaign.UpdateCampaign)
	r.GET("/campaigns/:orgId/all", middleware.JwtAuthMiddleware(), middleware.CheckOrgMiddleware(), campaign.GetAllCampaigns)

	r.POST("/templates/:orgId/email", middleware.JwtAuthMiddleware(), middleware.CheckOrgMiddleware(), templates.CreateEmailTemplate)
	r.GET("/templates/:orgId/email", middleware.JwtAuthMiddleware(), middleware.CheckOrgMiddleware(), templates.GetEmailTemplates)
	r.GET("/templates/:orgId/email/:id", middleware.JwtAuthMiddleware(), middleware.CheckOrgMiddleware(), templates.GetEmailTemplate)

	r.GET("/metrics/:orgId/email/:email_campaign_id", middleware.JwtAuthMiddleware(), middleware.CheckOrgMiddleware(), email3.GetEmailAnalytics)
	r.GET("/metrics/:orgId/totals", middleware.JwtAuthMiddleware(), middleware.CheckOrgMiddleware(), analytics.GetTotalRevenue)
	r.GET("/metrics/:orgId/item", middleware.JwtAuthMiddleware(), middleware.CheckOrgMiddleware(), analytics.GetReveneuByItem)

	r.GET("/organization/:orgId/items", middleware.JwtAuthMiddleware(), middleware.CheckOrgMiddleware(), items.GetAllItems)

	r.POST("/flows/:orgId", middleware.JwtAuthMiddleware(), middleware.CheckOrgMiddleware(), flows.CreateFlow)
	r.POST("/flows/:orgId/sms", middleware.JwtAuthMiddleware(), middleware.CheckOrgMiddleware(), apiFlows.ScheduleTextFlows)
	r.POST("/flows/:orgId/email", middleware.JwtAuthMiddleware(), middleware.CheckOrgMiddleware(), apiFlows.ScheduleEmailFlows)
	r.PUT("/flows/:orgId/status/:flowId", middleware.JwtAuthMiddleware(), middleware.CheckOrgMiddleware(), apiFlows.UpdateFlowStatus)
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
