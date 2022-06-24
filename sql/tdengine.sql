--创建库  参考：https://blog.csdn.net/john1337/article/details/120439636
-- days：数据文件存储数据的时间跨度，单位为天
-- keep：数据保留的天数
-- rows: 文件块中记录条数
-- comp: 文件压缩标志位，0：关闭，1:一阶段压缩，2:两阶段压缩
-- ctime：数据从写入内存到写入硬盘的最长时间间隔，单位为秒
-- clog：数据提交日志(WAL)的标志位，0为关闭，1为打开
-- tables：每个vnode允许创建表的最大数目
-- cache: 内存块的大小（字节数）
-- tblocks: 每张表最大的内存块数
-- ablocks: 每张表平均的内存块数
-- precision：时间戳为微秒的标志位，ms表示毫秒，us表示微秒
-- 数据保存90天  10天存一个文件块 内存块数 4个 允许修改
CREATE DATABASE dory_device KEEP 90 DAYS 10 BLOCKS 4 UPDATE 1;
SHOW VARIABLES;

USE dory_device;
--超级表的列分为两部分：动态部分，静态部分
-- 动态部分是采集数据，第一列为时间戳（ts）,其他列为采集的物理量
-- 静态部分指采集点的静态属性，一般作为标签。如采集点的地理位置、设备型号、设备组、管理员ID等
-- 创建正常上数超级表    字段：时间戳、数据、状态。 标签：产品ID、设备ID、属性ID、从机ID、
CREATE STABLE if not exists dory_device.device_data  (ts timestamp,gas float, status bool) TAGS(product_id binary(30), device_id binary(30),slave int,model_id binary(30));

CREATE STABLE if not exists dory_device.device_data  (ts timestamp,gas float) TAGS( device_id binary(30),slave int);


--根据超级表创建子表
CREATE TABLE if not exists  dory_device.device_data_01 USING dory_device.device_data TAGS("123","123",1,"123");

--插入时候自动创建子表
INSERT INTO dory_device.device_data_01 USING dory_device.device_data  TAGS("123","123",1,"123")VALUES (now,10.2,false);


--创建告警表  字段：时间戳、数据、告警级别  标签 产品ID、设备ID、属性ID、从机ID、告警ID
CREATE STABLE if not exists dory_device.device_alarm  (ts timestamp,gas float, alarm_level int) TAGS(product_id binary(30), device_id binary(30),slave int,model_id binary(30),alarm_id binary(30));

CREATE STABLE if not exists dory_device.device_alarm  (ts timestamp,gas float, alarm_level int) TAGS( device_id binary(30),slave int);


--创建通知表
CREATE STABLE if not exists dory_device.device_notify  (ts timestamp,notify_type binary(10),phone_Number binary(30)) TAGS(product_id binary(30), device_id binary(30),alarm_id binary(30),template_id binary(30));
