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
        "id": "z4u4s6ohz7ng9xe6azk3fbpqcw",
        "create_at": 1510639944500,
        "update_at": 1512527813516,
        "delete_at": 0,
        "username": "kenmy",
        "gender": "",
        "auth_data": "",
        "auth_service": "",
        "email": "",
        "email_verified": true,
        "nickname": "景雅曼",
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
            "mention_keys": "kenmy,@kenmy",
            "push": "mention"
        },
        "last_password_update": 1512527813516,
        "locale": "en",
        "mobile": "13544285662",
        "doctor_id": "59f6deef3f92a31b5752bc06",
        "region": "",
        "address": "",
        "birthDate": 0
    }
    
  
### Response Code
     HTTP/1.1 200 OK

### Example 手机+密码登录
    curl -X POST  "http://127.0.0.1:8065/sso/users/phone/login"  -i -d '{"mobile":"13544285662","password":"4**********011"}'
        HTTP/1.1 200 OK
        Content-Length: 619
        Authorization: BEARER eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiejR1NHM2b2h6N25nOXhlNmF6azNmYnBxY3ciLCJyb2xlcyI6Im5vcm1hbF91c2VyIiwicHJvcHMiOnsiYnJvd3NlciI6IkdvLWh0dHAtY2xpZW50LzEuMSIsIm9zIjoidW5rbm93biIsInBsYXRmb3JtIjoidW5rbm93biJ9LCJ0ZWFtX21lbWJlcnMiOm51bGwsImRldmljZV9pZCI6IiIsImlzX29hdXRoIjpmYWxzZSwiZXhwIjoxNTE0MzQ0MjU2LCJpYXQiOjE1MTI2MTYyNTYsImlzcyI6Ind3dy5hY2N1cm1lLmNvbSJ9.NcFz9hOv02JqJf3FqVupNXY2f-Def4qyo3cbNo1kV6E
        Content-Type: application/json
        Date: Thu, 07 Dec 2017 03:10:56 GMT
        Keep-Alive: timeout=38
        X-Request-Id: dz3othxrr7ysbcjk9jkxhef4bc
        X-Version-Id: 4.0.0.dev.8033da0d6ec9284484375ec6d73a4611

    {
        "id": "z4u4s6ohz7ng9xe6azk3fbpqcw",
        "create_at": 1510639944500,
        "update_at": 1512527813516,
        "delete_at": 0,
        "username": "kenmy",
        "gender": "",
        "auth_data": "",
        "auth_service": "",
        "email": "",
        "email_verified": true,
        "nickname": "景雅曼",
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
            "mention_keys": "kenmy,@kenmy",
            "push": "mention"
        },
        "last_password_update": 1512527813516,
        "locale": "en",
        "mobile": "13544285662",
        "doctor_id": "59f6deef3f92a31b5752bc06",
        "region": "",
        "address": "",
        "birthDate": 0
    }
### Example 手机 + 验证码
    curl -X POST  "http://127.0.0.1:8065/sso/users/phone/login"  -i -d '{"mobile":"13544285662","verification_code":"666666"}'
        HTTP/1.1 200 OK
        Content-Length: 619
        Authorization: BEARER eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiejR1NHM2b2h6N25nOXhlNmF6azNmYnBxY3ciLCJyb2xlcyI6Im5vcm1hbF91c2VyIiwicHJvcHMiOnsiYnJvd3NlciI6IkdvLWh0dHAtY2xpZW50LzEuMSIsIm9zIjoidW5rbm93biIsInBsYXRmb3JtIjoidW5rbm93biJ9LCJ0ZWFtX21lbWJlcnMiOm51bGwsImRldmljZV9pZCI6IiIsImlzX29hdXRoIjpmYWxzZSwiZXhwIjoxNTE0MzQ0NTQ3LCJpYXQiOjE1MTI2MTY1NDcsImlzcyI6Ind3dy5hY2N1cm1lLmNvbSJ9.yDE3N_Lv5yNrbRBXPPC7KYTd0HOEWbeWXG9qN2e22-Y
        Content-Type: application/json
        Date: Thu, 07 Dec 2017 03:15:47 GMT
        Keep-Alive: timeout=38
        X-Request-Id: hnszyrtiu3dm3je9qwgro3kchc
        X-Version-Id: 4.0.0.dev.8033da0d6ec9284484375ec6d73a4611

    {
        "id": "z4u4s6ohz7ng9xe6azk3fbpqcw",
        "create_at": 1510639944500,
        "update_at": 1512527813516,
        "delete_at": 0,
        "username": "kenmy",
        "gender": "",
        "auth_data": "",
        "auth_service": "",
        "email": "",
        "email_verified": true,
        "nickname": "景雅曼",
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
            "mention_keys": "kenmy,@kenmy",
            "push": "mention"
        },
        "last_password_update": 1512527813516,
        "locale": "en",
        "mobile": "13544285662",
        "doctor_id": "59f6deef3f92a31b5752bc06",
        "region": "",
        "address": "",
        "birthDate": 0
    }
