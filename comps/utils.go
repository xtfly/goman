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
	c.Data["import_css_files"] = []string{
		boot.SysSetting.Si.Static + "/front/css/" + boot.SysSetting.Ps.UiStyle + "/common.css",
		boot.SysSetting.Si.Static + "/front/css/" + boot.SysSetting.Ps.UiStyle + "/link.css",
	}
}

func AddCss(c *macaron.Context, css ...string) {
	for _, f := range css {
		c.Data["import_css_files"] = append(c.Data["import_css_files"].([]string), boot.SysSetting.Si.Static+"/front/css/"+f)
	}
}

func DefJs(c *macaron.Context) {
	c.Data["import_js_files"] = []string{
		boot.SysSetting.Si.Static + "/front/js/goman.js",
		boot.SysSetting.Si.Static + "/front/js/template.js",
		boot.SysSetting.Si.Static + "/front/js/app.js",
	}
}

func AddJs(c *macaron.Context, js ...string) {
	for _, f := range js {
		c.Data["import_js_files"] = append(c.Data["import_js_files"].([]string), boot.SysSetting.Si.Static+"/front/js/"+f)
	}
}

func DefCxt(c *macaron.Context, uid int64) {
	AddDatas(c, uid)
	DefCss(c)
	DefJs(c)
}

//格式化系统返回消息
//格式化系统返回的消息 json 数据包给前端进行处理
type RestErrResp struct {
	Rsm   string `json:"rsm"`
	Errno int    `json:"errno"`
	Err   string `json:"err"`
}

func NewRestErrResp(rsm string, errno int, err string) *RestErrResp {
	return &RestErrResp{
		Rsm:   rsm,
		Errno: errno,
		Err:   err,
	}
}
