/* VSIM Serial Product Toolkit
 * Author: daixijiang@gmail.com (2019)
 */

package main

import (
	"encoding/json"
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
	sim_pv1     devResPlainData
	sim_cv1     devResCipherData
	sim_ens     ENC_SIM_DATA
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

func serialProduce(portid int) int {
	ret := getDevInfo(portid)
	if ret != 0 {
		vlog.Info("!!!!Port[%d] => failed(%d) to getDevInfo ...", portid, ret)
		return ret
	}

	ret = getProToken(portid)
	if ret != 0 {
		vlog.Info("!!!!Port[%d] => failed(%d) to getProToken ...", portid, ret)
		return ret
	}

	ret = getServInfo(portid)
	if ret != 0 {
		vlog.Info("!!!!Port[%d] => failed(%d) to getServInfo ...", portid, ret)
		return ret
	}

	ret = cryptoVsim(portid)
	if ret != 0 {
		vlog.Info("!!!!Port[%d] => failed(%d) to cryptoVsim ...", portid, ret)
		return ret
	}

	ret = setVsimData(portid)
	if ret != 0 {
		vlog.Info("!!!!Port[%d] => failed(%d) to setVsimData ...", portid, ret)
		return ret
	}

	serial_port[portid].port_status = PORT_STATUS_OPEN
	return 0
}

func serialCheckDo(portid int) int {
	var result string
	vlog.Info("Port[%d] => check device produce ...", portid)

	//check1 ccid "AT+CCID"
	modEC20[Module_CMD3_CCID].CmdFunc(
		Module_CMD3_CCID, portid,
		serial_port[portid].comPort,
		&result)

	if result != serial_port[portid].sim_pv1.Iccid {
		vlog.Info("!!!!check ccid err: %s %s", result, serial_port[portid].sim_pv1.Iccid)
	} else {
		vlog.Info("    check ccid ok: %s", result)
	}

	//TODO, check2 creg "AT+CREG?"
	modEC20[Module_CMD3_CREG].CmdFunc(
		Module_CMD3_CREG, portid,
		serial_port[portid].comPort,
		&result)
	vlog.Info("    get creg: %s", result)
	return 0
}

func getDevInfo(portid int) int {
	vlog.Info("Port[%d] p(1.0)=> get device info ...", portid)
	time.Sleep(time.Duration(2) * time.Second)

	//get imei
	modEC20[Module_CMD2_IMEI].CmdFunc(
		Module_CMD2_IMEI, portid,
		serial_port[portid].comPort,
		&serial_port[portid].dev_data.Imei)

	//TODO, get chipid (del 9000)
	modEC20[Module_CMD2_CHIPID].CmdFunc(
		Module_CMD2_CHIPID, portid,
		serial_port[portid].comPort,
		&serial_port[portid].dev_data.Chipid)

	//TODO, get version (ascii to hex)
	modEC20[Module_CMD2_VER].CmdFunc(
		Module_CMD2_VER, portid,
		serial_port[portid].comPort,
		&serial_port[portid].dev_data.Ver)

	return 0
}

func getProToken(portid int) int {
	vlog.Info("Port[%d] p(2.0)=> get token info ...", portid)
	serial_port[portid].dev_data.Token = getToken(serial_port[portid].dev_data.Imei)
	return 0
}

func getServInfo(portid int) int {
	vlog.Info("Port[%d] p(3.0)=> get server info ...", portid)
	time.Sleep(time.Duration(2) * time.Second)

	var dev_data devReqPlainData
	var res []byte

	dev_data.Imei = serial_port[portid].dev_data.Imei
	dev_data.Token = serial_port[portid].dev_data.Token

	//TODO, test
	dev_data.Imei = "863412049788253"
	dev_data.Token = "YR0NI-259CE-R3JI5-01DJN-ENY2Z"

	reqSimServer(SERVER_PLAIN_v0, dev_data, &res)
	err := json.Unmarshal(res, &serial_port[portid].sim_pv1)
	if checkerr(err, 3, "Json parse server data") != 0 {
		return 3
	}

	vlog.Info("    Get siminfo: %+v", serial_port[portid].sim_pv1)
	serial_port[portid].strInfo = fmt.Sprintf("%s", "S")

	return 0
}

func cryptoVsim(portid int) int {
	vlog.Info("Port[%d] p(4.0)=> crypto vsim data ...", portid)

	srcsim := SRC_SIM_DATA{
		Imei:   serial_port[portid].dev_data.Imei,
		ChipID: serial_port[portid].dev_data.Chipid,
		CdmaData: CDMA_DATA{
			Imsi_m: serial_port[portid].sim_pv1.ImsiM,
			Uim_id: serial_port[portid].sim_pv1.Uimid,
			Hrdupp: serial_port[portid].sim_pv1.Hrpdupp,
		},
	}

	srcsim.VsimData[OPER_CN_MOBILE] = SIM_DATA{
		Iccid: serial_port[portid].sim_pv1.Iccid,
		Imsi:  serial_port[portid].sim_pv1.ImsiLTE,
		Ki:    serial_port[portid].sim_pv1.Ki,
		Opc:   serial_port[portid].sim_pv1.Opc,
	}

	Lib_vsim_encrypt(srcsim, &serial_port[portid].sim_ens)
	vlog.Info("    crypto de192: %s", serial_port[portid].sim_ens.EncData192)
	vlog.Info("    crypto de64: %s", serial_port[portid].sim_ens.EncData64)

	return 0
}

func setVsimData(portid int) int {
	var result string
	vlog.Info("Port[%d] p(5.0)=> do producing on ...", portid)

	//TODO, get result
	modEC20[Module_CMD2_SIM192].CmdFunc(
		Module_CMD2_SIM192, portid,
		serial_port[portid].comPort,
		&result)
	vlog.Info("    set de192 %s", result)

	//TODO, get result
	modEC20[Module_CMD2_SIM64].CmdFunc(
		Module_CMD2_SIM64, portid,
		serial_port[portid].comPort,
		&result)
	vlog.Info("    set de64 %s", result)

	vlog.Info("Port[%d] p(5.1)=> do producing ok!", portid)
	return 0
}

func serialList() {
}
