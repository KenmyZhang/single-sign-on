### 登录
  [手机号、邮箱、用户名 + 密码 ](http://120.24.62.67:10080/zhangkunming/sso-doc/src/master/login.md)

  [手机号 + 验证码  ](http://120.24.62.67:10080/zhangkunming/sso-doc/src/master/loginByMobile.md)

  微信登录

### 注册
  [判断手机是否已注册 ](http://120.24.62.67:10080/zhangkunming/sso-doc/src/master/isMobileExist.md)

  [判断邮箱是否已注册 ](http://120.24.62.67:10080/zhangkunming/sso-doc/src/master/isEmailExist.md)

  [发送手机短信验证码 ](http://120.24.62.67:10080/zhangkunming/sso-doc/src/master/sendSmsCode.md)

  [发送邮箱验证码 ](http://120.24.62.67:10080/zhangkunming/sso-doc/src/master/sendVerificationCodeEmail.md)

  [手机号码注册 ](http://120.24.62.67:10080/zhangkunming/sso-doc/src/master/signupByMobile.md)

  [邮箱注册 ](http://120.24.62.67:10080/zhangkunming/sso-doc/src/master/signupByEmail.md)

### 忘记密码 
  [邮件找回 ](http://120.24.62.67:10080/zhangkunming/sso-doc/src/master/resetPasswordByEmail.md)
  
  [手机找回 ](http://120.24.62.67:10080/zhangkunming/sso-doc/src/master/resetPasswordByMobile.md)


### 约束
  短信验证码有效期一分钟（MAX_SMS_TOKEN_EXIPRY_TIME），一分钟内只能发送一次验证码，24小时内只能发送60条短信（SEND_CODE_MAX，MAX_TOKEN_EXIPRY_TIME），以防止恶意用户
  
  邮件验证码有效期一分钟 (MAX_EMAIL_TOKEN_EXIPRY_TIME）,一分钟内只能发送一次验证码

  密码长度最小长度5，最大长度72（USER_PASSWORD_MAX_LENGTH）,必须包含大写字母、小写字母、数字
  
  昵称长度小于64（USER_NICKNAME_MAX_RUNES）
  
  邮箱长度小于128（USER_EMAIL_MAX_LENGTH）
  
  用户名长度3 ～ 64（USER_NAME_MIN_LENGTH、USER_NAME_MAX_LENGTH）
  
  用户名必须以字母开头,并且包含3到22个小写字母, '.', '-'和'_'.	
