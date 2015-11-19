package auth

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
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

// the type of token
type TokenFlag byte

// A simple Authorization Service
type AuthService interface {
	// generate a user token
	// clientip: the ip of the request
	// uid: the id of operate user
	// flag: the type of token
	GenUserToken(clientip string, uid int64, expire int, flag TokenFlag) (string, error)

	// generate a system token
	GenSysToken(clientip string, expire int) (string, error)

	// Check whether the user can access by the token
	Authenticate(clientip, token string) (*UserToken, bool)

	// Check the permissions of the user to operate the resource
	// utoken: the token info after decode
	// ownerUid: the owner of the resource
	// oid: the id of current operation
	Authorize(utoken *UserToken, ownerUid int64, oid OperationId) bool
}

// flag(1) + uid(8) + ip(16) + currenttime(8) + expire(4)
type UserToken struct {
	Flag     TokenFlag
	Uid      int64
	ClientIP string // the ip of the request
	GenTime  int64  // timestamp when generating token
	Expire   int    // unit: second
}

const (
	TokenLen             = 1 + 8 + 16 + 8 + 4
	TokenSys   TokenFlag = 0x00 // system general user
	TokenUser  TokenFlag = 0x01 // login user
	TokenAdmin TokenFlag = 0x01

	sysUserId int64 = int64(-8619820608)

	DummyOptId OperationId = 0
)

func (ut *UserToken) GenToken(crypto *kits.Crypto) (string, error) {
	bs := make([]byte, TokenLen)
	buf := bytes.NewBuffer(bs)

	// flag
	buf.WriteByte(byte(ut.Flag))

	// uid
	binary.Write(buf, binary.BigEndian, ut.Uid)

	// client ip
	ip := net.ParseIP(ut.ClientIP)
	buf.Write([]byte(ip))

	// current time
	binary.Write(buf, binary.BigEndian, time.Now().Unix())

	// current time
	binary.Write(buf, binary.BigEndian, ut.Expire)

	// encrypt
	if t, err := crypto.Encrypt(buf.Bytes()); err != nil {
		return "", err
	} else {
		return base64.URLEncoding.EncodeToString(t), nil
	}
}

func (ut *UserToken) DecodeToken(crypto *kits.Crypto, token string) bool {
	bs, err := base64.URLEncoding.DecodeString(token)
	if err != nil {
		return false
	}

	// decrypt
	var tbs []byte
	tbs, err = crypto.Decrypt(bs)
	if err != nil {
		return false
	}

	// check the len
	if len(tbs) != TokenLen {
		return false
	}

	ut.Flag = TokenFlag(tbs[0])
	// remain bytes
	buf := bytes.NewReader(tbs[1:])

	// read uid
	if err := binary.Read(buf, binary.BigEndian, &ut.Uid); err != nil {
		return false
	}

	// ip
	var ip net.IP
	buf.Read([]byte(ip))
	ut.ClientIP = ip.String()

	// gen time
	if err := binary.Read(buf, binary.BigEndian, &ut.GenTime); err != nil {
		return false
	}

	// gen time
	if err := binary.Read(buf, binary.BigEndian, &ut.Expire); err != nil {
		return false
	}

	return true
}

type authService struct {
	crypto *kits.Crypto
}

func (cas *authService) GenUserToken(clientip string, uid int64, expire int, flag TokenFlag) (string, error) {
	ut := &UserToken{
		Flag:     flag,
		Uid:      uid,
		ClientIP: clientip,
		GenTime:  time.Now().Unix(),
		Expire:   expire,
	}

	return ut.GenToken(cas.crypto)
}

func (cas *authService) GenSysToken(clientip string, expire int) (string, error) {
	ut := &UserToken{
		Flag:     TokenSys,
		Uid:      sysUserId,
		ClientIP: clientip,
		GenTime:  time.Now().Unix(),
		Expire:   expire,
	}

	return ut.GenToken(cas.crypto)
}

func (cas *authService) Authenticate(clientip, token string) (*UserToken, bool) {
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

	if time.Now().Unix()-ut.GenTime > int64(ut.Expire*60) {
		log.Errorf("token %s is expired", token)
		return nil, false
	}

	switch ut.Flag {
	case TokenSys:
		if ut.Uid != sysUserId {
			log.Errorf("Uid(%v!=%v) is invalid.", ut.Uid, sysUserId)
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

func (cas *authService) Authorize(ut *UserToken, ownerUid int64, oid OperationId) bool {
	if ut != nil {
		return false
	}

	// only support operate owner resource
	switch ut.Flag {
	case TokenUser:
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
