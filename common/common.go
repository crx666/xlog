package common

import (
	"errors"
	"fmt"
	"github.com/crx666/xlog/config"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"math/rand"
	"net"
	"os"
	"reflect"
	"strings"
	"time"
	"unsafe"
)

const (
	LogTemp   = ".temp"
	LogFormal = ".log"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func getClientIp() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "127.0.0.1"
	}
	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return "127.0.0.1"
}

func LogConfigCheck(config *config.LogConfig) error {
	if config.LogName == "" && config.LogDir == "" && !config.IsConsole {
		return errors.New("log config output set error")
	}

	if config.LogDir == "" && config.LogName != "" {
		config.LogDir = "./log"
	}
	if config.LogDir != "" && config.LogName == "" {
		config.LogName = "app"
	}

	//if config.LogDir != "" && config.LogName != "" {
	//	dir := ReplaceDir(config.LogDir)
	//	exist, err := IsPathExist(dir)
	//	if err != nil {
	//		return errors.New(fmt.Sprintf("Failed to check whether a path exist or not, DirPath=%s",
	//			dir))
	//	}
	//	if !exist {
	//		err = os.MkdirAll(dir, os.ModePerm)
	//		if err != nil {
	//			return errors.New(fmt.Sprintf("Failed to check whether a path exist or not,  DirPath=%s",
	//				dir))
	//		}
	//	}
	//	config.LogDir = dir
	//	//config.LogName = ReplaceName(config.LogName)
	//	//if config.ErrLogName != "" {
	//	//	config.ErrLogName = ReplaceName(config.ErrLogName)
	//	//}
	//}

	return nil
}

func ReplaceDir(dir string) string {
	now := time.Now()
	for strings.Contains(dir, "$") {
		if strings.Contains(dir, "$ip") {
			dir = strings.ReplaceAll(dir, "$ip", fmt.Sprintf("%s", getClientIp()))
		} else if strings.Contains(dir, "$rand") {
			dir = strings.ReplaceAll(dir, "$rand", fmt.Sprintf("%s", getRandString(6)))
		} else if strings.Contains(dir, "$date") {
			dir = strings.ReplaceAll(dir, "$date", fmt.Sprintf("%s", now.Format("2006-01-02")))
		} else {
			break
		}
	}
	return dir
}

func ReplaceName(name string) string {
	for strings.Contains(name, "$") {
		if strings.Contains(name, "$ip") {
			name = strings.ReplaceAll(name, "$ip", fmt.Sprintf("%s", getClientIp()))
		} else if strings.Contains(name, "$rand") {
			name = strings.ReplaceAll(name, "$rand", fmt.Sprintf("%s", getRandString(5)))
		} else if strings.Contains(name, "$ti") {
			name = strings.ReplaceAll(name, "$ti", fmt.Sprintf("%d", time.Now().Unix()))
		} else if strings.Contains(name, "$day") {
			name = strings.ReplaceAll(name, "$day", fmt.Sprintf("%s", "%Y_%m_%d"))
		} else if strings.Contains(name, "$hour") {
			name = strings.ReplaceAll(name, "$hour", fmt.Sprintf("%s", "%Y_%m_%d_%H"))
		} else if strings.Contains(name, "$minute") {
			name = strings.ReplaceAll(name, "$minute", fmt.Sprintf("%s", "%Y_%m_%d_%H_%M"))
		} else {
			break
		}
	}
	return name
}

func IsPathExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		} else {
			return false, err
		}
	}
	return true, nil
}

func ParserYamlData(path string, config interface{}) { //config 必须是个指针对象
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(yamlFile, config)
	if err != nil {
		panic(err)
	}
}

func String2Bytes(s string) []byte {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := reflect.SliceHeader{
		Data: sh.Data,
		Len:  sh.Len,
		Cap:  sh.Len,
	}
	return *(*[]byte)(unsafe.Pointer(&bh))
}

func Bytes2String(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

const (
	letterBytes   = "abcdefghijklmnopqrstuvwxyz1234567890"
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func getRandString(n int) string {
	b := make([]byte, n)
	// A rand.Int63() generates 63 random bits, enough for letterIdxMax letters!
	for i, cache, remain := n-1, rand.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = rand.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return string(b)
}

func ReplaceLogName(name string) error {
	if strings.Contains(name, LogTemp) {
		newName := strings.ReplaceAll(name, LogTemp, LogFormal)
		err := os.Rename(name, newName)
		if err != nil {
			return err
		}
	}
	return nil
}

func GetBaseLogConfig() *config.LogConfig {
	c := new(config.LogConfig)
	ParserYamlData("./config/base_log.yaml", c)
	return c
}
