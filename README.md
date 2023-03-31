# 业务-逻辑

这是一个业务系统的重构项目，展示了golang的显著优势：1.互联网的C语言 2.简单开发而高效使用  

参考文章:

【1】https://github.com/gzc426/Java-Interview/blob/master/%E9%A1%B9%E7%9B%AE%E6%8E%A8%E8%8D%90/%E7%A7%92%E6%9D%80.md  
【2】https://mp.weixin.qq.com/s?__biz=MzU0OTE4MzYzMw==&mid=2247517155&idx=4&sn=bf198afe7d2b498063a0416accffe74f&chksm=fbb10a1dccc6830be5c6b74cbf2de93c9fdd9dd12da4ae1c9b407d35949b92e8656856b7cd16&scene=27  

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

# 项目架构说明

## 系统的数据状态
本系统状态采用 启动预热 方式设置系统初始数据状态，即启动系统时，马上把数据库相应的数据预热到 redis中，后期将继续拓展 动态配置 方式，即通风后台管理端，把数据提交到数据库，然后马上缓存一份到redis中，让系统数据状态可控。

## 1.不使用orm
考虑到系统拓展、维护的复杂度，不使用orm

## 2.model 设计
本项目的 model 设计如下：
1.接口model围绕接口最终数据设计，不围绕基础表的结构设计。
2.表的mode围绕基础表的结构设计。

## 3.隔离级别
本系统设计上，抢购功能尽量采用消息队列，避免并发访问mysql，但是其他模块仍可能并发访问mysql，因此使用 事务隔离级别3。

读未提交 READ UNCOMMITTED | 0 : 存在脏读，不可重复读，幻读的问题。   
读已提交 READ COMMITTED | 1 : 解决脏读问题，存在不可重复读，幻读（幻行）的问题。  
可重复读 REPEATABLE READ | 2 : 解决脏读，不可重复读的问题，存在幻读（幻行），默认隔离级别，使用MVCC机制（多版本并发控制）实现可重复读。  
序列化 SERIALIZABLE | 3 : 解决脏读，不可重复读，幻读，可保证事务安全，但完全串行执行，性能最低。  

## 4.时间处理
在连接数据库时，开启parseTime=True，自动把datetime转换golang的time.Time。
https://github.com/go-sql-driver/mysql#columntype-support

## 5.请求参数的校验
关于请求参数的校验，参考以下：  
【1】https://gin-gonic.com/zh-cn/docs/examples/binding-and-validation/  
【2】https://pkg.go.dev/github.com/go-playground/validator/v10#hdr-One_Of  
【3】https://raw.githubusercontent.com/go-playground/validator/master/_examples/simple/main.go  
gin集成了很好使用的参数校验，包括缺失校验、类型校验等。另外，服务端只作简单校验，保证基本数据不缺失和类型正确，以减少服务复杂度，降低系统运维、拓展成本，业务校验如手机号校验，由前端来保证。

注意，如果一个参数可以为空，那么校验就不要使用required，如
```
Remarks  string `json:"remarks" binding:"required"`
```
一旦用户向remarks传入""，将会被校验为 remarks未传值而不是 remarks传了值但是是空值，这会使校验失败。因为，应该将这里参数设置为可不传。
```
Remarks  string `json:"remarks"`
```

## 6.分页查询
分页查询是后台管理最基本的功能，本系统仅作设计，暂不实现。设计如下：  
请求时，传入以下参数：  
```
per 每页记录数
total 总记录数
page 当前是第几页

total = page*per
```

## 7.时间的处理
注意这样一段程序：
```
package main

import (
	"fmt"
	"time"
)

func main() {
	location, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		panic(err)
	}

	timeInUTC := time.Date(2018, 8, 30, 12, 0, 0, 0, time.UTC)
	fmt.Println(timeInUTC.In(location))
}

```
Output:
2018-08-30 05:00:00 -0700 PDT

看起开很奇怪的时间格式实际上是通用格式，直接返回给前端即可。
PDT时间(Pacific Daylight Time)太平洋夏季时间，比UTC早7小时，而北京时间比UTC晚8小时，通常认为GMT和UTC一致，即GMT=UTC+0。

这种时间格式和我们国内应用常见的 2022-01-02 13:55:11 不太一样，后者通过前端JS的Date()或者使用moment.js可以很方便地得到
```
new Date("2018-08-30 05:00:00 -0700 PDT")
Thu Aug 30 2018 20:00:00 GMT+0800 (中国标准时间)

new Date("2023-03-27T16:15:47+08:00")
Mon Mar 27 2023 16:15:47 GMT+0800 (中国标准时间)
```

