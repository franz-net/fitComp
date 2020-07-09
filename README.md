<h1 align="center">
  fitComp
</h1>

I created this API project as I was learning go and to help my family keep track of a weight loss competition.
It was a great way of learning some basics, I hope it helps other people learn as well and get started with their own projects

## Libraries Used
- Gin and Gonic
- xid
- bcrypt
- MySql driver

## External components
- MySQL
- Docker

## Work in Progress
- Adding user management

## Docker Usage Instructions
* docker run --name mysql -e MYSQL_ROOT_PASSWORD=<root_pass> -p 3306:3306 -d mysql:latest
* On subsequent runs do docker start mysql
* docker exec -it mysql /bin/bash
* mysql -h localhost -u root -p
* USE fitComp or CREATE DATABASE fitComp;


## Things I would like to add
* Use SSL for the webserver
* Add tests
* Add a UI

## Installation (Building the code)

To build and run the program you need: 
* [Golang installed](https://golang.org/doc/install) and [GOPATH configured](https://golang.org/doc/gopath_code.html)
* Clone this repository
* Build the code with `go build...`

 Example of build for linux x86_64: `env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o fitComp main.go`