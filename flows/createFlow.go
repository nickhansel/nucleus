package flows

import (
	"github.com/gin-gonic/gin"
	"github.com/nickhansel/nucleus/config"
	"github.com/nickhansel/nucleus/model"
)

type FlowBody struct {
	Name            string `json:"name"`
	CustomerGroupID int32  `json:"customer_group_id"`
	TriggerEvent    string `json:"trigger_event"`
	ActionType      string `json:"action_type"`
	ActionWaitTime  string `json:"action_wait_time"`
	SmsBody         string `json:"sms_body"`
	EmailSubject    string `json:"email_subject"`
	EmailBody       string `json:"email_body"`
}

func CreateFlow(c *gin.Context) {
	org := c.MustGet("orgs").(model.Organization)
	var flowBody FlowBody
	if err := c.ShouldBindJSON(&flowBody); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var flow model.Flow
	flow.Name = flowBody.Name
	flow.CustomerGroupID = flowBody.CustomerGroupID
	flow.OrganizationID = org.ID

	if err := config.DB.Create(&flow).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var trigger model.Trigger
	trigger.FlowID = flow.ID
	trigger.Event = flowBody.TriggerEvent

	if err := config.DB.Create(&trigger).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var action model.Action
	action.TriggerID = trigger.ID
	action.Type = flowBody.ActionType
	action.WaitTime = flowBody.ActionWaitTime

	if err := config.DB.Create(&action).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if flowBody.ActionType == "SMS" {
		var sms model.SmsAction
		sms.ActionID = action.ID
		sms.Body = flowBody.SmsBody

		if err := config.DB.Create(&sms).Error; err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
	} else if flowBody.ActionType == "EMAIL" {
		var email model.EmailAction
		email.ActionID = action.ID
		email.Subject = flowBody.EmailSubject
		email.Body = flowBody.EmailBody

		if err := config.DB.Create(&email).Error; err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
	} else {
		c.JSON(400, gin.H{"error": "Invalid action type"})
		return
	}

	c.JSON(200, gin.H{
		"flow":    flow,
		"action":  action,
		"trigger": trigger,
	})

}
