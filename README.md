# Auth Proxy Container

Контейнер работающий как middleware для авторизации пользователя
при доступе к определенный маршрутам на сервере.

### Пример испольования в docker-compose

В данном примере контейнер будет блокировать доступ ко всем страницам habr.com,
при этом даже после авторизации увидет сайт только пользователь, чей ID задан в USER_ID.

```yaml
services:
  auth-proxy:
    image: ghcr.io/reansn0w/auth-proxy-container
    ports:
      - 8080:8080
    environment:
      OAUTH_PROVIDER: yandex
      OAUTH_CLIENT_ID: <client_id>
      OAUTH_CLIENT_SECRET: <client_secret>
      USER_ID: <user_id>
      PRIVATE: ^/(.*)
      OUTPUT: https://habr.com
```

### Параметры

```bash
Usage:
  app [OPTIONS]

Application Options:
      --debug                   enable debug mode [$DEBUG]
  -p, --port=                   application listen port (default: 8080) [$PORT]
  -o, --output=                 output URL (default: https://yandex.ru) [$OUTPUT]
      --private=                private path regexp (default: ^/(.*)) [$PRIVATE]
      --uid=                    user ID regexp (default: *) [$USER_ID]

oauth:
      --oauth.provider=[yandex] OAuth provider (default: yandex) [$OAUTH_PROVIDER]
      --oauth.client-id=        OAuth client ID [$OAUTH_CLIENT_ID]
      --oauth.client-secret=    OAuth client secret [$OAUTH_CLIENT_SECRET]

Help Options:
  -h, --help                    Show this help message
```

### Также важно знать

Сервер будет отправлять запросы к oauth провайдеру с redirect_uri равным URL сервера, без Query параметров и с Path = /oauth/authorize

```
# Пример:

https://mydomain.com -> https://mydomain.com/oauth/authorize
https://mydomain.com/pages/some -> https://mydomain.com/oauth/authorize
https://mydomain.com?some_parameter=some_value -> https://mydomain.com/oauth/authorize
https://mydomain.com:4444/pages/some -> https://mydomain.com:4444/oauth/authorize
```

Следует учитывать этот момент при настройке приложения на стороне oauth провайдера
