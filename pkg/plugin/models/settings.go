package models

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

type Settings struct {
	AccessToken string `json:"accessToken"`
	SiteId      string `json:"siteId"`
	AccountId   string `json:"accountId"`
	BaseUrl     string `json:"baseUrl"`
}

func LoadSettings(ctx context.Context, config backend.DataSourceInstanceSettings) (Settings, error) {
	s := Settings{}
	if err := json.Unmarshal(config.JSONData, &s); err != nil {
		return Settings{}, fmt.Errorf("failed to unmarshal settings JSONData: %w", err)
	}

	accessToken, ok := config.DecryptedSecureJSONData["accessToken"]
	if !ok {
		return Settings{}, fmt.Errorf("accessToken is missing")
	}

	s.AccessToken = accessToken

	return s, nil
}
