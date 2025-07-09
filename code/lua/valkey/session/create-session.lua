local argv = ARGV
local keys = KEYS
local server_call = server.call

local exists = tonumber(server_call("EXISTS", keys[1]))

if exists == 1 then
	return server.error_reply("KEY_EXISTS")
end

local session_field_name = "s"

local hmset = server_call(
	"HMSET",
	keys[1],
	session_field_name,
	argv[2]
)

if hmset.ok ~= "OK" then
	return server.error_reply("FAILED_TO_CREATE_SESSION")
end

server_call("EXPIRE", keys[1], argv[1])

return "OK"
