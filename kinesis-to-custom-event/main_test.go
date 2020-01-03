package main

import (
	"github.com/labstack/gommon/log"
	insights "github.com/newrelic/go-insights/client"
	"github.com/sirupsen/logrus"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestArnToEventType(t *testing.T) {
	assert.Equal(t, "test_stream", arnToEventType("arn:aws:kinesis:ap-northeast-1:398014708642:stream/test-stream"))
}

func TestToJson(t *testing.T) {
	a, b := toJson([]byte(`{"hello":"world", "hoge": 1, "fuga": "piyo", "array": [1,2,3,4,5]}`))
	assert.Equal(t, nil, a)
	assert.Equal(t, nil, b)
}

func TestHoge(t *testing.T) {
	insightInsertKey := os.Getenv("NEW_RELIC_INSIGHTS_INSERT_KEY")
	accountId := os.Getenv("NEW_RELIC_ACCOUNT_ID")
	insightClient = insights.NewInsertClient(insightInsertKey, accountId)

	// FIXME: set log level here
	insightClient.Logger.Level = logrus.DebugLevel
	log.SetLevel(log.DEBUG)

	err := insertEvent([]byte(`{"hello":"world", "hoge": 1, "fuga": [1,2,3,4,5], "hoge2": 2}`), "test")
	assert.NotEqual(t, nil, err)
}
