package comps

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
