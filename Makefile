FUNC_NAME ?= spotiot

.PHONY: build
build: # Build
	GOOS=linux go build -o dist/lambda ./cmd/lambda
	cd dist; zip lambda.zip lambda

.PHONY: deploy
deploy: # Deploy
	LAMBDA_VARS=$$(bash -c "cat <(grep '^SPOTIFY_.*' .env) <(go run ./cmd/auth) | tr '\n' ','"); \
	aws lambda update-function-code --function-name $(FUNC_NAME) --zip-file fileb://dist/lambda.zip; \
	aws lambda update-function-configuration --function-name $(FUNC_NAME) --environment Variables={$$LAMBDA_VARS}

.PHONY: clean
clean: # Clean
	rm -rf ./dist