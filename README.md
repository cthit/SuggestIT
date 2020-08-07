# SuggestIT

An application for sending suggestions about Hubben 2.1 to P.R.I.T..

## Setup development

Run the following commands to set up the development environment.

```
    git clone git@github.com:cthit/SuggestIT.git
    cd SuggestIT
    docker-compose up
```

The frontend is exposed to port 3000 and uses live reload.  
The backend is exposed to port 5000 but will not reload when the code change.  
The database can be administrated from `http://localhost:8080`, use the flowing credentials:

```
  System: PostgreSQL
  Server: db
  Username: suggestit
  Password: password
  Database: suggestit
```

## Backend documentation

The project is using [Swagger](http://swagger.io/) to document the restful api When the backend is running, the documentation. Copy the backend/static/swagger.yml to [swagger editor](http://editor.swagger.io/) or use the Swagger Viewer extension in VS Code to view it.
