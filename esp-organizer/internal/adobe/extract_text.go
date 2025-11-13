package adobe

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

type AdobePDFExtractor struct {
	ClientID     string
	ClientSecret string
	AccessToken  string
}

// ExtractTextFromDocumentCloudPDF extracts text from a PDF in Adobe Document Cloud
func (a *AdobePDFExtractor) ExtractTextFromDocumentCloudPDF(assetID string) (string, error) {
	// Step 1: Get access token
	if err := a.GetAccessToken(); err != nil {
		return "", err
	}

	// Step 2: Call Adobe PDF Extract API
	extractURL := "https://pdf-services.adobe.io/operation/extractpdf"

	payload := map[string]interface{}{
		"assetID": assetID,
		"renditions": []map[string]interface{}{
			{
				"type": "text",
			},
		},
	}

	jsonPayload, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", extractURL, bytes.NewBuffer(jsonPayload))
	req.Header.Set("Authorization", "Bearer "+a.AccessToken)
	req.Header.Set("x-api-key", a.ClientID)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	// Parse response to get extracted text
	var result struct {
		Text string `json:"text"`
	}
	json.Unmarshal(body, &result)

	return result.Text, nil
}

func (a *AdobePDFExtractor) GetAccessToken() error {
	// OAuth token retrieval (see earlier example)
	return nil
}
