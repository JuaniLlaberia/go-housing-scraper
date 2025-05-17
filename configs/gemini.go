package configs

import (
	"context"
	"zonaprop-scraper/utils"

	"google.golang.org/genai"
)

type GeminiConfig struct {
	APIKey          string
	ModelName       string
	Temperature     *float32
	TopP            *float32
	TopK            *int32
	MaxOutputTokens *int32
	StopSequences   []string
	SafeSettings    []*genai.SafetySetting
}

func DefaultSafetySettings() []*genai.SafetySetting {
	return []*genai.SafetySetting{
		{
			Category:  genai.HarmCategoryHateSpeech,
			Threshold: genai.HarmBlockThresholdBlockMediumAndAbove,
		},
		{
			Category:  genai.HarmCategoryDangerousContent,
			Threshold: genai.HarmBlockThresholdBlockMediumAndAbove,
		},
		{
			Category:  genai.HarmCategoryHarassment,
			Threshold: genai.HarmBlockThresholdBlockMediumAndAbove,
		},
		{
			Category:  genai.HarmCategorySexuallyExplicit,
			Threshold: genai.HarmBlockThresholdBlockMediumAndAbove,
		},
	}
}

func Gemini(contents []*genai.Content) (string, error) {
	apiKey := utils.ProcessEnv("GEMINI_API_KEY")
	ctx := context.Background()

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return "", err
	}

	temp := float32(0.1)
	topP := float32(0.3)
	topK := float32(10.0)
	maxOutputTokens := int32(8192)

	const systemInstructions = `
		You are a real state agent analyzing properties. Your task is to analyze certain properties and create analysis on them.

		General Behavior Guidelines:
		- Always return output in valid JSON.
		- Do not include any explanatory text or commentary outside the JSON.
		- When uncertain about missing context use default values based on the topics and prompt.
		- If you encounter a repeated property you can remove it.
		- If a property does not have a price, do not consider it.
		- Focus on accuracy, precision when analysing.

		Formatting Rules:
		- Always wrap your entire output inside a valid JSON object or array depending on the endpoint requirements.
		- Use snake_case for all keys.
		`

	config := &genai.GenerateContentConfig{
		Temperature:       &temp,
		TopP:              &topP,
		TopK:              &topK,
		MaxOutputTokens:   maxOutputTokens,
		ResponseMIMEType:  "application/json",
		SystemInstruction: genai.NewContentFromText(systemInstructions, genai.RoleUser),
		SafetySettings:    DefaultSafetySettings(),
	}

	result, err := client.Models.GenerateContent(ctx, "gemini-2.0-flash", contents, config)

	if err != nil {
		return "", err
	}

	return result.Text(), nil
}
