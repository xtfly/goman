package auth

import (
	"net"
	"net/http"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/xtfly/goman/boot"
	"github.com/xtfly/goman/kits"

	"gopkg.in/macaron.v1"
)

// the identifier of the operation which represents a api action
type OperationId int64

const DummyOptId OperationId = 0

// A simple Authorization Service
type AuthService interface {
	// generate a user token
	// clientip: the ip of the request
	// uid: the id of operate user
	// flag: the type of token
	GenUserToken(clientip string, uid int64, expire int, flag kits.TokenFlag) (string, error)

	// generate a system token
	GenSysToken(clientip string, expire int) (string, error)

	// Check whether the user can access by the token
	Authenticate(clientip, token string) (*kits.UserToken, bool)

	// Check the permissions of the user to operate the resource
	// utoken: the token info after decode
	// ownerUid: the owner of the resource
	// oid: the id of current operation
	Authorize(utoken *kits.UserToken, ownerUid int64, oid OperationId) bool
}

type authService struct {
	crypto *kits.Crypto
}

func (cas *authService) GenUserToken(clientip string, uid int64, expire int, flag kits.TokenFlag) (string, error) {
	ut := &kits.UserToken{
		Flag:     flag,
		Uid:      uid,
		ClientIP: clientip,
		GenTime:  time.Now().Unix(),
		Expire:   expire,
	}

	return ut.GenToken(cas.crypto)
}

func (cas *authService) GenSysToken(clientip string, expire int) (string, error) {
	ut := &kits.UserToken{
		Flag:     kits.TokenSys,
		Uid:      kits.SysUserId,
		ClientIP: clientip,
		GenTime:  time.Now().Unix(),
		Expire:   expire,
	}

	return ut.GenToken(cas.crypto)
}

func (cas *authService) Authenticate(clientip, token string) (*kits.UserToken, bool) {
	ut := &kits.UserToken{}
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

	if time.Now().Unix()-ut.GenTime > int64(ut.Expire*60) {
		log.Errorf("token %s is expired", token)
		return nil, false
	}

	switch ut.Flag {
	case kits.TokenSys:
		if ut.Uid != kits.SysUserId {
			log.Errorf("Uid(%v!=%v) is invalid.", ut.Uid, kits.SysUserId)
			return nil, false
		}
	case kits.TokenUser:
		if ut.Uid <= 0 {
			log.Errorf("Uid(%v) is invalid.", ut.Uid)
			return nil, false
		}
	default:
		break
	}

	return ut, true
}

func (cas *authService) Authorize(ut *kits.UserToken, ownerUid int64, oid OperationId) bool {
	if ut != nil {
		return false
	}

	// only support operate owner resource
	switch ut.Flag {
	case kits.TokenUser:
		if ut.Uid != ownerUid {
			log.Errorf("Uid(%v!=%v) is not equal.", ut.Uid, ownerUid)
			return false
		}
	default:
		break
	}

	return true
}

func Auther() macaron.Handler {
	as := &authService{crypto: boot.GetCrypto()}

	return func(res http.ResponseWriter, req *http.Request, c *macaron.Context) {
		c.Map(as)

		if !strings.HasPrefix(req.URL.Path, "/api/") {
			return
		}

		token := strings.TrimSpace(req.Header.Get("Authorization"))
		//log.Debugf("retrive token = % from header", token)

		if token == "" {
			http.Error(res, "Not Authorized", http.StatusUnauthorized)
			return
		}

		ut, ok := as.Authenticate(c.RemoteAddr(), token)
		if !ok {
			http.Error(res, "Not Authorized", http.StatusUnauthorized)
			return
		}

		c.Map(ut)
	}
}
