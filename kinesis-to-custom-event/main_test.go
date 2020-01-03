package main

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestArnToEventType(t *testing.T) {
	assert.Equal(t, "test_stream", arnToEventType("arn:aws:kinesis:ap-northeast-1:398014708642:stream/test-stream"))
}

func TestToJson(t *testing.T) {
	a, _ := mapEvent([]byte(`{"hello":"world", "hoge": 1, "fuga": "piyo", "array": [1,2,3,4,5], "nested": {"hoge": "fuga"}}`), "test")
	j, _ := json.Marshal(a)
	assert.Equal(t,
		`{"array.0":1,"array.1":2,"array.2":3,"array.3":4,"array.4":5,"eventType":"UnknownKinesisEvent","fuga":"piyo","hello":"world","hoge":1,"nested.hoge":"fuga"}`,
		string(j))
}
