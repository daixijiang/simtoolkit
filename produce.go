/* VSIM Serial Product Toolkit
 * Author: daixijiang@gmail.com (2019)
 */

package main

import (
	"fmt"
	"strings"
	"time"
	"vlog"

	"github.com/tarm/serial"
)

const COM_PREFIX string = "COM"
const APP_AT_OK string = "AT"

/* Port struct */
const SERAIL_PORT_MAX = 8

var close_port int = 0xFF
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
	sim_data    devResData
}

var serial_port [SERAIL_PORT_MAX]serial_port_info

/* Port end */

// serail function
func portIsOK(portid int) int {
	if portid == SERAIL_PORT_MAX {
		for index := 0; index < SERAIL_PORT_MAX; index++ {
			if serial_port[index].port_status == PORT_STATUS_CLOSE {
				close_port = index
				return 0
			}
		}
		return 1
	} else {
		if serial_port[portid].port_status == PORT_STATUS_CLOSE {
			close_port = portid
			return 0
		} else {
			return 1
		}
	}

	return 0
}

func serialWriteAndEcho(portid int, s *serial.Port, strCmd string) string {
	buf := make([]byte, 128)
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

func serialClose(portid int, strCom string) {
	vlog.Info("Port[%d] => close port %s", portid, strCom)
	if serial_port[portid].port_status != PORT_STATUS_CLOSE {
		serial_port[portid].comPort.Close()
		serial_port[portid].port_status = PORT_STATUS_CLOSE
	}
	serial_port[portid].strInfo = fmt.Sprintf("%s", "*")
}

func serialProduce(portid int) {
	getDevInfo(portid)
	getServInfo(portid)
	cryptoVSIM(portid)
	writeVsimData(portid)
	serial_port[portid].port_status = PORT_STATUS_OPEN
}

func getDevInfo(portid int) {
	vlog.Info("Port[%d] p(1.0)=> get device info ...", portid)
	time.Sleep(time.Duration(2) * time.Second)

	//serial_port[portid].dev_data.Imei = "868575021892064"
	modEC20[Module_CMD2_IMEI].CmdFunc(
		Module_CMD2_IMEI, portid,
		serial_port[portid].comPort,
		&serial_port[portid].dev_data.Imei)

	//serial_port[portid].dev_data.Chipid = "20171026050559A399032A3416886391"
	modEC20[Module_CMD2_CHIPID].CmdFunc(
		Module_CMD2_CHIPID, portid,
		serial_port[portid].comPort,
		&serial_port[portid].dev_data.Chipid)

	var version string
	modEC20[Module_CMD2_VER].CmdFunc(
		Module_CMD2_VER, portid, serial_port[portid].comPort, &version)

	//TODO
	serial_port[portid].dev_data.Token = getToken(serial_port[portid].dev_data.Imei)
}

func getServInfo(portid int) {
	vlog.Info("Port[%d] p(2.0)=> get server info ...", portid)
	time.Sleep(time.Duration(2) * time.Second)

	reqSimServer(serial_port[portid].dev_data, serial_port[portid].sim_data)
	serial_port[portid].strInfo = fmt.Sprintf("%s", "S")
}

func cryptoVSIM(portid int) {
	vlog.Info("Port[%d] p(3.0)=> crypto vsim data ...", portid)
	time.Sleep(time.Duration(2) * time.Second)
}

func writeVsimData(portid int) {
	vlog.Info("Port[%d] p(4.0)=> do producing on ...", portid)
	time.Sleep(time.Duration(2) * time.Second)
}

func serialList() {
}
