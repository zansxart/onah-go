package plugins

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"onah-go/config"
)

func init() {
	Register(Command{
		Name:  "ai",
		Tags:  []string{"ai"},
		Help:  "Bertanya kepada kecerdasan buatan (Gemini)",
		Limit: true,
		Execute: runGeminiAI,
	})

	Register(Command{
		Name:  "gemini",
		Tags:  []string{"ai"},
		Help:  "Bertanya kepada kecerdasan buatan (Gemini)",
		Limit: true,
		Execute: runGeminiAI,
	})
}

type GeminiRequest struct {
	Contents []Content `json:"contents"`
}

type Content struct {
	Parts []Part `json:"parts"`
}

type Part struct {
	Text string `json:"text"`
}

type GeminiResponse struct {
	Candidates []Candidate `json:"candidates"`
}

type Candidate struct {
	Content ResponseContent `json:"content"`
}

type ResponseContent struct {
	Parts []Part `json:"parts"`
}

func runGeminiAI(ctx *Context) error {
	if ctx.Query == "" {
		return ctx.Reply("⚠️ Harap berikan pertanyaan Anda setelah perintah! Contoh: `.ai siapa penemu listrik?`")
	}

	ctx.React("🕛")

	apiKey := config.ActiveConfig.ApiKeys["gemini"]
	if apiKey == "" || apiKey == "-" {
		return ctx.Reply("⚠️ API Key Gemini belum dikonfigurasi di config.json.")
	}

	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1/models/gemini-1.5-flash:generateContent?key=%s", apiKey)

	reqBody := GeminiRequest{
		Contents: []Content{
			{
				Parts: []Part{
					{Text: ctx.Query},
				},
			},
		},
	}

	jsonBytes, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var geminiResp GeminiResponse
	err = json.Unmarshal(bodyBytes, &geminiResp)
	if err != nil {
		return err
	}

	if len(geminiResp.Candidates) == 0 || len(geminiResp.Candidates[0].Content.Parts) == 0 {
		return ctx.Reply("❌ Maaf, Gemini tidak memberikan jawaban untuk pertanyaan tersebut.")
	}

	answer := geminiResp.Candidates[0].Content.Parts[0].Text
	ctx.React("✅")
	return ctx.Reply(answer)
}
