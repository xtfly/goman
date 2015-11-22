package models

//
import "github.com/astaxie/beego/orm"

//
type Pages struct {
	Id       int64  `json:"id" orm:"pk;auto"`                //ID
	UrlToken string `json:"url_token" orm:"size(32);unique"` //
	Title    string `json:"title" orm:"null"`                //
	Keywords string `json:"keywords" orm:"null"`             //
	Desc     string `json:"description" orm:"null"`          //
	Contents string `json:"contents" orm:"type(text);null"`  //
	Enabled  bool   `json:"enabled" orm:"default(0)"`        //

}

func init() {
	orm.RegisterModel(new(Pages))
}
