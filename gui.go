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

type TBaseInfo struct {
	Module   Module_cfg
	Message  string
	PortList []string
	Checkbox [SERIAL_PORT_MAX]bool

	nui *nucularUI
	wui *walkUI
}

var sysNeedSync = false

//////////////////////////////////////////////////////////////////////////////////
// show handle

func newTBaseInfo() (tbi *TBaseInfo) {
	newbi := &TBaseInfo{}
	newbi.Module = module_get()
	newbi.Message = ""
	newbi.PortList = ports_list

	myBtnTab[Btn_CMD_Produce] = BtnDoTable{Btn_CMD_Produce, "produce", serialProduce}
	myBtnTab[Btn_CMD_CheckDo] = BtnDoTable{Btn_CMD_CheckDo, "checkdo", serialCheckDo}
	myBtnTab[Btn_CMD_Close] = BtnDoTable{Btn_CMD_Close, "close", serialClose}

	return newbi
}

//////////////////////////////////////////////////////////////////////////////////
// btn handle
func (tbi *TBaseInfo) btnOpen(portid int, strCom string) {
	if serialOpen(portid, strCom) != 0 {
		msg := fmt.Sprintf("Filed to open the %s!", strCom)
		tbi.openMessage(msg)
	}
}

func (tbi *TBaseInfo) btnClose(portid int, strCom string) {
	serialClose(portid)
}

func (tbi *TBaseInfo) btnATSend(portid int, strCom string, strCmd string) string {
	if portIsOK(portid) == 0 {
		msg := fmt.Sprintf("Please open the port[%d]!", portid)
		tbi.openMessage(msg)
	} else {
		return serialATsendCmd(portid, strCom, strCmd)
	}

	return ""
}

func (tbi *TBaseInfo) btnProduce(portid int, strCom string) {
	if portIsOK(portid) == 0 {
		msg := fmt.Sprintf("Please open the port[%d]!", portid)
		tbi.openMessage(msg)
	} else {
		vlog.Info("Port[%d] => start produce %s", portid, strCom)

		if serial_port[portid].port_status != PORT_STATUS_PRODUCE {
			serial_port[portid].port_status = PORT_STATUS_PRODUCE
			taskChan := make(chan PortResult)
			go tbi.setTaskBtn(Btn_CMD_Produce, portid, taskChan)
			tbi.getTaskBtn(1, taskChan)
		}
		serial_port[portid].port_status = PORT_STATUS_OPEN
	}
}

func (tbi *TBaseInfo) btnHandleAll(oper int, check bool) {
	if check && (tbi.checkBoxAll() == false) {
		return
	}

	vlog.Info("start %s all", myBtnTab[oper].BtnStr)
	taskChan := make(chan PortResult)
	taskCnt := 0
	for port_id := 0; port_id < gConfig.Serial.Serial_max; port_id++ {
		if (tbi.Checkbox[port_id] || !check) && (portIsOK(port_id) != 0) {
			taskCnt++
			go tbi.setTaskBtn(oper, port_id, taskChan)
		}
	}
	tbi.getTaskBtn(taskCnt, taskChan)
}

func (tbi *TBaseInfo) btnExit() {
	os.Exit(1)
}

func (tbi *TBaseInfo) btnLoadToken() {
	loadTokenCfg(gConfig.Token.Cmcc_file, OPER_CN_MOBILE)
	loadTokenCfg(gConfig.Token.Uni_file, OPER_CN_UNICOM)
	loadTokenCfg(gConfig.Token.Tel_file, OPER_CN_TELECOM)
}

func (tbi *TBaseInfo) btnRefreshPort() {
	tbi.btnHandleAll(Btn_CMD_Close, false)

	tbi.PortList = serialList()

	for port_id := 0; port_id < gConfig.Serial.Serial_max; port_id++ {
		//TODO reset combo
		//pg.portCombo[port_id].resetItems()
	}

	vlog.Info("Portlists: %v", tbi.PortList)
}

//////////////////////////////////////////////////////////////////////////////////
// other handle

func (tbi *TBaseInfo) setTaskBtn(oper int, portid int, taskCH chan PortResult) {
	ret := myBtnTab[oper].BtnFunc(portid)
	resp := PortResult{
		Portid: portid,
		Oper:   oper,
		Result: ret,
	}
	vlog.Info("Handle-Put result of %s port[%d]: %d", myBtnTab[resp.Oper].BtnStr, resp.Portid, resp.Result)

	tbi.doResult(resp)

	if sysNeedSync {
		taskCH <- resp
	}
}

func (tbi *TBaseInfo) getTaskBtn(count int, taskCH chan PortResult) {
	if sysNeedSync {
		for i := 0; i < count; i++ {
			resp := <-taskCH
			vlog.Info("Handle-Get result of %s port[%d]: %d", myBtnTab[resp.Oper].BtnStr, resp.Portid, resp.Result)
			tbi.doResult(resp)
		}
		close(taskCH)
	}

	return
}

func (tbi *TBaseInfo) checkBoxAll() bool {
	cntCheck := 0
	portlist := ""

	for port_id := 0; port_id < gConfig.Serial.Serial_max; port_id++ {
		if tbi.Checkbox[port_id] {
			cntCheck++
			if portIsOK(port_id) == 0 {
				portlist += fmt.Sprintf("%d,", port_id)
			}
		}
	}

	vlog.Info("Check port(%d): %s", cntCheck, portlist)

	if cntCheck == 0 {
		tbi.openMessage("Please select(open) the ports!")
		return false
	} else if portlist != "" {
		msg := fmt.Sprintf("Please select(open) the ports: [%s]!", portlist)
		tbi.openMessage(msg)
		return false
	}

	return true
}

func (tbi *TBaseInfo) setModule(module Module_cfg) {
	if tbi.Module != module {
		vlog.Info("module from %d to %d", tbi.Module, module)
		tbi.btnHandleAll(Btn_CMD_Close, false)
		tbi.Module = module
		module_reinit(tbi.Module)
	}
}

//////////////////////////////////////////////////////////////////////////////////

func (tbi *TBaseInfo) openMessage(message string) {
	if tbi.nui != nil {
		tbi.nui.openMessage(message)
	}
	if tbi.wui != nil {
		tbi.wui.openMessage(message)
	}
}

func (tbi *TBaseInfo) doResult(resp PortResult) {
	if tbi.nui != nil {
		tbi.nui.doResult(resp)
	}
	if tbi.wui != nil {
		tbi.wui.doResult(resp)
	}
}
