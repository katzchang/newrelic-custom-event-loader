package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/labstack/gommon/log"
	"github.com/sirupsen/logrus"

	"github.com/jeremywohl/flatten"

	newrelic "github.com/newrelic/go-agent"
	"github.com/newrelic/go-agent/_integrations/nrlambda"
	insights "github.com/newrelic/go-insights/client"
)

var (
	insightClient *insights.InsertClient
)

func main() {
	initialize()

	cfg := nrlambda.NewConfig()
	app, err := newrelic.NewApplication(cfg)
	// TODO: エラーにはならないっぽいので別途チェックしないと
	if nil != err {
		lambda.Start(handler)
	} else {
		nrlambda.Start(handler, app)
	}
}

func initialize() {
	insightInsertKey := os.Getenv("NEW_RELIC_INSIGHTS_INSERT_KEY")
	accountId := os.Getenv("NEW_RELIC_ACCOUNT_ID")
	insightClient = insights.NewInsertClient(insightInsertKey, accountId)
	insightClient.Compression = insights.Deflate

	// FIXME: set log level
	insightClient.Logger.Level = logrus.InfoLevel
	log.SetLevel(log.INFO)

	for _, pair := range os.Environ() {
		log.Debug(pair)
	}
	log.Debug("insightInsertKey:", insightInsertKey)
	log.Debug("accountId:", accountId)
}

func handler(event events.KinesisEvent) (string, error) {
	var count, success, failure = 0, 0, 0
	var data []map[string]interface{}

	for _, record := range event.Records {
		if x, e := mapEvent(record.Kinesis.Data, record.EventSourceArn); e != nil {
			data = append(data, x)
			success += 1
		} else {
			// TODO: warn扱いにしたい
			failure += 1
			log.Errorf("Error: %v\n", e)
		}
		count += 1
	}

	if e := insightClient.PostEvent(data); e != nil {
		return "", e
	}

	msg := fmt.Sprintf(`{"count":%d,"success"%d,"failure":%d}`, count, success, failure)
	return msg, nil
}

// TODO: flatten https://docs.newrelic.co.jp/docs/insights/insights-data-sources/custom-data/introduction-event-api
// TODO: EventTypeが仕様と違っても、200 OKになるので注意。イベントは作成されない。 https://docs.newrelic.com/docs/insights/insights-data-sources/custom-data/insights-custom-data-requirements-limits
// TODO: AttributeNameも仕様を合わせないとだ https://docs.newrelic.com/docs/insights/insights-data-sources/custom-data/insights-custom-data-requirements-limits
func mapEvent(kinesisData []byte, eventSourceARN string) (map[string]interface{}, error) {
	var data map[string]interface{}


	if err := json.Unmarshal(kinesisData, &data); err != nil {
		return nil, err
	}
	data["eventType"] = arnToEventType(eventSourceARN)

	return flatten.Flatten(data, "", flatten.DotStyle)
}

/**
https://docs.newrelic.com/docs/insights/insights-data-sources/custom-data/insights-custom-data-requirements-limits
- Maximum name size: 255 bytes.
- Maximum total attributes per event: 254.
- Exception: If you use the APM agent API, the max is 64.
- Maximum total attributes per event type: 48,000.
- Event types (using the eventType attribute) can be a combination of alphanumeric characters, colons (:), and underscores (_).

https://docs.aws.amazon.com/kinesis/latest/APIReference/API_CreateStream.html#API_CreateStream_RequestParameters
- Length Constraints: Minimum length of 1. Maximum length of 128.
- Pattern: [a-zA-Z0-9_.-]+
*/
func arnToEventType(arn string) string {
	xx := strings.Split(arn, "/")
	if len(xx) != 2 {
		return "UnknownKinesisEvent"
	}
	return strings.ReplaceAll(xx[1], "-", "_")
}
