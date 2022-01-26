package controllers

import (
	"bytes"
	"fmt"
	"net/http"
	"text/template"
	"time"

	"github.com/StubbornYouth/goblog/app/models/passwordreset"
	"github.com/StubbornYouth/goblog/app/models/user"
	"github.com/StubbornYouth/goblog/app/requests"
	"gorm.io/gorm"

	"github.com/StubbornYouth/goblog/pkg/auth"
	"github.com/StubbornYouth/goblog/pkg/datetime"
	"github.com/StubbornYouth/goblog/pkg/flash"
	"github.com/StubbornYouth/goblog/pkg/logger"
	"github.com/StubbornYouth/goblog/pkg/rand"
	"github.com/StubbornYouth/goblog/pkg/route"
	"github.com/StubbornYouth/goblog/pkg/smtp"
	"github.com/StubbornYouth/goblog/pkg/view"
)

type AuthController struct {
}

// type userForm struct {
// 	Name            string `valid:"name"`
// 	Email           string `valid:"email"`
// 	Password        string `valid:password`
// 	PasswordConfirm string `valid:password_confirm`
// }

// 有两个问题
// 1.注册登录页面不需要左边导航栏，需使用不同的布局文件；
// 2.ArticlesFormData 只限于在文章控制器中使用，在此处显得格格不入，需使用更加通用的数据格式，以便在所有控制器中使用

// 注册页面
func (*AuthController) Register(w http.ResponseWriter, r *http.Request) {
	view.RenderSimple(w, view.D{}, "auth.register")
}

// 处理注册逻辑
// func (*AuthController) DoRegister(w http.ResponseWriter, r *http.Request) {
// 	// 初始化变量
// 	name := r.PostFormValue("name")
// 	email := r.PostFormValue("email")
// 	password := r.PostFormValue("password")

// 	// 表单验证
// 	// 验证通过入库 并跳转
// 	_user := user.User{
// 		Name:     name,
// 		Email:    email,
// 		Password: password,
// 	}

// 	_user.Create()
// 	if _user.ID > 0 {
// 		fmt.Fprint(w, "插入成功，新建得用户ID为"+_user.GetStringID())
// 	} else {
// 		w.WriteHeader(http.StatusInternalServerError)
// 		fmt.Fprint(w, "创建用户失败，请联系管理员")
// 	}
// 	// 验证不通过 返回表单
// }

// 处理注册逻辑
func (*AuthController) DoRegister(w http.ResponseWriter, r *http.Request) {
	// 初始化数据
	// _user := userForm{
	// 	Name:            r.PostFormValue("name"),
	// 	Email:           r.PostFormValue("email"),
	// 	Password:        r.PostFormValue("password"),
	// 	PasswordConfirm: r.PostFormValue("password_confirm"),
	// }
	_user := user.User{
		Name:            r.PostFormValue("name"),
		Email:           r.PostFormValue("email"),
		Password:        r.PostFormValue("password"),
		PasswordConfirm: r.PostFormValue("password_confirm"),
	}

	// // 表单规则
	// rules := govalidator.MapData{
	// 	// alpha_num 只允许英文字母和数字混合
	// 	"name":             []string{"required", "alpha_num", "between:3,20"},
	// 	"email":            []string{"required", "min:4", "max:30", "email"},
	// 	"password":         []string{"required", "min:6"},
	// 	"password_confirm": []string{"required"},
	// }

	// // 定制错误消息
	// messages := govalidator.MapData{
	// 	"name": []string{
	// 		"required:用户名为必填项",
	// 		"alpha_num:格式错误，只允许字母和数字",
	// 		"between:用户名长度需在3~20之间",
	// 	},
	// 	"email": []string{
	// 		"required:Email为必填项",
	// 		"min:长度必须大于4",
	// 		"max:长度必须小于30",
	// 		"email:格式错误，请提供有效的邮箱地址",
	// 	},
	// 	"password": []string{
	// 		"required:密码为必填项",
	// 		"min:密码长度必须大于6",
	// 	},
	// 	"password_confirm": []string{
	// 		"required:确认密码为必填项",
	// 	},
	// }

	// // 配置选项
	// opts := govalidator.Options{
	// 	Data:          &_user,
	// 	Rules:         rules,
	// 	TagIdentifier: "valid",  // Struct 标签标识符
	// 	Messages:      messages, // 增加自定义验证提示
	// }

	// // 开始验证
	// err := govalidator.New(opts).ValidateStruct()

	err := requests.ValidareRegistrationForm(_user)

	// 判断是否验证通过
	if len(err) > 0 {
		// json.MarshalIndent() 方法，一般此方法用来将 Go 对象格式成为 JSON 字符串，并加上合理的缩进
		// data, _ := json.MarshalIndent(err, "", " ")
		// fmt.Fprint(w, string(data))

		// 验证不通过 重新显示表单
		view.RenderSimple(w, view.D{
			"Errors": err,
			"User":   _user,
		}, "auth.register")
	} else {
		_user.Create()
		if _user.ID > 0 {
			// fmt.Fprint(w, "插入成功，新建得用户ID为"+_user.GetStringID())
			// 登录用户 并 跳转到首页
			auth.Login(_user)
			flash.Success("恭喜您注册成功！")
			http.Redirect(w, r, "/", http.StatusFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "创建用户失败，请联系管理员")
		}
	}
}

