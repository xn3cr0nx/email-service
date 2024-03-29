################################################################################
# Builder
################################################################################
# FROM dependencies as builder 
FROM golang:1.18 as builder
LABEL maintainer="Patrick Jusic <patrick.jusic@toggl.com>"

WORKDIR /service

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV GOPATH=""
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
RUN go build -trimpath -v -a -ldflags="-w -s" -o ./bin/service ./cmd/service

################################################################################
# Final image
################################################################################
FROM golang:alpine

WORKDIR /root/
COPY --from=builder /service/bin/service .
COPY --from=builder /service/templates/ ./templates/

CMD ./service