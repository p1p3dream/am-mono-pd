local argv = ARGV
local keys = KEYS
local server_call = server.call

local exists = tonumber(server_call("EXISTS", keys[1]))

if exists == 0 then
	return server.error_reply("KEY_NOT_FOUND")
end

local session_field_name = "s"

local hmget = server_call(
	"HMGET",
	keys[1],
	session_field_name
)

return {
	hmget[1],
}
