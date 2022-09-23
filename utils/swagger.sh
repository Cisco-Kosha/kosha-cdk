go install github.com/swaggo/swag/cmd/swag@latest
swag init -g main.go --parseDependency --parseInternal
npm i api-spec-converter
npx api-spec-converter --from=swagger_2 --to=openapi_3 --syntax=json ./docs/swagger.json > openapi.json