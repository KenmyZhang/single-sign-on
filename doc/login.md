# login
    
### Url
    /sso/users/login
    
### Method
    POST

### Request Payload
    {
        "login_id":"2224052849@qq.com",
        "password":"Zhang4***01"
    }

### Response Body
    {
        "id": "kyj99bt16ifimb973cbyegcs5o",
        "create_at": 1513171527086,
        "update_at": 1513171745524,
        "delete_at": 0,
        "username": "heskenmmy",
        "gender": "",
        "auth_service": "",
        "email": "2224052849@qq.com",
        "email_verified": true,
        "nickname": "lemnemynice",
        "first_name": "",
        "last_name": "",
        "position": "",
        "roles": "normal_user",
        "allow_marketing": true,
        "last_password_update": 1513171745524,
        "locale": "zh-CN",
        "mobile": "13544285662",
        "doctor_id": "",
        "region": "",
        "address": "",
        "birthDate": 0
    }
    
### Response Code
    HTTP/1.1 200 OK

### Example
    curl -X POST  "http://127.0.0.1:9966/sso/users/login"  -i -d '{"login_id":"2224052849@qq.com","password":"Zha****801"}'
        HTTP/1.1 200 OK
        Content-Length: 450
        Content-Type: application/json
        Date: Wed, 13 Dec 2017 13:51:51 GMT
        Keep-Alive: timeout=38
        Set-Cookie: AUTHTOKEN=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoia3lqOTlidDE2aWZpbWI5NzNjYnllZ2NzNW8iLCJyb2xlcyI6Im5vcm1hbF91c2VyIiwicHJvcHMiOnsiYnJvd3NlciI6ImN1cmwvNy40Ny4wIiwib3MiOiJ1bmtub3duIiwicGxhdGZvcm0iOiJ1bmtub3duIn0sImV4cCI6MTUxNDkwMTExMSwiaWF0IjoxNTEzMTczMTExLCJpc3MiOiJ3d3cuYWNjdXJtZS5jb20ifQ.lRfLBm9V7mwrrBrsbJA5P-kkJzmzu8z3U4GvoZdfL0E; Path=/; Expires=Tue, 02 Jan 2018 13:51:51 GMT; Max-Age=1728000; HttpOnly
        Set-Cookie: USERID=kyj99bt16ifimb973cbyegcs5o; Path=/; Expires=Tue, 02 Jan 2018 13:51:51 GMT; Max-Age=1728000
        Token: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoia3lqOTlidDE2aWZpbWI5NzNjYnllZ2NzNW8iLCJyb2xlcyI6Im5vcm1hbF91c2VyIiwicHJvcHMiOnsiYnJvd3NlciI6ImN1cmwvNy40Ny4wIiwib3MiOiJ1bmtub3duIiwicGxhdGZvcm0iOiJ1bmtub3duIn0sImV4cCI6MTUxNDkwMTExMSwiaWF0IjoxNTEzMTczMTExLCJpc3MiOiJ3d3cuYWNjdXJtZS5jb20ifQ.lRfLBm9V7mwrrBrsbJA5P-kkJzmzu8z3U4GvoZdfL0E
        X-Request-Id: u43hizkt6tfn38b9uq7qeh4bbe
        X-Version-Id: 4.0.0.dev.951781304cb60489c84a43745aae1df8
    {
        "id": "kyj99bt16ifimb973cbyegcs5o",
        "create_at": 1513171527086,
        "update_at": 1513171745524,
        "delete_at": 0,
        "username": "heskenmmy",
        "gender": "",
        "auth_service": "",
        "email": "2224052849@qq.com",
        "email_verified": true,
        "nickname": "lemnemynice",
        "first_name": "",
        "last_name": "",
        "position": "",
        "roles": "normal_user",
        "allow_marketing": true,
        "last_password_update": 1513171745524,
        "locale": "zh-CN",
        "mobile": "13544285662",
        "doctor_id": "",
        "region": "",
        "address": "",
        "birthDate": 0
    }
