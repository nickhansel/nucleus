package templates

import (
	"github.com/gin-gonic/gin"
	"github.com/nickhansel/nucleus/config"
	"github.com/nickhansel/nucleus/model"
)

type EmailTemplateBody struct {
	Name string `json:"name"`
	HTML string `json:"html"`
}

func CreateEmailTemplate(c *gin.Context) {
	org := c.MustGet("orgs").(model.Organization)

	var body EmailTemplateBody
	err := c.BindJSON(&body)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	var emailTemplate model.EmailTemplate
	emailTemplate.Name = body.Name
	emailTemplate.Html = body.HTML
	emailTemplate.OrganizationID = org.ID

	config.DB.Create(&emailTemplate)

	c.JSON(200, emailTemplate)
}

func GetEmailTemplates(c *gin.Context) {
	org := c.MustGet("orgs").(model.Organization)

	var emailTemplates []model.EmailTemplate
	config.DB.Where("\"organizationId\" = ?", org.ID).Find(&emailTemplates)

	c.JSON(200, emailTemplates)
}

func GetEmailTemplate(c *gin.Context) {
	org := c.MustGet("orgs").(model.Organization)

	id := c.Param("id")

	var emailTemplate model.EmailTemplate
	config.DB.Where("\"organizationId\" = ? AND \"id\" = ?", org.ID, id).Find(&emailTemplate)

	if emailTemplate.ID == 0 {
		c.JSON(404, gin.H{"error": "Email template not found"})
		return
	}

	c.JSON(200, emailTemplate)
}
