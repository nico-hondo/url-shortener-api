# ShortLink Backend
![MIT License](https://img.shields.io/badge/License-MIT-yellow?style=for-the-badge)

REST API server for the URL Shortener application, written in Go. It handles authentication, short link creation, and data storage for the frontend.

### Technology Stack
![Go](https://img.shields.io/badge/Go-1.27.1-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![Gin](https://img.shields.io/badge/Gin-1.12.0-008ECF?style=for-the-badge&logo=gin&logoColor=white)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-Database-4169E1?style=for-the-badge&logo=postgresql&logoColor=white)
![Redis](https://img.shields.io/badge/Redis-Cache-DC382D?style=for-the-badge&logo=redis&logoColor=white)
![JWT](https://img.shields.io/badge/JWT-Auth-000000?style=for-the-badge&logo=jsonwebtokens&logoColor=white)
![Swagger](https://img.shields.io/badge/Swagger-OpenAPI-85EA2D?style=for-the-badge&logo=swagger&logoColor=black)

### Features
- User registration, login, and logout.
- JWT-based authentication to secure protected endpoints.
- Short link creation from long URLs, with an optional custom slug or an automatically generated random slug.
- Soft deletion of links.
- Redirection from short links to their original URLs.
- Redis caching of link lookups for faster redirects.
- Interactive Swagger API documentation, available at `/swagger/index.html`.

### Prerequisites
Before you running this project locally, make sure you have:
- Git
- Go 1.26.3
- PostgreSQL
- Redis
- [golang-migrate](https://github.com/golang-migrate/migrate)
- `psql` command-line client for running seeders
- Optional: [Air](https://github.com/air-verse/air) for live reload
- Optional: [Swag CLI](https://github.com/swaggo/swag) for regenerating Swagger docs

## Setup Instruction

Clone the repository:

```bash
git clone https://github.com/nico-hondo/url-shortener-api.git
cd shortlink-backend
```

Copy the environment file and update the values:

```bash
cp .env.example .env
```

Configure these required values in `.env`:

```env
HOST=localhost
PORT=8080

DBUSER=user_db
DBPASS=pass_db
DBHOST=host_db
DBPORT=port_db
DBNAME=url_shortener

JWT_ISSUER=shortlink
JWT_SECRET=your-secret-key

RDB_ADDR=localhost
RDB_PORT=port_rdb
RDB_USER=user_rdb
RDB_PASS=pass_rdb
```

Install dependencies:

```bash
go mod download
```

Run database migrations:

```bash
make migrate-up
```

Run seeders if needed:

```bash
make seed
```

Start the application:

```bash
go run ./cmd/main.go
```

The API will run at:

```text
http://localhost:8080
```

Init Swagger documentation:

```bash
swag init -g cmd/main.go -o docs
```

Open Swagger documentation at:

```text
http://localhost:8080/swagger/index.html
```

For live reload development, install Air and run:

Here for [.air.toml](https://gist.github.com/nico-hondo/2c968070d9e3311cce91e66ebae7d48c) configuration file for Air live reload.

- then
```bash
air
```
## How to Contribute

1. Fork this repository.
2. Create a new branch:

   ```bash
   git checkout -b feature/your-feature-name
   ```

3. Make your changes.
4. Commit your changes:

   ```bash
   git commit -m "feat: add your feature"
   ```

5. Push your branch and open a pull request.

## Related Project

- [ShortLink Frontend](https://github.com/nico-hondo/url-shortener.git)

## License

This project is licensed under the [MIT License](LICENSE).