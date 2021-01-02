FROM golang:1.14.3-alpine AS build
WORKDIR /src
ENV CGO_ENABLED=0
COPY . .
ARG TARGETOS
ARG TARGETARCH
RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o /bin/goshort .


FROM redis:alpine
COPY --from=build /bin/goshort /bin/goshort
COPY redis-docker-commands.sh /scripts/redis-docker-commands.sh
RUN ["chmod", "+x", "/scripts/redis-docker-commands.sh"]
RUN ["chmod", "+x", "/bin/goshort"]
ENTRYPOINT ["sh", "/scripts/redis-docker-commands.sh"]