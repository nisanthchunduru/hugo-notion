FROM golang:1.22-alpine AS blog

ARG HUGO_NOTION_VERSION=latest

RUN go install github.com/nisanthchunduru/hugo-notion@${HUGO_NOTION_VERSION}

WORKDIR /opt/blog

CMD hugo-notion -r 5
