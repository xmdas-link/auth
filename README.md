# Auth
基于gin写的一套登录认证方式。认证模式为，登录成功后，生成JWT，传输给客户端。客户端在header中携带JWT作为token访问服务端。

##登录方式
当前版本支持的登录方式：
1. 账号密码登录
2. OAuth2.0登录（IM、微信扫码）

## Token
当前版本支持Token存储方式：
1. mysql数据库内存储
2. redis存放（推荐）

获取用户数据的方法
```go
func XXHandler (c *gin.Context) {
    user := c.GetStringMapString(auth.CtxKeyAuthUser) //CtxKeyAuthUser: AuthUser
}
```
user存放内容由provider自定义，必存字段包括,id, provider, name, role。为了方便获取，目前只支持map[string]string格式的user数据存储。


# 例子
见example目录