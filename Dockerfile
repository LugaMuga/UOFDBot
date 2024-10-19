# -----------------------------------------------------------------------------
#  Build Stage
# -----------------------------------------------------------------------------
FROM golang:alpine3.20 AS build

ENV CGO_ENABLED=1

RUN apk add --no-cache \
    gcc \
    musl-dev

WORKDIR /uofd

COPY . /uofd/

RUN \
    cd /uofd && \
    go mod tidy && \
    go build

# -----------------------------------------------------------------------------
#  Run Stage
# -----------------------------------------------------------------------------
FROM alpine:3.20.3

COPY --from=build /uofd/UOFDBot /opt/UOFDBot/UOFDBot

WORKDIR /opt/UOFDBot
ENV UOFD_DB_FILE_PATH="/opt/UOFDBot/default/uofd.db"
ENV UOFD_CONFIG_FILE_PATH="/opt/UOFDBot/default/config.yml"
ENV UOFD_LANG_DIR_PATH="/opt/UOFDBot/default/lang"

COPY ./configs/config.yml /opt/UOFDBot/default/config.yml
COPY ./lang/ /opt/UOFDBot/default/lang/

ENTRYPOINT [ "/opt/UOFDBot/UOFDBot" ]