FROM golang:1.19.0-alpine as builder
WORKDIR /
COPY . .
RUN apk add --no-cache gcc musl-dev
RUN GOPROXY=https://goproxy.io go build -o ./uask_node ./cmd/uask_node/full_node.go

FROM alpine:latest
WORKDIR /
COPY --from=builder /uask_node /
ENTRYPOINT ["/uask_node"]
