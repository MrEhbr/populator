# build
FROM            golang:1.15-alpine as builder
# dynamic config
ARG             BUILD_DATE
ARG             VCS_REF
ARG             VERSION

RUN             apk add --no-cache git gcc musl-dev make bash
ENV             GO111MODULE=on
WORKDIR         /go/src/github.com/MrEhbr/golang-repo-template
COPY            go.* ./
RUN             go mod download
COPY            . ./
RUN             make install VCS_REF=$VCS_REF VERSION=$VERSION BUILD_DATE=$BUILD_DATE

# minimalist runtime
FROM alpine:3.11
# dynamic config
ARG             BUILD_DATE
ARG             VCS_REF
ARG             VERSION

LABEL org.label-schema.build-date=$BUILD_DATE \
    org.label-schema.name="golang-repo-template" \
    org.label-schema.description="" \
    org.label-schema.url="" \
    org.label-schema.vcs-ref=$VCS_REF \
    org.label-schema.vcs-url="https://github.com/MrEhbr/golang-repo-template" \
    org.label-schema.vendor="Alexey Burmistrov" \
    org.label-schema.version=$VERSION \
    org.label-schema.schema-version="1.0" \
    org.label-schema.cmd="docker run -i -t --rm MrEhbr/golang-repo-template" \
    org.label-schema.help="docker exec -it $CONTAINER golang-repo-template --help"
COPY            --from=builder /go/bin/golang-repo-template /bin/
ENTRYPOINT      ["/bin/golang-repo-template"]
#CMD             []
