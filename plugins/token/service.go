package token

import (
	"net"
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/go-macaron/session"
	"github.com/xtfly/goman/boot"
	"github.com/xtfly/goman/kits"

	"gopkg.in/macaron.v1"
)

// A simple Authorization Service
type TokenService interface {
	// generate a user token
	// clientip: the ip of the request
	// uid: the id of operate user
	// flag: the type of token
	GenUserToken(clientip string, uid int64, expire int, flag TokenFlag) (string, error)

	// generate a system token
	GenSysToken(clientip string, expire int) (string, error)

	// Check whether the user can access by the token
	Authenticate(clientip, token string) (*UserToken, bool)

	ValidToken(clientip, token string) bool
}

type tokenService struct {
	crypto *kits.Crypto
}

func (cas *tokenService) GenUserToken(clientip string, uid int64, expire int, flag TokenFlag) (string, error) {
	ut := &UserToken{
		Flag:     flag,
		Uid:      uid,
		ClientIP: clientip,
		GenTime:  time.Now().Unix(),
		Expire:   expire,
	}

	return ut.GenToken(cas.crypto)
}

func (cas *tokenService) GenSysToken(clientip string, expire int) (string, error) {
	ut := &UserToken{
		Flag:     TokenSys,
		Uid:      SysUserId,
		ClientIP: clientip,
		GenTime:  time.Now().Unix(),
		Expire:   expire,
	}

	return ut.GenToken(cas.crypto)
}

func (cas *tokenService) Authenticate(clientip, token string) (*UserToken, bool) {
	ut := &UserToken{}
	if !ut.DecodeToken(cas.crypto, token) {
		log.Errorf("Decode token failed, token=%s", token)
		return nil, false
	}

	// check the ip
	tip := net.ParseIP(ut.ClientIP)
	rip := net.ParseIP(clientip)
	if !tip.Equal(rip) {
		log.Errorf("Check client failed, token ip=%s, real ip=%s", ut.ClientIP, clientip)
		return nil, false
	}

	switch ut.Flag {
	case TokenSys:
		if ut.Uid != SysUserId {
			log.Errorf("Uid(%v!=%v) is invalid.", ut.Uid, SysUserId)
			return nil, false
		}
	case TokenUser:
		if ut.Uid <= 0 {
			log.Errorf("Uid(%v) is invalid.", ut.Uid)
			return nil, false
		}
	default:
		break
	}

	return ut, true
}

func (cas *tokenService) ValidToken(clientip, token string) bool {
	if ut, ok := cas.Authenticate(clientip, token); !ok {
		return false
	} else {
		if ut.Expired() {
			return false
		}
	}
	return true
}

func Tokener() macaron.Handler {
	as := &tokenService{crypto: boot.GetCrypto()}

	return func(res http.ResponseWriter, c *macaron.Context, ss session.Store) {
		c.Map(as)

		// Cookie 清除则 Session 也清除
		utoken := c.GetCookie("utoken")
		if utoken == "" {
			ss.Delete("utoken")
			return
		}

		// 可能客户端变化，需要重新验证
		ut, ok := as.Authenticate(c.RemoteAddr(), utoken)
		if !ok || ut.Expired() {
			http.Error(res, "Not Authorized", http.StatusUnauthorized)
			return
		}

		// 服务端存在，但不相等
		stoken := ss.Get("utoken").(string)
		if stoken != "" {
			if utoken != stoken {
				http.Error(res, "Not Authorized", http.StatusUnauthorized)
				return
			}
		} else {
			// session被GC了，
			ss.Set("utoken", utoken)
		}

		if ut.NeedRefresh() {
			nt, err := as.GenUserToken(c.RemoteAddr(), ut.Uid, ut.Expire, ut.Flag)
			if err != nil {
				http.Error(res, "Not Authorized", http.StatusUnauthorized)
				return
			}
			ss.Set("utoken", nt)
			c.SetCookie("utoken", nt, ut.Expire*60)
		}

		c.Data["uid"] = ut.Uid
		c.Map(ut)
	}
}
