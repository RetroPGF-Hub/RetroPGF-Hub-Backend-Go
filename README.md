## Export path
export PATH="$PATH:$(go env GOPATH)/bin"

openssl genrsa -out cert/key_rsa 1024
openssl rsa -in cert/key_rsa -pubout -out cert/key_rsa.pub

## GRPC
### Export path
export PATH="$PATH:$(go env GOPATH)/bin"
```
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    modules/users/usersPb/usersPb.proto
```

## db diagram

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
  id              String   [primary key]
  name            String
  logo_url        String
  banner_url      String
  website_url     String
  crypto_category String
  description     String   
  reason          String   
  category        String
  contact         String
  // objectid
  create_by       String
  created_at      timestamp
  updated_at      timestamp
}

Table comment {
  id         String
  title      String
  content    String
  // objectid
  user_id    String
  project_id String

  created_at      timestamp
  updated_at      timestamp
}

Table favourite {
  project_id String
  // users is a array of user
  users String[]
  user_id String
  created_at      timestamp
}



Ref: comment.user_id > users.id
Ref: comment.project_id > projects.id
Ref: favourite.user_id > users.id
Ref: favourite.project_id > projects.id

```

