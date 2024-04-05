
# ลบ image ถ้า create project faild

## Documentation for this project
https://documenter.getpostman.com/view/17633056/2s9Yyth1cP

## Flow 
### Project Service
![alt text](./assets/projectSvc.png)
### User and Datacenter Service
![alt text](./assets/user&dcSvc.png)


## Run Service
### Service Users
```
make users
```

### Service Project
```
make project
```

### Service Datacenter
```
make datacenter
```

## Start Database
### This command will run Mongo and Redis container
```
make db-up
```
## Stop Database
```
make db-stop
```

## Grpc
### Export path
```
export PATH="$PATH:$(go env GOPATH)/bin"
```

## Migrate Grpc
### Migrate Grpc Users
```
make grpc-users
```

### Migrate Grpc Datacenter
```
make grpc-datacenter
```

### Migrate Grpc Favorite
```
make grpc-fav
```


## Migrate Database (Create Index)
### Migrate Project
```
make migrate-project
```

### Migrate Users
```
make migrate-users
```


## Random api
https://api.thecatapi.com/v1/images/search?limit=10

## Exec Redis cli
docker exec -it containerId redis-cli


## Generate Secret
```
openssl genrsa -out cert/key_rsa 1024
openssl rsa -in cert/key_rsa -pubout -out cert/key_rsa.pub
```


## Database Model

```
Table users {
  id         String    [primary key]
  email      String    
  password   String
  profile    String
  user_name  String
  first_name String
  last_name  String
  created_at      timestamp
  updated_at      timestamp
}

Table projects {
  // object id
  _id              String   [primary key]
  name            String
  logo_url        String
  banner_url      String
  website_url     String
  crypto_category String
  description     String   
  reason          String   
  category        String
  contact         String
  fav_count       Number
  comment_count   Number
  create_by       String
  created_at      timestamp
  updated_at      timestamp
}

Table comment {
  // _id is equal to projectId
  _id         String
  comments {
    _id             String
    title           String
    content         String
    created_by      String
    created_at      timestamp
    updated_at      timestamp
  }[]
  created_at      timestamp
  updated_at      timestamp
}

Table favorite {
  // _id is equal of userId
  _id  String
  project_id String[]
  created_at      timestamp
}

```

