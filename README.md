# 业务-逻辑

## 页面静态化
1.商品静态化。商品的名称、描述、图片等相对固定，把他们做成静态页面，后期可以通过管理端生成，通过nginx等进行动静分离，把对商品页面的刷新和加载和后端服务分离开，由nginx承担。为了保证数据完备性，后端仍要建立完整的商品表，并且也为未来拓展生成静态页面功能作基础。
2.库存需要发送请求获取，每刷新一次页面获取一次库存，库存请求通过 后端服务 打在redis上，避免访问mysql，以避免磁盘I/O。
3.点击秒杀按钮触发请求后端服务，进入抢购流程。

## 读多写少
只有少数人才能秒杀成功，把订单写入数据库，多数人只能读库存，然后秒杀失败，所以是 读多写少。
使用redis缓存解决方案：
1.缓存商品id、对应库存

生成缓存方式：
1.预热
开始秒杀时，把所有的商品id和库存，同步到缓存中
2.redis商品数据过期--缓存击穿
2.1 解决缓存击穿
预热的缓存数据有可能过期失效，缓存失效后，后端服务在redis查不到时一般回去mysql查，为了防止瞬间全部流量打在mysql上，使用lua脚本，把 查询到建立缓存 这段逻辑使用lua脚本代替。
3.查找不存在的商品id--缓存穿透
使用布隆过滤器，上架商品时，先在布隆过滤器中生成指纹，布隆过滤器的逻辑是 说一个指纹见过，可能实际是见过或没见过；说一个指纹没见过，实际就真的没见过。请求查询商品库存时，先访问布隆过滤器，如果查不到该商品，就表示mysql也没有该商品，直接返回抢购失败，避免缓存穿透。

## 库存扣减
要解决的问题：
1.并发扣减内存时，保证不超卖
2.并发时防止重复卖给同一用户

解决方案：使用原子操作 代替 锁。
1.先判断商品id是否存在，如果不存在则直接返回。（不管理商品缓存，只管理库存）
2.获取该商品id的库存，判断库存如果是-1，则直接返回，表示不限制库存。
3.如果库存大于0，则扣减库存。
4.如果库存等于0，是直接返回，表示库存不足。

## 使用Redis+Lua实现解决方案
https://help.aliyun.com/document_detail/92942.html

KEYS[1]是传入的key值，相当于函数的参数，可多次使用。
```
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
```

## mq（Message Queue）异步处理
秒杀->下单->支付，三者并发量不均等，秒杀最大，下单和支付很小。
支付模块需要对接第三方系统，本项目留好接口，不作实现。

1.通过 消息队列，让秒杀的下单解耦，下单和支付解耦。

基本架构：秒杀->mq服务器->下单->mq服务器->支付，使用异步消息。

常见问题：
如何保证消息被消费呢？
秒杀->mq服务器 之间增加 消息发送表，该表记录了已发送消息的状态：待处理/已处理，下单（消费者）消耗该消息，调用 消息发送表 修改状态为已处理，ack应答。

如何防止消息丢失呢？
为了保证 秒杀 发送的 消息 一定能到达 mq服务器，使用job，增加重试（重发）机制，类似TCP。
job每隔一段时间去查询消息发送表中状态为待处理的数据，然后重新发送mq消息。

如何防止重复消费呢？
当消费者消费完毕后，需要发送 ack消息 ，修改消息状态为 已处理。如果ack消息发送失败，将造成 消息的重复消费，加上 重试机制，消息重复消费概率增大。
在 mq服务器->下单 之间增加 消息处理表，标记哪些消息（id）是已经处理过的，再次遇到该消息时，直接不作处理。关键点，要保证 消息处理表 和 下单 是绝对一致的，放在同一个事务中，保证原子操作。

如何处理垃圾消息问题？
下单一直失败，一直没有调用 消息状态修改，job 会一直重试。这里直接设置 消息发送表 发生次数上限，达到上限，不再发该消息，未达上限，将次数加1，正常重试发消息。

如何处理延迟消费问题？
15分钟之内还未完成支付，订单取消，库存恢复。
下单->支付 之间的消息发 延迟消费消息，达到 延迟时间后，支付 消费该消息，如果 订单状态是 待支付，修改该订单状态为取消，库存恢复，否则，说明已支付。

2.怎么做消息队列和消费保证
本系统基于golang，为了保证部署的简单，这里使用redis.Stream作为消息队列，redis.Stream参考了kafka的设计。

通过PEL+ack保证消息已消费，通过redis本身的持久化能力，保证消息队列本身不会丢失。

## 限流
1.如果用户全部用手抢，也会不断 点击 秒杀按钮，一个用户会生成多条请求。
2.如果用户使用机器生成请求，1s可生成上千请求，而人手一秒只能生成一条请求。

常见方案：
对1，在前端JS中限制，每多少秒，可以发送1次请求，但是对2无效，2可以绕过JS限制。为了公平，需要在后端限流。

方案：限制同一个用户id，比如每分钟只能请求5次接口。
问题：请求方 模拟多个用户请求时无效

方案：限制同一个ip
问题：可能会导致公用ip的所有用户被连累。请求方 使用代理模拟请求ip时无效

推荐方案：移动滑块验证码
问题：影响用户体验，操作繁琐。解决方案是，提高业务门槛，使用用户画像，例如只有正式会员才能参与抽奖、等级到达3级以上的才可以参与、男性用户才可以参与。

本系统方案：
1.在几千-几万流量时，主要通过golang的rete包的 令牌桶方案 来限流，即业务层（网络模型的应用层）做限流
2.在十几万-几十万流量时，在nginx通过限制连接数来限制流量，即主机流量层（网络模型的网络层）做限流


封装MySQL https://www.liwenzhou.com/posts/Go/mysql/