package main

import (
	"encoding/json"
	"os"
	//"regexp"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/labstack/gommon/log"
	"github.com/sirupsen/logrus"

	newrelic "github.com/newrelic/go-agent"
	"github.com/newrelic/go-agent/_integrations/nrlambda"
	insights "github.com/newrelic/go-insights/client"
)

var (
	insightClient *insights.InsertClient
)

func main() {
	insightInsertKey := os.Getenv("NEW_RELIC_INSIGHTS_INSERT_KEY")
	accountId := os.Getenv("NEW_RELIC_ACCOUNT_ID")
	insightClient = insights.NewInsertClient(insightInsertKey, accountId)

	// FIXME: set log level here
	insightClient.Logger.Level = logrus.InfoLevel
	log.SetLevel(log.INFO)

	for _, pair := range os.Environ() {
		log.Debug(pair)
	}
	log.Debug("insightInsertKey:", insightInsertKey)
	log.Debug("accountId:", accountId)

	cfg := nrlambda.NewConfig()
	app, err := newrelic.NewApplication(cfg)
	// TODO: エラーにはならないっぽいので別途チェックしないと
	if nil != err {
		lambda.Start(handler)
	} else {
		nrlambda.Start(handler, app)
	}
}

// TODO: まとめて送る
// TODO: flattenが必要…
func handler(event events.KinesisEvent) (string, error) {
	var count = 0
	for _, record := range event.Records {
		kinesisData := record.Kinesis.Data
		eventType := record.EventSourceArn

		if err := insertEvent(kinesisData, eventType); err != nil {
			// TODO
			log.Errorf("Error: %v\n", err)
		}
		count += 1
	}
	// TODO: 処理した件数と最初のエラー
	return "", nil
}

func insertEvent(kinesisData []byte, arn string) error {
	var data map[string]interface{}
	if err := json.Unmarshal(kinesisData, &data); err != nil {
		return err
	}
	data["eventType"] = arnToEventType(arn)
	log.Debug(data)
	if err := insightClient.PostEvent(data); err != nil {
		return err
	}
	return nil
}

func toJson(kinesisData []byte) (map[string]interface{}, error) {
	var data map[string]interface{}
	if err := json.Unmarshal(kinesisData, &data); err != nil {
		return nil, err
	}
	return data, nil
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

TODO: EventTypeがしようと違っても、200 OKになるので注意。イベントは作成されない。
TODO: AttributeNameも仕様を合わせないとだ https://docs.newrelic.com/docs/insights/insights-data-sources/custom-data/insights-custom-data-requirements-limits
 */
func arnToEventType(arn string) string {
	defaultName := "UnknownKinesisEvent"

	xx := strings.Split(arn, "/")
	if len(xx) != 2 {
		return defaultName
	}

	return strings.ReplaceAll(xx[1], "-", "_")
}
