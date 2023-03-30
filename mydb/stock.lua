if (redis.call('HEXISTS',KEYS[1],KEYS[2])==1) then
    local stock = tonumber(redis.call('HGET',KEYS[1],KEYS[2]))
    if (stock == -1) then -- 不限库存
        return -1
    end
    if (stock > 0) then
        local d_stock = stock - 1
        redis.call('HSET',KEYS[1],KEYS[2],d_stock) -- 扣减库存
        return stock -- 返回本次消耗库存之后的库存
    end
    return 0 -- 库存不足
end
return -2 -- 不存在该商品

