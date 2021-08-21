REGION=asia-northeast1
TRIGGER_TOPIC=cost-alert

deploy:
	gcloud functions deploy CostAlert --env-vars-file env.yaml --trigger-topic $(TRIGGER_TOPIC) --region=$(REGION) --runtime=go113