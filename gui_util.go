/* VSIM Serial Product Toolkit
 * Author: daixijiang@gmail.com (2019)
 */

package main

import (
	"fmt"
	"os"
	"vlog"
)

const (
	Btn_CMD_Produce = 0
	Btn_CMD_CheckDo = 1
	Btn_CMD_Close   = 2
	Btn_CMD_MAX     = 3
)

type btnHandler func(portid int) int

type BtnDoTable struct {
	BtnID   int
	BtnStr  string
	BtnFunc btnHandler
}

var myBtnTab [Btn_CMD_MAX]BtnDoTable

type PortResult struct {
	Portid int
	Oper   int
	Result int
}

var sysNeedSync = false

//////////////////////////////////////////////////////////////////////////////////
// show handle

func newportGroup() (pg *portGroup) {
	pg = &portGroup{}

	pg.Message = ""
	pg.Module = module_get()
	pg.PortList = ports_list

	myBtnTab[Btn_CMD_Produce] = BtnDoTable{Btn_CMD_Produce, "produce", serialProduce}
	myBtnTab[Btn_CMD_CheckDo] = BtnDoTable{Btn_CMD_CheckDo, "checkdo", serialCheckDo}
	myBtnTab[Btn_CMD_Close] = BtnDoTable{Btn_CMD_Close, "close", serialClose}

	return pg
}

//////////////////////////////////////////////////////////////////////////////////
// btn handle
func (pg *portGroup) btnOpen(portid int, strCom string) {
	if serialOpen(portid, strCom) != 0 {
		msg := fmt.Sprintf("Filed to open the %s!", strCom)
		pg.openMessage(msg)
	}
}

func (pg *portGroup) btnClose(portid int, strCom string) {
	serialClose(portid)
}

func (pg *portGroup) btnATSend(portid int, strCom string, strCmd string) string {
	if portIsOK(portid) == 0 {
		msg := fmt.Sprintf("Please open the port[%d]!", portid)
		pg.openMessage(msg)
	} else {
		return serialATsendCmd(portid, strCom, strCmd)
	}

	return ""
}

func (pg *portGroup) btnProduce(portid int, strCom string) {
	if portIsOK(portid) == 0 {
		msg := fmt.Sprintf("Please open the port[%d]!", portid)
		pg.openMessage(msg)
	} else {
		vlog.Info("Port[%d] => start produce %s", portid, strCom)

		if serial_port[portid].port_status != PORT_STATUS_PRODUCE {
			serial_port[portid].port_status = PORT_STATUS_PRODUCE
			taskChan := make(chan PortResult)
			go pg.setTaskBtn(Btn_CMD_Produce, portid, taskChan)
			pg.getTaskBtn(1, taskChan)
		}
		serial_port[portid].port_status = PORT_STATUS_OPEN
	}
}

func (pg *portGroup) btnHandleAll(oper int, check bool) {
	if check && (pg.checkBox() == false) {
		return
	}

	vlog.Info("start %s all", myBtnTab[oper].BtnStr)
	taskChan := make(chan PortResult)
	taskCnt := 0
	for port_id := 0; port_id < gConfig.Serial.Serial_max; port_id++ {
		if (pg.Checkbox[port_id] || !check) && (portIsOK(port_id) != 0) {
			taskCnt++
			go pg.setTaskBtn(oper, port_id, taskChan)
		}
	}
	pg.getTaskBtn(taskCnt, taskChan)
}

func (pg *portGroup) btnExit() {
	os.Exit(1)
}

func (pg *portGroup) btnLoadToken() {
	loadTokenCfg(gConfig.Token.Cmcc_file, OPER_CN_MOBILE)
	loadTokenCfg(gConfig.Token.Uni_file, OPER_CN_UNICOM)
	loadTokenCfg(gConfig.Token.Tel_file, OPER_CN_TELECOM)
}

func (pg *portGroup) btnRefreshPort() {
	pg.btnHandleAll(Btn_CMD_Close, false)

	pg.PortList = serialList()

	for port_id := 0; port_id < gConfig.Serial.Serial_max; port_id++ {
		//TODO reset combo
		//pg.portCombo[port_id].resetItems()
	}

	vlog.Info("Portlists: %v", pg.PortList)
}

//////////////////////////////////////////////////////////////////////////////////
// other handle

func (pg *portGroup) setTaskBtn(oper int, portid int, taskCH chan PortResult) {
	ret := myBtnTab[oper].BtnFunc(portid)
	resp := PortResult{
		Portid: portid,
		Oper:   oper,
		Result: ret,
	}
	vlog.Info("Handle-Put result of %s port[%d]: %d", myBtnTab[resp.Oper].BtnStr, resp.Portid, resp.Result)

	pg.doResult(resp)

	if sysNeedSync {
		taskCH <- resp
	}
}

func (pg *portGroup) getTaskBtn(count int, taskCH chan PortResult) {
	if sysNeedSync {
		for i := 0; i < count; i++ {
			resp := <-taskCH
			vlog.Info("Handle-Get result of %s port[%d]: %d", myBtnTab[resp.Oper].BtnStr, resp.Portid, resp.Result)
			pg.doResult(resp)
		}
		close(taskCH)
	}

	return
}

func (pg *portGroup) checkBox() bool {
	cntCheck := 0
	portlist := ""

	for port_id := 0; port_id < gConfig.Serial.Serial_max; port_id++ {
		if pg.Checkbox[port_id] {
			cntCheck++
			if portIsOK(port_id) == 0 {
				portlist += fmt.Sprintf("%d,", port_id)
			}
		}
	}

	vlog.Info("Check port(%d): %s", cntCheck, portlist)

	if cntCheck == 0 {
		pg.openMessage("Please select(open) the ports!")
		return false
	} else if portlist != "" {
		msg := fmt.Sprintf("Please select(open) the ports: [%s]!", portlist)
		pg.openMessage(msg)
		return false
	}

	return true
}

func (pg *portGroup) setModule(module Module_cfg) {
	if pg.Module != module {
		vlog.Info("module from %d to %d", pg.Module, module)
		pg.btnHandleAll(Btn_CMD_Close, false)
		pg.Module = module
		module_reinit(pg.Module)
	}
}
