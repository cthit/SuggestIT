# SuggestIT

An application for sending suggestions about Hubben 2.1 to P.R.I.T..

## Setup development
Run the following commands to set up the development environment.  

```
    git clone git@github.com:cthit/SuggestIT.git
    cd SuggestIT/frontend
    yarn install
    docker-compose up --build
```
The frontend is exposed to port 3000 and uses live reload.  
The backend is exposed to port 5000 but will not reload when the code change.  
The database can be administrated from port 5432.

## Run prod docker
Use the ```prod.docker-compose.yml``` and change the passwords and the ```PRIT_AUTH_KEY```.

```
version: "2"
services:
  db:
    image: postgres
    restart: always
    environment:
      POSTGRES_USER: suggestit
      POSTGRES_DB: suggestit
      POSTGRES_PASSWORD: <your db password>
    ports:
      - 5432:5432

  backend:
    image: cthit/suggestit-backend
    environment: 
      SUGGESTIT_POSTGRES_USER: suggestit
      SUGGESTIT_POSTGRES_PASSWORD: <your db password>
      SUGGESTIT_POSTGRES_HOST: db
      SUGGESTIT_POSTGRES_PORT: 5432
      SUGGESTIT_POSTGRES_DB: suggestit
      
      AUTH_SECRET: <random long string>
    ports:
      - 5000:5000

  frontend:
    image: cthit/suggestit-frontend
    ports:
      - 8080:8080
    volumes:
      - ./frontend:/usr/src/app/
```

## Backend documentation
The project is using [Swagger](http://swagger.io/) to document the restful api When the backend is running, the documentation. Copy the backend/static/swagger.yml to [swagger editor](http://editor.swagger.io/) or use the Swagger Viewer extension in VS Code to view it.