package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestArnToEventType(t *testing.T) {
	assert.Equal(t, "test_stream", arnToEventType("arn:aws:kinesis:ap-northeast-1:398014708642:stream/test-stream"))
}

func TestToJson(t *testing.T) {
	a, b := mapEvent([]byte(`{"hello":"world", "hoge": 1, "fuga": "piyo", "array": [1,2,3,4,5]}`), "test")
	assert.Equal(t, nil, a)
	assert.Equal(t, nil, b)
}
