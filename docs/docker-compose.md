本节将讲述如何使用 docker-compose 启动整个项目

本功能模仿 [go–zero-tooklook](https://github.com/Mikaelemmmm/go-zero-looklook) ，晚上状态不佳，有疏漏可参考那个项目的文档，真的非常详细

首先请根据你的指令集选择对应的 `docker-compose` 文件！

![image-20230214002205723](https://pic-go-img.oss-cn-hangzhou.aliyuncs.com/202302140022813.png)

# 准备工作

## 准备配置文件

来到各个服务的配置目录，将示例配置文件复制一份，去掉example后缀

当然，OSS配置我不能告诉你（

![image-20230214003609781](https://pic-go-img.oss-cn-hangzhou.aliyuncs.com/202302140036810.png)

## 创建消息队列

先启动里面的卡夫卡

![image-20230213232725379](https://pic-go-img.oss-cn-hangzhou.aliyuncs.com/202302132327650.png)

卡夫卡启动后，进去创建用于传输日志的消息队列

```bash
$ docker exec -it kafka /bin/sh
$ cd /opt/kafka/bin/
$ ./kafka-topics.sh --create --zookeeper zookeeper:2181 --replication-factor 1 -partitions 1 --topic h68u-tiktok-log
```

或者你也可以直接在 Docker Desktop 启动它的终端

![image-20230213233046685](https://pic-go-img.oss-cn-hangzhou.aliyuncs.com/202302132330710.png)

## 初始化 MySQL

接着往下翻，找到 mysql 并启动

进入mysql，准备一下 root 用户

```bash
$ docker exec -it mysql mysql -uroot -p
##输入密码：PXDN93VRKUm8TeE7
$ use mysql;
$ update user set host='%' where user='root';
$ FLUSH PRIVILEGES;
```

![image-20230213233805227](https://pic-go-img.oss-cn-hangzhou.aliyuncs.com/202302132338255.png)

测试一下能否连接到数据库

![image-20230213233846140](https://pic-go-img.oss-cn-hangzhou.aliyuncs.com/202302132338164.png)

创建项目会访问到的库

![image-20230216223030583](https://pic-go-img.oss-cn-hangzhou.aliyuncs.com/202302162230942.png)

## 启动所有服务

![image-20230213234727364](https://pic-go-img.oss-cn-hangzhou.aliyuncs.com/202302132347411.png)

等待13个镜像全部拉取完成并启动

![image-20230213234242298](https://pic-go-img.oss-cn-hangzhou.aliyuncs.com/202302132342332.png)

项目本体使用 modd 热编译，任何修改都会自动重新编译并启动

![image-20230213234229785](https://pic-go-img.oss-cn-hangzhou.aliyuncs.com/202302132342841.png)

等待项目完全启动

![image-20230213234752966](https://pic-go-img.oss-cn-hangzhou.aliyuncs.com/202302132347005.png)

# 服务监控

[9090](http://127.0.0.1:9090/) 端口访问普罗米修斯，来到target页面，稍等片刻等待服务全部亮起

![image-20230213235050066](https://pic-go-img.oss-cn-hangzhou.aliyuncs.com/202302132350103.png)

![image-20230213235110737](https://pic-go-img.oss-cn-hangzhou.aliyuncs.com/202302132351768.png)

![image-20230213235157072](https://pic-go-img.oss-cn-hangzhou.aliyuncs.com/202302132351112.png)

# 链接追踪

去 apifox 上请求一下，然后来到 [16686](http://127.0.0.1:16686/) 端口访问 Jaeger，就能查看链路追踪

![image-20230214001801336](https://pic-go-img.oss-cn-hangzhou.aliyuncs.com/202302140018455.png)

![b696e8bb719a6d19fee3e168e9a76cc3](https://pic-go-img.oss-cn-hangzhou.aliyuncs.com/202302140025172.png)

# 日志收集

通过 [5601](http://127.0.0.1:5601/) 端口访问 Kibana



![image-20230214001900850](https://pic-go-img.oss-cn-hangzhou.aliyuncs.com/202302140019908.png)

![image-20230214001929216](https://pic-go-img.oss-cn-hangzhou.aliyuncs.com/202302140019273.png)

![image-20230214001950254](https://pic-go-img.oss-cn-hangzhou.aliyuncs.com/202302140019297.png)

![image-20230214001958514](https://pic-go-img.oss-cn-hangzhou.aliyuncs.com/202302140019553.png)

![image-20230214002030719](https://pic-go-img.oss-cn-hangzhou.aliyuncs.com/202302140020772.png)

![image-20230214002634676](https://pic-go-img.oss-cn-hangzhou.aliyuncs.com/202302140026725.png)

演示：根据请求ID查询日志排查错误

可以尝试把 mysql 关掉再请求

![ed2c40b507cf0fda52187cd346096a03](https://pic-go-img.oss-cn-hangzhou.aliyuncs.com/202302140025057.png)

![f09a295e995d99a8b00ff9b25d031a71](https://pic-go-img.oss-cn-hangzhou.aliyuncs.com/202302140025942.png)