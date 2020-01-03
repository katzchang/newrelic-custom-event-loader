package main

import (
	"encoding/json"
	"github.com/labstack/gommon/log"
	insights "github.com/newrelic/go-insights/client"
	"github.com/sirupsen/logrus"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestArnToEventType(t *testing.T) {
	assert.Equal(t, "test_stream", arnToEventType("arn:aws:kinesis:ap-northeast-1:398014708642:stream/test-stream"))
}

func TestToJson(t *testing.T) {
	a, _ := mapEvent([]byte(`{"hello":"world", "hoge": 1, "fuga": "piyo", "array": [1,2,3,4,5], "nested": {"hoge": "fuga"}}`),
		"arn:aws:kinesis:ap-northeast-1:398014708642:stream/test-stream")
	j, _ := json.Marshal(a)
	assert.Equal(t,
		`{"array.0":1,"array.1":2,"array.2":3,"array.3":4,"array.4":5,"eventType":"test_stream","fuga":"piyo","hello":"world","hoge":1,"nested.hoge":"fuga"}`,
		string(j))
}

func TestX(t *testing.T) {
	initialize()
	insightClient.Logger.Level = logrus.DebugLevel
	log.SetLevel(log.DEBUG)
	// TODO: Deflateつかえなくない？？
	insightClient.Compression = insights.Gzip
	arn := "arn:aws:kinesis:ap-northeast-1:398014708642:stream/test-stream"
	d, _ := mapEvent([]byte(`{"hello":"world"}`), arn)
	ddd := [...] map[string]interface{}{d, d, d}
	err := insightClient.PostEvent(ddd)
	assert.NotNil(t, err)
}
