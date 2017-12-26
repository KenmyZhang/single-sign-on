# resetPasswordByMobile 
    手机号找回密码

### Url
    /sso/users/phone/reset

### Method
    POST

### Request Payload
    {
        "mobile":"手机号码",
        "verification_code":"验证码",
        "new_password":"新密码"
    }


### Response Body
    {
        "status":"OK"
    }

### Response Code
    HTTP/1.1 200 OK


### Example
	curl -X POST  "http://127.0.0.1:8065/sso/users/phone/reset"  -i -d '{"mobile":"13544285662","verification_code":"666666","new_password":"4538***************1"}'
    HTTP/1.1 200 OK
    Content-Length: 15
    Content-Type: application/json
    Date: Thu, 07 Dec 2017 07:35:39 GMT
    Keep-Alive: timeout=38
    X-Request-Id: zufyhr19mjdytfkbumgc957ijh
    X-Version-Id: 4.0.0.dev.90a432f1facf65e607008acba4237acd

	{
		"status":"OK"
	}
