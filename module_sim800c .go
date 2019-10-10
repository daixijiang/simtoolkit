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

const head_ens128 string = "A0D600005001130303030303030303030303030303"

func module_sim800c_init() {
	modSIM800C[Module_CMD2_IMEI] = ModuleTable{
		Module_CMD2_IMEI,
		"AT+CGSN",
		sim800c_get_imei,
	}

	modSIM800C[Module_CMD2_VER] = ModuleTable{
		Module_CMD2_VER,
		"AT+CSIM=10,\"A0B009090C\"",
		sim800c_get_ver,
	}

	modSIM800C[Module_CMD2_CHIPID] = ModuleTable{
		Module_CMD2_CHIPID,
		"AT+CSIM=10,\"A0B0090910\"",
		sim800c_get_chipid,
	}

	modSIM800C[Module_CMD2_SIM64] = ModuleTable{
		Module_CMD2_SIM64,
		"AT+CSIM=170",
		sim800c_set_ens64,
	}

	modSIM800C[Module_CMD3_CCID] = ModuleTable{
		Module_CMD3_CCID,
		"AT+CCID",
		sim800c_get_ccid,
	}

	modSIM800C[Module_CMD3_CREG] = ModuleTable{
		Module_CMD3_CREG,
		"AT+CREG?",
		sim800c_get_creg,
	}
}

/* "AT+GSN\r\r\n
 * 862107043586551\r\n\r\n
 * OK\r\n"
 */
func sim800c_get_imei(cmdid int, portid int, s *serial.Port, reply *string) int {
	resp := serialWriteAndEcho(portid, s, modSIM800C[cmdid].CmdStr)
	rs := []byte(resp)
	length := len(rs)
	sublen := len(modSIM800C[cmdid].CmdStr)

	vlog.Info("    AT cmd(%d): %q", length, []byte(resp))
	pos1 := strings.Index(resp, modSIM800C[cmdid].CmdStr+"\r\r\n")
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
func sim800c_get_ccid(cmdid int, portid int, s *serial.Port, reply *string) int {
	resp := serialWriteAndEcho(portid, s, modSIM800C[cmdid].CmdStr)
	rs := []byte(resp)
	length := len(rs)
	sublen := len(modSIM800C[cmdid].CmdStr)

	vlog.Info("    AT cmd(%d): %q", length, []byte(resp))
	pos1 := strings.Index(resp, modSIM800C[cmdid].CmdStr+"\r\r\n")
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
func sim800c_get_creg(cmdid int, portid int, s *serial.Port, reply *string) int {
	resp := serialWriteAndEcho(portid, s, modSIM800C[cmdid].CmdStr)
	rs := []byte(resp)
	length := len(rs)
	sublen := len(modSIM800C[cmdid].CmdStr)

	vlog.Info("    AT cmd(%d): %q", length, []byte(resp))
	pos1 := strings.Index(resp, modSIM800C[cmdid].CmdStr+"\r\r\n")
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
func sim800c_get_ver(cmdid int, portid int, s *serial.Port, reply *string) int {
	resp := serialWriteAndEcho(portid, s, modSIM800C[cmdid].CmdStr)
	rs := []byte(resp)
	length := len(rs)
	sublen := len(modSIM800C[cmdid].CmdStr)

	vlog.Info("    AT cmd(%d): %q", length, []byte(resp))
	pos1 := strings.Index(resp, modSIM800C[cmdid].CmdStr+"\r\r\n")
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
func sim800c_get_chipid(cmdid int, portid int, s *serial.Port, reply *string) int {
	resp := serialWriteAndEcho(portid, s, modSIM800C[cmdid].CmdStr)
	rs := []byte(resp)
	length := len(rs)
	sublen := len(modSIM800C[cmdid].CmdStr)

	vlog.Info("    AT cmd(%d): %q", length, []byte(resp))
	pos1 := strings.Index(resp, modSIM800C[cmdid].CmdStr+"\r\r\n")
	pos2 := strings.Index(resp, "\r\n\r\nOK")
	if pos1 >= 0 && pos2 >= 0 {
		preresp := string(rs[(sublen + len("\r\r\n") + len("+CSIM: 36,\"")) : pos2-1-len("9000")])
		vlog.Info("    AT get(%d): %s", len(preresp), preresp)
		*reply = preresp
		return len(preresp)
	}

	return 0
}

/* "AT+CSIM=170,\"A0D600005001130303030303030303030303030303..."\r\r\n\
 * +CSIM: 4,"9000"\r\n\r\n
 * OK\r\n"
 */
func sim800c_set_ens64(cmdid int, portid int, s *serial.Port, reply *string) int {
	cmd_ens128 := fmt.Sprintf("AT+CSIM=170,\"%s%s\"", head_ens128, serial_port[portid].devInfo.sim_ens.EncData64)
	resp := serialWriteAndEcho(portid, s, cmd_ens128)
	rs := []byte(resp)
	length := len(rs)
	sublen := len(cmd_ens128)

	vlog.Info("    AT cmd(%d): %q", length, []byte(resp))
	pos1 := strings.Index(resp, cmd_ens128+"\r\r\n")
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
