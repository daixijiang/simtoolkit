/* VSIM Serial Product Toolkit
 * Author: daixijiang@gmail.com (2019)
 */

package main

import (
	"github.com/tarm/serial"
)

/* ModuleTable */
const (
	// prepare
	Module_CMD1_VER             = 10
	Module_CMD1_SOFTMODE        = 11
	Module_CMD1_REBOOT          = 12
	Module_PRE1_CMD_NETWORK_SET = 13
	Module_CMD1_URL_SET         = 14
	Module_CMD1_AUTOSWITCH_ON   = 15
	Module_CMD1_AUTOSWITCH_OFF  = 16
	Module_CMD1_CFG_BACKUP      = 17

	// produce
	Module_CMD2_IMEI   = 20
	Module_CMD2_VER    = 21
	Module_CMD2_CHIPID = 22
	Module_CMD2_SIM192 = 23
	Module_CMD2_SIM64  = 24

	// check
	Module_CMD3_CCID       = 30
	Module_CMD3_CIMI       = 31
	Module_CMD3_CREG       = 32
	Module_CMD3_CEREG      = 33
	Module_CMD3_COPS       = 34
	Module_CMD3_SWITCH_CM  = 35
	Module_CMD3_SWITCH_CU  = 36
	Module_CMD3_SWITCH_TEL = 37

	Module_TAB_AT_CMD_MAX = 40
)

type cmdHandler func(cmdid int, portid int, s *serial.Port, reply *string) int

type ModuleTable struct {
	CmdID   int
	CmdStr  string
	CmdFunc cmdHandler
}

var modEC20 [Module_TAB_AT_CMD_MAX]ModuleTable
var myMod *[Module_TAB_AT_CMD_MAX]ModuleTable
var myUrlVer int

func module_init() {
	module_ec20_init()
	// add list of module init

	// default is ec20
	myMod = &modEC20
	myUrlVer = SERVER_PLAIN_v0
}

/* ModuleTable end */
