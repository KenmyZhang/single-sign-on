# isEmailExist
    判断邮箱是否已经注册过
    
### Url
     /sso/users/email/exist
    
### Method
    POST

### Request Payload
    {
        "email":"2224052@qq.com"
    }

### Response Body
    {
        "status":"true"
    }
    
### Response Code
    HTTP/1.1 200 OK

### Example
     curl -X POST "http://127.0.0.1:8065/sso/users/email/exist" -d '{"email":"zhanhf@qq.com"}'  -i
        HTTP/1.1 200 OK
        Content-Length: 17
        Content-Type: application/json
        Date: Thu, 07 Dec 2017 07:01:47 GMT
        Keep-Alive: timeout=38
        X-Request-Id: o85ef6pq3fde7bx7hi3rkahhbr
        X-Version-Id: 4.0.0.dev.463aa9e9f0c1d9e0d9e24172a4bde3d8

    {
        "status":"true"
    }
