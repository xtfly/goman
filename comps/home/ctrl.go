package home

import (
	"github.com/xtfly/goman/comps/core"
	"gopkg.in/macaron.v1"
)

func HomeCtrl(c *macaron.Context) {
	r := core.NewRender(c)
	c.Data["app"] = "explore"

	r.AddJs("comps/home.js")
	r.RHTML(200, "home/index")
}

func GetFirstLoginCtrl(c *macaron.Context) {
	r := core.NewRender(c)
	c.Data["app"] = "explore"

	r.AddJs("comps/home.js")
	r.RHTML(200, "home/index")
}
