# ecom-go

## Go version - 1.22.5

## Please set up database before running the application - the db details are in the .env files

db creation -> (Please check the values in .env files for each service)

CREATE DATABASE "name of the db in .env file of each service";
CREATE USER "username"@"host" IDENTIFIED BY 'password';
GRANT ALL PRIVILEGES ON *.* TO "username"@"localhost" WITH GRANT OPTION;
FLUSH PRIVILEGES;

## Dependencies installation

in each service -> Ex -> go to the particualr service directory -> go to cart-service directory
run the following commands to install dependencies

### go mod tidy
### go mod download
### go mod vendor

## How to run the code (After setting up db)

in each service ->  go into `cmd` directory -> `main.go` -> run the following command

### go run main.go (run for all the services)
-> each service will run in a different port.

