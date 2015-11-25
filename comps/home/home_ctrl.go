package home

import (
	"github.com/xtfly/goman/comps"
	"gopkg.in/macaron.v1"
)

func HomeCtrl(c *macaron.Context) {
	r := comps.NewRender(c)
	r.DefCxt(0)
	c.Data["app"] = "explore"
	r.HTML(200, "home/index")
}
