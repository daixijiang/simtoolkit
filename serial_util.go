/* VSIM Serial Product Toolkit
 * Author: daixijiang@gmail.com (2019)
 */

package main

import (
	"fmt"
	"runtime"
	"strings"
	"time"
	"vlog"

	"github.com/tarm/serial"
)

/* Port struct */
const (
	PORT_STATUS_CLOSE   = 0
	PORT_STATUS_OPEN    = 1
	PORT_STATUS_PRODUCE = 2
)

type serial_port_info struct {
	port_status int
	comPort     *serial.Port
	portname    string
	strInfo     string
	devInfo     device_info
}

type device_info struct {
	version string
	token   [OPER_MAX]string
	sim_src SRC_SIM_DATA
	sim_ens ENC_SIM_DATA
	servde  string
}

const SERIAL_PORT_MAX = 16
const APP_AT_OK string = "AT"

var COM_RNAME_PREFIX string
var COM_SNAME_PREFIX string
var serial_port [SERIAL_PORT_MAX]serial_port_info
var at_reply [SERIAL_PORT_MAX]string
var ports_list []string

/* Port end */

func portIsOK(portid int) int {
	if serial_port[portid].port_status != PORT_STATUS_CLOSE {
		return 1
	}

	return 0
}

func serialWriteAndEcho(portid int, s *serial.Port, strCmd string, millsecond int) string {
	buf := make([]byte, 512)
	at_reply[portid] = ""

	n, err := s.Read(buf)
	if err != nil {
		vlog.Error("Port[%d] => read failed: %s", portid, serial_port[portid].portname, err)
		return at_reply[portid]
	}

	n, err = s.Write([]byte(strCmd + "\r\n"))
	if err != nil {
		vlog.Error("Port[%d] => write %s failed: %s", portid, serial_port[portid].portname, err)
		return at_reply[portid]
	}

	if millsecond == 0 {
		millsecond = int(gConfig.Serial.Cmd_timewait * 1000)
	} else if millsecond > gConfig.Serial.Cmd_timeout*1000 {
		millsecond = gConfig.Serial.Cmd_timeout * 1000
	}
	time.Sleep(time.Duration(millsecond) * time.Millisecond)

	at_reply[portid] = ""
	for i := 0; i < 10; i++ {
		n, err = s.Read(buf)
		if n > 0 {
			at_reply[portid] += fmt.Sprintf("%s", string(buf[:n]))
		}

		if strings.LastIndex(at_reply[portid], "\r\nOK\r\n") > 0 {
			break
		}
		if strings.LastIndex(at_reply[portid], "\r\nERROR\r\n") > 0 {
			break
		}
	}

	vlog.Debug("[Req]%s", strCmd)
	vlog.Debug("[Rly]%s", at_reply[portid])
	return at_reply[portid]
}

func serialOpen(portid int, strCom string) int {
	vlog.Info("Port[%d] => open port %s", portid, strCom)

	if serial_port[portid].port_status != PORT_STATUS_CLOSE {
		serial_port[portid].strInfo = fmt.Sprintf("%s", "O")
		return 0
	}
	c := &serial.Config{Name: strCom, Baud: 115200, ReadTimeout: time.Duration(gConfig.Serial.Cmd_timeout * 1000)}
	s, err := serial.OpenPort(c)
	if err != nil {
		vlog.Error("Port[%d] => open %s failed: %s", portid, strCom, err)
		return -1
	}
	serial_port[portid].portname = strCom

	resp := serialWriteAndEcho(portid, s, APP_AT_OK, 100)
	vlog.Info("%s", resp)

	serial_port[portid].port_status = PORT_STATUS_OPEN
	serial_port[portid].comPort = s
	serial_port[portid].strInfo = fmt.Sprintf("%s", "O")
	serial_port[portid].devInfo = device_info{}
	return 0
}

