FROM golang:1.22.4-bookworm

USER root

ARG SERVICE_NAME

COPY ./configs/${SERVICE_NAME}.json /configs/config.json
COPY ./${SERVICE_NAME} /app

ENTRYPOINT ["/app"]
