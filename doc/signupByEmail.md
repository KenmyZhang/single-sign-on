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
      "id": "gi6z9anrf7y9tqapirfpofs3uy",
      "create_at": 1512562562530,
      "update_at": 1512562562530,
      "delete_at": 0,
      "username": "heotestfdsaf",
      "gender": "",
      "auth_data": "",
      "auth_service": "",
      "email": "1027837952@qq.com",
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
          "mention_keys": "heotestfdsaf,@heotestfdsaf",
          "push": "mention"
      },
      "last_password_update": 1512562562530,
      "locale": "zh-CN",
      "mobile": "",
      "doctor_id": "",
      "region": "",
      "address": "",
      "birthDate": 0
  }

### Response Code
    201 Created

### Example 
	curl -X POST  "http://www.example.com:8065/sso/users/email/signup"  -i -d '{"username":"heotestfdsaf","nickname":"mynice","password":"12345678","verification_code": "206495","email":"1027837952@qq.com"}'
    HTTP/1.1 201 Created
    Content-Length: 600
    Token: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiYXQ1OWU3MWlndG5qcHB1ZGVnamRuam1pMXciLCJyb2xlcyI6Im5vcm1hbF91c2VyIiwicHJvcHMiOnsiYnJvd3NlciI6IkdvLWh0dHAtY2xpZW50LzEuMSIsIm9zIjoidW5rbm93biIsInBsYXRmb3JtIjoidW5rbm93biJ9LCJ0ZWFtX21lbWJlcnMiOm51bGwsImRldmljZV9pZCI6IiIsImlzX29hdXRoIjpmYWxzZSwiZXhwIjoxNTE0Mjg4MTkxLCJpYXQiOjE1MTI1NjAxOTEsImlzcyI6Ind3dy5hY2N1cm1lLmNvbSJ9.YdMwEqOK_61AeJtlVp-9f42C9Wy70slnL-1Uq5c7cq4
    Content-Type: application/json
    Date: Wed, 06 Dec 2017 11:36:31 GMT
    Keep-Alive: timeout=38
    X-Request-Id: dnqeow5syfrcucpisph1jyd3bh
    X-Version-Id: 4.0.0.dev.46e172c54ca01cf33a0f77d1b099c013

  {
      "id": "gi6z9anrf7y9tqapirfpofs3uy",
      "create_at": 1512562562530,
      "update_at": 1512562562530,
      "delete_at": 0,
      "username": "heotestfdsaf",
      "gender": "",
      "auth_data": "",
      "auth_service": "",
      "email": "1027837952@qq.com",
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
          "mention_keys": "heotestfdsaf,@heotestfdsaf",
          "push": "mention"
      },
      "last_password_update": 1512562562530,
      "locale": "zh-CN",
      "mobile": "",
      "doctor_id": "",
      "region": "",
      "address": "",
      "birthDate": 0
  }
