# single-sign-on

温馨提示：程序运行起来后，会判断是否存在数据表，如果不存在会自动创建，无需另外另外执行sql脚本来创建数据表(因比较多人找我要sql脚本，所以在这里提醒一下，我这没有sql脚本)

## 安装数据库（基于ubuntu）
* 利用apt-get install安装MySQl
 
    sudo apt-get install mysql-server

* 以root用户登录MySQL
  
  mysql -u root -p

* 创建sso用户'ssouser'
  
  mysql> create user 'ssouser'@'%' identified by 'ssouser-password'; 
   其中%表示网上的所有机器都可以连接上，使用具体的IP地址更安全点
  mysql> create user 'ssouser'@'10.10.10.2' identified by 'ssouser-password';


* 创建sso数据库

  mysql> create database sso


* 允许ssouser用户的访问权限

  mysql> grant all privileges on sso.* to 'ssouser'@'%';


* 退出MySQL

  mysql> exit


## 编译
  make build

### 打包 
  make package

### 运行
  ./single-sing-on



## api document
#### login
  [手机号、邮箱、用户名 + 密码 ](https://github.com/KenmyZhang/single-sign-on/blob/master/doc/login.md)

  [手机号 + 验证码  ](https://github.com/KenmyZhang/single-sign-on/blob/master/doc/loginByMobile.md)

  微信登录

#### sign up
  [判断手机是否已注册 ](https://github.com/KenmyZhang/single-sign-on/blob/master/doc/isMobileExist.md)

  [判断邮箱是否已注册 ](https://github.com/KenmyZhang/single-sign-on/blob/master/doc/isEmailExist.md)

  [发送手机短信验证码 ](https://github.com/KenmyZhang/single-sign-on/blob/master/sso-doc/sendSmsCode.md)

  [发送邮箱验证码 ](https://github.com/KenmyZhang/single-sign-on/blob/master/sso-doc/sendVerificationCodeEmail.md)

  [手机号码注册 ](https://github.com/KenmyZhang/single-sign-on/blob/master/doc/signupByMobile.md)

  [邮箱注册 ](https://github.com/KenmyZhang/single-sign-on/blob/master/doc/signupByEmail.md)

#### forget password 
  [邮件找回 ](https://github.com/KenmyZhang/single-sign-on/blob/master/doc/resetPasswordByEmail.md)
  
  [手机找回 ](https://github.com/KenmyZhang/single-sign-on/blob/master/doc/resetPasswordByMobile.md)


#### constraint
  短信验证码有效期一分钟（MAX_SMS_TOKEN_EXIPRY_TIME），一分钟内只能发送一次验证码，24小时内只能发送60条短信（SEND_CODE_MAX，MAX_TOKEN_EXIPRY_TIME），以防止恶意用户
  
  邮件验证码有效期一分钟 (MAX_EMAIL_TOKEN_EXIPRY_TIME）,一分钟内只能发送一次验证码

  密码长度最小长度5，最大长度72（USER_PASSWORD_MAX_LENGTH）,必须包含大写字母、小写字母、数字
  
  昵称长度小于64（USER_NICKNAME_MAX_RUNES）
  
  邮箱长度小于128（USER_EMAIL_MAX_LENGTH）
  
  用户名长度3 ～ 64（USER_NAME_MIN_LENGTH、USER_NAME_MAX_LENGTH）
  
  用户名必须以字母开头,并且包含3到22个小写字母, '.', '-'和'_'.   
  
### Contact
  kenmyzhang@gmail.com
