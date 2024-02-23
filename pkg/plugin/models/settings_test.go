package models

import (
	"context"
	"testing"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadSettings(t *testing.T) {
	t.Run("when settings are valid", func(t *testing.T) {
		t.Parallel()

		config := backend.DataSourceInstanceSettings{
			JSONData: []byte(`{"baseUrl":"http://localhost:3000", "siteId":"my-site-id"}`),
			DecryptedSecureJSONData: map[string]string{
				"accessToken": "my-access-token",
			},
		}

		settings, err := LoadSettings(context.Background(), config)
		require.NoError(t, err)
		assert.Equal(t, "my-access-token", settings.AccessToken)
		assert.Equal(t, "http://localhost:3000", settings.BaseUrl)
		assert.Equal(t, "my-site-id", settings.SiteId)
	})

	t.Run("returns error when missing accessToken", func(t *testing.T) {
		t.Parallel()

		config := backend.DataSourceInstanceSettings{
			JSONData: []byte(`{"baseUrl":"http://localhost:3000", "siteId":"my-site-id"}`),
		}

		_, err := LoadSettings(context.Background(), config)
		assert.Error(t, err, "accessToken is missing")
	})
}
