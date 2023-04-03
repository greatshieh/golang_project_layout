# Golang应用项目框架

参照[gin-vue-admin](https://github.com/flipped-aurora/gin-vue-admin)和[极客时间Go语言项目开发实战](https://time.geekbang.org/column/intro/100079601?tab=catalog)及其[源码](https://github.com/marmotedu/iam)实现

实现基本配置, 支持非web应用, 包括配置(viper)，log(zap), 数据库(mysql, pgsql)和邮件发送。

## 项目结构

```shell
├── cmd
│   └── application
├── configs
├── internal
│   └── application
└── pkg
    ├── app
    │   └── internal
    ├── db
    ├── global
    ├── options
    └── utils
```

| 文件夹          | 说明                    | 描述                |
|----------------|------------------------|--------------------|
| `cmd`          | 应用程序入口             | 应用程序主程序        |
| `-application` | application程序入口      | application应用主程序|
| `configs`      | 配置文件                 | yaml配置文件 |
| `internal`     | 应用程序实现              | 每个文件夹对应一个应用程序 |
| `-application` | 应用程序application的实现 | |
| `pkg`          | pkg包                   | 可以被应用程序和外部程序调用 |
| `-app`         | 系统级应用               | 系统及初始化函数           |
| `--internal`   | 初始化函数               | 只能被`app`层调用          |
| `-db`          | 数据库                   | 数据库连接和配置           |
| `-global`      | 全局对象                 | 全局对象                 |
| `-options`     | 组件配置                 | 组件结构体定义, 与yaml配置对应 |
| `-utils`       | 工具包                   | 工具函数的封装 |
