# Go authentication microservice
## Completed 
- Login  
- Registration  
- Getting users by username or uuid
- Banning users
- Role based authentication
- JWT authentication 
## TODO
- Captcha 
- Multi tenancy support
- Docker file (in progress)  
- Docker compose
## Setup

Download the repository and run `go mod tidy`

Create a postgres server and create the tables in auth.sql

Place the database URL and jwt secret inside a .env in the project directory

```
DATAGASE_URL=postgresql://postgres:password@ipaddr:5432/db
JWT_SECRET=token
```

run `go run .` or `go build -o auth` and `./auth` for a binary 