func (*AuthController) Login(w http.ResponseWriter, r *http.Request) {
	// 设置会话信息
	// session.Put("uid", 1)
	// 获取会话信息
	// session.Get("uid")
	// fmt.Fprint(w, session.Get("uid"))
	// 删除会话数据
	// session.Forget("uid")
	// 销毁整个会话
	// session.Flush()
	view.RenderSimple(w, view.D{}, "auth.login")
}

func (*AuthController) DoLogin(w http.ResponseWriter, r *http.Request) {
	email := r.PostFormValue("email")
	password := r.PostFormValue("password")

	if err := auth.Attempt(email, password); err == nil {
		// 登录成功
		flash.Success("欢迎回来！")
		http.Redirect(w, r, "/", http.StatusFound)
	} else {
		// 失败 页面提示错误信息
		view.RenderSimple(w, view.D{
			"Error":    err.Error(),
			"Email":    email,
			"Password": password,
		}, "auth.login")
	}
}

func (*AuthController) Logout(w http.ResponseWriter, r *http.Request) {
	auth.Loginout()
	// 跳转到首页
	flash.Success("您已退出登录")
	http.Redirect(w, r, "/", http.StatusFound)
}

func (*AuthController) Forget(w http.ResponseWriter, r *http.Request) {
	view.RenderSimple(w, view.D{}, "auth.forget")
}

func (*AuthController) DoForget(w http.ResponseWriter, r *http.Request) {
	email := r.PostFormValue("email")
	title := "GoBlog密码重置邮件"
	// body := goquery.ParseUrl()
	_reset := passwordreset.PasswordReset{
		Email: email,
		Token: rand.RandString(10),
	}

	err := requests.ValidarePasswordForm(_reset)

	// 判断是否验证通过
	if len(err) > 0 {
		// json.MarshalIndent() 方法，一般此方法用来将 Go 对象格式成为 JSON 字符串，并加上合理的缩进
		// data, _ := json.MarshalIndent(err, "", " ")
		// fmt.Fprint(w, string(data))

		// 验证不通过 重新显示表单
		view.RenderSimple(w, view.D{
			"Errors":        err,
			"PasswordReset": _reset,
		}, "auth.forget")
	} else {
		var buf bytes.Buffer

		// 判断用户是否存在
		// 根据email 获取用户
		_, err := user.GetByEmail(email)
		// 判断用户是否存在
		if err != nil {
			errors := make(map[string][]string)
			if err == gorm.ErrRecordNotFound {
				errors["email"] = []string{"当前邮箱账号不存在"}
			} else {
				errors["email"] = []string{"内部错误，请稍后尝试"}
			}

			view.RenderSimple(w, view.D{
				"Errors":        errors,
				"PasswordReset": _reset,
			}, "auth.forget")
		} else {
			_reset.Create()
			if _reset.ID > 0 {
				// 发送模板内容
				tmpl, err := template.New("reset_link.gohtml").ParseFiles("resources/views/auth/reset_link.gohtml")
				if err != nil {
					logger.LogError(err)
					w.WriteHeader(http.StatusInternalServerError)
					fmt.Fprint(w, "发送模板异常，请联系管理员")
				}

				//需要先导入bytes包
				// host := "localhost:3000"
				url := route.RouteNameToURL("auth.reset", "token", _reset.Token)
				//定义Buffer类型
				// var path bytes.Buffer
				// 向path中写入字符串
				// path.WriteString(host)
				// path.WriteString(url)
				//获得拼接后的字符串
				data := view.D{
					// "URL": path.String(),
					"URL": url,
				}
				// fmt.Printf(tmpl)
				execErr := tmpl.Execute(&buf, data)
				// execErr := view.RenderSimple(w, view.D{
				// 	"reset": _reset,
				// }, "auth.reset")

				if execErr != nil {
					logger.LogError(execErr)
					fmt.Fprint(w, "模板赋值异常，请联系管理员")
				}

				err = smtp.SendEmail(email, title, buf.String())
				if err == nil {
					flash.Success("邮件发送成功，请前往邮箱进行认证")
				} else {
					flash.Danger("邮件发送失败，请稍后再试")
				}

				// http.Redirect(w, r, route.RouteNameToURL("auth.forget"), http.StatusFound)
				view.RenderSimple(w, view.D{
					"PasswordReset": _reset,
				}, "auth.forget")
			} else {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprint(w, "系统异常，请稍后再试")
			}
		}
	}
}

