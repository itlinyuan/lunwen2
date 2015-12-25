#GoBlog 

基于Go语言和beego框架的简易博客系统

##编译安装说明：

修改数据库配置
	
	$ cd src
	$ vim github.com/lisijie/goblog/conf/app.conf
	
	appname = goblog
	httpport = 8080
	runmode = dev
	dbhost = localhost 
	dbport = 3306
	dbuser = root
	dbpassword = 123456
	dbname = goblog
	dbprefix = t_

导入MySQL

访问： 

http://localhost:8080

后台地址：

http://localhost:8080/admin

帐号：admin
密码：admin888

