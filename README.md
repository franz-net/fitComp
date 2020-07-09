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
* Run docker container for MySQL
  > ```docker run --name mysql -e MYSQL_ROOT_PASSWORD=<root_pass> -p 3306:3306 -d mysql:latest```
* Access the container to create the database
  > ```docker exec -it mysql /bin/bash```
* Connect to the database server
  > ```mysql -h localhost -u root -p```
* Create the database
  > ```USE fitComp or CREATE DATABASE fitComp;```

* If the machine is restarted or the container stopped, run the command below before starting the application
  > ```docker start mysql```


## Things I would like to add
* Use SSL for the requests
* Standarize the API schema
* Add tests
* Add a UI

## Installation (Building the code)

To build and run the program you need: 
* [Golang installed](https://golang.org/doc/install) and [GOPATH configured](https://golang.org/doc/gopath_code.html)
* Clone this repository
* Build the code with `go build .` or run with `go run .`
