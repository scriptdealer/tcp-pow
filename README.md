# Denial of denial
Faraway test to implement protection of a TCP server from DDOS attacks with the Proof of Work pattern

# Overview
For this implementation, we'll use the Hashcash algorithm. It's a simple yet effective PoW system that requires the client to find a partial hash collision. This algorithm is suitable because:
* It's computationally difficult but easy to verify
* It's adjustable in difficulty
* It's well-known and widely used

Base idea is to serve maximum number of users by denying only certain IP addresses. This is done by automatically adjusting the difficulty level based on the number of requests made from a particular IP address.

## Docker network
To create a common network for the server and client containers:

```bash
docker network create tcp-pow-network
```

## Server
```bash
docker build -t tcp-pow-server -f test/server/Dockerfile .
```   
then    
```bash
docker run -p 8080:8080 tcp-pow-server
```

## Client
```bash
docker build -t tcp-pow-client -f test/client/Dockerfile .
```
then
```bash
docker run -it tcp-pow-client
```
