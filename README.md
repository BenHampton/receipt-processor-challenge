# receipt-processor-challenge

Steps:
1. build docker image: `docker build -t go-docker-demo .`
2. run docker container: `docker run -d -p 8080:8080 go-docker-demo`
3. using an api testing tool like postman call the `POST` endpoint and use the `ID` from that response in the  `GET` endpoint. 

Postman collection can be found at the root level of the project or reference the Endpoints section of the readme

## Docker:
### build image:
`docker build -t go-receipt-processor-challenge .`

### Run container:
`docker run -d -p 8080:8080 go-receipt-processor-challenge`


## Endpoints
#### 2 endpoints running on port 8080:

1: Method: `POST` <br>
Url: `http://localhost:8080/receipts/process`

2: Method: `GET` <br>
Url: `http://localhost:8080/receipts/{id}/points`
