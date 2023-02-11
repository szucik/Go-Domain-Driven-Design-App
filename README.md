# Simple implementation of Domain-Driven Design

### Installation 
More info: [https://go.dev/doc/manage-install](https://go.dev/doc/manage-install)

## Architecture
This application is based on Domain-Driven Design (DDD). 

Terms used:
* Entities - a struct that has a unique identifier and is mutable.
* Aggregate - a collection of combined entities and value objects. It has a life cycle and is stored in the DB.
* Repository - is used to store and manage aggregates.
* Service - the combination of business logic and repositories together.

### Routing
Application implements Gorilla mux router.</br>
More info: [https://github.com/gorilla/mux]( https://github.com/gorilla/mux)