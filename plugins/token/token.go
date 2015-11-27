package token

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"net"
	"time"

	"github.com/xtfly/goman/kits"
)

// the type of token
type TokenFlag byte

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

	SysUserId int64 = int64(-8619820608)
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

func (ut *UserToken) Expired() bool {
	if time.Now().Unix()-ut.GenTime > int64(ut.Expire*60) {
		return true
	}

	return false
}

// 每隔2分刷新token
func (ut *UserToken) NeedRefresh() bool {
	ivl := time.Now().Unix() - ut.GenTime
	if ivl < int64(ut.Expire*60) && ivl > 2*60 {
		return true
	}

	return false
}
