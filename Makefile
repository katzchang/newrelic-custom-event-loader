.PHONY: deps clean build

main=./kinesis-to-custom-event/kinesis-to-custom-event

clean:
	rm -rf $(main)

test:
	go test ./kinesis-to-custom-event

build:
	GOOS=linux GOARCH=amd64 go build -o $(main) ./kinesis-to-custom-event \

deploy: clean build
	sam deploy\
	 --parameter-overrides "NewRelicInsightsInsertKey=$(NEW_RELIC_INSIGHTS_INSERT_KEY) NewRelicAccountId=$(NEW_RELIC_ACCOUNT_ID)"

put-record:
	aws kinesis put-record --stream-name test-stream --data '{"hello":"world", "hoge": 1, "fuga": "piyo", "array": [1,2,3,4,5], "nested": {"hoge": "fuga"}}' --partition-key 0

put-record2:
	aws kinesis put-record --stream-name test-stream --data '{"hello":"world"}' --partition-key 0
