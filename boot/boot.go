package boot

import (
	"os"
	"path/filepath"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/xtfly/goman/kits"
	"github.com/xtfly/goman/models"
	"gopkg.in/ini.v1"
	"gopkg.in/macaron.v1"
)

var (
	cfg *ini.File
)

var (
	//
	WebListenIP = ""
	WebPort     = 8080

	// app work path
	AppPath = ""

	// System setting
	SysSetting *models.GlobalSetting

	// a global AES crypte object
	crypto *kits.Crypto
)

// initialize by the secure config
func init() {
	macaron.SetConfig(getCfgFile())
	cfg = macaron.Config()

	getWebCfg()
	getSecureCfg()

	models.InitDB(crypto)
}

func BootStrap() {
	models.ConnectDB()
	SysSetting = models.NewGlobalSetting()
	SysSetting.LoadAll()
}

func getWebCfg() {
	web, err := cfg.GetSection("web")
	if err != nil {
		panic(err)
	}

	WebListenIP = web.Key("ip").MustString("0.0.0.0")
	WebPort = web.Key("port").MustInt(8080)
}

func getSecureCfg() {
	secure, err := cfg.GetSection("secure")
	if err != nil {
		panic(err)
	}

	factor := secure.Key("factor").String()
	crc := secure.Key("crc").String()

	crypto, err = kits.NewCrypto(factor, crc)
	if err != nil {
		panic(err)
	}
}

// Get the config path
func getCfgFile() string {
	workPath, _ := os.Getwd()
	workPath, _ = filepath.Abs(workPath)

	// initialize default configurations
	AppPath, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	configPath := filepath.Join(AppPath, "conf", "app.ini")

	if workPath != AppPath {
		if kits.FileExists(configPath) {
			os.Chdir(AppPath)
		} else if strings.HasSuffix(workPath, "goman") {
			configPath = filepath.Join(workPath, "conf", "app.ini")
		} else {
			configPath = filepath.Join(workPath, "../conf", "app.ini")
		}
	}

	log.Debug("config path=", configPath)
	return configPath
}

// retrive the Crypto using defalt config
func GetCrypto() *kits.Crypto {
	return crypto
}
