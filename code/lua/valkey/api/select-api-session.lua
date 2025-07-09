local argv = ARGV
local keys = KEYS
local server_call = server.call

local exists = tonumber(server_call("EXISTS", keys[1]))

if exists == 0 then
	return server.error_reply("KEY_NOT_FOUND")
end

local invalid_field_name = "i"
local quota_exhausted_field_name = "q"
local session_field_name = "s"

local hmget = server_call(
	"HMGET",
	keys[1],
	invalid_field_name,
	quota_exhausted_field_name,
	session_field_name
)

return hmget
