# loginByMobile
    手机验证码或者密码登录
    
### Url
    /sso/users/phone/login
    
### Method
    POST

### Request Payload
    {
        "mobile":"手机号码",
        "verification_code":"验证码"，
        "password":"密码"
    }

### Response Body
    {
        "id": "74zat1zt97b9uyhm7bfbuprbje",
        "create_at": 1514278513093,
        "update_at": 1514278513093,
        "delete_at": 0,
        "username": "hesdfdsfotstni",
        "gender": "",
        "auth_service": "",
        "email": "2224052849@qq.com",
        "nickname": "mynice",
        "first_name": "",
        "last_name": "",
        "position": "",
        "roles": "system_admin",
        "allow_marketing": true,
        "locale": "zh-CN",
        "mobile": "13544285662"
    }
    
  
### Response Code
     HTTP/1.1 200 OK

### Example 手机 + 验证码
    curl -X POST  "http://127.0.0.1:9966/sso/users/phone/login"  -i -d '{"mobile":"13544285662","verification_code":"666666"}'
        HTTP/1.1 200 OK
        Content-Type: application/json
        Set-Cookie: AUTHTOKEN=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiNzR6YXQxenQ5N2I5dXlobTdiZmJ1cHJiamUiLCJyb2xlcyI6InN5c3RlbV9hZG1pbiIsInByb3BzIjp7ImJyb3dzZXIiOiJjdXJsLzcuNDcuMCIsIm9zIjoidW5rbm93biIsInBsYXRmb3JtIjoidW5rbm93biJ9LCJleHAiOjE1MTYwMDcwNTAsImlhdCI6MTUxNDI3OTA1MCwiaXNzIjoid3d3LmFjY3VybWUuY29tIn0.LnNXdBeeuGlal697KOTdbd6ebvCr0NQKxgoasFPnvLs; Path=/; Expires=Mon, 15 Jan 2018 09:04:10 GMT; Max-Age=1728000; HttpOnly
        Set-Cookie: USERID=74zat1zt97b9uyhm7bfbuprbje; Path=/; Expires=Mon, 15 Jan 2018 09:04:10 GMT; Max-Age=1728000
        Token: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiNzR6YXQxenQ5N2I5dXlobTdiZmJ1cHJiamUiLCJyb2xlcyI6InN5c3RlbV9hZG1pbiIsInByb3BzIjp7ImJyb3dzZXIiOiJjdXJsLzcuNDcuMCIsIm9zIjoidW5rbm93biIsInBsYXRmb3JtIjoidW5rbm93biJ9LCJleHAiOjE1MTYwMDcwNTAsImlhdCI6MTUxNDI3OTA1MCwiaXNzIjoid3d3LmFjY3VybWUuY29tIn0.LnNXdBeeuGlal697KOTdbd6ebvCr0NQKxgoasFPnvLs
        X-Request-Id: iy3h7p6z33y3bfr6de9a61rgsc
        X-Version-Id: 4.0.0.dev.f3384cff7a166fbeb1a6d90426f62fbd
        Date: Tue, 26 Dec 2017 09:04:10 GMT
        Content-Length: 338

    {
        "id": "74zat1zt97b9uyhm7bfbuprbje",
        "create_at": 1514278513093,
        "update_at": 1514278513093,
        "delete_at": 0,
        "username": "hesdfdsfotstni",
        "gender": "",
        "auth_service": "",
        "email": "2224052849@qq.com",
        "nickname": "mynice",
        "first_name": "",
        "last_name": "",
        "position": "",
        "roles": "system_admin",
        "allow_marketing": true,
        "locale": "zh-CN",
        "mobile": "13544285662"
    }
