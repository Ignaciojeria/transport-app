go test ./app/adapter/out/tidbrepository/... -coverprofile=coverage.out && go tool cover -html=coverage.out -o coverage.html

