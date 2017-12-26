# login
    login by mobile or email or username
### Url
    /sso/users/login
    
### Method
    POST

### Request Payload
    {
        "login_id":"2224052849@qq.com",
        "password":"Zht*****3t401"
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
        "mobile": "135***662"
    }
    
### Response Code
    HTTP/1.1 200 OK

### Example
    curl -X POST  "http://127.0.0.1:9966/sso/users/login"  -i -d '{"login_id":"13544285662","password":"Zhtreter****45678"}'
        HTTP/1.1 200 OK
        Content-Type: application/json
        Set-Cookie: AUTHTOKEN=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiNzR6YXQxenQ5N2I5dXlobTdiZmJ1cHJiamUiLCJyb2xlcyI6InN5c3RlbV9hZG1pbiIsInByb3BzIjp7ImJyb3dzZXIiOiJjdXJsLzcuNDcuMCIsIm9zIjoidW5rbm93biIsInBsYXRmb3JtIjoidW5rbm93biJ9LCJleHAiOjE1MTYwMDY3OTAsImlhdCI6MTUxNDI3ODc5MCwiaXNzIjoid3d3LmFjY3VybWUuY29tIn0.si7qPIYV4n3xRzyDgeVWRFCiHnfkkwGbRsFJQpSlF_o; Path=/; Expires=Mon, 15 Jan 2018 08:59:50 GMT; Max-Age=1728000; HttpOnly
        Set-Cookie: USERID=74zat1zt97b9uyhm7bfbuprbje; Path=/; Expires=Mon, 15 Jan 2018 08:59:50 GMT; Max-Age=1728000
        Token: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiNzR6YXQxenQ5N2I5dXlobTdiZmJ1cHJiamUiLCJyb2xlcyI6InN5c3RlbV9hZG1pbiIsInByb3BzIjp7ImJyb3dzZXIiOiJjdXJsLzcuNDcuMCIsIm9zIjoidW5rbm93biIsInBsYXRmb3JtIjoidW5rbm93biJ9LCJleHAiOjE1MTYwMDY3OTAsImlhdCI6MTUxNDI3ODc5MCwiaXNzIjoid3d3LmFjY3VybWUuY29tIn0.si7qPIYV4n3xRzyDgeVWRFCiHnfkkwGbRsFJQpSlF_o
        X-Request-Id: s1h5dofupt8kpyz1i1rayebngc
        X-Version-Id: 4.0.0.dev.f3384cff7a166fbeb1a6d90426f62fbd
        Date: Tue, 26 Dec 2017 08:59:50 GMT
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
        "mobile": "13566*****2"
    }
