# hello-api

A simple hello world service written in go. Intent is to build docker image and use it while learning stuff in k8s.


## Docker Image

```
docker build -t hello-api .
docker tag hello-api:latest <docker location>
docker push <docker location>
```
