# go-graphql

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



