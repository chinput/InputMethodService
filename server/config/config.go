package config

import (
	"encoding/base64"
	"log"
	"os"
	"path/filepath"
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

	// 目录格式
	// 比如 data/2006/01/02/
	TimePattern string `toml:"time_pattern"`

	// des 加密模式 ECB/CBC
	EncType string `toml:"enc_type"`
}

var cfg *Config

const (
	kECB = "ECB"
	kCBC = "CBC"
)

var (
	kSecretKey []byte
)

func init() {
	cfg = new(Config)
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	configPath := filepath.Join(dir, "config.toml")
	toml.DecodeFile(configPath, cfg)

	kSecretKey = []byte(cfg.SecretKey)
}

func DBUrl() string {
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
