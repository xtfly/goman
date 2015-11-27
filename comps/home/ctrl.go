package home

import (
	"github.com/xtfly/goman/comps/core"
	"gopkg.in/macaron.v1"
)

func HomeCtrl(c *macaron.Context) {
	r := core.NewRender(c)
	c.Data["app"] = "explore"
	r.HTML(200, "home/index")
}
