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

var COM_PREFIX string
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

/* Port end */

func portIsOK(portid int) int {
	if portid == SERAIL_PORT_MAX {
		for index := 0; index < SERAIL_PORT_MAX; index++ {
			if serial_port[index].port_status == PORT_STATUS_CLOSE {
				return 0
			}
		}
		return 1
	} else {
		if serial_port[portid].port_status == PORT_STATUS_CLOSE {
			return 0
		} else {
			return 1
		}
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
	}
	serial_port[portid].strInfo = fmt.Sprintf("%s", "*")
	return 0
}

func serialList() {
	//TODO, list ports
}

func serial_util_init() {
	sysType := runtime.GOOS

	if sysType == "linux" {
		COM_PREFIX = "/dev/ttyUSB"
	} else if sysType == "windows" {
		COM_PREFIX = "COM"
	} else {
		COM_PREFIX = "COM"
	}
}
