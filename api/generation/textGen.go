package generation

import (
	"github.com/gin-gonic/gin"
	"github.com/nickhansel/nucleus/openai"
)

type TextGenBody struct {
	Length int    `json:"length"`
	Prompt string `json:"prompt"`
}

func GenerateText(c *gin.Context) {
	var body TextGenBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	length := body.Length
	prompt := body.Prompt

	prompt += " make the response" + string(rune(length)) + "characters long"

	text, err := openai.GenerateText(prompt)
	if err != nil {
		return
	}

	c.JSON(200, gin.H{"text": text})
}
