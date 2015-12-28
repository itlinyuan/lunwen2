package admin

import (
	"github.com/lisijie/goblog/models"
	"github.com/lisijie/goblog/util"
	"regexp"
	"strconv"
	"strings"
)

type AccountController struct {
	baseController
}

//登录
func (this *AccountController) Login() {
	//>0:已经登录了，重定向到用户管理界面
	if this.userid > 0 {
		this.Redirect("/admin", 302)
	}

	if this.GetString("dosubmit") == "yes" {
		account := strings.TrimSpace(this.GetString("account"))
		password := strings.TrimSpace(this.GetString("password"))
		remember := this.GetString("remember")
		if account != "" && password != "" {
			var user models.User
			user.UserName = account
			if user.Read("user_name") != nil || user.Password != util.Md5([]byte(password)) {
				this.Data["errmsg"] = "帐号或密码错误"
			} else if user.Active == 0 {
				this.Data["errmsg"] = "该帐号未激活"
			} else {
				user.LoginCount += 1
				user.LastIp = this.getClientIp()
				user.LastLogin = this.getTime()
				user.Update()
				authkey := util.Md5([]byte(this.getClientIp() + "|" + user.Password))
				if remember == "yes" {
					this.Ctx.SetCookie("auth", strconv.Itoa(user.Id)+"|"+authkey, 7*86400)
				} else {
					this.Ctx.SetCookie("auth", strconv.Itoa(user.Id)+"|"+authkey)
				}

				this.Redirect("/admin", 302)
			}
		}
	}
	this.TplNames = this.moduleName + "/account/login.html"
}

//退出登录
func (this *AccountController) Logout() {
	this.Ctx.SetCookie("auth", "")
	this.Redirect("/admin/login", 302)
}

//资料修改
func (this *AccountController) Profile() {
	user := models.User{Id: this.userid}
	if err := user.Read(); err != nil {
		this.showmsg(err.Error())
	}
	//	user.LastIp = this.getClientIp()
	//	user.Update("lastip")
	if this.isPost() {
		errmsg := make(map[string]string)
		email := strings.TrimSpace(this.GetString("email"))
		password := strings.TrimSpace(this.GetString("password"))
		newpassword := strings.TrimSpace(this.GetString("newpassword"))
		newpassword2 := strings.TrimSpace(this.GetString("newpassword2"))
		updated := false
		if email == user.Email {
			errmsg["email"] = ""
		} else if email == "" {
			errmsg["email"] = "请输入邮箱"
		} else if m, _ := regexp.MatchString(`^([\w\.\_]{2,10})@(\w{1,}).([a-z]{2,4})$`, email); !m {
			errmsg["email"] = "不是有效的邮箱"
		} else {
			user.Email = email
			user.Update("email")
			updated = true
		}

		if newpassword != "" {
			if password == "" || util.Md5([]byte(password)) != user.Password {
				errmsg["password"] = "当前密码错误"
			}
			if len(newpassword) < 6 {
				errmsg["newpassword"] = "密码长度不能少于6个字符"
			} else if newpassword != newpassword2 {
				errmsg["newpassword2"] = "两次输入的密码不一致"
			}
		}

		if len(errmsg) == 0 {
			user.Password = util.Md5([]byte(newpassword))
			user.Update("email")
			updated = true
		}
		this.Data["updated"] = updated
		this.Data["errmsg"] = errmsg
	}
	this.Data["user"] = user
	this.display()
}
