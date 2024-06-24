FROM golang:1.22-alpine AS blog

RUN go install github.com/nisanthchunduru/hugo-notion@latest

WORKDIR /opt/blog

CMD hugo-notion -r 5
