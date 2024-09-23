package banner

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/lcrownover/banner-tools/internal/keys"
)

type DuckIDResponseData struct {
	BannerID string `json:"bannerID,omitempty"`
	DuckID   string `json:"duckID,omitempty"`
}

type DuckIDResponse struct {
	Message    string             `json:"message,omitempty"`
	Data       DuckIDResponseData `json:"data,omitempty"`
	StatusCode int                `json:"statusCode,omitempty"`
}

func BannerIDToDuckID(ctx context.Context, bannerID string) (string, error) {
	apiKey := ctx.Value(keys.ApiKeyKey).(string)
	apiURL := ctx.Value(keys.ApiURL).(string)

	url := fmt.Sprintf("https://%s/person/uo/duckid/%s", apiURL, bannerID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Ocp-Apim-Subscription-Key", apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response: %v", err)
	}

	var duckIdResp DuckIDResponse
	err = json.Unmarshal(body, &duckIdResp)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal json: %v", err)
	}

	status := duckIdResp.StatusCode
	if status != 200 {
		return "", fmt.Errorf("banner response code %d: %s", status, duckIdResp.Message)
	}

	duckid := duckIdResp.Data.DuckID
	return duckid, nil
}
