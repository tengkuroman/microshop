local http = require "resty.http"
local cjson = require "cjson"

local access_token = kong.request.get_header("Authorization")

if not access_token then
    kong.response.exit(401)  -- Unauthorized
end

local httpc = http:new()

-- Check user validation and get user info
local res, err = httpc:request_uri("http://user-srv:8082/auth/validate", {
    method = "POST",
    ssl_verify = false,
    headers = {
        ["Content-Type"] = "application/x-www-form-urlencoded",
        ["Authorization"] = access_token }
})

if not res then
    kong.response.exit(500)  -- Internal error
end

user_info = cjson.decode(res.body)

-- Delete bearer token and set user info to request header
kong.service.request.clear_header("Authorization")
kong.service.request.set_header("X-User-ID", user_info.user_id)
kong.service.request.set_header("X-User-Role", user_info.role)