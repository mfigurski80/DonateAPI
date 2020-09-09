# Donate API

Central API for a larger **Donate** project aimed at allowing people to donae compute time from their home computers or laptops to a chosen computational task. These 'jobs' will be specified by researchers and will allow them to perform high quality computation on a budget.

Works on a distributed system built around transferring docker images.

## Building

`go build -o main .`

## Dockerization

`docker build --build-arg salt=<RANDOM SALT HERE> -t donate-api`

`docker run -it -p 8080:8080 donate-api`
