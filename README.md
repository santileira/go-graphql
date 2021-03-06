# go-graphql

This repository has a proof of concept about graphql.

I will use a video publishing site as an example in which a user can publish a video, add users and get videos and users.

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

## Data loaders

I am loading data when needed. Clients have the control of the data, there is no underfetching and no overfetching but everything comes with a cost.

I will show an example to understand the data cost:

```
query{
  Videos{
      name
      user{
        name
      }
  }
}
```

If we have 10 videos entries and only 5 differents users. We will do 11 queries, 1 to videos table and 10 to users table.

This is known as the `N+1` problem. There will be one query to get all the data and for each data (N) there will be another database query.

Dataloaders will solve this problem. I will use dataloaden (https://github.com/vektah/dataloaden) library from the author of gqlgen.

I will generate the user data loader with the next command: `go run github.com/vektah/dataloaden UserLoader int github.com/santileira/go-graphql/api/models.User`

This will generate a dataloader called `UserLoader` that looks for `github.com/santileira/go-graphql/api/models.User` objects based on an `int` key.

I need to define the `Fetch` method to get the result in bulk. I am waiting for 1ms for a user to load queries and i have kept a maximum batch of 10 queries. So now, instead of firing a query for each user, dataloader will wait for either 1 millisecond for 10 users before hitting the database.

## Query complexity

In GraphQL you are giving a powerful way for the client to fetch whatever they need, but this exposes you to the risk of denial of service attacks.

For example, if i have a related field in video type which returns related videos. And each related video is of the graphql video type so they all have related videos too… and this goes on.

Library gqlgen assigns fix complexity weight for each field so it will consider struct, array, and string all as equals.

For this example, complexity limit is 100.

## Library used

gqlgen (https://github.com/99designs/gqlgen)

## Other libraries

Thunder, graphql and graphql-go

## Bibliography

- https://www.freecodecamp.org/news/deep-dive-into-graphql-with-golang-d3e02a429ac3/
- https://gqlgen.com/
- https://github.com/99designs/gqlgen
- https://github.com/vektah/dataloaden

