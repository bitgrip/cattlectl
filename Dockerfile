# STEP 1 build executable binary
FROM golang:alpine as builder
# Install git + SSL ca certificates
RUN apk update && apk add git && apk add ca-certificates
# Create user cattlectl
RUN adduser -D -g '' cattlectl \
&& mkdir /cattlectl \
&& touch /cattlectl/.keep \
&& chown -R cattlectl:cattlectl /cattlectl
COPY . $GOPATH/src/github.com/bitgrip/cattlectl/
WORKDIR $GOPATH/src/github.com/bitgrip/cattlectl/

ENV GO111MODULE=on \
CGO_ENABLED=0 \
GOOS=linux \
GOARCH=amd64

RUN go test -mod=vendor -v ./...

#build the binary
RUN echo "Building VERSION=$(git describe --tags)" 1>&2 && \
go build \
-ldflags "-X github.com/bitgrip/cattlectl/internal/pkg/ctl.Version=$(git describe --tags) -d -s -w -extldflags \"-static\"" \
-a -tags netgo -installsuffix netgo \
-mod=vendor \
-o /go/bin/cattlectl

# STEP 2 package the result image
FROM scratch

ARG BUILD_DATE
ARG VCS_REF

LABEL org.label-schema.schema-version="1.0" \
    org.label-schema.build-date=$BUILD_DATE \
    org.label-schema.name="bitgrip/cattlectl" \
    org.label-schema.vendor="Bitgrip GmbH" \
    org.label-schema.license="Apache 2.0" \
    org.label-schema.vcs-url="https://github.com/bitgrip/cattlectl.git" \
    org.label-schema.vcs-ref=$VCS_REF

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /go/bin/cattlectl /bin/cattlectl
COPY --from=builder /cattlectl /cattlectl
COPY sample/data /data
ENV HOME=/cattlectl
USER cattlectl
WORKDIR /data
ENTRYPOINT ["/bin/cattlectl"]
