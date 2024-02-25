![finding-prices](assets/banner.gif)

# Finding prices

This is a project for finding prices from websites. 

I wanted to buy a macbook and I need to understand how often the price changes.

## Getting started

You need: 

1. download go, current version 1.18.5 https://go.dev/dl/
2. download postgresql https://www.postgresql.org/download/
3. create local database and insert tables, data from migration.sql
4. create config.yml at the root of the project and transfer structure from config.template.uml
5. run the file ./cmd/main/maing.go with command ```go run main.go```


## What's contained in this project

- main.go - is the main definition of the service
- internal/customers - contains search function for each customers
- internal/link - contains functions for finding prices for items
- internal/prices - functions for creating new prices
- pkg/postgresql/postgresql - create connection to db

## Dependencies

Install the following

- [goquery](https://github.com/PuerkitoBio/goquery)
- [cleanenv](https://github.com/ilyakaznacheev/cleanenv)
- [pgx4](https://github.com/jackc/pgx/v4)