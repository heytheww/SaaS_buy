if (redis.call('exists',KEYS[1])==1) then
    local stock = tonumber(redis.call('get',KEYS[1]))
    if (stock == -1) then -- 不限库存
        return -1
    end
    if (stock>0) then
        redis.call('incrby',KEYS[1],-1) -- 扣减库存
        return stock -- 返回本次消耗库存之前的库存
    end
    return 0 -- 库存不足
end
return -2 -- 不存在该商品