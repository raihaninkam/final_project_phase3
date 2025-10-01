# Social Media App Backend

![badge golang](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![badge postgresql](https://img.shields.io/badge/PostgreSQL-316192?style=for-the-badge&logo=postgresql&logoColor=white)
![badge redis](https://img.shields.io/badge/redis-%23DD0031.svg?&style=for-the-badge&logo=redis&logoColor=white)

Welcome!
A social media platform built to make it easy for you to share moments, connect with friends, and express yourself anytime, anywhere. With a fast and reliable system, you’ll enjoy a smooth experience whether posting updates, exploring content, or engaging with your community.

# System Design Social Media App

Ekspektasi awal yakni aplikasi kecil. Start with Monolith. langsung menerapkan microservice pada semua aplikasi (overkill) terurama pada aplikasi baru yang trafficnya masih kecil. terlalu mengeluarkan resource tanpa memperhatikan kebutuhan.

![alt text](/image.png)

Ketika Kebutuhan aplikasi semakin kompleks. Lakukan Scalling Up salah satunya dengan menerapkan MicroService, yakni kita pecah aplikasi kita menjadi ke beberapa server dalam servicenya. menggunakan CDN untuk konten yang sifatnya statis dan jika user kita sudah ada di beberapa region

![alt text](/image-1.png)

## ERD

![alt text](/social_media.png)

### Relationship Type Cardinality

- User - Posts One-to-Many 1:N
- User - Comments One-to-Many 1:N
- User - Notifications One-to-Many 1:N
- User - Follows (as follower) Many-to-Many M:N
- User - Follows (as following) Many-to-Many M:N
- User - Likes Many-to-Many M:N
- Post - Comments One-to-Many 1:N
- Post - Likes Many-to-Many M:N

## 🔧 Tech Stack

- [Go](https://go.dev/dl/)
- [PostgreSQL](https://www.postgresql.org/download/)
- [Redis](https://redis.io/docs/latest/operate/oss_and_stack/install/archive/install-redis/install-redis-on-windows/)
- [JWT](https://github.com/golang-jwt/jwt)
- [argon2](https://pkg.go.dev/golang.org/x/crypto/argon2)
- [migrate](https://github.com/golang-migrate/migrate)
- [Docker](https://docs.docker.com/engine/install/ubuntu/#install-using-the-repository)
- [Swagger for API docs](https://swagger.io/) + [Swaggo](https://github.com/swaggo/swag)

## 🗝️ Environment

```bash
# database
DBUSER=<your_database_user>
DBPASS=<your_database_password>
DBNAME=<your_database_name
DBHOST=<your_database_host>
DBPORT=<your_database_port>

# JWT hash
JWT_SECRET=<your_secret_jwt>
JWT_ISSUER=<your_jwt_issuer>

# Redish
RDB_HOST=<your_redis_host>
RDB_PORT=<your_redis_port>
RDB_USER=<your_redis_user>
RDB_PWD=<your_redis_password>


```

## ⚙️ Installation

1. Clone the project

```sh
$ https://github.com/raihaninkam/final_project_phase3
```

2. Navigate to project directory

```sh
$ cd final_project_phase3
```

3. Install dependencies

```sh
$ go mod tidy
```

4. Setup your [environment](##-environment)

5. Install [migrate](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate#installation) for DB migration

6. Do the DB Migration

```sh
$ migrate -database YOUR_DATABASE_URL -path ./db/migrations up
```

or if u install Makefile run command

```sh
$ make migrate-createUp
```

7. Run the project

```sh
$ go run ./cmd/main.go
```

## 🚧 API Documentation

| Method | Endpoint               | Body                                                       | Description                      |
| ------ | ---------------------- | ---------------------------------------------------------- | -------------------------------- |
| POST   | /auth/register         | email:string, password:string                              | Register                         |
| POST   | /auth/login            | email:string, password:string                              | Login                            |
| POST   | /post                  | header: Authorization (token jwt) content:form, image:form | Create Post                      |
| GET    | /folllowing            | header: Authorization (token jwt)                          | Get Following List               |
| POST   | /follow/:user_id       | header: Authorization (token jwt)                          | Follow Some User                 |
| GET    | /POST                  | header: Authorization (token jwt)                          | Get Following Post               |
| POST   | /post/:post_id/like    | header: Authorization (token jwt)                          | Like Some Post by Id             |
| POST   | /post/:post_id/comment | header: Authorization (token jwt)                          | Comment Some Post by Id          |
| GET    | /post/:post_id/comment | header: Authorization (token jwt)                          | Get Comment of a post by post_id |
| PATCH  | /auth/profile          | header: Authorization (token jwt)                          | Update Profile                   |
| GET    | /post/popular          | header: Authorization (token jwt)                          | Get a Popular post               |
| DELETE | /post/:post_id/like    | header: Authorization (token jwt)                          | Unlike Some Post by Post_id      |

## 📄 LICENSE

MIT License

Copyright (c) 2025 raihaninkam

## 📧 Contact Info & Contributor

raihankamil37@gmail.com

[https://www.linkedin.com/in/raihan-insan-kamil/](https://github.com/raihaninkam)
