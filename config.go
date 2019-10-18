/* VSIM Serial Product Toolkit
 * Author: daixijiang@gmail.com (2019)
 */
package main

import (
	"fmt"
	"github.com/koding/multiconfig"
	"vlog"
)

/*
1. config.json:
--------------------------------------------------------
{
  "scaling": 1.5,
  "verbose": 0,
  "module": "sim800c",
  "simfake": 0,
  "uistyle": "nucular",
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
    "conn_timeout": 25,
    "rw_timeout": 20,
    "plain_url": "https://rdp.showmac.cn/api/v1/profile/clear/get",
    "cipher_url": "https://ldp.showmac.cn/api/openluat/profile",
    "cipherv1_url": "https://rdp.showmac.cn/api/v1/profile/get",
    "cipherv3_url": "https://rdp.showmac.cn/api/v3/profile/get"
  },
  "serial": {
    "serial_max": 10,
    "cmd_timeout": 3,
    "cmd_timewait": 0.5
  },
  "produce": {
    "cold_reset_timeout": 0,
    "hot_reset_timeout": 0,
    "creg_timeout": 3,
    "common_timeout": 1
  }
}
--------------------------------------------------------

2. config.toml:
--------------------------------------------------------
scaling = 1.5
verbose = 1
#module = sim800c
#simfake = 1
uistyle = nucular

[log]
level = info
file = vsimkit.log
maxday = 7

[token]
cmcc_file = token.cfg
uni_file = token_uni.cfg
tel_file = token_tel.cfg

[server]
conn_timeout = 25
rw_timeout = 20
plain_url = https://rdp.showmac.cn/api/v1/profile/clear/get
cipher_url = https://ldp.showmac.cn/api/openluat/profile
cipherv1_url = https://rdp.showmac.cn/api/v1/profile/get
cipherv3_url = https://rdp.showmac.cn/api/v3/profile/get

[serial]
serial_max = 8
cmd_timeout = 3
cmd_timewait = 0.5

[produce]
cold_reset_timeout = 0
hot_reset_timeout = 0
creg_timeout = 3
common_timeout = 1
--------------------------------------------------------
*/

const CONFIG_PATH string = "./"
const CONFIG_NAME string = "simconfig.toml"

type config_log struct {
	Level  string `default:"info"`
	File   string `default:"vsimkit.log"`
	Maxday int    `default:"7"`
}

type config_token struct {
	Max       int    `default:"100"`
	Cmcc_file string `default:"token.cfg"`
	Uni_file  string `default:"token_uni.cfg"`
	Tel_file  string `default:"token_tel.cfg"`
}

type config_server struct {
	Conn_timeout int    `default:"25"`
	Rw_timeout   int    `default:"20"`
	Plain_url    string `default:"https://rdp.showmac.cn/api/v1/profile/clear/get"`
	Cipher_url   string `default:"https://ldp.showmac.cn/api/openluat/profile"`
	Cipherv1_url string `default:"https://rdp.showmac.cn/api/v1/profile/get"`
	Cipherv3_url string `default:"https://rdp.showmac.cn/api/v3/profile/get"`
}

type config_serial struct {
	Serial_max   int     `default:"8"`
	Cmd_timeout  int     `default:"3"`
	Cmd_timewait float64 `default:"0.5"`
}

type config_produce struct {
	Cold_reset_timeout int `default:"30"`
	Hot_reset_timeout  int `default:"5"`
	Creg_timeout       int `default:"3"`
	Common_timeout     int `default:"1"`
}

type SysConfig struct {
	Scaling float64 `default:"1.3"`
	Verbose int     `default:"0"`
	Simfake int     `default:"0"`
	Module  string  `default:"sim800c"`
	UIstyle string  `default:"nucular"`
	Log     config_log
	Token   config_token
	Server  config_server
	Produce config_produce
	Serial  config_serial
}

var gConfig SysConfig

func config_print_value(gConfig *SysConfig) {
	fmt.Printf("---------------------------------\n")
	fmt.Printf("version:                \t%s\n", szVersion)
	fmt.Printf("scaling:                \t%f\n", gConfig.Scaling)
	fmt.Printf("verbose:                \t%d\n", gConfig.Verbose)
	fmt.Printf("simfake:                \t%d\n", gConfig.Simfake)
	fmt.Printf("module:                 \t%s\n", gConfig.Module)
	fmt.Printf("uistyle:                \t%s\n", gConfig.UIstyle)

	fmt.Printf("log.level:              \t%s\n", gConfig.Log.Level)
	fmt.Printf("log.file:               \t%s\n", gConfig.Log.File)
	fmt.Printf("log.maxday:             \t%d\n", gConfig.Log.Maxday)

	fmt.Printf("token.max:              \t%d\n", gConfig.Token.Max)
	fmt.Printf("token.cmcc_file:        \t%s\n", gConfig.Token.Cmcc_file)
	fmt.Printf("token.uni_file:         \t%s\n", gConfig.Token.Uni_file)
	fmt.Printf("token.tel_file:         \t%s\n", gConfig.Token.Tel_file)

	fmt.Printf("server.conn_timeout:     \t%d\n", gConfig.Server.Conn_timeout)
	fmt.Printf("server.rw_timeout:       \t%d\n", gConfig.Server.Rw_timeout)
	if false {
		fmt.Printf("server.plain_url:       \t%s\n", gConfig.Server.Plain_url)
		fmt.Printf("server.cipher_url:      \t%s\n", gConfig.Server.Cipher_url)
		fmt.Printf("server.cipherv1_url:    \t%s\n", gConfig.Server.Cipherv1_url)
		fmt.Printf("server.cipherv3_url:    \t%s\n", gConfig.Server.Cipherv3_url)
	}

	fmt.Printf("serial.serial_max:      \t%d\n", gConfig.Serial.Serial_max)
	fmt.Printf("serial.cmd_timeout:     \t%d\n", gConfig.Serial.Cmd_timeout)
	fmt.Printf("serial.cmd_timewait:    \t%f\n", gConfig.Serial.Cmd_timewait)

	fmt.Printf("produce.cold_reset_timeout: \t%d\n", gConfig.Produce.Cold_reset_timeout)
	fmt.Printf("produce.hot_reset_timeout: \t%d\n", gConfig.Produce.Hot_reset_timeout)
	fmt.Printf("produce.creg_timeout:   \t%d\n", gConfig.Produce.Creg_timeout)
	fmt.Printf("produce.common_timeout: \t%d\n", gConfig.Produce.Common_timeout)

	fmt.Printf("---------------------------------\n")
}

func config_init() {
	m := multiconfig.NewWithPath(CONFIG_PATH + CONFIG_NAME)

	err := m.Load(&gConfig)
	if err != nil {
		fmt.Printf("Load configure %s error: %s\n", CONFIG_PATH+CONFIG_NAME, err)
	} else {
		m.MustLoad(&gConfig)
		//fmt.Printf("%+v\n", gConfig)
	}

	config_print_value(&gConfig)
}

func log_init() {
	vlog.InitLog("file", gConfig.Log.File, gConfig.Log.Level, int64(gConfig.Log.Maxday))
}
