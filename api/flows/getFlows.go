package flows

import (
	"github.com/gin-gonic/gin"
	"github.com/nickhansel/nucleus/config"
	"github.com/nickhansel/nucleus/model"
)

type FlowData struct {
	Flow    model.Flow      `json:"flow"`
	Action  model.Action    `json:"action"`
	Trigger model.Trigger   `json:"trigger"`
	FlowRan []model.FlowRan `json:"flow_ran"`
}

func GetFlows(c *gin.Context) {
	org := c.MustGet("orgs").(model.Organization)

	var flows []model.Flow
	config.DB.Where(&model.Flow{OrganizationID: org.ID}).Find(&flows)

	if len(flows) == 0 {
		c.JSON(200, gin.H{"data": []model.Flow{}})
		return
	}

	flowIds := make([]int64, len(flows))
	for i, flow := range flows {
		flowIds[i] = flow.ID
	}

	var triggers []model.Trigger
	config.DB.Where("\"flowId\" IN ?", flowIds).Find(&triggers)

	triggerIds := make([]int64, len(triggers))
	for i, trigger := range triggers {
		triggerIds[i] = trigger.ID
	}

	var actions []model.Action
	config.DB.Where("\"triggerId\" IN ?", triggerIds).Find(&actions)

	var ran []model.FlowRan
	config.DB.Where("\"flowId\" IN ?", flowIds).Find(&ran)

	var flowData []FlowData
	for _, flow := range flows {
		var trigger model.Trigger
		var action model.Action
		var flowRan []model.FlowRan

		for _, t := range triggers {
			if t.FlowID == flow.ID {
				trigger = t
			}
		}

		for _, a := range actions {
			if a.TriggerID == trigger.ID {
				action = a
			}
		}

		for _, r := range ran {
			if r.FlowID == flow.ID {
				flowRan = append(flowRan, r)
			}
		}

		if len(flowRan) == 0 {
			flowRan = append(flowRan, model.FlowRan{FlowID: flow.ID})
		}

		flowData = append(flowData, FlowData{Flow: flow, Trigger: trigger, Action: action, FlowRan: flowRan})
	}

	c.JSON(200, gin.H{"data": flowData})

}
