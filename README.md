# ecom-go

Go version - 1.22.5

Please run the following commands to run the project successfully

in each service - run the following commands
go mod tidy
go mod vendor

db details are in the .env file of each service

please check the env files to set-up database (MySQL) with correct names.

after setting up the databases
cd cmd (in each service) -> go run main.go

each service will be running in different ports.
 
 
