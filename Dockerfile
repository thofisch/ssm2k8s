FROM golang:1.13 as builder
LABEL maintainer="Thomas Fischer <thfis@dfds.com>"
ARG LDFLAGS

WORKDIR ${GOPATH}/src/github.com/thofisch/ssm2k8s

ENV GO111MODULE=on

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "${LDFLAGS}" -a -installsuffix cgo -o /go/bin/mysticod ./cmd/mysticod
# RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go install -a -tags netgo -ldflags '-w -extldflags "-static"' ./cmd/weaviate-server

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /go/bin/mysticod .
#
#EXPOSE 8080
#
CMD ["./mysticod"]
