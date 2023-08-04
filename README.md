# Payment Gateway

The Payment Gateway is a solution to the ProcessOut take-home technical assesment. This solution is a RESTful
HTTP server, built using Gin framework for the server, and Swaggo for generating an API spec and means of testing via the Swagger UI.

## Technologies Used

The server uses the following technologies:

-   [Gin](https://github.com/gin-gonic/gin) - a web framework for Go
-   [Swaggo](https://github.com/swaggo/swag) - an auto-generated API documentation tool for Go

## Installation

Prerequisites
To run this server, you need to have [Go](https://go.dev/doc/install) 1.20 or later installed on your machine.

1. Clone this repository using the following command:

`https://github.com/KevGow/payment-gateway.git`

2.  Navigate to the project directory:

`cd PATHTODIRECTORY/payment-gateway`

3. Install the project:

`go install payment-gateway`

4. Run the server:

`payment-gateway`

## API Documentation

There are two endpoints for this server:

#### POST /pay

and

#### GET /findpayment/{uuid}

As server is documented using Swaggo, you can view the full API specs by viewing the Swagger documentation. To view the Swagger documentation, open a browser and navigate to [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html) once the server is built and running.

## Running Tests

To run the tests, use the following command:

`go test ./...` 

The tests are located in `main_test.go`.

These tests test at the API level, with no other testing. This allows for testing all components from the moment a request is received, to the reponse generated. 


## Further Improvements

While this solution aims to be through, it also aims to be simple and scalable. There are many improvements that could be made, but here are some of the main ones that have been identified, in no particular order.

#### Additional Validation: 
While there are already basic validation checks in place such as Luhn's Algorithm for the card details, the world of payment validation and fraud checking is vast. Due to the way the solution is built, more checks would be trivial to emplement going forward.

#### Data Persistence: 
Currently, the payments are stored in memory using a map. This creates a few problems, one being data persistance, as data will only live as long as the server does. Another issue is that in memory data storage is not scalable. However, due to the way data is accessed via methods, the implementation of these methods easily be changed to either push or fetch data from a database, another RESTful server, gRPC, etc in production.

#### Input Sanitization: 
Go Gin performs minimal to no input santisation. Input santisation would be a nessesity in production to ensure that the data received from clients is safe, and does not lead to security vulnerabilities such as SQL injection or cross-site scripting.

#### Authentication and Authorization: 
There is no authenticaion or authorisation in this solution. Authenticaion and authorisation mechanisms would be needed to secure the API in production. 

#### Configuration Management: 
Move hard-coded configuration (e.g. server port) to a configuration file or environment variables for easier deployment and management.

#### Improved Logging and Monitoring: 
Implement more verbose and useful logging as well as mechanisms for monitoring the status and performance of the server.

#### Graceful Shutdown: 
Gin does not handle graceful shutdown out of the box. Implementmenting a graceful shutdown mechanism to handle server shutdowns properly would be beneficial.
