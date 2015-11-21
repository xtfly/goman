package home

import (
	"github.com/xtfly/goman/comps"
	"gopkg.in/macaron.v1"
)

func HomeCtrl(c *macaron.Context) {
	comps.DefCxt(c, 0)
	c.Data["app"] = "explore"
	comps.HTML(c, 200, "home/index")
}
