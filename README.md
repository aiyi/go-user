#Overview
go-user is a micro service written in go to provide essential user management features such as:
- Registration and confirmation via email/phone.
- Login and logout.
- Auth token / session management.
- Password reset.  

# need go1.5.2+ database/sql
see  
https://github.com/golang/go/commit/40457745e51eb327751de5be4c69c9079db69f66

## api 返回
成功时返回的样式：  
```json
{
    "err_code": 0, 
    "err_msg": "success", 
    "token": "eyJzaWQiOiJ0ZW1wLlZFUF9Lam8zdHU5bWhrSmRVRVNDTFh6b0NLMmU0NTlWIiwidG9rZW5faWQiOiJmZjZkNGE2NzE1NDExZWJkOTFlNmYxNTE1ODQxOTA2ZSIsImF1dGhfdHlwZSI6Imd1ZXN0IiwiZXhwX2FjY2VzcyI6MCwiZXhwX3JlZnJlc2giOjB9.f0b03ca7af190925c9ed12d555690a8b81aa9c2c691ff6942ad689cf99b9b865"
}
```

失败时： 
```json 
{
    "err_code": 20000, 
    "err_msg": "token missing"
}
```
通过 err_code 来判断，详细的 err_code 请移步到:  
https://github.com/aiyi/go-user/blob/master/frontend/errors/list.go
## token

eyJzaWQiOiJVLUhFVlZiSnR1OW1oa0pkV3h6XzBDbEo4TDFEemllYSIsInRva2VuX2lkIjoiODExZTNmMGJhYTc1ZmNjOGRjMWVkOGU3MDY4YzA5ZjAiLCJhdXRoX3R5cGUiOiJwaG9uZV9wYXNzd29yZCIsImV4cF9hY2Nlc3MiOjE0NDQ3ODk1NDksImV4cF9yZWZyZXNoIjoxNDc2MzAyMjQ1fQ.a3850335036e76c200aef18ce0280d9807fb4a6ee96e19a9c9792666146dbfb0

获取了token之后，以后所有的api都要携带这个token去调用！  
token 记得保存在客户端，可以作为免登凭证。  

## token 解码

点号 '.' 前面是 base64 编码（最后的 '=' 已经去掉，解码的时候要加上），后面是签名，客户端不用理会；  
把点号前面的base64编码解码，得到token的json：

```json
{
    "sid": "U-HEVVbJtu9mhkJdWxz_0ClJ8L1Dziea", 
    "token_id": "811e3f0baa75fcc8dc1ed8e7068c09f0", 
    "auth_type": "phone_password", 
    "exp_access": 1444789549, 
    "exp_refresh": 1476302245
}
```

sid，token_id 客户端不用理会，  
auth_type 是token的认证类型:  

```go
const (
	AuthTypeGuest          = "guest"           // 游客
	AuthTypeEmailPassword  = "email_password"  // 邮箱+密码
	AuthTypeEmailCheckCode = "email_checkcode" // 邮箱+校验码, 校验码推送到邮箱
	AuthTypePhonePassword  = "phone_password"  // 手机+密码
	AuthTypePhoneCheckCode = "phone_checkcode" // 手机+校验码, 校验码短信推送给手机
	AuthTypeOAuthQQ        = "oauth_qq"        // QQ oauth
	AuthTypeOAuthWechat    = "oauth_wechat"    // 微信 oauth
	AuthTypeOAuthWeibo     = "oauth_weibo"     // 微博 oauth
)
```

exp_access 是token的过期时间，unixtime，过期了需要刷新 token  
exp_refresh 是token的刷新过期时间，unixtime，过了这个时间，需要重新认证  

对于 auth_type == guest 的 token，没有过期时间一说，永不过期！  

## 认证

