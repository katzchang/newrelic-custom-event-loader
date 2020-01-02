.PHONY: deps clean build

main=./kinesis-to-custom-event/kinesis-to-custom-event

clean:
	rm -rf $(main)
	
build:
	GOOS=linux GOARCH=amd64 go build -o $(main) ./kinesis-to-custom-event \

deploy: clean build
	sam deploy\
	 --parameter-overrides "NewRelicInsightsInsertKey=$(NEW_RELIC_INSIGHTS_INSERT_KEY) NewRelicAccountId=$(NEW_RELIC_ACCOUNT_ID)"

put-record:
	aws kinesis put-record --stream-name test-stream --data '{"hello":"world"}' --partition-key 0
