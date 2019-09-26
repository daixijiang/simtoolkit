/* VSIM Serial Product Toolkit
 * Author: daixijiang@gmail.com (2019)
 */

package main

import (
	"strings"
	"vlog"

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

var modEC20 [100]ModuleTable
var modSIM800C [100]ModuleTable

func module_init() {
	module_ec20_init()
	module_sim800c_init()
}

func module_ec20_init() {
	modEC20[Module_CMD2_IMEI] = ModuleTable{
		Module_CMD2_IMEI,
		"AT+GSN",
		ec20_get_imei,
	}

	modEC20[Module_CMD2_VER] = ModuleTable{
		Module_CMD2_VER,
		"AT+CSIM=10,\"A0B009090C\"",
		ec20_get_ver,
	}

	modEC20[Module_CMD2_CHIPID] = ModuleTable{
		Module_CMD2_CHIPID,
		"AT+CSIM=10,\"A0B0090910\"",
		ec20_get_chipid,
	}
}

func module_sim800c_init() {

}

/* "AT+GSN\r\r\n
 * 862107043586551\r\n\r\n
 * OK\r\n"
 */
func ec20_get_imei(cmdid int, portid int, s *serial.Port, reply *string) int {
	resp := serialWriteAndEcho(portid, s, modEC20[cmdid].CmdStr)
	rs := []rune(resp)
	length := len(rs)
	sublen := len(modEC20[cmdid].CmdStr)

	vlog.Info("%d %q", length, []byte(resp))
	pos1 := strings.Index(resp, modEC20[cmdid].CmdStr+"\r\r\n")
	pos2 := strings.Index(resp, "\r\n\r\nOK")
	if pos1 >= 0 && pos2 >= 0 {
		preresp := string(rs[sublen+len("\r\r\n") : pos2])
		vlog.Info("%d %s", len(preresp), preresp)
		*reply = preresp
		return len(*reply)
	}

	return 0
}

/* "AT+CSIM=10,"A0B009090C"\r\r\n
 * +CSIM: 28,"436F735665725F312E312E349000"\r\n\r\n
 * OK\r\n"
 */

func ec20_get_ver(cmdid int, portid int, s *serial.Port, reply *string) int {
	resp := serialWriteAndEcho(portid, s, modEC20[cmdid].CmdStr)
	rs := []rune(resp)
	length := len(rs)
	sublen := len(modEC20[cmdid].CmdStr)

	vlog.Info("%d %q", length, []byte(resp))
	pos1 := strings.Index(resp, modEC20[cmdid].CmdStr+"\r\r\n")
	pos2 := strings.Index(resp, "\r\n\r\nOK")
	if pos1 >= 0 && pos2 >= 0 {
		preresp := string(rs[(sublen + len("\r\r\n") + len("+CSIM: 36,\"")) : pos2-1])
		vlog.Info("%d %s", len(preresp), preresp)
		*reply = preresp
		return len(preresp)
	}

	return 0
}

/* "AT+CSIM=10,"A0B0090910"\r\r\n
 * +CSIM: 36,"3934363531303236320A3A373B3C3A3B9000"\r\n\r\n
 * OK\r\n"
 */
func ec20_get_chipid(cmdid int, portid int, s *serial.Port, reply *string) int {
	resp := serialWriteAndEcho(portid, s, modEC20[cmdid].CmdStr)
	rs := []rune(resp)
	length := len(rs)
	sublen := len(modEC20[cmdid].CmdStr)

	vlog.Info("%d %q", length, []byte(resp))
	pos1 := strings.Index(resp, modEC20[cmdid].CmdStr+"\r\r\n")
	pos2 := strings.Index(resp, "\r\n\r\nOK")
	if pos1 >= 0 && pos2 >= 0 {
		preresp := string(rs[(sublen + len("\r\r\n") + len("+CSIM: 36,\"")) : pos2-1])
		vlog.Info("%d %s", len(preresp), preresp)
		*reply = preresp
		return len(preresp)
	}

	return 0
}

/* ModuleTable end */
