go test ./app/adapter/out/tidbrepository -coverprofile=coverage.out && go tool cover -html=coverage.out -o coverage.html

go test ./app/domain/... -coverprofile=coverage.out && go tool cover -html=coverage.out -o coverage.html

gqlgen generate
