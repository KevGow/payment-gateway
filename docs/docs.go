// Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/findpayment/{uuid}": {
            "get": {
                "description": "Get payment information by UUID",
                "produces": [
                    "application/json"
                ],
                "summary": "Get payment information by UUID",
                "operationId": "get-payment-by-uuid",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Payment UUID",
                        "name": "uuid",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/api.GetResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/api.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/api.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/pay": {
            "post": {
                "description": "Make a payment",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Make a payment",
                "operationId": "make-payment",
                "parameters": [
                    {
                        "description": "Payment Data",
                        "name": "paymentData",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/api.PostJsonRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/api.PostResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/api.ErrorResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "api.ErrorResponse": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                }
            }
        },
        "api.GetResponse": {
            "type": "object",
            "properties": {
                "amount": {
                    "type": "number",
                    "example": 100
                },
                "bank-payment-status": {
                    "type": "string",
                    "example": "Success"
                },
                "card-number-masked": {
                    "type": "string",
                    "example": "****5070"
                },
                "currency": {
                    "type": "string",
                    "example": "GBP"
                },
                "expiry-date": {
                    "type": "string",
                    "example": "11/26"
                }
            }
        },
        "api.PostJsonRequest": {
            "type": "object",
            "required": [
                "amount",
                "card-number",
                "currency",
                "cvv",
                "expiry-date"
            ],
            "properties": {
                "amount": {
                    "type": "number",
                    "example": 100
                },
                "card-number": {
                    "type": "string",
                    "example": "4032 0341 3083 5070"
                },
                "currency": {
                    "type": "string",
                    "example": "GBP"
                },
                "cvv": {
                    "type": "string",
                    "example": "975"
                },
                "expiry-date": {
                    "type": "string",
                    "example": "11/26"
                }
            }
        },
        "api.PostResponse": {
            "type": "object",
            "properties": {
                "uuid": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "Payment Gateway API",
	Description:      "This is a simple Payment Gateway API for the ProcessOut take-home technical assessment.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
