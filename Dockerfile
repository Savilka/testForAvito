FROM golang:latest

WORKDIR /app

COPY ./ /app
ENV CONFIG_PATH=config/cfg.yaml
RUN go mod download
RUN go build cmd/app/main.go
ENTRYPOINT ./main
EXPOSE 8000