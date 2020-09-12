# Donate API

Central API for a larger **Donate** project aimed at allowing people to donae compute time from their home computers or laptops to a chosen computational task. These 'jobs' will be specified by researchers and will allow them to perform high quality computation on a budget.

Works on a distributed system built around transferring docker images.

## Explanation of Common Concepts in final concept

A singular job is a datum with:

* unique id
* pointer to original dockerhub image (containing the job)
* pointer to completed dockerhub image
* list of pointers to partially completed images, along with total work
* title
* description
* creation timestamp
* mark for completion
* mark for allowing multiple runners
* author
* runner list

A singlar user is a datum with:

* username
* hashed password
* list of authored jobs
* list of currently running jobs

## Paths Supported

### Users

* [x] `POST /register`
* [x] `GET /user`
* [x] `PUT /user`
* [x] `DELETE /user`

### Jobs

* [x] `GET /job`
* [x] `POST /job`
* [x] `GET /job/{id}`
* [x] `DELETE /job/{id}`
* [ ] `PUT /job/{id}/take`
* [ ] `PUT /job/{id}/return`

## Building

`go build -o main .`

## Dockerization

`docker build --build-arg salt=<RANDOM SALT HERE> -t donate-api .`

`docker run -it -p 8080:8080 donate-api`
