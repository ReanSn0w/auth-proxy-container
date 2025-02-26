# Auth Proxy Container

Контейнер работающий как middleware для авторизации пользователя
при доступе к определенный маршрутам на сервере.

### Пример испольования в docker-compose

В данном примере контейнер будет блокировать доступ ко всем страницам habr.com,
при этом даже после авторизации увидет сайт только пользователь, чей ID задан в USER_ID.

```yaml
services:
  auth-proxy:
    image: auth-proxy
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
