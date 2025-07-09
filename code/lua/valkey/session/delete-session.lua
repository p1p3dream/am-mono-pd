local keys = KEYS
local server_call = server.call

return server_call("DEL", keys[1]) or 0
