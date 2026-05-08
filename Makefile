swag:
	@swag init --parseDependency --parseInternal

dev: swag
	@air serve

go: swag
	@go run main.go serve

migrate:
	@go run main.go migrate

migrate-force:
	@go run main.go migrate --drop

# bunx --bun shadcn@latest init --preset b2p8PagwS --template start --pointer