# go-graphql

This repository has a proof of concept about graphql.

I will use a video publishing site as an example in which a user can publish a video, add users, add screenshots and get videos and users.

## Getting started

1. Clone repository: `git clone https://github.com/santileira/go-graphql.git`
2. Run `dep ensure`
3. Run `go run server/server.go`
4. Open `http://localhost:8080/` for GraphQL Playground

## Definitions

- Video schema has an user but i only fetch the user if it is requested.
- Scalar type `DateTime` with format `2006-01-02 15:04:05`.

## Users

### Query

Gets users in database.

`URL:`
`http://localhost:8080/query`

`BODY:`
```
query{
  Users{
    id
    name
    email
  }
}
```

`RESPONSE:`
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

`URL:`
`http://localhost:8080/query`

`BODY:`
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

`RESPONSE:`
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

`URL:`
`http://localhost:8080/query`

`BODY:`
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
  }
}
```

`RESPONSE:`
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
        "createdAt": 1564489100000
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

`URL:`
`http://localhost:8080/query`

`BODY:`
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
  }
}
```

`HEADERS:`
```
{
  "Role": "ADMIN"
}
```


`RESPONSE:`
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
      "createdAt": "2019-07-31 19:42:34"
    }
  }
}
```

If user id not exists this is the response:

```
{
  "errors": [
    {
      "message": "user not exists",
      "path": [
        "createVideo"
      ]
    }
  ],
  "data": null
}
```

### Subscription

I subscribe to videos creation, when video is created i receive it.
The library (gqlgen) provides web socket-based real-time subscription events.

`URL`
`http://localhost:8080/query`

`BODY:`
```
subscription{
  videoCreated{
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
  }
}
```

`REPSONSE:`
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
      "createdAt": "2019-08-02 20:59:10"
    }
  }
}
```

## Authentication

In GraphQL only one endpoint is exposed, in this case `http://localhost:8080/query`, so you can achieve authorization with schema directives.

`createVideo` mutation has directive `hasRole(ADMIN, USER)`, so only roles admin and user can create videos, else returns not authorized error.

To execute `createVideo` you must pass header `"Role": "ADMIN"` or `"Role": "ADMIN"`. This verification is very easy and silly, but you can do a very complexity authentication.

`URL:`
`http://localhost:8080/query`

`BODY:`
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
  }
}
```

`HEADERS:`
```
{
  "Role": "TEST"
}
```


`RESPONSE:`
```
{
  "errors": [
    {
      "message": "you aren't authorized to perform this action",
      "path": [
        "createVideo"
      ]
    }
  ],
  "data": null
}
```

## Data Loader

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


