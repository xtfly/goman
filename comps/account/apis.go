package account

import (
	"bytes"
	"image"
	"image/png"
	"os"
	"strconv"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/Unknwon/com"
	"github.com/nfnt/resize"

	"github.com/go-macaron/captcha"
	"github.com/go-macaron/session"
	"github.com/xtfly/gokits"
	"github.com/xtfly/goman/boot"
	"github.com/xtfly/goman/comps"
	"github.com/xtfly/goman/comps/core"

	"github.com/xtfly/goman/models"
	"github.com/xtfly/goman/plugins/token"
	"gopkg.in/macaron.v1"
)

//----------------------------------------------------------
// POST /api/account/check/
func ApiCheckUserName(c *macaron.Context) {
	un := strings.TrimSpace(c.QueryEscape("username"))
	s := NewService()
	if err, ok := s.CheckUsernameChar(un); !ok {
		c.JSON(200, comps.NewRestErrResp(-1, err))
		return
	}

	if s.CheckUsernameSensitiveWords(un) || models.UserExistedByName(un) {
		c.JSON(200, comps.NewRestErrResp(-1, "用户名已被注册"))
		return
	}

	c.JSON(200, comps.NewRestErrResp(1, ""))
}

//----------------------------------------------------------
// POST /api/account/signup/
func ApiUserSignup(f SignupForm, c *macaron.Context, cpt *captcha.Captcha,
	a token.TokenService, ss session.Store) {
	if !a.ValidToken(c.RemoteAddr(), f.CsrfToken) {
		c.JSON(200, comps.NewRestErrResp(-1, "非法的跨站请求"))
		return
	}

	if !cpt.VerifyReq(c.Req) {
		c.JSON(200, comps.NewRestResp(comps.NewCaptcha(cpt), -1, "请填写正确的验证码"))
		return
	}

	s := NewService()
	u, msg, ok := s.Signup(f, c.RemoteAddr())
	if !ok {
		c.JSON(200, comps.NewRestResp(comps.NewCaptcha(cpt), -1, msg))
		return
	}

	// 如果不需要email验证
	if boot.SysSetting.Ra.RegisterValidType == models.RegValidNone ||
		u.GroupId != models.GroupNotValidated ||
		u.ValidEmail {
		SetSigninCookies(c, u, a, ss)
		c.JSON(200, comps.NewRestRedirectResp("/h/firstlogin"))
		return
	}

	ss.Set("validemail", u.Email)
	if !models.NewValidByEmail(models.NewTr(), u.Id, u.Email) {
		c.JSON(200, comps.NewRestErrResp(-1, "内部系统错误"))
		return
	}

	SetSigninCookies(c, u, a, ss)
	c.JSON(200, comps.NewRestRedirectResp("/a/validemail/"))
	return
}

//----------------------------------------------------------
// POST /api/account/signin/
func ApiSignin(c *macaron.Context, f SigninForm, a token.TokenService, ss session.Store) {
	u := &models.Users{}
	if !u.CheckSignin(f.Input, f.Password) {
		c.JSON(200, comps.NewRestErrResp(-1, "输入正确的帐号或密码"))
		return
	}

	s := NewService()
	if err, ok := s.CheckSignin(u); !ok {
		c.JSON(200, comps.NewRestErrResp(-1, err))
		return
	}

	// 需要审批
	if u.GroupId == models.GroupNotValidated &&
		boot.SysSetting.Ra.RegisterValidType == models.RegValidApproval {
		c.JSON(200, comps.NewRestRedirectResp("/a/validapproval/"))
		return
	}

	//
	u.LastLogin = time.Now()
	u.LastIp = c.RemoteAddr()
	u.LoginCount = u.LoginCount + 1
	if _, ok := models.NewTr().Update(u, "LastLogin", "LastIp", "LoginCount"); !ok {
		// todo log
	}

	CleanCookies(c, ss)
	SetSigninCookies(c, u, a, ss)

	url := ""
	if !u.ValidEmail && boot.SysSetting.Ra.RegisterValidType == models.RegValidEmail {
		ss.Set("validemail", u.Email)
		url = "/a/validemail/"
	} else if u.FirstLogin {
		url = "/h/firstlogin/"
	} else if f.ReturnUrl != "" {
		url = f.ReturnUrl
	}

	c.JSON(200, comps.NewRestRedirectResp(url))
}

