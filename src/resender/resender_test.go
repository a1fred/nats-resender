package resender

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func GetEnv(key, fallback string) string {
	val := os.Getenv(key)

	if val == "" {
		return fallback
	}
	return val
}

func GetNatsUrls(t *testing.T) (string, string) {
	fromUrl, toUrl := GetEnv("NATS_FROM", "nats://127.0.0.1:4222"), GetEnv("NATS_TO", "nats://127.0.0.1:4223")
	t.Logf("Nats connections: from: %s, to: %s\n", fromUrl, toUrl)
	return fromUrl, toUrl
}

func TestResender(t *testing.T) {
	fromUrl, toUrl := GetNatsUrls(t)

	assert.NotEqual(t, fromUrl, toUrl, "Nats from/to must be different")

	// Connection tests
	nn := NewResender(fromUrl, toUrl)
	assert.NotNil(t, nn.ncFrom)
	assert.NotNil(t, nn.ncTo)
	assert.Nil(t, nn.sub)
	assert.Equal(t, nn.processedCounter, uint(0))

	sub, err := nn.Subscribe("*", "", "to", 10, 10, true)
	assert.Nil(t, err)
	assert.NotNil(t, sub)
	assert.Equal(t, nn.processedCounter, uint(0))

	err = nn.ncFrom.Publish("test", []byte("test"))
	time.Sleep(1 * time.Second)
	assert.Nil(t, err)
	assert.Equal(t, nn.processedCounter, uint(1))

	// Close no errors
	assert.Nil(t, nn.Close())
}