## redis的准备
参考：  
【1】https://hub.docker.com/_/redis  
【2】https://redis.uptrace.dev/zh/guide/go-redis.html  
【3】https://redis.io/commands/xtrim/  
【4】https://redis.io/commands/xadd/  

本系统将使用redis的两项功能，一是持久化缓存，二是异步消息队列。

持久化缓存：
```
docker run --name buy -d redis redis-server --save 300 1 --loglevel warning
```
容器后台运行，每300秒若至少有1次写错字，就进行一次快照保存，只记录warning级别的log

Redis进程运行日志的级别优先级从高到低分别是warning、notice、verbose、debug，程序会打印高于或等于所设置级别的日志，设置的日志等级越高，打印出来的日志就越少。

运行日志：
1.warning warning表示只打印非常重要的信息。  
2.notice notice表示打印适当的详细信息，适用于生产环境。  
3.verbose verbose表示记录系统及各事件正常运行状态信息。  
4.debug debug表示记录系统及系统的调试信息。  

【注意】
在开发和测试阶段，应采用本地目录；生成环境，使用创建的卷。

```
#dev
docker run -p 6379:6379 -v C:/Users/Administrator/Desktop/v:/data --name buy -d redis redis-server --save 300 1 --loglevel warning 

# prod
docker volume create v1
docker run -p 6379:6379 -v v1:/data --name buy -d redis redis-server --save 300 1 --loglevel warning 
```

异步消息队列：
使用以下命令测试redis.Streams并与本系统行为进行对比，检查正确性
1.创建stream
MAXLEN ~ 1000 限定长度约是1000，可能多几十条，或MAXLEN = 1000，精确控制数量，这些策略是插入新的消息，驱逐旧的消息，应用在本系统时，可能会造成消息被驱逐而丢失，进而导致订单生成数据丢失。
这里存在一个问题，消息队列是存在于内存的，为避免内存消耗过大，应使用XLEN判断异步消息队列长度，若长度超过一定数，则停止插入消息，等待一定时间后，再尝试插入消息。
```
XADD key [NOMKSTREAM] [<MAXLEN | MINID> [= | ~] threshold
  [LIMIT count]] <* | id> field value [field value ...]
```
```
XADD ww * user 1001 product_id 1001
XADD ww * user 1002 product_id 1002
XADD ww * user 1003 product_id 1003
XADD ww * user 1004 product_id 1004
XADD ww * user 1005 product_id 1005

XRANGE ww - + #查看stream
XLEN ww #查看stream长度
```


2.创建消费组
规定组内消费者从第一条消息开始消费：0-0
```
XGROUP CREATE ww cg1 0-0
```
获取stream各消费组的详情
name表示组名，entries-read表示已被读取数，lag表示未被读取数
```
XINFO GROUPS ww 
```

3.组内消费
```
XREADGROUP GROUP group consumer [COUNT count] [BLOCK milliseconds]
  [NOACK] STREAMS key [key ...] id [id ...]
```
一条一条消费，>指定读取 从未被消费过的消息，0指定当前消费者消费了但未ack的消息
```
XREADGROUP GROUP cg1 c1 COUNT 1 BLOCK 0 STREAMS ww >
XREADGROUP GROUP cg1 c1 COUNT 1 BLOCK 0 STREAMS ww 0
```

4.ack确认消费
```
XACK ww cg1 1680100221976-0
```

## 8.redis扣减库存的设计
参考：
【1】https://redis.com/redis-best-practices/lua-helpers/  

```
HGET stock 1001
HEXISTS stock 1001
HSET stock 1001 100
```

使用lua脚本做到原子操作
```
if (redis.call('HEXISTS',KEYS[1],KEYS[2])==1) then
    local stock = tonumber(redis.call('HGET',KEYS[1],KEYS[2]))
    if (stock == -1) then -- 不限库存
        return -1
    end
    if (stock > 0) then
        stock = stock - 1
        redis.call('HSET',KEYS[1],KEYS[2],stock) -- 扣减库存
        return stock -- 返回本次消耗库存之后的库存
    end
    return 0 -- 库存不足
end
return -2 -- 不存在该商品
```

## 9.限流
context的超时时长不能小于一个令牌生成的时长，否则，只有一开始的请求可以拿到Bursts个令牌，后来的请求，失败率很高。

## 10.异步生成订单

排队号
订单号