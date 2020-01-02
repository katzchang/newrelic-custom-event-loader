package main

import (
	"encoding/json"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/labstack/gommon/log"

	newrelic "github.com/newrelic/go-agent"
	"github.com/newrelic/go-agent/_integrations/nrlambda"
	insights "github.com/newrelic/go-insights/client"
)

var (
	insightClient    *insights.InsertClient
)

func main() {
	insightInsertKey := os.Getenv("NEW_RELIC_INSIGHTS_INSERT_KEY")
	accountId        := os.Getenv("NEW_RELIC_ACCOUNT_ID")
	insightClient = insights.NewInsertClient(insightInsertKey, accountId)

	cfg := nrlambda.NewConfig()
	app, err := newrelic.NewApplication(cfg)
	if nil != err {
		lambda.Start(handler)
	} else {
		nrlambda.Start(handler, app)
	}
}

func handler(event events.KinesisEvent) (string, error) {
	for _, record := range event.Records {
		kinesisData := record.Kinesis.Data
		eventType := record.EventSourceArn

		if err := insertEvent(kinesisData, eventType); err != nil {
			// TODO
			log.Errorf("Error: %v\n", err)
		}
	}
	// TODO: 処理した件数と最初のエラー
	return "", nil
}

func insertEvent(kinesisData []byte, eventType string) error {
	var data map[string]interface {}
	if err := json.Unmarshal(kinesisData, &data); err != nil {
		return err
	}
	data["eventType"] = eventType
	if err := insightClient.PostEvent(data); err != nil {
		return err
	}
	return nil
}
