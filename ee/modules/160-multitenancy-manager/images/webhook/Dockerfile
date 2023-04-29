ARG BASE_ALPINE
ARG BASE_GOLANG_16_ALPINE

FROM $BASE_GOLANG_16_ALPINE as builder
WORKDIR /src/multitenancy-manager-webhook

COPY main.go go.mod go.sum /src/multitenancy-manager-webhook/
COPY web /src/multitenancy-manager-webhook/web/

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go test ./...
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -o multitenancy-manager-webhook main.go

FROM $BASE_ALPINE
COPY --from=builder /src/multitenancy-manager-webhook/multitenancy-manager-webhook /multitenancy-manager-webhook
ENTRYPOINT [ "/multitenancy-manager-webhook" ]
