package segmentQL

import (
	"fmt"
	"github.com/nickhansel/nucleus/config"
	"strconv"
	"time"
)

//SELECT "customerId" from "Purchase" INNER JOIN "purchased_item" ON "Purchase".purchase_id = purchased_item."purchaseId"
//INNER JOIN "Items" ON purchased_item."itemId" = "Items".id WHERE "Items".id = 839081377060323329
//AND "Items"."organizationId" = 838194565431033857 AND cast("Purchase".created_at as DATE) BETWEEN '2023-01-15' AND '2023-02-15'
//AND "Purchase".amount_money BETWEEN 0 AND 1000;

func Parse(item int64, orgId int64, startDate string, endDate string, minPurchasePrice float64, maxPurchasePrice float64) []int64 {
	//	make a hashmap of string,string
	commands := make(map[string]string)
	commands["dates"] = "cast(\"Purchase\".created_at as DATE)"
	commands["price"] = "\"Purchase\".amount_money"
	commands["item"] = "\"Items\".id = " + fmt.Sprintf("%d", item)

	query := fmt.Sprintf("SELECT DISTINCT \"customerId\" from \"Purchase\" INNER JOIN \"purchased_item\" ON \"Purchase\".purchase_id = purchased_item.\"purchaseId\"\nINNER JOIN \"Items\" ON purchased_item.\"itemId\" = \"Items\".id WHERE \"Items\".\"organizationId\" = %s", strconv.FormatInt(orgId, 10))

	if startDate != "" && endDate == "" {
		endDate = time.Now().Format(time.RFC3339Nano)
	}

	//	build the query based on the parameters
	if item != 0 {
		query += " AND " + commands["item"]
	}
	if startDate != "" && endDate != "" {
		commandWithParams := commands["dates"]
		// add "BETWEEN startDate AND endDate"

		t, _ := time.Parse(time.RFC3339Nano, startDate)
		startParsed := t.Format("2006-01-02")

		t, _ = time.Parse(time.RFC3339Nano, endDate)
		endParsed := t.Format("2006-01-02")

		// add quotes to the dates
		startParsed = "'" + startParsed + "'"
		endParsed = "'" + endParsed + "'"

		commandWithParams += " BETWEEN " + startParsed + " AND " + endParsed
		query += " AND " + commandWithParams
	}
	if maxPurchasePrice != 0 {
		commandWithParams := commands["price"]
		minPurchasePrice = 0
		// add "BETWEEN minPurchasePrice AND maxPurchasePrice" to commandWithParams
		commandWithParams += " BETWEEN " + fmt.Sprintf("%f", minPurchasePrice) + " AND " + fmt.Sprintf("%f", maxPurchasePrice)
		query += " AND " + commandWithParams
	}

	query = query + " AND \"customerId\" != 0"

	var ids []int64
	// execute the raw query
	config.DB.Raw(query).Pluck("customerId", &ids)
	fmt.Println(query)
	return ids
}
