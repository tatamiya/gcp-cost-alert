REGION=asia-northeast1
TRIGGER_TOPIC=cost-alert

local-run:
	godotenv -f ./.env go test -run TestCostAlert,TestNotificationIsNotSentWhenPayloadIsEmpty

test:
	godotenv -f ./.env go test ./... -cover -short

test-all:
	godotenv -f ./.env go test ./... -cover

deploy:
	gcloud functions deploy CostAlert --env-vars-file env.yaml --trigger-topic $(TRIGGER_TOPIC) --region=$(REGION) --runtime=go113