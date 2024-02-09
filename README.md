# Project bikesRentalAPI

Simple API used to complete the **Take-Home Test: Backend Engineer - Go Developer** application at [Forest](https://www.humanforest.co.uk).

## Getting Started

Copy the env.example file to add your environment variables.

```bash
cp env.example .env
```

For manage migrations please use [golang-migrate CLI](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)

For mocking this project please use [GoMock](https://github.com/uber-go/mock)
To install
```bash
go install go.uber.org/mock/mockgen@latest
```
Example of use
```bash
mockgen -source=internal/users/repository/repository.go -destination=internal/users/repository/mocks/repository_mock.go -package=mocks
```

## MakeFile

There is a Make file to make easy to run this project.

run all make commands with clean tests

```bash
make all build
```

build the application

```bash
make build
```

run the application

```bash
make run
```

This project could use live-reloading for development with [Air](https://github.com/cosmtrek/air) tool.

```bash
make watch
```

run the test suite

```bash
make test
```

clean up binary from the last build

```bash
make clean
```

---

## Test Definition - Take-Home Test: Backend Engineer - Go Developer

You are tasked with creating a simple RESTful API service for a bike rental company. The service will be used by the company’s mobile app to handle user authentication and bike rental operations. The service will also be used by the company’s admin dashboard to manage bikes, users, and rentals.

### Requirements

1. The service should be written in Go version 1.21 or higher. Use Go modules for dependency management.
2. The service should use the [Go-chi router](https://github.com/go-chi/chi), or the [standard library](https://golang.org/pkg/net/http/) for routing. Do not use frameworks like Gin, or Echo for this test.
3. The service should use [JWT](https://jwt.io/) for user authentication. You can use whatever JWT library you prefer or write your own.
4. The service should user Basic Auth for admin authentication.
5. The service should use [SQLite](https://www.sqlite.org/index.html) for data storage. Do not use an ORM library for this test, instead use the [database/sql](https://golang.org/pkg/database/sql/) package or [sqlx](https://github.com/jmoiron/sqlx), and any SQLite driver of your choice (for example, [go-sqlite3](https://github.com/mattn/go-sqlite3)).
6. You can use any other third-party libraries or tools you need.

#### Optional requirements

1. Unit tests.
2. Containerization.
3. API documentation.
4. Logging.
5. Configuration management.
6. Migration scripts.

### API Specification

#### Database Schema

The service will handle user authentication and will interact with a SQLite database. The database should contain three tables: `users`, `bikes`, and `rentals`.

##### Users

| Column          | Description                                 |
|-----------------|---------------------------------------------|
| id              | Identifier of the user                      |
| email           | Email address of the user                   |
| hashed_password | Hashed password of the user                 |
| first_name      | First name of the user                      |
| last_name       | Last name of the user                       |
| created_at      | Timestamp of when the user was created      |
| updated_at      | Timestamp of when the user was last updated |

##### Bikes

|    column    |                    description                   |
|--------------|--------------------------------------------------|
| id           | Identifier of the bike                           |
| is_available | Indicates whether the bike is available for rent |
| latitude     | Latitude of the bike                             |
| longitude    | Longitude of the bike                            |
| created_at   | Timestamp of when the bike was created           |
| updated_at   | Timestamp of when the bike was last updated      |
| updated_at   | Timestamp of when the user was last updated      |

##### Rentals

|      column     |                    description                   |
|-----------------|--------------------------------------------------|
| id              | Identifier of the rental                         |
| user_id         | Identifier of the user who rented the bike       |
| bike_id         | Identifier of the bike that was rented           |
| status          | Status of the rental (running, ended)            |
| start_time      | Start time of the rental                         |
| end_time        | End time of the rental                           |
| start_latitude  | Latitude of the bike at the start of the rental  |
| start_longitude | Longitude of the bike at the start of the rental |
| end_latitude    | Latitude of the bike at the end of the rental    |
| end_longitude   | Longitude of the bike at the end of the rental   |

![ERD](erd_diagram.svg)

Feel free to extend the database schema if necessary.

### Authentication

The service should have two types of authentication: user authentication and admin authentication.

#### User Authentication

The service will handle user authentication. The mobile app will use the service to register new users, login existing users, and retrieve the profile of the logged-in user. For user authentication, the service should use JWT.

The JWT claims should contain the following information:

```json
{
    "sub": 1, // user ID
    "exp": 1620000000,
    "email": "name@example.com",
    "first_name": "John",
    "last_name": "Doe"
}
```

The `exp` claim should be set to the expiry time of the JWT. The expiry time should be 30 days from the time of issue. The JWT should be signed using the HMAC SHA256 algorithm. The secret key used to sign the JWT should be stored in an environment variable named `JWT_SECRET`, you can use any value for the secret key for this test.

#### Admin Authentication

The service will handle admin authentication. The admin dashboard will use the service to login admins. For admin authentication, the service should use Basic Auth.

The admin credentials should be stored in an environment variable named `ADMIN_CREDENTIALS`. The value of the environment variable should be a Base64 encoded string of the admin username and password, separated by a colon. For example, if the admin username is `admin` and the password is `password`, the value of the environment variable should be `YWRtaW46cGFzc3dvcmQ=`.

### Endpoints

#### User-Related Endpoints

All the following endpoints require user authentication, except for registration and login.

1. **User Authentication and Management**
    * `POST /users/register`: Register a new user.
    * `POST /users/login`: Authenticate a user and return a JWT.
    * `GET /users/profile`: Retrieve the profile of the logged-in user.
    * `PATCH /users/profile`: Update user profile details.
2. **Bike Rental Operations**
    * `GET /bikes/available`: List all available bikes for rent.
    * `POST /rentals/start`: Start a bike rental.
    * `POST /rentals/end`: End a bike rental and return the bike.
    * `GET /rentals/history`: Retrieve the rental history of the logged-in user.

#### Administrative Endpoints

All the following endpoints require admin authentication.

1. **Bike Management**
    * `POST /admin/bikes`: Add a new bike to the system.
    * `PATCH /admin/bikes/{bike_id}`: Update details of a specific bike.
    * `GET /admin/bikes`: List all bikes in the system.
2. **User Management**
    * `GET /admin/users`: List all registered users.
    * `GET /admin/users/{user_id}`: Retrieve details of a specific user.
    * `PATCH /admin/users/{user_id}`: Update user details.
3. **Rental Management**
    * `GET /admin/rentals`: List all bike rentals.
    * `GET /admin/rentals/{rental_id}`: Get details of a specific rental.
    * `PATCH /admin/rentals/{rental_id}`: Update rental details.

#### Utility Endpoints

* `GET /status`: Health check endpoint to ensure the API is running.

### Business Logic

1. **Concurrent Bike Rentals**: A bike that is currently rented cannot be rented by another user. The system should enforce this rule and return an appropriate error message if there’s an attempt to rent an already rented bike.
2. **Bike Availability**: A bike that is currently rented cannot be listed as available.
3. **Simultaneous Bike Rentals**: A user can only rent one bike at a time. The system should enforce this rule and return an appropriate error message if there’s an attempt to rent more than one bike at a time.
4. **Start and End Location**: The system should record the start and end location of a bike rental. The start location should be the current location of the bike at the time of the rental. And for this test, the end location should be a random location within a 5km radius of the start location.
5. **Unique Email Address**: The system should enforce the uniqueness of email addresses. The system should return an appropriate error message if there’s an attempt to register a user with an email address that already exists in the system.

#### Optional bussines logic

1. **Bike Rental Duration**: The system should calculate the duration of a bike rental in minutes, rounded up to the nearest minute. This data should be stored in the `rentals` table.
2. **Bike Rental Cost**: You could add a price per minute to the `bikes` table, and calculate the cost of a bike rental based on the duration of the rental and the price per minute and store the cost in the `rentals` table. The price per minute could be fixed for all bikes, or it could be different for each bike. This data should be returned in the response of the `GET /rentals/history` endpoint and the admin endpoints for rental management.

### Submission Instructions

Submit your code in a zip file containing your source code, along with the SQLite database file. Ensure that any necessary setup instructions are included in your README file. Send the zip file attached to the email you received with the test instructions. You have **7 days** to complete the test from the time you received the email.

Good luck!
