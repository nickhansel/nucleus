package groups

import (
	"github.com/gin-gonic/gin"
	"github.com/nickhansel/nucleus/model"
	"github.com/nickhansel/nucleus/segmentQL"
	"net/http"
)

func SegmentCustomers(c *gin.Context) {
	org := c.MustGet("orgs").(model.Organization)
	if org.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	query := SegmentQLRequestBody{}
	err := c.BindJSON(&query)
	if err != nil {
		return
	}

	ql := segmentQL.Parse(int64(query.ItemID), org.ID, query.StartDate, query.EndDate, query.MinPurchasePrice, query.MaxPurchasePrice)

	c.JSON(http.StatusOK, gin.H{"customers": ql})
}

