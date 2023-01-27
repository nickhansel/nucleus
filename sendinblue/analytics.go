package sendinblue

import (
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Analytics struct {
	Reports []struct {
		Date         string `json:"date"`
		Requests     int    `json:"requests"`
		Delivered    int    `json:"delivered"`
		HardBounces  int    `json:"hardBounces"`
		SoftBounces  int    `json:"softBounces"`
		Clicks       int    `json:"clicks"`
		UniqueClicks int    `json:"uniqueClicks"`
		Opens        int    `json:"opens"`
		UniqueOpens  int    `json:"uniqueOpens"`
		SpamReports  int    `json:"spamReports"`
		Blocked      int    `json:"blocked"`
		Invalid      int    `json:"invalid"`
		Unsubscribed int    `json:"unsubscribed"`
	} `json:"reports"`
}

func GetEmailAnalytics(campaignId int32) (Analytics, error) {
	idToString := strconv.Itoa(int(campaignId))

	reqUrl := "https://api.sendinblue.com/v3/smtp/statistics/reports?limit=10&offset=0&days=1&tag=" + idToString + "&sort=desc"
	client := &http.Client{Timeout: 10 * time.Second}
	err := godotenv.Load("../.env")

	apiKey := os.Getenv("SENDINBLUE_API_KEY")

	req, err := http.NewRequest("GET", reqUrl, nil)
	if err != nil {
		fmt.Println(err)
		return Analytics{}, nil
	}

	req.Header.Add("api-key", apiKey)
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return Analytics{}, nil
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
			return
		}
	}(resp.Body)

	var result map[string]interface{}
	result = make(map[string]interface{})

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		fmt.Println(err)
		return Analytics{}, nil
	}

	//	convert to Analytics struct
	var analytics Analytics
	analytics.Reports = make([]struct {
		Date         string `json:"date"`
		Requests     int    `json:"requests"`
		Delivered    int    `json:"delivered"`
		HardBounces  int    `json:"hardBounces"`
		SoftBounces  int    `json:"softBounces"`
		Clicks       int    `json:"clicks"`
		UniqueClicks int    `json:"uniqueClicks"`
		Opens        int    `json:"opens"`
		UniqueOpens  int    `json:"uniqueOpens"`
		SpamReports  int    `json:"spamReports"`
		Blocked      int    `json:"blocked"`
		Invalid      int    `json:"invalid"`
		Unsubscribed int    `json:"unsubscribed"`
	}, len(result["reports"].([]interface{})))

	for i, v := range result["reports"].([]interface{}) {
		analytics.Reports[i].Date = v.(map[string]interface{})["date"].(string)
		analytics.Reports[i].Requests = int(v.(map[string]interface{})["requests"].(float64))
		analytics.Reports[i].Delivered = int(v.(map[string]interface{})["delivered"].(float64))
		analytics.Reports[i].HardBounces = int(v.(map[string]interface{})["hardBounces"].(float64))
		analytics.Reports[i].SoftBounces = int(v.(map[string]interface{})["softBounces"].(float64))
		analytics.Reports[i].Clicks = int(v.(map[string]interface{})["clicks"].(float64))
		analytics.Reports[i].UniqueClicks = int(v.(map[string]interface{})["uniqueClicks"].(float64))
		analytics.Reports[i].Opens = int(v.(map[string]interface{})["opens"].(float64))
		analytics.Reports[i].UniqueOpens = int(v.(map[string]interface{})["uniqueOpens"].(float64))
		analytics.Reports[i].SpamReports = int(v.(map[string]interface{})["spamReports"].(float64))
		analytics.Reports[i].Blocked = int(v.(map[string]interface{})["blocked"].(float64))
		analytics.Reports[i].Invalid = int(v.(map[string]interface{})["invalid"].(float64))
		analytics.Reports[i].Unsubscribed = int(v.(map[string]interface{})["unsubscribed"].(float64))
	}

	return analytics, nil
}
