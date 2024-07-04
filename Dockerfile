FROM golang:1.22-alpine AS hugo-notion

RUN apk add --no-cache git

ARG HUGO_NOTION_VERSION=latest

RUN go install github.com/nisanthchunduru/hugo-notion@${HUGO_NOTION_VERSION}

RUN mkdir /opt/hugo-site

WORKDIR /opt/hugo-site

CMD hugo-notion -r 5
