# ETL server

The ETL server is built to facilitate API development based on SQL
queries, and also allows Web development by defining a Vuego template in
the definition.

## Simple API

![](diagrams/api.svg?v1)

A simple request will run one or more queries against the database.

## Web development

![](diagrams/webdev-db.svg?v1)

The data returned from the database can be used for front-end rendering.

## Web development with external source

![](diagrams/webdev-http.svg?v1)

ETL would allow a distributed service configuration. Service discovery
could look up on which servers the API is accessible, and use a HTTP
request to fetch the JSON. This allows to consider network topology.
