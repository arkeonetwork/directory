# Arkeo Directory Service

The directory is an off-chain API/service (utilizing openapi/swagger) that makes it easier for users and client to discover data providers.


## go-swagger

Currently, this API uses go-swagger to generate the API spec and can be used to also generate a client in the future. For now this can be installed using brew or by following alternate 
instructions [here](https://goswagger.io/install.html). 

```
brew tap go-swagger/go-swagger
brew install go-swagger
```

Once installed you can use the following make commands to either generate a swagger.yaml or serve the .yaml.

```
make swagger
make swagger-serve
```

An important note for VS code users, leave a single blank linke between your code annotations and the function definition.  This will enable the go formatter to not re-format the needed `+` into a `-`

## Database
- Postgres 15
### Migrations
[Tern](https://github.com/jackc/tern) is used to manage database structure and initialize state. For now
tern is using the default configuration mechanism, see db/tern.conf.

With postgres running and pointed to by tern.conf:
```
make db-migrate
```

Create a new migration:
```
tern new -m db useful_name

ls db/007_useful_name.sql
db/007_useful_name.sql
```

Undo the most recently applied migration:
```
tern migrate -c db/tern.conf -m db --destination -1
```
