# chirpy

Chirpy is a [Boot.dev](https://www.boot.dev) project!

The goal of this project was to improve upon my Backend development knowledge by building an HTTP server in [Go](https://go.dev/).

The course that I completed can be viewed [here](https://www.boot.dev/courses/learn-http-servers-golang).

## Overview
Chirpy is a social network application that allows for posting and removing messages (chirps) by an authenticated user. 

Additional features will be outlined in the API Routes section.

## Setup
1) [Clone repo locally](https://docs.github.com/en/repositories/creating-and-managing-repositories/cloning-a-repository)
2) Ensure that [Go 1.20+](https://go.dev/dl/) is installed
3) Ensure that [PostgreSQL 15+](https://www.postgresql.org/download/) is installed (or via [brew](https://formulae.brew.sh/formula/postgresql@18) if using a Mac computer)
4) Create a local `.env` file with the following values:
```
DB_URL="{YOUR_POSTGRES_CONNECTION_STRING}/chirpy?sslmode=disable"
PLATFORM="dev"
SECRET="{YOUR_SECRET_STRING}"
POLKA_KEY="{YOUR_POLKA_KEY}"
```
5) Start PostgreSQL
```
IE: brew services start postgresql@18
```
6) Run migrations (Note: install [goose](https://pressly.github.io/goose/installation/) if not currently on your machine)
```
Pouplate database:
goose postgres "postgres://postgres:@localhost:5432/chirpy" up

Rollback population(s):
goose postgres "postgres://postgres:@localhost:5432/chirpy" down
```
7) Start the app locally
```
Option 1: go run .

Option 2: go build -o out && ./out
```
8) You should receive the following notice upon successful app start: `"Serving on port: 8080"`

## API Routes
**Note:** All routes served on http://localhost:8080/{Route}

### ## GET:

#### **Route:** `/admin/metrics`
**Description**: Returns Admin page with number of times Chirpy has been visited

#### **Route:** `/api/healthz`
**Description**: Returns a page with status page 200 if the app is running

#### **Route:** `/api/chirps`
**Description**: Get all chirps that are currently in the database

#### **Route:** `/api/chirps/{chirpID}`
**Description**: Get a single chirp by its ID

#### **Route:** `/api/chirps?author_id=${userID}`
**Description**: Get all chirps by author

#### **Route:** `/api/chirps?sort=asc`
**Description**: Get all chirps in ascending order

#### **Route:** `/api/chirps?sort=desc`
**Description**: Get all chirps in descending order

### ## POST:
#### **Route:** `/admin/reset`
**Description**: Delete all users from the database and reset number of times Chirpy has been visited

#### **Route:** `/api/chirps`
**Description**: Creates a new chirp for authenticated user

**Body**:
```
TBD
```

#### **Route:** `/api/users`
**Description**: Creates new user

**Body**:
```
TBD
```

#### **Route:** `/api/login`
**Description**: Logs in a valid user

**Body**:
```
TBD
```

#### **Route:** `/api/refresh`
**Description**: Refreshes access token for a user


#### **Route:** `/api/revoke`
**Description**: Revoked access token for a user

#### **Route:** `/api/polka/webhooks`
**Description**: Webhook for mock payment processor - grants access to premium Chirpy Red Membership

**Body**:
```
TBD
```

### ## PUT:
#### **Route:** `/api/users`
**Description**: Updates credentials (username and password) for a user

**Body**:
```
TBD
```

### ## DELETE:
#### **Route:** `/api/chirps/{chirpID}`
**Description**: Deletes a chirp for a user

**Body**:
```
TBD
```

