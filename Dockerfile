# -----------------------------------------------------------------------------
#  Build Stage
# -----------------------------------------------------------------------------
FROM --platform=$BUILDPLATFORM tonistiigi/xx AS xx
FROM --platform=$BUILDPLATFORM golang:alpine3.20 AS build

ARG TARGETOS
ARG TARGETARCH

RUN apk add clang lld
COPY --from=xx / /
RUN xx-apk add --no-cache \
 gcc \
 musl-dev

ENV CGO_ENABLED=1

WORKDIR /uofd

COPY . /uofd/

RUN cd /uofd && \
    xx-go mod tidy
RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} xx-go build -tags musl --ldflags "-extldflags -static"

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