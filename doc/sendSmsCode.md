# sendSmsCode
    发送手机验证码
    
### Url
    /sso/users/sendsms
    
### Method
    POST

### Request Payload
    {
        "mobile":"手机号码"
    }
    
### Response Body
    {"status":"OK"}
    
### Response Code
    HTTP/1.1 200 OK

### Example
    curl -X POST  "http://127.0.0.1:8065/sso/users/sendsms"  -d '{"mobile":"13544285662"}'  -D cookie-file -i
        HTTP/1.1 200 OK
        Content-Length: 15
        Content-Type: application/json
        Date: Wed, 06 Dec 2017 08:50:48 GMT
        Keep-Alive: timeout=38
        X-Request-Id: gspaaa38zidzucciep6rgtyxao
        X-Version-Id: 4.0.0.dev.b88ebfe669ef663bb26f25ac85d6bf0d

    {
        "status":"OK"
    }
