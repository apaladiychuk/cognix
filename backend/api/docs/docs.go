// Package docs GENERATED BY SWAG; DO NOT EDIT
// This file was generated by swaggo/swag
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "API Support"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/auth/google/invite": {
            "get": {
                "description": "join user to tenant using invitation link",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Auth"
                ],
                "summary": "join user to tenant using invitation link",
                "operationId": "auth_join_to_team",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            },
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "create invitation for user",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Auth"
                ],
                "summary": "create invitation for user",
                "operationId": "auth_invitation",
                "parameters": [
                    {
                        "description": "invitation  parameter",
                        "name": "params",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/parameters.InviteParam"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/auth/google/login": {
            "get": {
                "description": "register new user and tenant using google auth",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Auth"
                ],
                "summary": "register new user and tenant using google auth",
                "operationId": "auth_login",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/chats/create-chat-session": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "creates new chat session",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Chat"
                ],
                "summary": "creates new chat session",
                "operationId": "chat_create_session",
                "parameters": [
                    {
                        "description": "create session parameters",
                        "name": "payload",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/parameters.CreateChatSession"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.ChatSession"
                        }
                    }
                }
            }
        },
        "/chats/get-chat-session/{id}": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "return chat session with messages by given id",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Chat"
                ],
                "summary": "return chat session with messages by given id",
                "operationId": "chat_get_by_id",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "session id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.ChatSession"
                        }
                    }
                }
            }
        },
        "/chats/get-user-chat-sessions": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "return list of chat sessions for current user",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Chat"
                ],
                "summary": "return list of chat sessions for current user",
                "operationId": "chat_get_sessions",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/model.ChatSession"
                            }
                        }
                    }
                }
            }
        },
        "/manage/connector": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "return list of allowed connectors",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Connectors"
                ],
                "summary": "return list of allowed connectors",
                "operationId": "connectors_get_all",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/model.Connector"
                            }
                        }
                    }
                }
            }
        },
        "/manage/connector/": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "creates connector",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Connectors"
                ],
                "summary": "creates connector",
                "operationId": "connectors_create",
                "parameters": [
                    {
                        "description": "connector create parameter",
                        "name": "params",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/parameters.CreateConnectorParam"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/model.Connector"
                        }
                    }
                }
            }
        },
        "/manage/connector/{id}": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "return list of allowed connectors",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Connectors"
                ],
                "summary": "return list of allowed connectors",
                "operationId": "connectors_get_by_id",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.Connector"
                        }
                    }
                }
            },
            "put": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "updates connector",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Connectors"
                ],
                "summary": "updates connector",
                "operationId": "connectors_update",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "connector id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "connector update parameter",
                        "name": "params",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/parameters.UpdateConnectorParam"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.Connector"
                        }
                    }
                }
            }
        },
        "/manage/credential": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "return list of allowed credentials",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Credentials"
                ],
                "summary": "return list of allowed credentials",
                "operationId": "credentials_get_all",
                "parameters": [
                    {
                        "type": "string",
                        "description": "source of credentials",
                        "name": "source",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/model.Credential"
                            }
                        }
                    }
                }
            },
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "creates new credential",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Credentials"
                ],
                "summary": "creates new credential",
                "operationId": "credentials_create",
                "parameters": [
                    {
                        "description": "credential create parameter",
                        "name": "params",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/parameters.CreateCredentialParam"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/model.Credential"
                        }
                    }
                }
            }
        },
        "/manage/credential/{id}": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "return list of allowed credentials",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Credentials"
                ],
                "summary": "return list of allowed credentials",
                "operationId": "credentials_get_by_id",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "credential id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.Credential"
                        }
                    }
                }
            },
            "put": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "updates credential",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Credentials"
                ],
                "summary": "updates credential",
                "operationId": "credentials_update",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "credential id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "credential update parameter",
                        "name": "params",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/parameters.UpdateCredentialParam"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.Credential"
                        }
                    }
                }
            }
        },
        "/manage/persona": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "return list of allowed personas",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Connectors"
                ],
                "summary": "return list of allowed personas",
                "operationId": "personas_get_all",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/model.Persona"
                            }
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "model.ChatMessage": {
            "type": "object",
            "properties": {
                "chat_session_id": {
                    "type": "integer"
                },
                "citations": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                },
                "error": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "latest_child_message": {
                    "type": "integer"
                },
                "message": {
                    "type": "string"
                },
                "message_type": {
                    "type": "string"
                },
                "parent_message": {
                    "type": "integer"
                },
                "rephrased_query": {
                    "type": "string"
                },
                "time_sent": {
                    "type": "string"
                },
                "token_count": {
                    "type": "integer"
                }
            }
        },
        "model.ChatSession": {
            "type": "object",
            "properties": {
                "created_date": {
                    "type": "string"
                },
                "deleted_date": {
                    "$ref": "#/definitions/pg.NullTime"
                },
                "description": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "messages": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/model.ChatMessage"
                    }
                },
                "one_shot": {
                    "type": "boolean"
                },
                "persona_id": {
                    "type": "integer"
                },
                "user_id": {
                    "type": "string"
                }
            }
        },
        "model.Connector": {
            "type": "object",
            "properties": {
                "connector_specific_config": {
                    "$ref": "#/definitions/model.JSONMap"
                },
                "created_date": {
                    "type": "string"
                },
                "credential_id": {
                    "type": "integer"
                },
                "deleted_date": {
                    "$ref": "#/definitions/pg.NullTime"
                },
                "disabled": {
                    "type": "boolean"
                },
                "id": {
                    "type": "integer"
                },
                "input_type": {
                    "type": "string"
                },
                "last_attempt_status": {
                    "type": "string"
                },
                "last_successful_index_time": {
                    "$ref": "#/definitions/pg.NullTime"
                },
                "name": {
                    "type": "string"
                },
                "refresh_freq": {
                    "type": "integer"
                },
                "shared": {
                    "type": "boolean"
                },
                "source": {
                    "type": "string"
                },
                "tenant_id": {
                    "type": "string"
                },
                "total_docs_indexed": {
                    "type": "integer"
                },
                "updated_date": {
                    "$ref": "#/definitions/pg.NullTime"
                },
                "user_id": {
                    "type": "string"
                }
            }
        },
        "model.Credential": {
            "type": "object",
            "properties": {
                "created_date": {
                    "type": "string"
                },
                "credential_json": {
                    "$ref": "#/definitions/model.JSONMap"
                },
                "deleted_date": {
                    "$ref": "#/definitions/pg.NullTime"
                },
                "id": {
                    "type": "integer"
                },
                "shared": {
                    "type": "boolean"
                },
                "source": {
                    "type": "string"
                },
                "tenant_id": {
                    "type": "string"
                },
                "updated_date": {
                    "$ref": "#/definitions/pg.NullTime"
                },
                "user_id": {
                    "type": "string"
                }
            }
        },
        "model.JSONMap": {
            "type": "object",
            "additionalProperties": true
        },
        "model.Persona": {
            "type": "object",
            "properties": {
                "default_persona": {
                    "type": "boolean"
                },
                "description": {
                    "type": "string"
                },
                "display_priority": {
                    "type": "integer"
                },
                "id": {
                    "type": "integer"
                },
                "is_visible": {
                    "type": "boolean"
                },
                "llm_id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "starter_messages": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                },
                "tenant_id": {
                    "type": "string"
                }
            }
        },
        "parameters.CreateChatSession": {
            "type": "object",
            "properties": {
                "description": {
                    "type": "string"
                },
                "one_shot": {
                    "type": "boolean"
                },
                "persona_id": {
                    "type": "integer"
                }
            }
        },
        "parameters.CreateConnectorParam": {
            "type": "object",
            "properties": {
                "connector_specific_config": {
                    "$ref": "#/definitions/model.JSONMap"
                },
                "credential_id": {
                    "type": "integer"
                },
                "disabled": {
                    "type": "boolean"
                },
                "input_type": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "refresh_freq": {
                    "type": "integer"
                },
                "shared": {
                    "type": "boolean"
                },
                "source": {
                    "type": "string"
                }
            }
        },
        "parameters.CreateCredentialParam": {
            "type": "object",
            "properties": {
                "credential_json": {
                    "$ref": "#/definitions/model.JSONMap"
                },
                "shared": {
                    "type": "boolean"
                },
                "source": {
                    "type": "string"
                }
            }
        },
        "parameters.InviteParam": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string"
                },
                "role": {
                    "type": "string"
                }
            }
        },
        "parameters.UpdateConnectorParam": {
            "type": "object",
            "properties": {
                "connector_specific_config": {
                    "$ref": "#/definitions/model.JSONMap"
                },
                "credential_id": {
                    "type": "integer"
                },
                "disabled": {
                    "type": "boolean"
                },
                "input_type": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "refresh_freq": {
                    "type": "integer"
                },
                "shared": {
                    "type": "boolean"
                }
            }
        },
        "parameters.UpdateCredentialParam": {
            "type": "object",
            "properties": {
                "credential_json": {
                    "$ref": "#/definitions/model.JSONMap"
                },
                "shared": {
                    "type": "boolean"
                }
            }
        },
        "pg.NullTime": {
            "type": "object",
            "properties": {
                "time.Time": {
                    "type": "string"
                }
            }
        }
    },
    "securityDefinitions": {
        "ApiKeyAuth": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "Cognix API",
	Description:      "This is Cognix Golang API Documentation",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
