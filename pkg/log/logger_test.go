package log

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLogger_NewLogger(t *testing.T) {
	loggerConfig := Config{
		AppName: "network-stats",
		Mode:    "dev",
		LokiURL: "https://logs-prod-017.grafana.net",
	}

	logger, err := NewLogger(loggerConfig)
	assert.Nil(t, err)
	logger.Info("dm thang phat")
	defer logger.Sync()

	//time.Sleep(30 * time.Minute)
}
