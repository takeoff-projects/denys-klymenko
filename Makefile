SHELL := /bin/bash

VERSION := 1.0

api-swagger:
	swag init --g api/pets/api.go -o ./docs