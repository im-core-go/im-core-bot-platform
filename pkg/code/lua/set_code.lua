local val = ARGV[1]
local key = KEYS[1]
local cntKey = key..":cnt"
if not redis.call("get", key) then
    redis.call("set", key, val)
    redis.call("expire", key, 300)
    redis.call("set", cntKey, 5)
    redis.call("expire", cntKey, 300)
    return 0
end

local ttl = tonumber(redis.call("ttl", key))
if ttl == -1 then
    return -1
elseif ttl == -2 or ttl < 240 then
    redis.call("set", key, val)
    redis.call("expire", key, 300)
    redis.call("set", cntKey, 5)
    redis.call("expire", cntKey, 300)
    return 0
else
    return -2
end