func serialATsendCmd(portid int, strCom string, strCmd string) string {
	vlog.Info("Port[%d] => AT send cmd[%s] port %s", portid, strCmd, strCom)
	resp := serialWriteAndEcho(portid, serial_port[portid].comPort, strCmd, int(gConfig.Serial.Cmd_timewait*1000))
	vlog.Info("%s", resp)
	return resp
}

/* general get info */
/* eg: "AT+GSN\r\r\n
 * 862107043586551\r\n\r\n
 * OK\r\n"
 */
func serial_atget_info(cmdid int, cmdstr string, portid int, s *serial.Port, reply *string) int {
	resp := serialWriteAndEcho(portid, s, cmdstr, 0)
	rs := []byte(resp)
	length := len(rs)
	sublen := len(cmdstr)

	vlog.Info("    AT cmd(%d): %q", length, []byte(resp))
	pos1 := strings.Index(resp, cmdstr+"\r\r\n")
	pos2 := strings.Index(resp, "\r\n\r\nOK")
	if pos1 >= 0 && pos2 >= 0 {
		preresp := string(rs[sublen+len("\r\r\n") : pos2])
		vlog.Info("    AT get(%d): %s", len(preresp), preresp)
		*reply = preresp
		return len(*reply)
	}

	return 0
}

func serial_atget2_info(cmdid int, cmdstr string, portid int, s *serial.Port, reply *string) int {
	resp := serialWriteAndEcho(portid, s, cmdstr, gConfig.Serial.Cmd_timeout*1000)
	rs := []byte(resp)
	length := len(rs)
	sublen := len(cmdstr)

	vlog.Info("    AT cmd(%d): %q", length, []byte(resp))
	pos1 := strings.Index(resp, cmdstr+"\r\r\n")
	pos2 := strings.Index(resp, "\r\n\r\nOK")
	if pos1 >= 0 && pos2 >= 0 {
		preresp := string(rs[sublen+len("\r\r\n") : pos2])
		vlog.Info("    AT get(%d): %s", len(preresp), preresp)
		*reply = preresp
		return len(*reply)
	}

	return 0
}
func serialClose(portid int) int {
	vlog.Info("Port[%d] => close port %s", portid, serial_port[portid].portname)
	if serial_port[portid].port_status != PORT_STATUS_CLOSE {
		serial_port[portid].comPort.Close()
		serial_port[portid].port_status = PORT_STATUS_CLOSE
		serial_port[portid] = serial_port_info{}
	}
	serial_port[portid].strInfo = fmt.Sprintf("%s", "*")
	return 0
}

func serialProduce(portid int) int {
	return thisModule.GoProduce(portid)
}

func serialCheckDo(portid int) int {
	return thisModule.GoCheck(portid)
}

func serialList() []string {
	portlist := make([]string, 0)
	for id := 0; id < 128; id++ {
		strCom := fmt.Sprintf("%s%s%d", COM_RNAME_PREFIX, COM_SNAME_PREFIX, id)
		strSCom := fmt.Sprintf("%s%d", COM_SNAME_PREFIX, id)
		c := &serial.Config{Name: strCom, Baud: 115200, ReadTimeout: 100}
		s, err := serial.OpenPort(c)
		if err == nil {
			s.Close()
			portlist = append(portlist, strSCom)
		}
	}

	if len(portlist) == 0 {
		portlist = append(portlist, "null")
	}

	return portlist
}

func serial_util_init() {
	sysType := runtime.GOOS

	if sysType == "linux" {
		COM_SNAME_PREFIX = "USB"
		COM_RNAME_PREFIX = "/dev/tty"
	} else if sysType == "windows" {
		COM_SNAME_PREFIX = "com"
		COM_RNAME_PREFIX = ""
	} else {
		COM_SNAME_PREFIX = "USB"
		COM_RNAME_PREFIX = "/dev/tty"
	}

	ports_list = serialList()
	vlog.Info("Portlists: %v", ports_list)
}
