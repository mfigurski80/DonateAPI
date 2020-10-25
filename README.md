# Donate API

![Go Tests](https://github.com/DonateComputing/DonateAPI/workflows/Go%20Tests/badge.svg?branch=master)
![Docker Image CI](https://github.com/mfigurski80/DonateAPI/workflows/Docker%20Image%20CI/badge.svg)

Central API for a larger **Donate** project aimed at allowing people to donae compute time from their home computers or laptops to a chosen computational task. These 'jobs' will be specified by researchers and will allow them to perform high quality computation on a budget.

Works on a distributed system built around transferring docker images.

## Paths Supported

### Users

* [x] `POST /register`
* [x] `GET /{userId}`
* [x] `PUT /{userId}`
* [x] `DELETE /{userId}`

### Jobs

* [x] `GET /job`
* [x] `POST /job`
* [x] `GET /{userId}/{jobId}`
* [x] `DELETE /{userId}/{jobId}`
* [x] `PUT /{userId}/{jobId}/take`
* [x] `PUT /{userId}/{jobId}/return`

## Testing

`go test ./...`

## Building

`go build -o main .`

## Dockerization

`docker build -t donate-api -t mfigurski80/donate-api .`

`docker-compose up`
