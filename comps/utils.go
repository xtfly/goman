package comps

import (
	"github.com/xtfly/goman/boot"
	"gopkg.in/macaron.v1"
)

func AddDatas(c *macaron.Context, uid int64) {
	c.Data["X_UA_COMPATIBLE"] = boot.X_UA_COMPATIBLE
	c.Data["VERSION_BUILD"] = boot.VERSION_BUILD
	c.Data["sys"] = boot.SysSetting
	c.Data["uid"] = uid
}

func HTML(c *macaron.Context, sc int, tn string) {
	us := boot.SysSetting.Ps.UiStyle
	c.HTML(sc, us+"/"+tn)
}

func DefCss(c *macaron.Context) {
	s := boot.SysSetting.Si.Static
	u := boot.SysSetting.Ps.UiStyle
	c.Data["css_files"] = []string{
		s + "/front/css/" + u + "/common.css",
		s + "/front/css/" + u + "/link.css",
	}
}

func AddCss(c *macaron.Context, css ...string) {
	for _, f := range css {
		c.Data["css_files"] = append(c.Data["css_files"].([]string), boot.SysSetting.Si.Static+"/front/css/"+f)
	}
}

func DefJs(c *macaron.Context) {
	s := boot.SysSetting.Si.Static
	c.Data["js_files"] = []string{
		s + "/front/js/goman.js",
		s + "/front/js/template.js",
		s + "/front/js/app.js",
	}
}

func AddJs(c *macaron.Context, js ...string) {
	for _, f := range js {
		c.Data["js_files"] = append(c.Data["js_files"].([]string), boot.SysSetting.Si.Static+"/front/js/"+f)
	}
}

func DefCxt(c *macaron.Context, uid int64) {
	AddDatas(c, uid)
	DefCss(c)
	DefJs(c)
}

//格式化系统返回消息
//格式化系统返回的消息 json 数据包给前端进行处理
type RestResp struct {
	Rsm   interface{} `json:"rsm"`
	Errno int         `json:"errno"`
	Err   string      `json:"err"`
}

func NewRestErrResp(errno int, err string) *RestResp {
	return &RestResp{
		Rsm:   nil,
		Errno: errno,
		Err:   err,
	}
}

func NewRestRedirectResp(url string) *RestResp {
	return &RestResp{
		Rsm:   map[string]string{"url": url},
		Errno: 1,
		Err:   "",
	}
}
