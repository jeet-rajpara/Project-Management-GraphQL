
# Project Management System

Project Management system for create/update/delete/share Projects and related Screenshot.

# Features
- There are four different roles Owner(creator of project),Admin,Editor and Viewer.
- Owner can create,update and delete projects
- Owner can create and delete project Screenshots.
- User can fetch all projects(based on filter if provided) and single project 
- Owner can share project with another user with roles like Admin,Editor and Viewer
- User can get real time update when he becomes part of any project with specific role

# Tech Stack 
- GO 1.21
- CockroachDB 23.1
- JWT (json web token)

## Run Locally

Prerequisites you need to set up on your local computer:

- [Golang](https://go.dev/doc/install)
- [Cockroach](https://www.cockroachlabs.com/docs/releases/)
- [Dbmate](https://github.com/amacneil/dbmate#installation)

1. Clone the project

```bash
  git clone https://github.com/zuru2024/Project-Management-GraphQL.git
  cd Project-Management-GraphQL
```

2. Copy the .env.example file to new .env file inside root directory and set env variables in .env:

```bash
  cd . > .env
  cp .env.example .env
```

3. Run `dbmate up` to create database schema or Run `dbmate migrate` to migrate database schema.
4. Run `go run main.go` to run the programme.

## API Documentation:

After executing run command, open your favorite browser and type below URL to open graphQL playground.
```
http://localhost:8010/
```