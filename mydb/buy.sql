CREATE DATABASE IF NOT EXISTS wechat_saas;

USE wechat_saas;

CREATE TABLE user(
	id INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '用户表id',
    username VARCHAR(20) COMMENT '用户名',
    password VARCHAR(40) NOT NULL COMMENT '密码，MD5加密',
    phone VARCHAR(11) COMMENT '用户电话',
    role TINYINT NOT NULL COMMENT '角色：1普通用户 3管理员 7超级管理员',
    grade TINYINT NOT NULL COMMENT '用户等级',
    del_flag TINYINT NOT NULL DEFAULT(1) COMMENT '假删除标志 1正常 2删除',
    create_time DATETIME NOT NULL COMMENT '数据创建时间',
    update_time DATETIME NOT NULL COMMENT '数据最近修改时间',
    PRIMARY KEY (id)
)ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE product(
	id INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '产品表id',
    name VARCHAR(20) COMMENT '产品名',
    sub_name VARCHAR(20) COMMENT '产品副标题',
    main_image VARCHAR(255) NOT NULL COMMENT '产品主图url',
    detail_image  VARCHAR(255) NOT NULL COMMENT '产品详情图url',
    price MEDIUMINT UNSIGNED NOT NULL COMMENT '产品价格',
    del_flag TINYINT NOT NULL DEFAULT(1) COMMENT '假删除标志 1正常 2删除',
    create_time DATETIME NOT NULL COMMENT '数据创建时间',
    update_time DATETIME NOT NULL COMMENT '数据最近修改时间',
    PRIMARY KEY (id)
)ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE buy_order(
	id INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '订单表id',
    user_id INT UNSIGNED NOT NULL COMMENT '用户id',
    product_id INT UNSIGNED NOT NULL COMMENT '产品id',
    pay_id VARCHAR(255) COMMENT '支付凭证：微信',
    status TINYINT UNSIGNED NOT NULL COMMENT '订单状态：1待支付 3已支付 7取消',
    name VARCHAR(10) NOT NULL COMMENT '收货人',
    phone VARCHAR(11) NOT NULL COMMENT '收货电话',
    address VARCHAR(100) NOT NULL COMMENT '收货地址',
    remarks VARCHAR(255) COMMENT '订单备注',
    del_flag TINYINT NOT NULL DEFAULT(1) COMMENT '假删除标志 1正常 2删除',
    create_time DATETIME NOT NULL COMMENT '数据创建时间',
    update_time DATETIME NOT NULL COMMENT '数据最近修改时间',
    PRIMARY KEY (id),
    CONSTRAINT u_id FOREIGN KEY (user_id) REFERENCES user(id),
    CONSTRAINT p_id FOREIGN KEY (product_id) REFERENCES product(id)
)ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE activities(
	id INT UNSIGNED NOT NULL COMMENT '活动表id',
    product_id INT UNSIGNED NOT NULL COMMENT '产品id',
    burst INT NOT NULL COMMENT '令牌桶大小',
    limt INT NOT NULL COMMENT '是否开启令牌桶限流：0关闭 >0开启limit/s的限流',
    stock MEDIUMINT UNSIGNED NOT NULL COMMENT '产品库存',
    name VARCHAR(20) NOT NULL COMMENT '活动名称',
    sub_name VARCHAR(20) COMMENT '活动副标题',
    start_time DATETIME NOT NULL COMMENT '活动开始时间',
    ground TINYINT UNSIGNED NOT NULL COMMENT '活动上架：1上架 2下架',
    del_flag TINYINT NOT NULL DEFAULT(1) COMMENT '假删除标志 1正常 2删除',
    create_time DATETIME NOT NULL COMMENT '数据创建时间',
    update_time DATETIME NOT NULL COMMENT '数据最近修改时间',
    PRIMARY KEY (id),
   CONSTRAINT p_id2 FOREIGN KEY (product_id) REFERENCES product(id)
)ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- INSERT INTO user(id,username,password,phone,role,grade,create_time,update_time) VALUES(1,"zhangsan","123456","13211223366",1,11,UTC_DATE(),UTC_DATE());

-- INSERT INTO product(id,name,sub_name,main_image,detail_image,price,create_time,update_time) VALUES(1,"N95口罩","医用口罩","https://aliyun.com/main_image.jpg","https://aliyun.com/detail_image.jpg",9900,UTC_DATE(),UTC_DATE());

-- INSERT INTO buy_order(id,user_id,product_id,pay_id,status,name,phone,address,remarks,create_time,update_time) VALUES(1,1,1,"10001223324aashhhrf00001",1,"张三","13244554455","广东省深圳市宝安区123号","",UTC_DATE(),UTC_DATE());

-- INSERT INTO activities(id,product_id,burst,limt,stock,name,sub_name,start_time,ground,create_time,update_time) VALUES(1,1,10,10,9999,"抗疫惠民物质派发活动","抗疫专项行动",UTC_DATE(),1,UTC_DATE(),UTC_DATE());