//----------------------------------------------------------
// POST /api/account/setting/profile
func ApiSettingProfile(c *macaron.Context, f UserSettingForm, ss session.Store) {
	r := core.NewRender(c)

	if msg, ok := r.CheckUser(); !ok {
		c.JSON(200, comps.NewRestErrResp(-1, msg))
		return
	}

	s := NewService()
	u := r.UserInfo
	nu := &models.Users{Id: u.Id}
	t := models.NewTr().Begin()
	defer t.End()

	// 如果原来是采用Email注册，默认使用Email做为username
	nu.UserName = u.UserName
	if f.UserName != "" {
		if msg, ok := s.CheckUsernameChar(f.UserName); !ok {
			c.JSON(200, comps.NewRestErrResp(-1, msg))
			return
		}
		if u.UserName != f.UserName && models.UserExistedByName(f.UserName) {
			c.JSON(200, comps.NewRestErrResp(-1, "已经存在相同的姓名, 请重新填写"))
			return
		}
		nu.UserName = f.UserName
	}

	//
	nu.UrlToken = u.UrlToken
	if f.UrlToken != "" && f.UrlToken != u.UrlToken {
		if msg, ok := s.CheckUrlToken(&u.Users, f.UrlToken); !ok {
			c.JSON(200, comps.NewRestErrResp(-1, msg))
			return
		}
		nu.UrlToken = f.UrlToken
	}

	nu.Email = u.Email
	if f.Email != "" {
		if !gokits.IsEmail(f.Email) {
			c.JSON(200, comps.NewRestErrResp(-1, "请输入正确的 E-Mail 地址"))
			return
		}
		if !models.UserExistedByEmail(f.Email) {
			c.JSON(200, comps.NewRestErrResp(-1, "邮箱已经存在, 请使用新的邮箱"))
			return
		}
		nu.Email = f.Email
		models.NewValidByEmail(t, u.Id, nu.Email)
	}

	nu.CommonEmail = u.CommonEmail
	if f.CommonEmail != "" {
		if !gokits.IsEmail(f.CommonEmail) {
			c.JSON(200, comps.NewRestErrResp(-1, "请输入正确的常用邮箱地址"))
			return
		}
		nu.CommonEmail = f.CommonEmail
	}

	nu.Gender = f.Gender
	nu.Province = gokits.IfEmpty(f.Province, u.Province)
	nu.City = gokits.IfEmpty(f.City, u.City)

	nu.Birthday = u.Birthday
	if f.Birthday != "" {
		nu.Birthday, _ = time.Parse("19801010", f.Birthday)
	}

	nu.Signature = u.Signature
	if f.Signature != "" {
		nu.Signature = f.Signature
		if !models.IntegralLogExistByUidAction(u.Id, models.IntegralUpdateUserSignature) {
			models.AddIntegralLog(t, u.Id, models.IntegralUpdateUserSignature, int64(float64(boot.SysSetting.Ir.FinishProfile)*0.1), "完善一句话介绍")
		}
	}

	nu.JobId = u.JobId
	if f.JobId != 0 {
		nu.JobId = f.JobId
	}
	nu.Mobile = gokits.IfEmpty(f.Mobile, u.Mobile)

	if boot.SysSetting.Cs.AutoCreateSocialTopic {
		if f.Province != "" {
			models.AddTopic(t, f.Province)
		}
		if f.City != "" {
			models.AddTopic(t, f.City)
		}
	}

	if _, ok := t.Update(nu, "UserName", "Gender", "Province", "Province", "JobId",
		"Signature", "Email", "Signature", "UrlToken", "CommonEmail",
		"Birthday", "Mobile"); !ok {
		c.JSON(200, comps.NewRestErrResp(-1, "个人资料保存成功失败"))
	} else {
		c.JSON(200, comps.NewRestErrResp(1, "个人资料保存成功"))
	}
}

