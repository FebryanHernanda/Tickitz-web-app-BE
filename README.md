# Tickitz Backend

![badge golang](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![badge postgresql](https://img.shields.io/badge/PostgreSQL-316192?style=for-the-badge&logo=postgresql&logoColor=white)
![badge redis](https://img.shields.io/badge/Redis-DC382D?style=for-the-badge&logo=redis&logoColor=white)
![badge docker](https://img.shields.io/badge/Docker-2CA5E0?style=for-the-badge&logo=docker&logoColor=white)
![badge swagger](https://img.shields.io/badge/Swagger-85EA2D?style=for-the-badge&logo=swagger&logoColor=black)

Backend project for [frontend Tickitz Web App](https://github.com/FebryanHernanda/tickitz-web-app-react). This service powers a cinema ticketing web app: browse movies, find cinemas and schedules, pick seats in real time, and place secure orders. Built with Gin (Go), PostgreSQL, and Redis.

## 🔧 Tech Stack

- [Go](https://go.dev/dl/)
- [Gin](https://gin-gonic.com/)
- [PostgreSQL](https://www.postgresql.org/download/)
- [Redis](https://redis.io/)
- [JWT](https://github.com/golang-jwt/jwt)
- [Bcrypt](https://pkg.go.dev/golang.org/x/crypto/bcrypt)
- [golang-migrate](https://github.com/golang-migrate/migrate)
- [Docker](https://docs.docker.com/get-docker/)
- [Swagger + Swaggo](https://github.com/swaggo/swag)

## 🗝️ Environment

```bash
# PostgreSQL
DBHOST=<your_db_host>
DBPORT=<your_db_port>
DBUSER=<your_db_user>
DBPASSWORD=<your_db_password>
DBNAME=<your_db_name>

# Redis
RDBHOST=<your_redis_host>
RDBPORT=<your_redis_port>

# JWT
JWTKEY=<your_jwt_secret>

```

## ⚙️ Installation

1. Clone the project

```sh
git clone https://github.com/FebryanHernanda/Tickitz-web-app-BE.git
cd Tickitz-web-app-BE
```

2. Install dependencies

```sh
go mod tidy
```

3. Configure your environment

4. Install migrate (for DB migration)

- Follow official install guide for your OS: https://github.com/golang-migrate/migrate/tree/master/cmd/migrate#installation

5. Run DB migrations

```sh
make migrate-up
```

6. Run the server

```sh
go run ./cmd/main.go
```

7. Optional: Docker Compose

```sh
docker compose up -d
```

## 🚧 API Documentation

Swagger UI will be served when the app is running:

```
http://localhost:8080/swagger/index.html
```

## 🧭 REST Endpoints

### Authentication

| Method | Endpoint       | Body / Headers                | Description              |
| ------ | -------------- | ----------------------------- | ------------------------ |
| POST   | /auth/register | email, password               | Register new user        |
| POST   | /auth/login    | email, password               | Login and get JWT        |
| POST   | /auth/logout   | Authorization: Bearer <token> | Logout + blacklist token |

### Movies

| Method | Endpoint             | Query / Body                               | Description                      |
| ------ | -------------------- | ------------------------------------------ | -------------------------------- |
| GET    | /movies              | page:int, search:string, genres:[]string   | List movies with filters         |
| GET    | /movies/popular      |                                            | Popular movies                   |
| GET    | /movies/upcoming     |                                            | Upcoming movies                  |
| GET    | /movies/genres       |                                            | List available genres            |
| GET    | /movies/casts        |                                            | List casts                       |
| GET    | /movies/directors    |                                            | List directors                   |
| GET    | /movies/{id}/details | path: id:int                               | Movie details by ID              |
| GET    | /movies/schedules    |                                            | Aggregated schedules for a movie |

### Cinemas

| Method | Endpoint                                      | Query / Body                 | Description                    |
| ------ | --------------------------------------------- | ---------------------------- | ------------------------------ |
| GET    | /cinemas/list                                 |                              | List cinemas                   |
| GET    | /cinemas/location                             |                              | Cinemas by location            |
| GET    | /cinemas/{movieId}                            | path: movieId:int            | Cinemas showing specific movie |
| GET    | /cinemas/available-seats/{cinema_schedule_id} |                              | Available seats for a schedule |

### Orders

| Method | Endpoint        | Headers / Body                                                                                  | Description            |
| ------ | --------------- | ----------------------------------------------------------------------------------------------- | ---------------------- |
| POST   | /orders         | Authorization: Bearer <token>, schedule_id:int, payment_id:int, seats:[]string, total_price:int | Create new order       |
| GET    | /orders/history | Authorization: Bearer <token>                                                                   | Get user order history |

### Profile

| Method | Endpoint              | Headers / Body                                              | Description      |
| ------ | --------------------- | ----------------------------------------------------------- | ---------------- |
| GET    | /profile              | Authorization: Bearer <token>                               | Get user profile |
| PATCH  | /profile/edit         | Authorization: Bearer <token>, first_name, last_name, phone, etc | Update profile   |
| PATCH  | /profile/editpassword | Authorization: Bearer <token>, password                     | Change password  |

### Admin

| Method | Endpoint                             | Headers / Body                                                                                                          | Description                 |
| ------ | ------------------------------------ | ----------------------------------------------------------------------------------------------------------------------- | --------------------------- |
| GET    | /admin/movies                        | Authorization: Bearer <admin_token>, page:int                                                                           | Admin movie list            |
| POST   | /admin/movies/add                    | Authorization: Bearer <admin_token>, title, poster_path, backdrop_path, overview, duration, casts[], director, genres[] | Create movie                |
| PATCH  | /admin/movies/edit/{id}              | Authorization: Bearer <admin_token>, path: id:int, fields to update                                                     | Update a movie              |
| DELETE | /admin/movies/delete/{id}            | Authorization: Bearer <admin_token>, path: id:int                                                                       | Hard delete movie           |
| POST   | /admin/movies/cinemaschedule/add     | Authorization: Bearer <admin_token>, movie_id, cinema_id, room, date, time, price                                       | Add cinema schedule         |
| GET    | /admin/movies/schedule               | Authorization: Bearer <admin_token>, movie_id:int                                                                       | List schedules (admin view) |
| GET    | /admin/movies/{movieId}/edit-details | Authorization: Bearer <admin_token>, path: movieId:int                                                                  | Get editable movie details  |

Notes:

- All protected endpoints require Authorization header with a valid Bearer token.
- Seat arrays should be sent as JSON arrays of seat codes (e.g., ["A1","A2"]).
- Dates/times use ISO-8601 where applicable.

## 📄 License

MIT License

Copyright (c) 2025 Febryan Hernanda

## 🎯 Related Project

[Frontend Tickitz](https://github.com/FebryanHernanda/tickitz-web-app-react)
