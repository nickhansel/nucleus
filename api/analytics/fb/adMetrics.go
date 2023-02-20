package fb

import (
	"encoding/json"
	"fmt"
	"github.com/nickhansel/nucleus/model"
	"io"
	"net/http"
)

func GetMetricsFromFB(org model.Organization, level string, adId string) (error, map[string]interface{}) {

	accessToken := org.FbAccessToken

	reqUrl := fmt.Sprintf("https://graph.facebook.com/v15.0/%s/insights?fields=clicks,reach,spend,account_id,adset_id,campaign_id,ad_id,impressions,unique_clicks,cpc,conversions,cpm,ctr,frequency&level=%s&access_token=%s", adId, level, accessToken)

	req, err := http.NewRequest("GET", reqUrl, nil)
	if err != nil {
		return err, nil
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err, nil
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	// convert the response to json
	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return err, nil
	}

	return nil, result
}
