gen:
	go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen -o internal/gen/types.go -generate chi-server,types -package gen api/api.yaml