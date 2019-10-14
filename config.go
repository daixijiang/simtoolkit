/* VSIM Serial Product Toolkit
 * Author: daixijiang@gmail.com (2019)
 */
package main

import (
	"fmt"
	"github.com/spf13/viper"
	"vlog"
)

/*

{
  "scaling": 1.5,
  "testflag": 0,
  "log": {
    "level": "info",
    "file": "vsimkit.log",
    "maxday": 7
  },
  "token": {
    "cmcc_file": "token.cfg",
    "uni_file": "token_uni.cfg",
    "tel_file": "token_tel.cfg"
  },
  "server": {
    "plain_url": "https://rdp.showmac.cn/api/v1/profile/clear/get",
    "cipher_url": "https://ldp.showmac.cn/api/openluat/profile",
    "cipherv1_url": "https://rdp.showmac.cn/api/v1/profile/get",
    "cipherv3_url": "https://rdp.showmac.cn/api/v3/profile/get"
  },
  "serial": {
    "serial_max": 10,
    "serial_timeout": 3000,
    "serial_timewait": 200
  },
  "produce": {
    "timeout_cold_reset": 1,
    "timeout_hot_reset": 1,
    "timeout_creg": 3,
    "timeout_common": 1
  }
}

*/

type config_log struct {
	Level  string
	File   string
	Maxday int
}

type config_token struct {
	Max       int
	Cmcc_file string
	Uni_file  string
	Tel_file  string
}

type config_server struct {
	Plain_url    string
	Cipher_url   string
	Cipherv1_url string
	Cipherv3_url string
}

type config_produce struct {
	Timeout_cold_reset int
	Timeout_hot_reset  int
	Timeout_creg       int
	Timeout_common     int
}

type config_serial struct {
	Serial_max      int
	Serial_timeout  int
	Serial_timewait int
}

type SysConfig struct {
	Testflag int
	Scaling  float64
	Log      config_log
	Token    config_token
	Server   config_server
	Produce  config_produce
	Serial   config_serial
}

var gConfig SysConfig

const CONFIG_PATH string = "./"
const CONFIG_NAME string = "simconfig"
const CONFIG_TYPE string = "json"

func config_init() {
	config := viper.New()
	config.AddConfigPath(CONFIG_PATH)
	config.SetConfigName(CONFIG_NAME)
	config.SetConfigType(CONFIG_TYPE)

	//set default value
	config.SetDefault("scaling", "1.3")
	config.SetDefault("testflag", "0")

	config.SetDefault("log.level", "info")
	config.SetDefault("log.file", CONFIG_PATH+"vsimkit.log")
	config.SetDefault("log.maxday", "7")

	config.SetDefault("token.max", "100")
	config.SetDefault("token.cmcc_file", CONFIG_PATH+"token.cfg")
	config.SetDefault("token.uni_file", CONFIG_PATH+"token_uni.cfg")
	config.SetDefault("token.tel_file", CONFIG_PATH+"token_tel.cfg")

	config.SetDefault("server.plain_url", "https://rdp.showmac.cn/api/v1/profile/clear/get")
	config.SetDefault("server.cipher_url", "https://ldp.showmac.cn/api/openluat/profile")
	config.SetDefault("server.cipherv1_url", "https://rdp.showmac.cn/api/v1/profile/get")
	config.SetDefault("server.cipherv3_url", "https://rdp.showmac.cn/api/v3/profile/get")

	config.SetDefault("produce.timeout_cold_reset", "30")
	config.SetDefault("produce.timeout_hot_reset", "5")
	config.SetDefault("produce.timeout_creg", "3")
	config.SetDefault("produce.timeout_common", "1")

	config.SetDefault("serial.serial_max", "8")
	config.SetDefault("serial.serial_timeout", "3000")
	config.SetDefault("serial.serial_timewait", "200")

	// read config
	if err := config.ReadInConfig(); err != nil {
		fmt.Printf("Faile to read %s%s.%s !\n", CONFIG_PATH, CONFIG_NAME, CONFIG_TYPE)
	}

	{
		//get value
		fmt.Printf("scaling: %f\n", config.GetFloat64("scaling"))
		fmt.Printf("testflag: %d\n", config.GetInt("testflag"))

		fmt.Printf("log.level: %s\n", config.GetString("log.level"))
		fmt.Printf("log.file: %s\n", config.GetString("log.file"))
		fmt.Printf("log.maxday: %d\n", config.GetInt("log.maxday"))

		fmt.Printf("token.max: %d\n", config.GetInt("token.max"))
		fmt.Printf("token.cmcc_file: %s\n", config.GetString("token.cmcc_file"))
		fmt.Printf("token.uni_file: %s\n", config.GetString("token.uni_file"))
		fmt.Printf("token.tel_file: %s\n", config.GetString("token.tel_file"))

		fmt.Printf("server.plain_url: %s\n", config.GetString("server.plain_url"))
		fmt.Printf("server.cipher_url: %s\n", config.GetString("server.cipher_url"))
		fmt.Printf("server.cipherv1_url: %s\n", config.GetString("server.cipherv1_url"))
		fmt.Printf("server.cipherv3_url: %s\n", config.GetString("server.cipherv3_url"))

		fmt.Printf("serial.serial_max: %d\n", config.GetInt("serial.serial_max"))
		fmt.Printf("serial.serial_timeout: %d\n", config.GetInt("serial.serial_timeout"))
		fmt.Printf("serial.serial_timewait: %d\n", config.GetInt("serial.serial_timewait"))

		fmt.Printf("produce.timeout_cold_reset: %d\n", config.GetInt("produce.timeout_cold_reset"))
		fmt.Printf("produce.timeout_hot_reset: %d\n", config.GetInt("produce.timeout_hot_reset"))
		fmt.Printf("produce.timeout_creg: %d\n", config.GetInt("produce.timeout_creg"))
		fmt.Printf("produce.timeout_common: %d\n", config.GetInt("produce.timeout_common"))
	}

	//直接反序列化为Struct
	if err := config.Unmarshal(&gConfig); err != nil {
		fmt.Println(err)
	}

	//fmt.Println(gConfig)
}

func log_init() {
	vlog.InitLog("file", gConfig.Log.File, gConfig.Log.Level, int64(gConfig.Log.Maxday))
}
