// middleware to get the user id from the previous middleware and check if the user is apart of the org that is in the url

package middleware

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nickhansel/nucleus/config"
	"github.com/nickhansel/nucleus/model"
)

func CheckOrgMiddleware() gin.HandlerFunc {
	// get the members of the org
	// check if the user is apart of the org
	return func(c *gin.Context) {
		// get the org id from the url
		orgId := c.Param("orgId")

		// get the id from the previous middleware
		id, _ := c.Get("id")

		orgIdInt, err := strconv.Atoi(orgId)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Invalid URL parameter",
			})
			return
		}

		// query all the data from the organizationtouser table
		// OrganizationsToUsers is a model where A is the organization id and B is the user id
		// query all organizations that the user is apart of
		orgs := model.Organization{}
		// find the members of the org and find the member where the id is the same as the id from the previous middleware
		res := config.DB.Preload("Members").Find(&orgs, orgIdInt).Error

		if res != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": res.Error(),
			})
			return
		}

		// find if the user is apart of the org
		for _, member := range orgs.Members {
			if member.ID == id {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		c.Abort()
	}
}
