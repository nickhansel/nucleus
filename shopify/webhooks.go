package shopify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/nickhansel/nucleus/model"
	"github.com/pkg/errors"
	"io"
	"log"
	"net/http"
	"os"
)

func RegisterHooks(shopifyPos model.Pos, shopifyDomain string, org model.Organization) error {
	accessToken := shopifyPos.AccessToken

	topics := []string{
		"orders/create",
		"products/create",
		"products/update",
		"products/delete",
		"customers/create",
		"customers/update",
		"customers/delete",
		"locations/create",
		"locations/update",
		"locations/delete",
		"domains/create",
		"domains/update",
		"domains/destroy",
		"app/uninstalled",
	}

	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file for sendgrid")
	}

	arn := os.Getenv("SHOPIFY_ARN")
	if arn == "" {
		return errors.New("invalid arn")
	}

	if accessToken == "" || shopifyDomain == "" || org.ID == 0 {
		return errors.New("invalid parameters")
	}

	reqUrl := "https://" + shopifyDomain + "/admin/api/2023-01/webhooks.json"

	for _, topic := range topics {
		reqBody := fmt.Sprintf(`{
			"webhook": {
				"topic": "%s",
				"address": "%s",
				"format": "json"
			}
		}`, topic, arn)

		fmt.Println(reqBody)

		req, err := http.NewRequest("POST", reqUrl, bytes.NewBuffer([]byte(reqBody)))
		if err != nil {
			return errors.Wrap(err, "error creating request")
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Shopify-Access-Token", accessToken)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			continue
		}

		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				return
			}
		}(resp.Body)

		var result map[string]interface{}

		err = json.NewDecoder(resp.Body).Decode(&result)
		if err != nil {
			fmt.Println(err)
			continue
		}

		fmt.Println(result)

	}

	return nil

}
