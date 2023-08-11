# TinyTik

## 1、相关资料

**github仓库地址:** https://github.com/JiangH156/TinyTik


## 2、技术栈

Go版本：1.20

| 技术     | 功能             | 官网                                         |
|--------|----------------| -------------------------------------------- |
| Gin    | Web 框架，路由注册    | https://gin-gonic.com/zh-cn/                 |
| Gorm   | ORM框架，用于对象关系映射 | https://gorm.io/zh_CN/                       |
| MySQL  | 关系型数据库         | https://www.mysql.com/                       |
| Redis  | 缓存数据库          | https://redis.uptrace.dev/zh/guide/ |
| JWT    | 跨域认证，生成和验证令牌   | https://jwt.io/                              |
| Viper  | 配置文件           | https://github.com/spf13/viper               |
| Gomail | 邮件服务           | https://github.com/go-gomail/gomail          |
| Bcrypt | 密码加密服务         | https://godoc.org/golang.org/x/crypto/bcrypt |
|        |                |



## 3、团队分工

| 成员   | 介绍                 | 内容                        |
| ------ | -------------------- |---------------------------|
| 黄江   | 福建农林大学准大三   | 用户模块（注册、登录、信息管理）          |
| 周帅鹏 | 湖南大学软件研二在读 | 社交模块（用户关注、粉丝、好友）          |
| 周灿   |                      | 评论和聊天模块（评论、消息管理）          |
| 殷家豪 |                      | 喜欢和发布模块（视频流、点赞、喜欢、视频投稿发布） |

## 4、项目管理

github仓库地址：https://github.com/JiangH156/TinyTik 

主分支：master 

开发分支：dev

## 5、项目结构

```
go_lib
├─cmd             -- 后台启动管理  
├─common          -- 通用的代码或功能
├─config          -- 配置文件
├─controller      -- 控制器（Controller）层，接受请求并处理响应
├─dto             -- 数据传输（DTO）层，处理实体传输和响应数据的转换
├─middleware      -- 中间件层，处理请求处理前后的逻辑
├─model           -- 数据库实体（Model）层，处理数据库相关操作
├─public          -- 静态资源文件
├─resp            -- 响应处理层，格式化响应数据和处理异常
├─router          -- 路由（Router）层，负责处理路由和中间件的注册
├─service         -- 服务（Service）层，处理业务逻辑
├─test            -- 项目测试文件
├─utils           -- 工具类和通用函数
└─vo              -- 视图对象（VO）层，封装视图数据传输格式和管理响应数据格式
```

## 5、其他

1. **git规范**

> 表示代码提交的类型，可以是以下之一
> 示例： fix: update the bug
- feat: 新功能（feature）
- fix: 修复 Bug
- docs: 文档更新
- style: 代码格式（不影响功能）
- refactor: 重构代码
- test: 添加或修改测试代码
- chore: 构建过程或辅助工具的变动
- perf: 改进性能的代码更改
