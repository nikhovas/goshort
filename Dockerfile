FROM golang:1.14.3-alpine AS build
WORKDIR /src
ENV CGO_ENABLED=0
COPY . .
ARG TARGETOS
ARG TARGETARCH
RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o /bin/goshort .

FROM golang:1.14.3-alpine
COPY --from=build /bin/UrlShortener /bin/goshort
ENTRYPOINT ["/bin/goshort"]