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

func module_sim800c_init(myModCmd *[Module_TAB_AT_CMD_MAX]ModCmdTable) {
	////cmd for prepare
	myModCmd[Module_CMD1_SYSVER] = ModCmdTable{
		Module_CMD1_SYSVER,
		"ATI",
		serial_atget_info,
	}

	myModCmd[Module_CMD1_SOFTMODE] = ModCmdTable{
		Module_CMD1_SOFTMODE,
		"AT+SSIM=1",
		serial_atget_info,
	}

	myModCmd[Module_CMD1_RESET0] = ModCmdTable{
		Module_CMD1_RESET0,
		"AT+CFUN=0",
		serial_atget_info,
	}

	myModCmd[Module_CMD1_RESET1] = ModCmdTable{
		Module_CMD1_RESET1,
		"AT+CFUN=1",
		serial_atget_info,
	}

	myModCmd[Module_CMD1_RESET2] = ModCmdTable{
		Module_CMD1_RESET2,
		"AT+CFUN=1,1",
		serial_atget_info,
	}

	////cmd for produce
	myModCmd[Module_CMD2_IMEI] = ModCmdTable{
		Module_CMD2_IMEI,
		"AT+CGSN",
		serial_atget_info,
	}

	myModCmd[Module_CMD2_COSVER] = ModCmdTable{
		Module_CMD2_COSVER,
		"AT+CSIM=10,\"A0B009090C\"",
		sim800c_get_ver,
	}

	myModCmd[Module_CMD2_CHIPID] = ModCmdTable{
		Module_CMD2_CHIPID,
		"AT+CSIM=10,\"A0B0090910\"",
		sim800c_get_chipid,
	}

	myModCmd[Module_CMD2_SIM64] = ModCmdTable{
		Module_CMD2_SIM64,
		"AT+CSIM=170",
		sim800c_set_ens64,
	}

	////cmd for check
	myModCmd[Module_CMD3_CCID] = ModCmdTable{
		Module_CMD3_CCID,
		"AT+CCID",
		sim800c_get_ccid,
	}

	myModCmd[Module_CMD3_CIMI] = ModCmdTable{
		Module_CMD3_CIMI,
		"AT+CIMI",
		serial_atget_info,
	}

	myModCmd[Module_CMD3_CREG] = ModCmdTable{
		Module_CMD3_CREG,
		"AT+CREG?",
		serial_atget_info,
	}

	myModCmd[Module_CMD3_CEREG] = ModCmdTable{
		Module_CMD3_CEREG,
		"AT+CEREG?",
		serial_atget_info,
	}

	myModCmd[Module_CMD3_COPS] = ModCmdTable{
		Module_CMD3_COPS,
		"AT+COPS?",
		serial_atget_info,
	}
}

/* "AT+CCID\r\r\n
 * +CCID: 898602B2211790026229\r\n\r\n
 * OK\r\n"
 */
func sim800c_get_ccid(cmdid int, cmdstr string, portid int, s *serial.Port, reply *string) int {
	resp := serialWriteAndEcho(portid, s, cmdstr, 0)
	rs := []byte(resp)
	length := len(rs)
	sublen := len(cmdstr)

	vlog.Info("    AT cmd(%d): %q", length, []byte(resp))
	pos1 := strings.Index(resp, cmdstr+"\r\r\n")
	pos2 := strings.Index(resp, "\r\n\r\nOK")
	if pos1 >= 0 && pos2 >= 0 {
		posBgn := sublen + len("\r\r\n") + len("+CCID: ")
		posEnd := pos2
		if posEnd < posBgn {
			posEnd = posBgn
		}
		preresp := string(rs[posBgn:posEnd])

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
func sim800c_get_chipid(cmdid int, cmdstr string, portid int, s *serial.Port, reply *string) int {
	resp := serialWriteAndEcho(portid, s, cmdstr, 0)
	rs := []byte(resp)
	length := len(rs)
	sublen := len(cmdstr)

	vlog.Info("    AT cmd(%d): %q", length, []byte(resp))
	pos1 := strings.Index(resp, cmdstr+"\r\r\n")
	pos2 := strings.Index(resp, "\r\n\r\nOK")
	if pos1 >= 0 && pos2 >= 0 {
		posBgn := sublen + len("\r\r\n") + len("+CSIM: 28,\"")
		posEnd := pos2 - 1 - len("9000")
		if posEnd < posBgn {
			posEnd = posBgn
		}
		preresp := string(rs[posBgn:posEnd])

		vlog.Info("    AT get(%d): %s", len(preresp), preresp)
		*reply = preresp
		return len(preresp)
	}

	return 0
}

/* "AT+CSIM=10,"A0B009090C"\r\r\n
 * +CSIM: 28,"436F735665725F312E312E349000"\r\n\r\n
 * OK\r\n"
 */
func sim800c_get_ver(cmdid int, cmdstr string, portid int, s *serial.Port, reply *string) int {
	resp := serialWriteAndEcho(portid, s, cmdstr, 0)
	rs := []byte(resp)
	length := len(rs)
	sublen := len(cmdstr)

	vlog.Info("    AT cmd(%d): %q", length, []byte(resp))
	pos1 := strings.Index(resp, cmdstr+"\r\r\n")
	pos2 := strings.Index(resp, "\r\n\r\nOK")
	if pos1 >= 0 && pos2 >= 0 {
		posBgn := sublen + len("\r\r\n") + len("+CSIM: 28,\"")
		posEnd := pos2 - 1 - len("9000")
		if posEnd < posBgn {
			posEnd = posBgn
		}
		preresp := string(rs[posBgn:posEnd])
		hexb := Ascii2Hex([]byte(preresp))

		vlog.Info("    AT get(%d): %s to %s", len(preresp), preresp, string(hexb[:]))
		*reply = preresp
		return len(preresp)
	}

	return 0
}

/* "AT+CSIM=170,\"A0D600005001130303030303030303030303030303..."\r\r\n\
 * +CSIM: 4,"9000"\r\n\r\n
 * OK\r\n"
 */
func sim800c_set_ens64(cmdid int, cmdstr string, portid int, s *serial.Port, reply *string) int {
	cmd_ens128 := fmt.Sprintf("AT+CSIM=170,\"%s%s\"", head_ens128, serial_port[portid].devInfo.sim_ens.EncData64)
	resp := serialWriteAndEcho(portid, s, cmd_ens128, 0)
	rs := []byte(resp)
	length := len(rs)
	sublen := len(cmd_ens128)

	vlog.Info("    AT cmd(%d): %q", length, []byte(resp))
	pos1 := strings.Index(resp, cmd_ens128+"\r\r\n")
	pos2 := strings.Index(resp, "\r\n\r\nOK")
	if pos1 >= 0 && pos2 >= 0 {
		posBgn := sublen + len("\r\r\n") + len("+CSIM: 4,\"")
		posEnd := pos2 - 1
		if posEnd < posBgn {
			posEnd = posBgn
		}
		preresp := string(rs[posBgn:posEnd])
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