func (*AuthController) Reset(w http.ResponseWriter, r *http.Request) {
	token := route.GetRouteVarible("token", r)
	view.RenderSimple(w, view.D{"Token": token}, "auth.reset")
}

func (*AuthController) DoReset(w http.ResponseWriter, r *http.Request) {
	email := r.PostFormValue("email")
	password := r.PostFormValue("password")
	token := r.PostFormValue("token")
	_user := user.User{
		Email:           email,
		Password:        password,
		PasswordConfirm: r.PostFormValue("password_confirm"),
	}

	err := requests.ValidareResetForm(_user)

	// 判断是否验证通过
	if len(err) > 0 {
		// json.MarshalIndent() 方法，一般此方法用来将 Go 对象格式成为 JSON 字符串，并加上合理的缩进
		// data, _ := json.MarshalIndent(err, "", " ")
		// fmt.Fprint(w, string(data))

		// 验证不通过 重新显示表单
		view.RenderSimple(w, view.D{
			"Errors":        err,
			"PasswordReset": _user,
			"Token":         token,
		}, "auth.reset")
	} else {
		// 判断 token是否正确且 未过期
		_passwordreset, err := passwordreset.GetByEmail(email)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				flash.Danger("未获取到重置记录")
			} else {
				flash.Danger("内部错误，请稍后尝试")
			}

			view.RenderSimple(w, view.D{
				"PasswordReset": _user,
			}, "auth.reset")
		} else {
			expires := 60 * 10 // 过期时间 10分钟
			time := time.Now().Unix()
			createTime := datetime.DateToTime(_passwordreset.CreatedAt.Format("2006-01-02 15:04:05"))
			if int(time-createTime) > expires {
				flash.Danger("链接已过期，请重新尝试")
				view.RenderSimple(w, view.D{
					"PasswordReset": _user,
					"Token":         token,
				}, "auth.reset")
			} else {
				if _passwordreset.Token != token {
					flash.Danger("令牌错误")
					view.RenderSimple(w, view.D{
						"PasswordReset": _user,
						"Token":         token,
					}, "auth.reset")
				} else {
					// 一切正常 更新密码
					_, err := user.UpdatePassword(_user)
					if err != nil {
						w.WriteHeader(http.StatusInternalServerError)
						flash.Danger("500 服务器错误")
						view.RenderSimple(w, view.D{
							"PasswordReset": _user,
							"Token":         token,
						}, "auth.reset")
					} else {
						flash.Success("密码重置成功，请使用新密码登录")
						http.Redirect(w, r, route.RouteNameToURL("auth.login"), http.StatusFound)
					}
				}
			}
		}
	}
}
