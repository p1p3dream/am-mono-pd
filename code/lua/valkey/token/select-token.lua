local argv = ARGV
local keys = KEYS
local server_call = server.call

local exists = tonumber(server_call("EXISTS", keys[1]))

if exists == 0 then
	return server.error_reply("KEY_NOT_FOUND")
end

local payload_field_name = "p"
local quota_field_name = "q"

local hmget = server_call(
	"HMGET",
	keys[1],
	payload_field_name,
	quota_field_name
)

local decreaseBy = tonumber(argv[1])
local quota = tonumber(hmget[2]) - decreaseBy

if quota < 0 then
	return server.error_reply("INSUFFICIENT_QUOTA")
end

if quota == 0 then
	server_call("DEL", keys[1])
else
	server_call("HSET", keys[1], quota_field_name, quota)
end

return {
	hmget[1],
	quota,
}
