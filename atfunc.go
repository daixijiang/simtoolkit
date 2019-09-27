/* VSIM Serial Product Toolkit
 * Author: daixijiang@gmail.com (2019)
 */

package main

import (
	"fmt"
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

var modEC20 [Module_TAB_AT_CMD_MAX]ModuleTable
var modSIM800C [Module_TAB_AT_CMD_MAX]ModuleTable

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

	modEC20[Module_CMD2_SIM192] = ModuleTable{
		Module_CMD2_CHIPID,
		"AT+CSIM=426",
		ec20_set_ens192,
	}

	modEC20[Module_CMD2_SIM64] = ModuleTable{
		Module_CMD2_CHIPID,
		"AT+CSIM=170",
		ec20_set_ens64,
	}

	modEC20[Module_CMD3_CCID] = ModuleTable{
		Module_CMD3_CCID,
		"AT+CCID",
		ec20_get_ccid,
	}

	modEC20[Module_CMD3_CREG] = ModuleTable{
		Module_CMD3_CREG,
		"AT+CREG?",
		ec20_get_ccid,
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

	vlog.Info("    AT cmd(%d): %q", length, []byte(resp))
	pos1 := strings.Index(resp, modEC20[cmdid].CmdStr+"\r\r\n")
	pos2 := strings.Index(resp, "\r\n\r\nOK")
	if pos1 >= 0 && pos2 >= 0 {
		preresp := string(rs[sublen+len("\r\r\n") : pos2])
		vlog.Info("    AT get(%d): %s", len(preresp), preresp)
		*reply = preresp
		return len(*reply)
	}

	return 0
}

/* "AT+CCID\r\r\n
 * 89860317492034500726\r\n\r\n
 * OK\r\n"
 */
func ec20_get_ccid(cmdid int, portid int, s *serial.Port, reply *string) int {
	resp := serialWriteAndEcho(portid, s, modEC20[cmdid].CmdStr)
	rs := []rune(resp)
	length := len(rs)
	sublen := len(modEC20[cmdid].CmdStr)

	vlog.Info("    AT cmd(%d): %q", length, []byte(resp))
	pos1 := strings.Index(resp, modEC20[cmdid].CmdStr+"\r\r\n")
	pos2 := strings.Index(resp, "\r\n\r\nOK")
	if pos1 >= 0 && pos2 >= 0 {
		preresp := string(rs[sublen+len("\r\r\n") : pos2])
		vlog.Info("    AT get(%d): %s", len(preresp), preresp)
		*reply = preresp
		return len(*reply)
	}

	return 0
}

/* "AT+CREG?\r\r\n
 * xxxxxxxx\r\n\r\n
 * OK\r\n"
 */
func ec20_get_creg(cmdid int, portid int, s *serial.Port, reply *string) int {
	resp := serialWriteAndEcho(portid, s, modEC20[cmdid].CmdStr)
	rs := []rune(resp)
	length := len(rs)
	sublen := len(modEC20[cmdid].CmdStr)

	vlog.Info("    AT cmd(%d): %q", length, []byte(resp))
	pos1 := strings.Index(resp, modEC20[cmdid].CmdStr+"\r\r\n")
	pos2 := strings.Index(resp, "\r\n\r\nOK")
	if pos1 >= 0 && pos2 >= 0 {
		preresp := string(rs[sublen+len("\r\r\n") : pos2])
		vlog.Info("    AT get(%d): %s", len(preresp), preresp)
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

	vlog.Info("    AT cmd(%d): %q", length, []byte(resp))
	pos1 := strings.Index(resp, modEC20[cmdid].CmdStr+"\r\r\n")
	pos2 := strings.Index(resp, "\r\n\r\nOK")
	if pos1 >= 0 && pos2 >= 0 {
		preresp := string(rs[(sublen + len("\r\r\n") + len("+CSIM: 28,\"")) : pos2-1])
		vlog.Info("    AT get(%d): %s", len(preresp), preresp)
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

	vlog.Info("    AT cmd(%d): %q", length, []byte(resp))
	pos1 := strings.Index(resp, modEC20[cmdid].CmdStr+"\r\r\n")
	pos2 := strings.Index(resp, "\r\n\r\nOK")
	if pos1 >= 0 && pos2 >= 0 {
		preresp := string(rs[(sublen + len("\r\r\n") + len("+CSIM: 36,\"")) : pos2-1])
		vlog.Info("    AT get(%d): %s", len(preresp), preresp)
		*reply = preresp
		return len(preresp)
	}

	return 0
}

/* "AT+CSIM=426,"A0D60000D002FD5646696C6531312E62696E000000...."\r\r\n
 * OK\r\n"
 */
const head_ens192 string = "A0D60000D002FD5646696C6531312E62696E000000"

func ec20_set_ens192(cmdid int, portid int, s *serial.Port, reply *string) int {
	cmd_ens192 := fmt.Sprintf("AT+CSIM=426,\"%s%s\"", head_ens192, serial_port[portid].sim_ens.EncData192)
	resp := serialWriteAndEcho(portid, s, cmd_ens192)
	rs := []rune(resp)
	length := len(rs)
	sublen := len(cmd_ens192)

	vlog.Info("    AT cmd(%d): %q", length, []byte(resp))
	pos1 := strings.Index(resp, cmd_ens192+"\r\r\n")
	pos2 := strings.Index(resp, "\r\n\r\nOK")
	if pos1 >= 0 && pos2 >= 0 {
		preresp := string(rs[(sublen + len("\r\r\n") + len("+CSIM:\"")) : pos2-1])
		vlog.Info("    AT get(%d): %s", len(preresp), preresp)
		*reply = preresp
		return len(preresp)
	}

	*reply = "ERROR"
	return 0
}

/* "AT+CSIM=170,\"A0D600015002FD5646696C6531312E62696E000000...\x00\"\r\r\n
 * ERROR\r\n"
 */
const head_ens64 string = "A0D600015002FD5646696C6531312E62696E000000"

func ec20_set_ens64(cmdid int, portid int, s *serial.Port, reply *string) int {
	cmd_ens64 := fmt.Sprintf("AT+CSIM=170,\"%s%s\"", head_ens64, serial_port[portid].sim_ens.EncData64)
	resp := serialWriteAndEcho(portid, s, cmd_ens64)
	rs := []rune(resp)
	length := len(rs)
	sublen := len(cmd_ens64)

	vlog.Info("    AT cmd(%d): %q", length, []byte(resp))
	pos1 := strings.Index(resp, cmd_ens64+"\r\r\n")
	pos2 := strings.Index(resp, "\r\n\r\nOK")
	if pos1 >= 0 && pos2 >= 0 {
		preresp := string(rs[(sublen + len("\r\r\n") + len("+CSIM:\"")) : pos2-1])
		vlog.Info("    AT get(%d): %s", len(preresp), preresp)
		*reply = preresp
		return len(preresp)
	}

	*reply = "ERROR"
	return 0
}

/* ModuleTable end */
