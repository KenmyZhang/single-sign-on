# isMobileExist
    判断手机号是否已经注册过
    
### Url
    /sso/users/phone/exist
    
### Method
    POST

### Request Payload
    {
        "mobile":"13544285663"
    }

### Response Body
    {
        "status":"true"
    }
    
### Response Code
    HTTP/1.1 200 OK

### Example(存在)
    curl -X POST  "http://127.0.0.1:8065/sso/users/phone/exist"  -i -d '{"mobile":"13544285663"}'
        HTTP/1.1 200 OK
        Content-Length: 17
        Content-Type: application/json
        Date: Thu, 07 Dec 2017 03:34:12 GMT
        Keep-Alive: timeout=38
        X-Request-Id: ka16q5isi7rs98yqytsjzk3ype
        X-Version-Id: 4.0.0.dev.0d1bdaabb8272d8016bd57816dad483f

    {
        "status":"true"
    }
