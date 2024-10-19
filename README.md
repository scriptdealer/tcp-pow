# Faraway test
Proof-of-Work over TCP

## Algorithm choice
For this implementation, we'll use the Hashcash algorithm. It's a simple yet effective PoW system that requires the client to find a partial hash collision. This algorithm is suitable because:
* It's computationally difficult but easy to verify
* It's adjustable in difficulty
* It's well-known and widely used

## Local server
`docker build -t tcp-pow-server -f test/server/Dockerfile .`   
then    
`docker run -p 8080:8080 tcp-pow-server`

## Local client
`docker build -t tcp-pow-client -f test/client/Dockerfile .`
then
`docker run -it tcp-pow-client`

## Docker network
To create a common network for the server and client containers:

```bash
docker network create tcp-pow-network
```