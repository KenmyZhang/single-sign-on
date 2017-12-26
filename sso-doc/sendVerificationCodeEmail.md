# sendVerificationCodeEmail
    发送邮件验证
    
### Url
    /sso/users/email/verify/code/send
    
### Method
    POST

### Request Payload
    {
        email: "邮箱地址"
    }
    
### Response Body
    {
        "status":"OK"
    }
    
### Response Code
    HTTP/1.1 200 OK

### Example
    curl -X POST "http://127.0.0.1:8065/sso/users/email/verify/code/send" -d '{"email":"1027837952@qq.com"}' -i
        HTTP/1.1 200 OK
        Content-Length: 15
        Content-Type: application/json
        Date: Wed, 06 Dec 2017 09:53:40 GMT
        Keep-Alive: timeout=38
        X-Request-Id: iuhi5ybfk7n8jmyooof5u56yrh
        X-Version-Id: 4.0.0.dev.298c10b2fa4782ea6ac753859868a67f


    {
        "status":"OK"
    }
