package token

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"net"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/xtfly/gokits"
)

// the type of token
type TokenFlag byte

// flag(1) + uid(8) + ip(4) + currenttime(8) + expire(4)
type UserToken struct {
	Flag     TokenFlag
	Uid      int64
	ClientIP string // the ip of the request
	GenTime  int64  // timestamp when generating token
	Expire   int32  // unit: second
}

const (
	TokenLen             = 1 + 8 + 4 + 8 + 4
	TokenSys   TokenFlag = 0x00 // system general user
	TokenUser  TokenFlag = 0x01 // login user
	TokenAdmin TokenFlag = 0x01

	SysUserId int64 = int64(-8619820608)
)

func (ut *UserToken) GenToken(crypto *gokits.Crypto) (string, error) {
	bs := make([]byte, 0)
	buf := bytes.NewBuffer(bs)

	// flag
	buf.WriteByte(byte(ut.Flag))

	// uid
	binary.Write(buf, binary.BigEndian, ut.Uid)

	// client ip
	ip := net.ParseIP(ut.ClientIP)
	//log.Info("token length, ip  ", ut.ClientIP)
	buf.Write([]byte(ip.To4()))

	// current time
	binary.Write(buf, binary.BigEndian, time.Now().Unix())

	// expire
	binary.Write(buf, binary.BigEndian, ut.Expire)

	// encrypt
	//log.Info("token length, expire ", buf.Len())
	if t, err := crypto.Encrypt(buf.Bytes()); err != nil {
		return "", err
	} else {
		return base64.URLEncoding.EncodeToString(t), nil
	}
}

func (ut *UserToken) DecodeToken(crypto *gokits.Crypto, token string) bool {
	bs, err := base64.URLEncoding.DecodeString(token)
	if err != nil {
		return false
	}

	// decrypt
	var tbs []byte
	tbs, err = crypto.Decrypt(bs)
	if err != nil {
		log.Error("Decrypt failed ", err.Error())
		return false
	}

	// check the len
	if l := len(tbs); l != TokenLen {
		log.Error("Invalid length ", l)
		return false
	}

	ut.Flag = TokenFlag(tbs[0])
	// remain bytes
	buf := bytes.NewReader(tbs[1:])

	// read uid
	if err := binary.Read(buf, binary.BigEndian, &ut.Uid); err != nil {
		log.Error("Read uid failed ", err.Error())
		return false
	}

	// ip
	ip := make([]byte, 4)
	buf.Read(ip)
	ut.ClientIP = net.IP(ip).String()
	//log.Info("token ip  ", ut.ClientIP)

	// gen time
	if err := binary.Read(buf, binary.BigEndian, &ut.GenTime); err != nil {
		log.Error("Read generate time failed ", err.Error())
		return false
	}

	// expire time
	if err := binary.Read(buf, binary.BigEndian, &ut.Expire); err != nil {
		log.Error("Read expire time failed ", err.Error())
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