// /api/account/avatar/upload/
func ApiUploadAvatar(c *macaron.Context) {
	r := core.NewRender(c)

	if msg, ok := r.CheckUser(); !ok {
		c.JSON(200, comps.NewRestErrResp(-1, msg))
		return
	}

	f, h, err := c.GetFile("upload_file")
	if err != nil {
		log.Errorln("not find image, ", err.Error())
		c.JSON(200, comps.NewUploadFileErrRsp("你没有上传文件"))
		return
	}
	//log.Infoln("filename=", h.Filename)
	defer f.Close()

	// FIXME: the workdir should app root path
	ext := strings.ToLower(h.Filename[strings.LastIndex(h.Filename, ".")+1:])
	if !strings.Contains("jpg,jpeg,png,gif", ext) {
		c.JSON(200, comps.NewUploadFileErrRsp("文件类型无效"))
		return
	}

	// 上传的路径
	path := boot.SysSetting.Si.UploadDir + "/avatar/" + com.DateT(time.Now(), "YY/MM/DD/")
	os.MkdirAll(path, 0744)
	path = path + strconv.Itoa(int(r.Uid)) + "_"

	img, _, err := image.Decode(f)
	if err != nil {
		log.Errorln("decode image failed, ", err.Error())
		c.JSON(200, comps.NewUploadFileErrRsp("生成文件失败"))
		return
	}

	// 缩放图片
	sizes := []uint{32, 50, 100}
	names := []string{"min", "mid", "max"}
	tfn := ""
	for i, v := range sizes {
		tfn = path + names[i] + ".png"
		//log.Infoln("new image name, ", tfn)
		tf, err := os.OpenFile(tfn, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			log.Errorln("write failed, ", err.Error())
			c.JSON(200, comps.NewUploadFileErrRsp("打开目标文件失败"))
			return
		}
		defer tf.Close()

		timg := resize.Thumbnail(v, v, img, resize.Lanczos3)
		var buf bytes.Buffer
		if err := png.Encode(&buf, timg); err != nil {
			log.Errorln("encode png failed, ", err.Error())
			c.JSON(200, comps.NewUploadFileErrRsp("编码PNG文件失败"))
			return
		}

		if _, err := tf.Write(buf.Bytes()); err != nil {
			log.Errorln("write png failed, ", err.Error())
			c.JSON(200, comps.NewUploadFileErrRsp("写入PNG文件失败"))
			return
		}
	}

	// 更新数据库
	t := models.NewTr().Begin()
	defer t.End()
	r.UserInfo.Avatar = "/" + tfn
	if _, ok := t.Update(&r.UserInfo.Users, "Avatar"); !ok {
		c.JSON(200, comps.NewUploadFileErrRsp("内部错误"))
		return
	}

	// 增加积分
	if !models.IntegralLogExistByUidAction(r.Uid, models.IntegralUploadUserAvatar) {
		models.AddIntegralLog(t, r.Uid, models.IntegralUploadUserAvatar, int64(float64(boot.SysSetting.Ir.FinishProfile)*0.2), "上传头像")
	}

	// c.JSON(200, comps.NewUploadFileRsp(r.UserInfo.Avatar))
	json, _ := c.JSONString(comps.NewUploadFileRsp(r.UserInfo.Avatar))
	json = strings.Replace(json, "/", "\\/", -1)
	c.Resp.Header().Set("Content-Type", "text/html; charset=utf-8")
	c.Write([]byte(json))
}

// /api/account/firstlogin/clean/
func ApiCleanFirstLogin(c *macaron.Context) {
	r := core.NewRender(c)

	if msg, ok := r.CheckUser(); !ok {
		c.JSON(200, comps.NewRestErrResp(-1, msg))
		return
	}

	r.UserInfo.FirstLogin = false

	if _, ok := models.NewTr().Update(&r.UserInfo.Users, "FirstLogin"); !ok {
		r.PlainText(200, []byte("failed"))
	} else {
		r.PlainText(200, []byte("success"))
	}
}
