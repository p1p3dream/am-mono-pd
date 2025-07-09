local argv = ARGV
local keys = KEYS
local server_call = server.call

local invalid_field_name = "i"
local quota_exhausted_field_name = "q"
local session_field_name = "s"

local hmset = server_call(
	"HMSET",
	keys[1],
	invalid_field_name, argv[1],
	quota_exhausted_field_name, argv[2],
	session_field_name, argv[3]
)

if hmset.ok ~= "OK" then
	return server.error_reply("FAILED_TO_UPDATE_SESSION")
end

local expire = server_call(
	"EXPIRE",
	keys[1],
	argv[4]
)

if expire ~= 1 then
	return server.error_reply("FAILED_TO_EXPIRE_SESSION")
end

return "OK"
