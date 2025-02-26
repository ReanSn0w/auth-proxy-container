# Сборка приложения
FROM golang:1.23-alpine AS application
ARG TAG
ADD . /bundle
WORKDIR /bundle

RUN apk --no-cache add ca-certificates

RUN \
  version=${TAG} && \
  echo "Building service. Version: ${version}" && \
  go build -ldflags "-X main.build=${version}" -o /srv/app ./cmd/app/main.go

# Финальная сборка образа
FROM scratch
COPY --from=application /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=application /srv /srv
EXPOSE 8080
WORKDIR /srv
ENTRYPOINT ["/srv/app"]
