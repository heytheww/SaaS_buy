if (redis.call('exists',KEYS[1])==1) then
    local stock = tonumber(redis.call('get',KEYS[1]))
    if (stock == -1) then
        return -1
    end
    if (stock > 0) then
        redis.call('incrby',KEYS[1],-1)
        return stock --返回开始的库存量，实际上的库存已经扣1
    end
    return 0
end
return -1