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

func module_ec20_init() *[Module_TAB_AT_CMD_MAX]ModCmdTable {
	var myModCmd [Module_TAB_AT_CMD_MAX]ModCmdTable

	////cmd for prepare
	myModCmd[Module_CMD1_SYSVER] = ModCmdTable{
		Module_CMD1_SYSVER,
		"ATI",
		serial_atget_info,
	}

	myModCmd[Module_CMD1_SOFTMODE] = ModCmdTable{
		Module_CMD1_SOFTMODE,
		"AT+QCFG=\"sim/softsimmode\",2",
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

	myModCmd[Module_PRE1_SET_NETWORK0] = ModCmdTable{
		Module_PRE1_SET_NETWORK0,
		"AT+QMBNCFG=\"autosel\",0",
		serial_atget_info,
	}

	myModCmd[Module_PRE1_SET_NETWORK1] = ModCmdTable{
		Module_PRE1_SET_NETWORK1,
		"AT+QMBNCFG=\"select\",\"ROW_Generic_3GPP\"",
		serial_atget_info,
	}

	myModCmd[Module_CMD1_SET_SERVURL] = ModCmdTable{
		Module_CMD1_SET_SERVURL,
		"AT+CSIM=72,\"A0D600021F5CD9767557F5A4F56CD1A19ACD3DBB0C6573696D2E73686F776D61632E636E\"",
		serial_atget_info,
	}

	myModCmd[Module_CMD1_AUTOSWITCH_ON] = ModCmdTable{
		Module_CMD1_AUTOSWITCH_ON,
		"AT+CSIM=42,\"A0D6000410C7FD5646696C6535362E62696EDA96A3\"",
		serial_atget_info,
	}

	myModCmd[Module_CMD1_AUTOSWITCH_OFF] = ModCmdTable{
		Module_CMD1_AUTOSWITCH_OFF,
		"AT+CSIM=42,\"A0D6000310C7FD5646696C6535362E62696EDA96A3\"",
		serial_atget_info,
	}

	myModCmd[Module_CMD1_BACKUP_CONFIG] = ModCmdTable{
		Module_CMD1_BACKUP_CONFIG,
		"AT+QPRTPARA=1",
		serial_atget2_info,
	}

	////cmd for produce
	myModCmd[Module_CMD2_IMEI] = ModCmdTable{
		Module_CMD2_IMEI,
		"AT+GSN",
		serial_atget_info,
	}

	myModCmd[Module_CMD2_COSVER] = ModCmdTable{
		Module_CMD2_COSVER,
		"AT+CSIM=10,\"A0B009090C\"",
		ec20_get_ver,
	}

	myModCmd[Module_CMD2_CHIPID] = ModCmdTable{
		Module_CMD2_CHIPID,
		"AT+CSIM=10,\"A0B0090910\"",
		ec20_get_chipid,
	}

	myModCmd[Module_CMD2_SIM192] = ModCmdTable{
		Module_CMD2_SIM192,
		"AT+CSIM=426",
		ec20_set_ens192,
	}

	myModCmd[Module_CMD2_SIM64] = ModCmdTable{
		Module_CMD2_SIM64,
		"AT+CSIM=170",
		ec20_set_ens64,
	}

	////cmd for check
	myModCmd[Module_CMD3_CCID] = ModCmdTable{
		Module_CMD3_CCID,
		"AT+CCID",
		ec20_get_ccid,
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

	myModCmd[Module_CMD3_SWITCH_TEL] = ModCmdTable{
		Module_CMD3_SWITCH_TEL,
		"AT+CSIM=12,\"00B001090103\"",
		serial_atget_info,
	}

	myModCmd[Module_CMD3_SWITCH_CU] = ModCmdTable{
		Module_CMD3_SWITCH_CU,
		"AT+CSIM=12,\"00B001090102\"",
		serial_atget_info,
	}

	myModCmd[Module_CMD3_SWITCH_CM] = ModCmdTable{
		Module_CMD3_SWITCH_CM,
		"AT+CSIM=12,\"00B001090101\"",
		serial_atget_info,
	}

	return &myModCmd
}

/* "AT+CCID\r\r\n
 * +CCID: 898602B2211790026229\r\n\r\n
 * OK\r\n"
 */
func ec20_get_ccid(cmdid int, cmdstr string, portid int, s *serial.Port, reply *string) int {
	resp := serialWriteAndEcho(portid, s, cmdstr, 0)
	rs := []byte(resp)
	length := len(rs)
	sublen := len(cmdstr)

	vlog.Info("    AT cmd(%d): %q", length, []byte(resp))
	pos1 := strings.Index(resp, cmdstr+"\r\r\n")
	pos2 := strings.Index(resp, "\r\n\r\nOK")
	if pos1 >= 0 && pos2 >= 0 {
		posBgn := sublen + len("\r\r\n") + len("+CCID: ")
		posEnd := pos2 - 1
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
func ec20_get_chipid(cmdid int, cmdstr string, portid int, s *serial.Port, reply *string) int {
	resp := serialWriteAndEcho(portid, s, cmdstr, 0)
	rs := []byte(resp)
	length := len(rs)
	sublen := len(cmdstr)

	vlog.Info("    AT cmd(%d): %q", length, []byte(resp))
	pos1 := strings.Index(resp, cmdstr+"\r\r\n")
	pos2 := strings.Index(resp, "\r\n\r\nOK")
	if pos1 >= 0 && pos2 >= 0 {
		posBgn := sublen + len("\r\r\n") + len("+CSIM: 36,\"")
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
func ec20_get_ver(cmdid int, cmdstr string, portid int, s *serial.Port, reply *string) int {
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

/* "AT+CSIM=426,"A0D60000D002FD5646696C6531312E62696E000000...."\r\r\n\
 * +CSIM: 4,"9000"\r\n\r\n
 * OK\r\n"
 */
func ec20_set_ens192(cmdid int, cmdstr string, portid int, s *serial.Port, reply *string) int {
	cmd_ens192 := fmt.Sprintf("AT+CSIM=426,\"%s%s\"", head_ens192, serial_port[portid].devInfo.sim_ens.EncData192)
	resp := serialWriteAndEcho(portid, s, cmd_ens192, 0)
	rs := []byte(resp)
	length := len(rs)
	sublen := len(cmd_ens192)

	vlog.Info("    AT cmd(%d): %q", length, []byte(resp))
	pos1 := strings.Index(resp, cmd_ens192+"\r\r\n")
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

/* "AT+CSIM=170,\"A0D600015002FD5646696C6531312E62696E000000..."\r\r\n\
 * +CSIM: 4,"9000"\r\n\r\n
 * OK\r\n"
 */
func ec20_set_ens64(cmdid int, cmdstr string, portid int, s *serial.Port, reply *string) int {
	cmd_ens64 := fmt.Sprintf("AT+CSIM=170,\"%s%s\"", head_ens64, serial_port[portid].devInfo.sim_ens.EncData64)
	resp := serialWriteAndEcho(portid, s, cmd_ens64, 0)
	rs := []byte(resp)
	length := len(rs)
	sublen := len(cmd_ens64)

	vlog.Info("    AT cmd(%d): %q", length, []byte(resp))
	pos1 := strings.Index(resp, cmd_ens64+"\r\r\n")
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
