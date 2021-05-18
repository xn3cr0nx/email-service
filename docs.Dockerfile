# FROM dependencies as builder 
FROM golang:1.16.2 as builder
LABEL maintainer="Patrick Jusic <patrick.jusic@docs.com>"

WORKDIR /docs

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go get -u github.com/swaggo/swag/cmd/swag@v1.6.7
RUN swag init -g cmd/api/main.go

RUN CGO_ENABLED=0 GOOS=linux go build -mod=mod -o ./bin/docs -tags docs ./cmd/docs

################################################################################
# Final image
################################################################################
FROM golang:alpine

WORKDIR /root/
COPY --from=builder /docs/bin/docs .

EXPOSE 8085

CMD ["./docs"]