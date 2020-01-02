package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestArnToEventType(t *testing.T) {
	assert.Equal(t, "test_stream", arnToEventType("arn:aws:kinesis:ap-northeast-1:398014708642:stream/test-stream"))
}