### guest
```http
curl -X GET \    
  -G \  
  --data-urlencode 'auth_type=guest' \    
  https://xxx.xxx.xxx/user/auth 
```
返回：
```json   
{
    "err_code": 0, 
    "err_msg": "success", 
    "token": "eyJzaWQiOiJ0ZW1wLlZFUF9Lam8zdHU5bWhrSmRVRVNDTFh6b0NLMmU0NTlWIiwidG9rZW5faWQiOiJmZjZkNGE2NzE1NDExZWJkOTFlNmYxNTE1ODQxOTA2ZSIsImF1dGhfdHlwZSI6Imd1ZXN0IiwiZXhwX2FjY2VzcyI6MCwiZXhwX3JlZnJlc2giOjB9.f0b03ca7af190925c9ed12d555690a8b81aa9c2c691ff6942ad689cf99b9b865"
} 
``` 
### 邮箱+密码
```http
curl -X GET \   
  -G \  
  --data-urlencode 'auth_type=email_password' \   
  --data-urlencode 'email=test@qq.com' \   
  --data-urlencode 'password=123456' \   
  https://xxx.xxx.xxx/user/auth  
```
返回：
```json  
{
    "err_code": 0, 
    "err_msg": "success", 
    "token": "eyJzaWQiOiJWRVEwVTJMcHR1OW1oa0pkVUVTQ0xuMFJlR1BTRnlwRiIsInRva2VuX2lkIjoiMjhlMjBkNDIwNzhhNWM5M2YwYjA1MDM4NTVkYjAzNjEiLCJhdXRoX3R5cGUiOiJlbWFpbF9wYXNzd29yZCIsImV4cF9hY2Nlc3MiOjE0NDQ3OTQ3NzIsImV4cF9yZWZyZXNoIjoxNDc2MzQ0NTI0fQ.1dc8aa84064a902de74d4f582a8e41d01d00c690888bc958ad4d8ac9aebad805"
} 
```
### 手机+密码
```http
curl -X GET \    
  -G \   
  --data-urlencode 'auth_type=phone_password' \   
  --data-urlencode 'phone=18888888888' \   
  --data-urlencode 'password=123456' \   
  https://xxx.xxx.xxx/user/auth  
```
返回：
```json  
{
    "err_code": 0, 
    "err_msg": "success", 
    "token": "eyJzaWQiOiJWRVJOQy1fdnR1OW1oa0pkVUVTQ0wxUG1xZWwwcVdSSSIsInRva2VuX2lkIjoiY2E5YmJhNjM2YmYxMzZiNmQxZjU4Njg3NzUyNzcwZTMiLCJhdXRoX3R5cGUiOiJwaG9uZV9wYXNzd29yZCIsImV4cF9hY2Nlc3MiOjE0NDQ3OTQ4MTMsImV4cF9yZWZyZXNoIjoxNDc2MzQ0NTY1fQ.ff5582e8680d30fc9e073ef3f4db94f152386dbe072b00fa0307818e5d8bf38e"
}
```
## 第三方认证

