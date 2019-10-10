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

const head_ens192 string = "A0D60000D002FD5646696C6531312E62696E000000"
const head_ens64 string = "A0D600015002FD5646696C6531312E62696E000000"

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
		ec20_get_creg,
	}
}

/* "AT+GSN\r\r\n
 * 862107043586551\r\n\r\n
 * OK\r\n"
 */
func ec20_get_imei(cmdid int, portid int, s *serial.Port, reply *string) int {
	resp := serialWriteAndEcho(portid, s, modEC20[cmdid].CmdStr)
	rs := []byte(resp)
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
	rs := []byte(resp)
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
	rs := []byte(resp)
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
	rs := []byte(resp)
	length := len(rs)
	sublen := len(modEC20[cmdid].CmdStr)

	vlog.Info("    AT cmd(%d): %q", length, []byte(resp))
	pos1 := strings.Index(resp, modEC20[cmdid].CmdStr+"\r\r\n")
	pos2 := strings.Index(resp, "\r\n\r\nOK")
	if pos1 >= 0 && pos2 >= 0 {
		posBgn := sublen + len("\r\r\n") + len("+CSIM: 28,\"")
		posEnd := pos2 - 1 - len("9000")
		preresp := string(rs[posBgn:posEnd])
		hexb := Ascii2Hex([]byte(preresp))

		vlog.Info("    AT get(%d): %s to %s", len(preresp), preresp, string(hexb[:]))
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
	rs := []byte(resp)
	length := len(rs)
	sublen := len(modEC20[cmdid].CmdStr)

	vlog.Info("    AT cmd(%d): %q", length, []byte(resp))
	pos1 := strings.Index(resp, modEC20[cmdid].CmdStr+"\r\r\n")
	pos2 := strings.Index(resp, "\r\n\r\nOK")
	if pos1 >= 0 && pos2 >= 0 {
		preresp := string(rs[(sublen + len("\r\r\n") + len("+CSIM: 36,\"")) : pos2-1-len("9000")])
		vlog.Info("    AT get(%d): %s", len(preresp), preresp)
		*reply = preresp
		return len(preresp)
	}

	return 0
}

/* "AT+CSIM=426,"A0D60000D002FD5646696C6531312E62696E000000...."\r\r\n\
 * +CSIM: 4,"9000"\r\n\r\n
 * OK\r\n"
 */
func ec20_set_ens192(cmdid int, portid int, s *serial.Port, reply *string) int {
	cmd_ens192 := fmt.Sprintf("AT+CSIM=426,\"%s%s\"", head_ens192, serial_port[portid].devInfo.sim_ens.EncData192)
	resp := serialWriteAndEcho(portid, s, cmd_ens192)
	rs := []byte(resp)
	length := len(rs)
	sublen := len(cmd_ens192)

	vlog.Info("    AT cmd(%d): %q", length, []byte(resp))
	pos1 := strings.Index(resp, cmd_ens192+"\r\r\n")
	pos2 := strings.Index(resp, "\r\n\r\nOK")
	if pos1 >= 0 && pos2 >= 0 {
		preresp := string(rs[(sublen + len("\r\r\n") + len("+CSIM: 4,\"")) : pos2-1])
		vlog.Info("    AT get(%d): %s", len(preresp), preresp)
		if preresp == "9000" {
			*reply = "OK"
			return len("OK")
		} else {
			*reply = "ERROR"
			return 0
		}
	}

	*reply = "ERROR"
	return 0
}

/* "AT+CSIM=170,\"A0D600015002FD5646696C6531312E62696E000000..."\r\r\n\
 * +CSIM: 4,"9000"\r\n\r\n
 * OK\r\n"
 */
func ec20_set_ens64(cmdid int, portid int, s *serial.Port, reply *string) int {
	cmd_ens64 := fmt.Sprintf("AT+CSIM=170,\"%s%s\"", head_ens64, serial_port[portid].devInfo.sim_ens.EncData64)
	resp := serialWriteAndEcho(portid, s, cmd_ens64)
	rs := []byte(resp)
	length := len(rs)
	sublen := len(cmd_ens64)

	vlog.Info("    AT cmd(%d): %q", length, []byte(resp))
	pos1 := strings.Index(resp, cmd_ens64+"\r\r\n")
	pos2 := strings.Index(resp, "\r\n\r\nOK")
	if pos1 >= 0 && pos2 >= 0 {
		preresp := string(rs[(sublen + len("\r\r\n") + len("+CSIM: 4,\"")) : pos2-1])
		vlog.Info("    AT get(%d): %s", len(preresp), preresp)
		if preresp == "9000" {
			*reply = "OK"
			return len("OK")
		} else {
			*reply = "ERROR"
			return 0
		}
	}

	*reply = "ERROR"
	return 0
}
