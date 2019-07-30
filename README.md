# go-graphql

This repository has a proof of concept about graphql.

I will use a video publishing site as an example in which a user can publish a video, add users, add screenshots and get videos and users.

## Getting started

1. Clone repository: `git clone https://github.com/santileira/go-graphql.git`
2. Run `dep ensure`
3. Run `go run server/server.go`
4. Open `http://localhost:8080/` for GraphQL Playground

## Definitions

Video schema has an user but i only fetch the user if it is requested.

## Users

### Query

Gets users in database.

`http://localhost:8080/query`

```
query{
  Users{
    id
    name
    email
  }
}
```

```
{
  "data": {
    "Users": [
      {
        "id": "5577006791947779410",
        "name": "Claudio Leira_5577006791947779410",
        "email": "claudio@gmail.com_5577006791947779410"
      }
    ]
  }
}
```

### Mutation

Inserts user in database with random id.

Parameters:
- name (string): required
- email (string): required

`http://localhost:8080/query`

```
mutation{
     createUser(user:{
       name: "Claudio Leira",
       email: "claudio@gmail.com"
     }) {
       id
       name
       email
     }
   }
```

```
{
  "data": {
    "createUser": {
      "id": "5577006791947779410",
      "name": "Claudio Leira_5577006791947779410",
      "email": "claudio@gmail.com_5577006791947779410"
    }
  }
}
```

## Videos

### Query

Gets videos in database.

`http://localhost:8080/query`

```
query{
  Videos{
    id
    name
    description
    user {
      id
      name
      email
    }
    url
    createdAt
    related{
      id
    }
  }
}
```

```
{
  "data": {
    "Videos": [
      {
        "id": "8674665223082153551",
        "name": "Video Claudio Leira",
        "description": "Description Video Claudio Leira",
        "user": {
          "id": "5577006791947779410",
          "name": "Claudio Leira_5577006791947779410",
          "email": "claudio@gmail.com_5577006791947779410"
        },
        "url": "Url Video Claudio Leira",
        "createdAt": 1564489100000,
        "related": []
      }
    ]
  }
}
```

### Mutation

Inserts video in database with random id. If user id isn't exist so returns error.

Parameters:
- name (string): required
- email (string): required

`http://localhost:8080/query`

```
mutation{
  createVideo(input:{
    name: "Video Claudio Leira",
    description: "Description Video Claudio Leira"
    userId: 5577006791947779410
		url:"Url Video Claudio Leira"
  }) {
    id
    name
    description
    user {
      id
      name
      email
    }
    url
    createdAt
    related {
      id
    }
  }
}
```

```
{
  "data": {
    "createVideo": {
      "id": "8674665223082153551",
      "name": "Video Claudio Leira",
      "description": "Description Video Claudio Leira",
      "user": {
        "id": "5577006791947779410",
        "name": "Claudio Leira_5577006791947779410",
        "email": "claudio@gmail.com_5577006791947779410"
      },
      "url": "Url Video Claudio Leira",
      "createdAt": 1564489100000,
      "related": []
    }
  }
}
```

### Subscription

## Library used

gqlgen (https://github.com/99designs/gqlgen)

## Other libraries

Thunder, graphql and graphql-go

## Steps

dep init

go run scripts/gqlgen.go init (generates all the boring (but interesting for a few) skeleton code)

which will create the following files:

gqlgen.yml — Config file to control code generation.
generated.go — The generated code which you might not want to see.
models_gen.go — All the models for input and type of your provided schema.
resolver.go — You need to write your implementations.
server/server.go — entry point with an http.Handler to start the GraphQL server.

generates new models and marshal and unmarshal

modify gqlgen.yml

Regenerate the code by running: go run scripts/gqlgen.go -v.


generate dataloader: `go run github.com/vektah/dataloaden UserLoader string github.com/santileira/go-graphql/api/models.User`


