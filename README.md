# Donate API

![Go Tests](https://github.com/DonateComputing/DonateAPI/workflows/Go%20Tests/badge.svg?branch=master)
![Docker Image CI](https://github.com/mfigurski80/DonateAPI/workflows/Docker%20Image%20CI/badge.svg)

Central API for a larger **Donate** project aimed at allowing people to donae compute time from their home computers or laptops to a chosen computational task. These 'jobs' will be specified by researchers and will allow them to perform high quality computation on a budget.

Works on a distributed system built around transferring docker images.

## Paths Supported

### Users

* [x] `POST /register` <- {username, password}
	- creates a new user entry
* [x] `GET /{userId}`
	- gets user if matches auth
	- requires auth
* [x] `PUT /{userId}` <- {password}
	- updates password of user
	- requires auth 
* [x] `DELETE /{userId}`
	- removes user entry
	- requires auth

### Jobs

* [x] `GET /job`
	- gets list of jobs without runner
* [x] `POST /job` <- {title, description, image}
	- creates new job entry
	- requires auth
* [x] `GET /{userId}/{jobId}`
	- gets job entry by given user and id
* [x] `DELETE /{userId}/{jobId}`
	- removes job entry by given user and id
	- requires auth
* [x] `PUT /{userId}/{jobId}/take` <- {}
	- marks self as runner of job entry by given user and id
	- requires auth
* [x] `PUT /{userId}/{jobId}/return` <- {image}
	- removes runner reference in job entry
	- sets completed image in job entry if image ref is given
	- requires auth

## Testing

`go test ./...`

## Building

`go build -o main .`

## Dockerization

`docker build -t donate-api -t mfigurski80/donate-api .`

`docker-compose up`
