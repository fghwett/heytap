# HeyTap
欢太商城每日签到Go脚本，项目 [https://github.com/hwkxk/HeytapTask](https://github.com/hwkxk/HeytapTask) 的Go版本

作者python版本已经够用了，但是因为我在群晖上的pip环境又问题，导致每次都只能手动运行项目，故写了顺手的版本


## 使用方法
第一次使用需要下载或编译、并填写相关配置，之后就可以直接使用

### 直接下载使用


### 从源码编译使用
1. 下载源码
    ```shell
    git clone https://github.com/fghwett/heytap.git
    ```
   
2. 复制配置文件并更改为自己的
    ```shell
    cp config.yml.example config.yml
    ```

3. 编译项目
    ```shell
    go build ./
    ```

4. 给编译后的文件添加运行权限（windows平台不用）
   ```shell
   chmod a+x ./heytap
   ```

5. 运行做任务
   ```shell
   ./heytap 
   ```
   
## 注意事项
cookie和ua需要手机抓包获取，两个要一起获取。尽量使用家庭环境网络，有些接口可能屏蔽国外IP

更多注意事项可参考 [原项目](https://github.com/hwkxk/HeytapTask)

## 开发
**项目文件：**

```log
.
├── LICENSE
├── README.md  # 文档
├── config     # 配置文件读取
│   ├── config.go
│   └── config_test.go
├── config.yml.example  # 示例配置文件
├── go.mod
├── go.sum
├── main.go    # 入口函数
├── notify     # 第三方通知 可自行添加其他通知方式
│   └── notify.go
├── task       # 任务中心
│   ├── model.go    # 数据接口
│   ├── task.go     # 任务逻辑
│   └── task_test.go
└── util       # 工具函数
    ├── gzip.go           # gzip函数解码
    ├── http.go           # http处理结果封装
    └── rand.go           # 随机函数相关
```

## 申明
- 本项目仅用于学习研究，禁止任何人用于商业用途，本人不保证其合法性，准确性，完整性和有效性
- 本人无法100%保证使用本项目之后不会造成账号异常问题，请根据情况自行判断再下载执行！否则请勿下载运行！
- 如果任何单位或个人认为该项目的脚本可能涉及侵犯其权利，则应及时通知并提供相关证明，我将在收到认证文件后删除相关脚本.