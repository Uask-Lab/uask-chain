FROM golang:1.19.0-alpine as builder
WORKDIR /
COPY . .
RUN apk add --no-cache gcc musl-dev
RUN go build -o ./uask_node ./cmd/uask_node/full_node.go

FROM alpine:latest
WORKDIR /
COPY --from=builder /uask_node /
EXPOSE 7999/tcp 8999/tcp
ENTRYPOINT ["/uask_node"]
