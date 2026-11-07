# SuggestIT

A suggestion box for P.R.I.T., the IT student division's committee at Chalmers.
Anyone can post a suggestion without logging in. Members who log in with
[Gamma](https://auth.chalmers.it) can read, delete and clear suggestions.

- `backend/`: Go (standard library `net/http`) with
  [zitadel/oidc](https://github.com/zitadel/oidc) for Gamma login and
  [lib/pq](https://github.com/lib/pq) for PostgreSQL.
- `frontend/`: React 19 + TypeScript, built with Vite+ and styled with
  Tailwind CSS and shadcn/ui components.

## Requirements

- Go (version in `backend/go.mod`)
- Node.js 24 and pnpm (version in `frontend/package.json`; `corepack enable` installs it)
- Docker with Docker Compose
- make

## Development

### 1. Create a Gamma client

Create an OAuth client in Gamma with the redirect URI
`http://localhost:3000/callback`. Then copy `.env.example` to `.env` and fill
in the values:

```sh
cp .env.example .env
```

| Variable              | Description                                                                 |
| --------------------- | --------------------------------------------------------------------------- |
| `OAUTH_CLIENT_ID`     | Gamma client ID                                                             |
| `OAUTH_CLIENT_SECRET` | Gamma client secret                                                         |
| `SESSION_SECRET`      | At least 32 bytes, used to sign and encrypt session cookies (`openssl rand -base64 48`). Optional in development: if it is empty, a random one is generated at startup and everyone is logged out when the backend restarts. |

`docker compose` reads `.env` automatically. `.env` is gitignored.

### 2. Start the backend and database

```sh
make up
```

This builds the backend binary to `backend/bin/suggestit` and starts:

| Service | URL                                                 |
| ------- | --------------------------------------------------- |
| backend | http://localhost:3001                               |
| db      | `localhost:5432` (user, password and database: `suggestit` / `password` / `suggestit`) |
| adminer | http://localhost:8080 (System: PostgreSQL, Server: `db`) |

The binary is bind-mounted into the backend container. After changing Go code,
run:

```sh
make restart   # rebuild the binary and restart the backend container
```

Other targets:

```sh
make logs      # follow container logs
make down      # stop the containers (the database volume is kept)
make build     # only build the backend binary
```

### 3. Start the frontend

```sh
make frontend  # same as: cd frontend && pnpm dev
```

On the first run, run `pnpm install` in `frontend/` first. Then open
http://localhost:3000. The Vite dev server proxies `/api` to the backend
(`http://127.0.0.1:3001` by default; override it with `SUGGESTIT_BACKEND_URL`).

### Checks

```sh
make -C backend check   # gofmt, go vet, go test
cd frontend && pnpm typecheck && pnpm check && pnpm build
```

These are the checks CI runs.

## Configuration

The backend is configured with environment variables:

| Variable                      | Required | Default         | Description |
| ----------------------------- | -------- | --------------- | ----------- |
| `SUGGESTIT_POSTGRES_HOST`     | yes      |                 | PostgreSQL host |
| `SUGGESTIT_POSTGRES_PORT`     | no       | `5432`          | PostgreSQL port |
| `SUGGESTIT_POSTGRES_USER`     | yes      |                 | PostgreSQL user |
| `SUGGESTIT_POSTGRES_PASSWORD` | no       | (empty)         | PostgreSQL password |
| `SUGGESTIT_POSTGRES_DB`       | yes      |                 | PostgreSQL database |
| `SUGGESTIT_POSTGRES_SSLMODE`  | no       | `disable`       | PostgreSQL `sslmode` |
| `OAUTH_URL`                   | yes      |                 | OIDC issuer, `https://auth.chalmers.it` |
| `OAUTH_CLIENT_ID`             | yes      |                 | Gamma client ID |
| `OAUTH_CLIENT_SECRET`         | yes      |                 | Gamma client secret |
| `CALLBACK_URL`                | yes      |                 | Redirect URI registered in Gamma |
| `SESSION_SECRET`              | yes*     |                 | At least 32 bytes. *May be empty only when `CALLBACK_URL` is `http://` on `localhost`, `127.0.0.1` or `::1` (a random secret is then generated at startup) |
| `COOKIE_DOMAIN`               | no       | (host-only)     | Domain attribute for the session cookie |
| `LISTEN_ADDR`                 | no       | `:3001`         | Address the HTTP server listens on |

### Gamma redirect URI

`CALLBACK_URL` must match a redirect URI registered on the Gamma client. Two
forms work:

- `https://suggestit.chalmers.it/callback` (the frontend route, used in
  production today): the frontend forwards the browser, including the query
  string, to `/api/auth/callback`.
- `https://suggestit.chalmers.it/api/auth/callback`: the backend directly.

For local development use `http://localhost:3000/callback`.

## API

All endpoints except `/healthz` are under `/api` and use JSON. Errors are returned as
`{"error": "..."}`. Endpoints marked "auth" require a session cookie and
return `401` without one.

| Method   | Path                      | Auth | Description |
| -------- | ------------------------- | ---- | ----------- |
| `GET`    | `/api/auth/login`         |      | Redirect to Gamma to log in |
| `GET`    | `/api/auth/callback`      |      | OIDC callback; sets the session cookie and redirects to `/` (`/?login=failed` on failure) |
| `POST`   | `/api/auth/logout`        |      | Clear the session cookie (`204`) |
| `GET`    | `/api/me`                 | auth | The logged-in user: `{cid, nickname, name, picture}` |
| `GET`    | `/api/suggestions`        | auth | All suggestions, newest first |
| `POST`   | `/api/suggestions`        |      | Create a suggestion from `{title, text, author}` (`201`) |
| `GET`    | `/api/suggestions/{id}`   | auth | One suggestion (`404` if it does not exist) |
| `DELETE` | `/api/suggestions/{id}`   | auth | Delete one suggestion (`204`) |
| `POST`   | `/api/suggestions/delete` | auth | Delete several suggestions from `{ids: [...]}` (`204`) |
| `GET`    | `/healthz`                |      | Health check, including the database |

A suggestion is `{id, timestamp, title, text, author}`, where `timestamp` is
RFC 3339 in UTC. The full OpenAPI description is in
[`backend/static/swagger.yaml`](backend/static/swagger.yaml).

## Production

GitHub Actions checks, builds and, on `master`, pushes
`ghcr.io/cthit/suggestit-backend:latest` and
`ghcr.io/cthit/suggestit-frontend:latest`. `prod.docker-compose.yml` runs them
together with PostgreSQL. A reverse proxy must serve both on the same origin:
`/api/*` to the backend (host port 5000, as before; the container listens on
3001) and everything else to the frontend (host port 8080). Set
`OAUTH_CLIENT_ID`, `OAUTH_CLIENT_SECRET` and `SESSION_SECRET` in the
environment.

### Upgrading from the previous stack

The previous `prod.docker-compose.yml` ran an untagged `postgres` image with
its data in an anonymous volume. The new file pins `postgres:18-alpine` and
stores data in the named volume `db-data`, which starts empty. Move the data
with a dump and restore before switching:

```sh
# 1. With the OLD compose file, while the old stack is still running:
docker compose -f prod.docker-compose.yml exec -T db \
  pg_dump -U suggestit -d suggestit --clean --if-exists > suggestit.sql

# 2. Remove the old containers (the old anonymous volume is kept as a backup):
docker compose -f prod.docker-compose.yml down

# 3. With the NEW compose file, start only the database and restore:
docker compose -f prod.docker-compose.yml up -d --wait db
docker compose -f prod.docker-compose.yml exec -T db \
  psql -U suggestit -d suggestit -v ON_ERROR_STOP=1 < suggestit.sql

# 4. Start everything:
docker compose -f prod.docker-compose.yml up -d
```

Do not skip step 2. If the old `db` container is recreated in place, Compose
carries its anonymous `/var/lib/postgresql/data` volume over, and
`postgres:18` refuses to start with data at that path. The backend keeps the
existing `suggestion` table as it is, so no schema migration is needed.
