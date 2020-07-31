FROM golang:1.14-alpine as build
WORKDIR /src/wooc
COPY . .
RUN go mod vendor
RUN go build -mod=vendor

FROM alpine:latest
COPY --from=build /src/wooc/wooc .
COPY --from=build /src/wooc/config.toml ./config.toml
COPY --from=build /src/wooc/swaggerui ./swaggerui
CMD ["./wooc"]
EXPOSE 8080
