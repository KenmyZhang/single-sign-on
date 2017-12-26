# signupByEmail 
    

### Url
    /sso/users/email/signup

### Method
    POST

### Request Payload
	{
		"username":"heotest",
		"nickname":"mynice",
		"password":"12345678",
		"verification_code": "164285",
		"email":"1027837952@qq.com"
	}


### Response Body
    {
        "id": "ish14g5zftgftnqruubw3qmt1a",
        "create_at": 1514281415487,
        "update_at": 1514281415487,
        "delete_at": 0,
        "username": "heotf",
        "gender": "",
        "auth_service": "",
        "email": "1027837952@qq.com",
        "nickname": "mynice",
        "first_name": "",
        "last_name": "",
        "position": "",
        "roles": "system_admin",
        "allow_marketing": true,
        "locale": "zh-CN",
        "mobile": ""
    }

### Response Code
    201 Created

### Example 
  curl -X POST  "http://www.example.com:9968/sso/users/email/signup"  -i -d '{"username":"heotf","nickname":"mynice","password":"Zhang12****5678","verification_code": "491969","email":"1027837952@qq.com"}'
    HTTP/1.1 201 Created
    Content-Type: application/json
    Set-Cookie: AUTHTOKEN=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiaXNoMTRnNXpmdGdmdG5xcnV1YnczcW10MWEiLCJyb2xlcyI6InN5c3RlbV9hZG1pbiIsInByb3BzIjp7ImJyb3dzZXIiOiJjdXJsLzcuNDcuMCIsIm9zIjoidW5rbm93biIsInBsYXRmb3JtIjoidW5rbm93biJ9LCJleHAiOjE1MTYwMDk0MTUsImlhdCI6MTUxNDI4MTQxNSwiaXNzIjoid3d3LmFjY3VybWUuY29tIn0.CxPnCTLvJFi9xXZjK_eysmXa2SREgxm49d7i16Y6xSA; Path=/; Expires=Mon, 15 Jan 2018 09:43:35 GMT; Max-Age=1728000; HttpOnly
    Set-Cookie: USERID=ish14g5zftgftnqruubw3qmt1a; Path=/; Expires=Mon, 15 Jan 2018 09:43:35 GMT; Max-Age=1728000
    Token: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiaXNoMTRnNXpmdGdmdG5xcnV1YnczcW10MWEiLCJyb2xlcyI6InN5c3RlbV9hZG1pbiIsInByb3BzIjp7ImJyb3dzZXIiOiJjdXJsLzcuNDcuMCIsIm9zIjoidW5rbm93biIsInBsYXRmb3JtIjoidW5rbm93biJ9LCJleHAiOjE1MTYwMDk0MTUsImlhdCI6MTUxNDI4MTQxNSwiaXNzIjoid3d3LmFjY3VybWUuY29tIn0.CxPnCTLvJFi9xXZjK_eysmXa2SREgxm49d7i16Y6xSA
    X-Request-Id: 9coeyobq4fg33eccfesyynnznh
    X-Version-Id: 4.0.0.dev.6b8de684aa2e0d9c94b645e75c772902
    Date: Tue, 26 Dec 2017 09:43:35 GMT
    Content-Length: 318

    {
        "id": "ish14g5zftgftnqruubw3qmt1a",
        "create_at": 1514281415487,
        "update_at": 1514281415487,
        "delete_at": 0,
        "username": "heotf",
        "gender": "",
        "auth_service": "",
        "email": "1027837952@qq.com",
        "nickname": "mynice",
        "first_name": "",
        "last_name": "",
        "position": "",
        "roles": "system_admin",
        "allow_marketing": true,
        "locale": "zh-CN",
        "mobile": ""
    }