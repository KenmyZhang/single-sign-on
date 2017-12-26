# resetPasswordByEmail 
    邮箱找回密码

### Url
    /sso/users/email/reset

### Method
    POST

### Request Payload
    {
        "email":"邮箱地址",
        "verification_code": "",
        "new_password":""        
    }


### Response Body
    {
        "status":"OK"
    }


### Response Code
    HTTP/1.1 200 OK


### Example
    curl -X POST http://www.example.com:8866/sso/users/email/reset -d '{"email":"1027837952@qq.com","verification_code":"733644","new_password":"Zh***g12345678"}' -i
        HTTP/1.1 200 OK
        Content-Length: 15
        Content-Type: application/json
        Date: Thu, 14 Dec 2017 09:38:01 GMT
        Keep-Alive: timeout=38
        X-Request-Id: 1qztm788rtbh3jomuhsg4jwg3e
        X-Version-Id: 4.0.0.dev.bae7f642f417866b946b24c3b5acf6fb

    {
        "status":"OK"
    }
