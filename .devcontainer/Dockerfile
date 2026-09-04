FROM golang:1.26-alpine

# Setup bash and other utilities
RUN apk add --no-cache bash git make curl gcc musl-dev

# Setup Installing development tools
RUN go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.62.0 && \
    go install github.com/daixiang0/gci@latest && \
    go install golang.org/x/tools/cmd/goimports@latest

WORKDIR /workspace

RUN adduser -D vscode
RUN chown -R vscode:vscode /workspace
USER vscode

ENTRYPOINT [ "sleep", "infinity" ]