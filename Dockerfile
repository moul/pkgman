# dynamic config
ARG             BUILD_DATE
ARG             VCS_REF
ARG             VERSION

# build
FROM            golang:1.22.1-alpine as builder
RUN             apk add --no-cache git gcc musl-dev make
ENV             GO111MODULE=on
WORKDIR         /go/src/moul.io/pkgman
COPY            go.* ./
RUN             go mod download
COPY            . ./
RUN             make install

# minimalist runtime
FROM alpine:3.19.1
LABEL           org.label-schema.build-date=$BUILD_DATE \
                org.label-schema.name="pkgman" \
                org.label-schema.description="" \
                org.label-schema.url="https://moul.io/pkgman/" \
                org.label-schema.vcs-ref=$VCS_REF \
                org.label-schema.vcs-url="https://github.com/moul/pkgman" \
                org.label-schema.vendor="Manfred Touron" \
                org.label-schema.version=$VERSION \
                org.label-schema.schema-version="1.0" \
                org.label-schema.cmd="docker run -i -t --rm moul/pkgman" \
                org.label-schema.help="docker exec -it $CONTAINER pkgman --help"
COPY            --from=builder /go/bin/pkgman /bin/
ENTRYPOINT      ["/bin/pkgman"]
#CMD             []
