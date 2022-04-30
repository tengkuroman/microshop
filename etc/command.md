This is a command documentation to set custom-auth (pre-function serverless) plugin in a service.

curl -i -X POST http://localhost:8001/services/SERVICE_NAME/plugins -F "name=pre-function" -F "config.access[1]=@custom-auth.lua"
