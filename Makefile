.PHONY: deps clean build

deps:
	go get -u ./...

clean: 
	rm -rf ./kinesis-to-custom-event/kinesis-to-custom-event
	
build:
	GOOS=linux GOARCH=amd64 go build -o kinesis-to-custom-event/kinesis-to-custom-event ./kinesis-to-custom-event \

deploy:
	sam deploy \
	 --parameter-overrides "NewRelicInsightsInsertKey=$(NEW_RELIC_INSIGHTS_INSERT_KEY) NewRelicAccountId=$(NEW_RELIC_ACCOUNT_ID)"
