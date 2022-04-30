This is a command documentation to set custom-auth (pre-function serverless) plugin in a service.

curl -i -X POST http://localhost:8001/services/SERVICE_NAME/plugins -F "name=pre-function" -F "config.access[1]=@custom-auth.lua"

curl -i -X POST http://localhost:8001/services/order-service-auth/plugins -F "name=pre-function" -F "config.access[1]=@custom-auth.lua"

curl -i -X POST http://localhost:8001/services/payment-service-auth/plugins -F "name=pre-function" -F "config.access[1]=@custom-auth.lua"

curl -i -X POST http://localhost:8001/services/product-service-auth/plugins -F "name=pre-function" -F "config.access[1]=@custom-auth.lua"

curl -i -X POST http://localhost:8001/services/shopping-service-auth/plugins -F "name=pre-function" -F "config.access[1]=@custom-auth.lua"

curl -i -X POST http://localhost:8001/services/user-service-auth/plugins -F "name=pre-function" -F "config.access[1]=@custom-auth.lua"
