local key = KEYS[1]
local cnt = redis.call("get", key..":cnt")

if not cnt then
    return -3
end
cnt = tonumber(cnt)

local expected = ARGV[1]
if cnt <=0 then
    return -1
elseif expected == redis.call("get", key) then
        redis.call("del", key..":cnt")
        redis.call("del", key)
        return 0
    else
        redis.call("decr", key..":cnt")
        return -2
    end