### 微信公众号
#### 获取认证url
```http
curl -X GET \      
  -H "x-token: {{token(normally guest)}}" \   
  -G \   
  --data-urlencode 'redirect_uri=/xxx/yyy' \    
  https://xxx.xxx.xxx/oauth/wechat/mp/auth_url  
```
返回：
```json 
{
    "err_code": 0, 
    "err_msg": "success", 
    "url": "https://open.weixin.qq.com/connect/oauth2/authorize?appid=appid&redirect_uri=http%3A%2F%2Flocalhost%3A8080%2Fcallback.html&response_type=code&scope=snsapi_userinfo&state=28c10d7a840cc25a988b282614b02179#wechat_redirect"
}
```
#### 认证
```http
curl -X GET \      
  -H "x-token: {{token(normally guest)}}" \   
  -G \   
  --data-urlencode 'code=XXXXXX' \   
  --data-urlencode 'state=XXXXXX' \   
  https://xxx.xxx.xxx/oauth/wechat/mp/auth  
```
返回：
```json  
{
    "err_code": 0, 
    "err_msg": "success", 
    "token": "eyJzaWQiOiJWRVJOQy1fdnR1OW1oa0pkVUVTQ0wxUG1xZWwwcVdSSSIsInRva2VuX2lkIjoiY2E5YmJhNjM2YmYxMzZiNmQxZjU4Njg3NzUyNzcwZTMiLCJhdXRoX3R5cGUiOiJwaG9uZV9wYXNzd29yZCIsImV4cF9hY2Nlc3MiOjE0NDQ3OTQ4MTMsImV4cF9yZWZyZXNoIjoxNDc2MzQ0NTY1fQ.ff5582e8680d30fc9e073ef3f4db94f152386dbe072b00fa0307818e5d8bf38e"
}
```
### 微信开放平台--网站应用
#### 获取认证url
```http
curl -X GET \     
  -H "x-token: {{token(normally guest)}}" \  
  -G \  
  --data-urlencode 'redirect_uri=/xxx/yyy' \   
  https://xxx.xxx.xxx/oauth/wechat/open/web/auth_url  
```
返回：
```json  
{
    "err_code": 0, 
    "err_msg": "success", 
    "url": "https://open.weixin.qq.com/connect/qrconnect?appid=appid&redirect_uri=http%3A%2F%2Flocalhost%3A8080%2Fcallback.html&response_type=code&scope=snsapi_login&state=5126d8ddd24ee49ef88329a1bb55aec1#wechat_redirect"
}
```
#### 认证
```http
curl -X GET \     
  -H "x-token: {{token(normally guest)}}" \  
  -G \  
  --data-urlencode 'code=XXXXXX' \  
  --data-urlencode 'state=XXXXXX' \  
  https://xxx.xxx.xxx/oauth/wechat/open/web/auth  
```
返回：
```json  
{
    "err_code": 0, 
    "err_msg": "success", 
    "token": "eyJzaWQiOiJWRVJOQy1fdnR1OW1oa0pkVUVTQ0wxUG1xZWwwcVdSSSIsInRva2VuX2lkIjoiY2E5YmJhNjM2YmYxMzZiNmQxZjU4Njg3NzUyNzcwZTMiLCJhdXRoX3R5cGUiOiJwaG9uZV9wYXNzd29yZCIsImV4cF9hY2Nlc3MiOjE0NDQ3OTQ4MTMsImV4cF9yZWZyZXNoIjoxNDc2MzQ0NTY1fQ.ff5582e8680d30fc9e073ef3f4db94f152386dbe072b00fa0307818e5d8bf38e"
}
```

### 微信开放平台--移动应用
#### 获取认证参数
```http
curl -X GET \     
  -H "x-token: {{token(normally guest)}}" \    
  https://xxx.xxx.xxx/oauth/wechat/open/app/auth_para  
```
返回：
```json  
{
    "err_code": 0, 
    "err_msg": "success", 
    "appid": "appid",
	"state": "fasdfasdfhashdflkasjdhfjlkasdhf",
	"scope": "snsapi_userinfo"
}
```
#### 认证
```http
curl -X GET \     
  -H "x-token: {{token(normally guest)}}" \  
  -G \  
  --data-urlencode 'code=XXXXXX' \  
  --data-urlencode 'state=XXXXXX' \  
  https://xxx.xxx.xxx/oauth/wechat/open/web/auth  
```
返回：
```json  
{
    "err_code": 0, 
    "err_msg": "success", 
    "token": "eyJzaWQiOiJWRVJOQy1fdnR1OW1oa0pkVUVTQ0wxUG1xZWwwcVdSSSIsInRva2VuX2lkIjoiY2E5YmJhNjM2YmYxMzZiNmQxZjU4Njg3NzUyNzcwZTMiLCJhdXRoX3R5cGUiOiJwaG9uZV9wYXNzd29yZCIsImV4cF9hY2Nlc3MiOjE0NDQ3OTQ4MTMsImV4cF9yZWZyZXNoIjoxNDc2MzQ0NTY1fQ.ff5582e8680d30fc9e073ef3f4db94f152386dbe072b00fa0307818e5d8bf38e"
}
```

## 刷新token
```http
curl -X GET \     
  -H "x-token: {{token}}" \  
  https://xxx.xxx.xxx/token/refresh  
```
返回：
```json
{
    "err_code": 0, 
    "err_msg": "success", 
    "token": "eyJzaWQiOiJ0ZW1wLlU1OU5MUV9SdHU5bWhrSmRhaWktelIzWGZmX0ZESzIzIiwidG9rZW5faWQiOiJkMGQxY2U0ZGYyYzcyYmJhMDdmNGVlMTNmZDMxMWFlNyIsImF1dGhfdHlwZSI6Imd1ZXN0IiwiZXhwX2FjY2VzcyI6MCwiZXhwX3JlZnJlc2giOjB9.a18ea69f1936b87fae1a2750ce9820a9d17b16625ac14a4dfcde7cb1955b7f9a"
}
```