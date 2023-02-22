package openai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"io"
	"net/http"
	"os"
)

type GenerateTextBody struct {
	Model       string  `json:"model"`
	Prompt      string  `json:"prompt"`
	Temperature float64 `json:"temperature"`
	MaxTokens   int     `json:"max_tokens"`
}

//curl https://api.openai.com/v1/completions \
//-H "Content-Type: application/json" \
//-H "Authorization: Bearer YOUR_API_KEY" \
//-d '{"model": "text-davinci-003", "prompt": "Say this is a test", "temperature": 0, "max_tokens": 7}'

func GenerateText(prompt string) ([]string, error) {
	err := godotenv.Load("../.env")
	if err != nil {
		fmt.Println("Error loading .env file")
		return nil, err
	}

	reqUrl := "https://api.openai.com/v1/completions"

	client := &http.Client{}

	req, err := http.NewRequest("POST", reqUrl, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+os.Getenv("OPENAI_API_KEY"))
	var body GenerateTextBody
	body.Model = "text-davinci-003"
	body.Prompt = prompt
	body.Temperature = 0
	body.MaxTokens = 100

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req.Body = io.NopCloser(bytes.NewBuffer(jsonBody))

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
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
		return nil, err
	}

	//loop through choices and return them
	choices := result["choices"].([]interface{})
	var responses []string

	for _, choice := range choices {
		response := choice.(map[string]interface{})["text"].(string)
		responses = append(responses, response)
	}
	return responses, nil

}
