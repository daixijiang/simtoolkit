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
  "verbose": 0,
  "simfake": 0,
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
	Verbose int
	Simfake int
	Scaling float64
	Log     config_log
	Token   config_token
	Server  config_server
	Produce config_produce
	Serial  config_serial
}

var gConfig SysConfig

const CONFIG_PATH string = "./"
const CONFIG_NAME string = "simconfig"
const CONFIG_TYPE string = "json"

func config_set_default(config *viper.Viper) {
	config.SetDefault("scaling", "1.3")
	config.SetDefault("verbose", "0")
	config.SetDefault("simfake", "0")

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

	config.SetDefault("serial.serial_max", "8")
	config.SetDefault("serial.serial_timeout", "3000")
	config.SetDefault("serial.serial_timewait", "200")

	config.SetDefault("produce.timeout_cold_reset", "30")
	config.SetDefault("produce.timeout_hot_reset", "5")
	config.SetDefault("produce.timeout_creg", "3")
	config.SetDefault("produce.timeout_common", "1")
}

func config_get_value(config *viper.Viper) {
	config.GetFloat64("scaling")
	config.GetInt("verbose")
	config.GetInt("simfake")

	config.GetString("log.level")
	config.GetString("log.file")
	config.GetInt("log.maxday")

	config.GetInt("token.max")
	config.GetString("token.cmcc_file")
	config.GetString("token.uni_file")
	config.GetString("token.tel_file")

	config.GetString("server.plain_url")
	config.GetString("server.cipher_url")
	config.GetString("server.cipherv1_url")
	config.GetString("server.cipherv3_url")

	config.GetInt("serial.serial_max")
	config.GetInt("serial.serial_timeout")
	config.GetInt("serial.serial_timewait")

	config.GetInt("produce.timeout_cold_reset")
	config.GetInt("produce.timeout_hot_reset")
	config.GetInt("produce.timeout_creg")
	config.GetInt("produce.timeout_common")
}

func config_print_value(gConfig SysConfig) {
	fmt.Printf("---------------------------------\n")
	fmt.Printf("version:                \t%s\n", szVersion)
	fmt.Printf("scaling:                \t%f\n", gConfig.Scaling)
	fmt.Printf("verbose:                \t%d\n", gConfig.Verbose)
	fmt.Printf("simfake:                \t%d\n", gConfig.Simfake)

	fmt.Printf("log.level:              \t%s\n", gConfig.Log.Level)
	fmt.Printf("log.file:               \t%s\n", gConfig.Log.File)
	fmt.Printf("log.maxday:             \t%d\n", gConfig.Log.Maxday)

	fmt.Printf("token.max:              \t%d\n", gConfig.Token.Max)
	fmt.Printf("token.cmcc_file:        \t%s\n", gConfig.Token.Cmcc_file)
	fmt.Printf("token.uni_file:         \t%s\n", gConfig.Token.Uni_file)
	fmt.Printf("token.tel_file:         \t%s\n", gConfig.Token.Tel_file)

	//fmt.Printf("server.plain_url:       \t%s\n", gConfig.Server.Plain_url)
	//fmt.Printf("server.cipher_url:      \t%s\n", gConfig.Server.Cipher_url)
	//fmt.Printf("server.cipherv1_url:    \t%s\n", gConfig.Server.Cipherv1_url)
	//fmt.Printf("server.cipherv3_url:    \t%s\n", gConfig.Server.Cipherv3_url)

	fmt.Printf("serial.serial_max:      \t%d\n", gConfig.Serial.Serial_max)
	fmt.Printf("serial.serial_timeout:  \t%d\n", gConfig.Serial.Serial_timeout)
	fmt.Printf("serial.serial_timewait: \t%d\n", gConfig.Serial.Serial_timewait)

	fmt.Printf("produce.timeout_cold_reset: \t%d\n", gConfig.Produce.Timeout_cold_reset)
	fmt.Printf("produce.timeout_hot_reset: \t%d\n", gConfig.Produce.Timeout_hot_reset)
	fmt.Printf("produce.timeout_creg:   \t%d\n", gConfig.Produce.Timeout_creg)
	fmt.Printf("produce.timeout_common: \t%d\n", gConfig.Produce.Timeout_common)

	fmt.Printf("---------------------------------\n")
}

func config_init() {
	config := viper.New()
	config.AddConfigPath(CONFIG_PATH)
	config.SetConfigName(CONFIG_NAME)
	config.SetConfigType(CONFIG_TYPE)

	//set default value
	config_set_default(config)

	// read config
	if err := config.ReadInConfig(); err != nil {
		fmt.Printf("Faile to read %s%s.%s\n", CONFIG_PATH, CONFIG_NAME, CONFIG_TYPE)
	} else {
		//get value
		config_get_value(config)
	}

	// json to struct
	if err := config.Unmarshal(&gConfig); err != nil {
		fmt.Println(err)
	}

	config_print_value(gConfig)
}

func log_init() {
	vlog.InitLog("file", gConfig.Log.File, gConfig.Log.Level, int64(gConfig.Log.Maxday))
}
