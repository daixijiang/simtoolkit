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

const APP_AT_OK string = "AT"

/* Port struct */
const SERAIL_PORT_MAX = 8

var COM_RNAME_PREFIX string
var COM_SNAME_PREFIX string
var at_reply [SERAIL_PORT_MAX]string

const (
	PORT_STATUS_CLOSE   = 0
	PORT_STATUS_OPEN    = 1
	PORT_STATUS_PRODUCE = 2
)

type serial_port_info struct {
	port_status int
	portname    string
	strInfo     string
	comPort     *serial.Port
	dev_data    devReqData
	sim_pv1     devResPlainData
	sim_cv1     devResCipherData
	sim_ens     ENC_SIM_DATA
}

var serial_port [SERAIL_PORT_MAX]serial_port_info
var ports_list []string

/* Port end */

func portIsOK(portid int) int {
	if serial_port[portid].port_status != PORT_STATUS_CLOSE {
		return 1
	}

	return 0
}

func serialWriteAndEcho(portid int, s *serial.Port, strCmd string) string {
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

	time.Sleep(time.Duration(100) * time.Millisecond)

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

	//vlog.Info("[Req]%s", strCmd)
	//vlog.Info("[Rly]%s", at_reply[portid])
	return at_reply[portid]
}

func serialOpen(portid int, strCom string) int {
	vlog.Info("Port[%d] => open port %s", portid, strCom)

	if serial_port[portid].port_status != PORT_STATUS_CLOSE {
		return 0
	}
	c := &serial.Config{Name: strCom, Baud: 115200, ReadTimeout: 100}
	s, err := serial.OpenPort(c)
	if err != nil {
		vlog.Error("Port[%d] => open %s failed: %s", portid, strCom, err)
		return -1
	}
	serial_port[portid].portname = strCom

	resp := serialWriteAndEcho(portid, s, APP_AT_OK)
	vlog.Info("%s", resp)

	serial_port[portid].port_status = PORT_STATUS_OPEN
	serial_port[portid].comPort = s
	serial_port[portid].strInfo = fmt.Sprintf("%s", "o")
	return 0
}

func serialATsendCmd(portid int, strCom string, strCmd string) {
	vlog.Info("Port[%d] => AT send cmd[%s] port %s", portid, strCmd, strCom)
	resp := serialWriteAndEcho(portid, serial_port[portid].comPort, strCmd)
	vlog.Info("%s", resp)

	serial_port[portid].strInfo = fmt.Sprintf("%s", "T")
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
