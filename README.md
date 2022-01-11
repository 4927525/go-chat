# 基于 WEBSOCKET + MONGODB 的IM即时聊天DEMO

写在前面

> 这个项目是基于WebSocket + MongoDB + MySQL + Redis。  
> 业务逻辑很简单，只是两人的聊天。

- `MySQL` 用来存储用户基本信息
- `MongoDB` 用来存放用户聊天信息
- `Redis` 用来存储处理过期信息

## 项目功能实现

- 登录鉴权（jwt-go）

- 清单 crud

- 分页

## 项目架构

```
TodoList/

├── api

├── conf

├── e

├── middleware

├── model

├── pkg

│  ├── util

├── routes

├── serializer

└── service
```

- api : 用于定义接口函数

- conf : 用于存储配置文件
- e: 错误代码文件

- middleware : 应用中间件

- model : 应用数据库模型

- pkg / util : 工具函数

- routers : 路由逻辑处理

- serializer : 将数据序列化为 json 的函数

- service : 接口函数的实现

## 项目依赖

- go1.16

- gin

- gorm

- mysql
- redis
- mongodb

## 项目接口

```bash
[GIN-debug] GET    /ping                     --> chat/routers.NewRouters.func1 (5 handlers)
[GIN-debug] POST   /user/register            --> chat/api.UserRegister (5 handlers)
[GIN-debug] GET    /ws                       --> chat/service.Handler (5 handlers)
```

## 项目运行

```bash
go mod tidy
```

```bash
go run main.go
```