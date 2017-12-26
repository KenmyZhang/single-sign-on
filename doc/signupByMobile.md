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
        "id": "bnnkzm6mtfrd9exg9hsncjgr1y",
        "create_at": 1512559515149,
        "update_at": 1512559515149,
        "delete_at": 0,
        "username": "hesdfdsfotstni",
        "gender": "",
        "auth_data": "",
        "auth_service": "",
        "email": "",
        "nickname": "mynice",
        "first_name": "",
        "last_name": "",
        "position": "",
        "roles": "normal_user",
        "allow_marketing": true,
        "notify_props": {
            "channel": "true",
            "desktop": "all",
            "desktop_sound": "true",
            "email": "true",
            "first_name": "false",
            "mention_keys": "hesdfdsfotstni,@hesdfdsfotstni",
            "push": "mention"
        },
        "last_password_update": 1512559515149,
        "locale": "zh-CN",
        "mobile": "13544285684",
        "doctor_id": "",
        "region": "",
        "address": "",
        "birthDate": 0
    }

    
### Response Code
    HTTP/1.1 201 Created

### Example
    curl -X POST "http://127.0.0.1:8065/sso/users/phone/signup"  -i -d '{"username":"hesdfdsfotstni","nickname":"mynice","password":"12345678","verification_code": "666666","mobile":"13544285684"}'  -b  ~/go/src/github.com/mattermost/example/bin/cookie-file
        HTTP/1.1 201 Created
        Content-Length: 600
        Token: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiYXQ1OWU3MWlndG5qcHB1ZGVnamRuam1pMXciLCJyb2xlcyI6Im5vcm1hbF91c2VyIiwicHJvcHMiOnsiYnJvd3NlciI6IkdvLWh0dHAtY2xpZW50LzEuMSIsIm9zIjoidW5rbm93biIsInBsYXRmb3JtIjoidW5rbm93biJ9LCJ0ZWFtX21lbWJlcnMiOm51bGwsImRldmljZV9pZCI6IiIsImlzX29hdXRoIjpmYWxzZSwiZXhwIjoxNTE0Mjg4MTkxLCJpYXQiOjE1MTI1NjAxOTEsImlzcyI6Ind3dy5hY2N1cm1lLmNvbSJ9.YdMwEqOK_61AeJtlVp-9f42C9Wy70slnL-1Uq5c7cq4
        Content-Type: application/json
        Date: Wed, 06 Dec 2017 11:36:31 GMT
        Keep-Alive: timeout=38
        X-Request-Id: dnqeow5syfrcucpisph1jyd3bh
        X-Version-Id: 4.0.0.dev.46e172c54ca01cf33a0f77d1b099c013

    {
        "id": "bnnkzm6mtfrd9exg9hsncjgr1y",
        "create_at": 1512559515149,
        "update_at": 1512559515149,
        "delete_at": 0,
        "username": "hesdfdsfotstni",
        "gender": "",
        "auth_data": "",
        "auth_service": "",
        "email": "",
        "nickname": "mynice",
        "first_name": "",
        "last_name": "",
        "position": "",
        "roles": "normal_user",
        "allow_marketing": true,
        "notify_props": {
            "channel": "true",
            "desktop": "all",
            "desktop_sound": "true",
            "email": "true",
            "first_name": "false",
            "mention_keys": "hesdfdsfotstni,@hesdfdsfotstni",
            "push": "mention"
        },
        "last_password_update": 1512559515149,
        "locale": "zh-CN",
        "mobile": "13544285684",
        "doctor_id": "",
        "region": "",
        "address": "",
        "birthDate": 0
    }
