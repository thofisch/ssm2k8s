FROM golang:1.12 as builder
LABEL maintainer="Thomas Fischer <thfis@dfds.com>"

WORKDIR ${GOPATH}/src/github.com/thofisch/ssm2k8s

ENV GO111MODULE=on

COPY go.mod .
COPY go.sum .

# COPY Gopkg.toml Gopkg.lock ./
# COPY vendor vendor
# ARG DEP_ENSURE=""
# RUN if [ ! -z "${DEP_ENSURE}" ]; then \
#       go get -u github.com/golang/dep/cmd/dep && \
#       dep ensure --vendor-only; \
#     fi

RUN go mod download

COPY . .

# Download all the dependencies
#RUN go get -d -v ./...

# Install the package
#RUN go install -v ./...

RUN CGO_ENABLED=0 GOOS=linux go build -a -v -installsuffix cgo -o /go/bin/mysticod ./cmd/mysticod
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
