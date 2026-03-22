# -----------------------------------------------------------------------------
#  Build Stage
# -----------------------------------------------------------------------------
FROM --platform=$BUILDPLATFORM tonistiigi/xx AS xx
FROM --platform=$BUILDPLATFORM dhi.io/golang:1.26.1-alpine3.23-dev AS build

ARG TARGETPLATFORM
ARG TARGETOS
ARG TARGETARCH

RUN apk add --update --no-cache zig
COPY --from=xx / /
RUN xx-apk add --no-cache musl-dev 

ENV CGO_ENABLED=1
ENV XX_GO_PREFER_C_COMPILER=zig

WORKDIR /uofd

COPY . /uofd/

RUN cd /uofd && \
    xx-go mod tidy
RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} xx-go build -tags musl --ldflags "-extldflags -static"

# -----------------------------------------------------------------------------
#  Run Stage
# -----------------------------------------------------------------------------
FROM dhi.io/alpine-base:3.23

COPY --from=build /uofd/UOFDBot /opt/UOFDBot/UOFDBot

WORKDIR /opt/UOFDBot
ENV UOFD_DB_FILE_PATH="/opt/UOFDBot/default/uofd.db"
ENV UOFD_CONFIG_FILE_PATH="/opt/UOFDBot/default/config.yml"
ENV UOFD_LANG_DIR_PATH="/opt/UOFDBot/default/lang"

COPY ./configs/config.yml /opt/UOFDBot/default/config.yml
COPY ./lang/ /opt/UOFDBot/default/lang/

ENTRYPOINT [ "/opt/UOFDBot/UOFDBot" ]