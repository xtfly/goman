package comps

import (
	"fmt"

	"github.com/go-macaron/captcha"
)

//----------------------------------------------------------
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

func NewRestResp(rsm interface{}, errno int, err string) *RestResp {
	return &RestResp{
		Rsm:   rsm,
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

//----------------------------------------------------------
// 验证码信息
type CaptchaInfo struct {
	CaptchaId  string `json:"captcha_id"`
	CaptchaUrl string `json:"captcha_url"`
}

func NewCaptcha(cpt *captcha.Captcha) *CaptchaInfo {
	cptvalue, err := cpt.CreateCaptcha()
	if err != nil {
		return nil
	}

	return &CaptchaInfo{
		CaptchaId:  cptvalue,
		CaptchaUrl: fmt.Sprintf("%s%s%s.png", cpt.SubURL, cpt.URLPrefix, cptvalue),
	}
}

//----------------------------------------------------------
type UploadFileErrRsp struct {
	Error string `json:"error"`
}

type UploadFileRsp struct {
	Success bool   `json:"success"`
	Thumb   string `json:"thumb"`
}

func NewUploadFileErrRsp(err string) *UploadFileErrRsp {
	return &UploadFileErrRsp{
		Error: err,
	}
}

func NewUploadFileRsp(thumb string) *UploadFileRsp {
	return &UploadFileRsp{
		Success: true,
		Thumb:   thumb,
	}
}
