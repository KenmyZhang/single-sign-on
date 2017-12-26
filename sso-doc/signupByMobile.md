# signupByMobile
    校验手机验证码注册用户
    
### Url
    /sso/users/phone/signup
    
### Method
    POST

### Request Payload
    {
        "username": "用户名",
        "password": "密码",
        "nickname": "昵称",                    //optional 选填
        "verification_code":"手机验证码",
        "mobile": "手机号码"
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
        "email": "",
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
    HTTP/1.1 201 Created

### Example
    curl -X POST "http://127.0.0.1:9966/sso/users/phone/signup"  -i -d '{"username":"hesdfdsfotstni","nickname":"mynice","password":"Zhang12345678","verification_code": "666666","mobile":"13544285662"}'
        HTTP/1.1 201 Created
        Content-Type: application/json
        Set-Cookie: AUTHTOKEN=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiNzR6YXQxenQ5N2I5dXlobTdiZmJ1cHJiamUiLCJyb2xlcyI6InN5c3RlbV9hZG1pbiIsInByb3BzIjp7ImJyb3dzZXIiOiJjdXJsLzcuNDcuMCIsIm9zIjoidW5rbm93biIsInBsYXRmb3JtIjoidW5rbm93biJ9LCJleHAiOjE1MTYwMDY1MTMsImlhdCI6MTUxNDI3ODUxMywiaXNzIjoid3d3LmFjY3VybWUuY29tIn0.MWQpddSiOMUYJycaSDn0feB48aQ248LEI-VKPAJJ714; Path=/; Expires=Mon, 15 Jan 2018 08:55:13 GMT; Max-Age=1728000; HttpOnly
        Set-Cookie: USERID=74zat1zt97b9uyhm7bfbuprbje; Path=/; Expires=Mon, 15 Jan 2018 08:55:13 GMT; Max-Age=1728000
        Token: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiNzR6YXQxenQ5N2I5dXlobTdiZmJ1cHJiamUiLCJyb2xlcyI6InN5c3RlbV9hZG1pbiIsInByb3BzIjp7ImJyb3dzZXIiOiJjdXJsLzcuNDcuMCIsIm9zIjoidW5rbm93biIsInBsYXRmb3JtIjoidW5rbm93biJ9LCJleHAiOjE1MTYwMDY1MTMsImlhdCI6MTUxNDI3ODUxMywiaXNzIjoid3d3LmFjY3VybWUuY29tIn0.MWQpddSiOMUYJycaSDn0feB48aQ248LEI-VKPAJJ714
        X-Request-Id: utgtu9xs5jyomnsq3jw4qnga8r
        X-Version-Id: 4.0.0.dev.f3384cff7a166fbeb1a6d90426f62fbd
        Date: Tue, 26 Dec 2017 08:55:13 GMT
        Content-Length: 321

    {
        "id": "74zat1zt97b9uyhm7bfbuprbje",
        "create_at": 1514278513093,
        "update_at": 1514278513093,
        "delete_at": 0,
        "username": "hesdfdsfotstni",
        "gender": "",
        "auth_service": "",
        "email": "",
        "nickname": "mynice",
        "first_name": "",
        "last_name": "",
        "position": "",
        "roles": "system_admin",
        "allow_marketing": true,
        "locale": "zh-CN",
        "mobile": "13544285662"
    }
