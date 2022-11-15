swagger:
	swagger generate spec -o ./swagger.yaml --scan-models
swagger-serve: swagger
	swagger serve -F=swagger swagger.yaml