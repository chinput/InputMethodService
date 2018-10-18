package config

import (
	"encoding/base64"
	"time"

	"github.com/BurntSushi/toml"
)

type Config struct {
	DBHost string `toml:"db_host"`
	DBName string `toml:"db_name"`
	DBUser string `toml:"db_user"`
	DBPass string `toml:"db_pass"`

	// 用来加密用的秘钥
	SecretKey string `toml:"secret_key"`
	// des 加密模式 ECB/CBC
	EncType string `toml:"enc_type"`

	// 目录格式
	// 比如 data/2006/01/02/
	TimePattern string `toml:"time_pattern"`

	// 邮箱相关
	EmailHost string `toml:"email_host"`
	EmailUser string `toml:"email_user"`
	EmailPass string `toml:"email_pass"`
	// 邮件发送最短间隔，单位 ms
	EmailTimeout int64 `toml:"email_timeout"`

	// 存储用户数据文件的目录
	DataPath string `toml:"data_path"`

	// 用于下载文件的临时目录，5分钟清理一次
	TmpPath string `toml:"tmp_path"`
}

var cfg *Config

const (
	kECB = "ECB"
	kCBC = "CBC"
)

var (
	kSecretKey []byte
)

func Init(configPath string) {
	cfg = new(Config)
	toml.DecodeFile(configPath, cfg)
	kSecretKey = []byte(cfg.SecretKey)
}

func DBUrl() string {
	if cfg.DBUser == "" || cfg.DBPass == "" {
		return cfg.DBHost
	}
	return cfg.DBUser + ":" + cfg.DBPass + "@" + cfg.DBHost
}

func DBName() string {
	return cfg.DBName
}

func PatternTime(pre string, t time.Time, ext string) string {
	return pre + t.Format(cfg.TimePattern) + ext
}

func Encode(data []byte) ([]byte, error) {

	if cfg.EncType == kECB {
		return DesECBEncrypt(data, kSecretKey)
	}
	return DesCBCEncrypt(data, kSecretKey)
}

func Decode(data []byte) ([]byte, error) {
	if cfg.EncType == kECB {
		return DesECBDecrypt(data, kSecretKey)
	}

	return DesCBCDecrypt(data, kSecretKey)
}

func EncodeB64(data []byte) (string, error) {
	tmp, err := Encode(data)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(tmp), nil

}

func DecodeB64(data string) ([]byte, error) {
	tmp, err := base64.URLEncoding.DecodeString(data)
	if err != nil {
		return nil, err
	}

	return Decode(tmp)
}

func EmailConfig() (string, string, string, int64) {
	if cfg == nil {
		return "", "", "", 0
	}
	return cfg.EmailHost, cfg.EmailUser, cfg.EmailPass, cfg.EmailTimeout
}

func TmpPath() string {
	return cfg.TmpPath
}

func DataPath() string {
	return cfg.DataPath
}
