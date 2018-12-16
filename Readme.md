
```bash

========================================================================================  
                           888                        888      
                           888                        888      
                           888                        888      
.d8888b   8888b.  88888b.  88888b.   .d88b.   .d88b.  888  888 
88K          "88b 888 "88b 888 "88b d88""88b d88""88b 888 .88P 
"Y8888b. .d888888 888  888 888  888 888  888 888  888 888888K  
     X88 888  888 888  888 888  888 Y88..88P Y88..88P 888 "88b 
 88888P' "Y888888 888  888 888  888  "Y88P"   "Y88P"  888  888
========================================================================================
- port    : 8080
- log     : logs
--------------------------------------------------------
```

#### Structures

```bash
.
+-- api
   ├── swagger-iam.json
   ├── swagger-activity-log.json
+-- configs
   ├── app.dev.yaml
   ├── app.yaml.dist
+-- cmd
   ├── app
      ├── main.go
+-- docs
+-- deployments
   ├── Dockerfile
   ├── docker-compose.yaml
+-- infrastructures
   ├── mariadb.go
   ├── postgres.go
   ├── sentry.go
+-- internal
   ├── pkg
   +-- http
      ├── routes.go
      ├── user.go
   +-- user
      ├── model.go
      +-- schema
         ├── mariadb.go
         ├── sql.go
      ├── repository.go
      ├── response.go
      ├── request.go
      ├── service.go
+-- pkg
   ├── utils
├── .dockerignore
├── .gitignore
├── Makefile``
├── Readme.md
```


## Pre Test Backend Simple API

Requirement: 
- Go 1.11+

### Setup

```bash
go mod tidy

```

### Running Service

* Http Service

```bash
go run cmd/sanhook/main.go http

```

* Nats Pub-Sub Service

```bash
go run cmd/sanhook/main.go nats

```

#### Create a simple API


##### API for sending a message 

> GET: /v1/api/message/publish/{subject:[a-z-]+}


 ```bash
 curl -i http://localhost:8080/v1/api/message/publish/nats-topic?message=makan+daging+kambing+pakai+sayur+lodeh
 
 ```
 
##### API for collect message that has been sent out

> GET: /v1/api/message/inbox/{id:[a-z-0-9]+}

```bash
curl -i http://localhost:8080/v1/api/message/inbox/824bed3e-5073-4ebe-b3a0-833aa1247372
```
 
- Get All Inbox 

> GET: /v1/api/message/inbox

```bash
curl -i http://localhost:8080/v1/api/message/inbox
```

##### API for display message in real time

- For Subscribe Service Real Time


> Service: -s [subject-for-subscribe]

```bash
go run cmd/sanhook/main.go subscribe -s nats-topic 
```
