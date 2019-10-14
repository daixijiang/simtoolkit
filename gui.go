/* VSIM Serial Product Toolkit
 * Author: daixijiang@gmail.com (2019)
 */

package main

import (
	"fmt"
	"image/color"
	"os"
	"vlog"

	"github.com/aarzilli/nucular"
	"github.com/aarzilli/nucular/label"
	"github.com/aarzilli/nucular/rect"
	nstyle "github.com/aarzilli/nucular/style"
)

type portGroup struct {
	Theme         nstyle.Theme
	Module        Module_cfg
	Message       string
	Checkbox      [SERIAL_PORT_MAX]bool
	TestCmdEditor [SERIAL_PORT_MAX]nucular.TextEditor
	CheckValues   [SERIAL_PORT_MAX]int
	CurrentPortId [SERIAL_PORT_MAX]int
	ComboPorts    []string
}

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

//////////////////////////////////////////////////////////////////////////////////
// show handle

func newportGroup() (pg *portGroup) {
	pg = &portGroup{}

	pg.Message = ""
	pg.Theme = mytheme
	pg.Module = SIM800C
	pg.ComboPorts = ports_list

	for port_id := 0; port_id < SERIAL_PORT_MAX; port_id++ {
		//pg.TestCmdEditor[port_id].Flags = nucular.EditSelectable
		pg.TestCmdEditor[port_id].Flags |= nucular.EditBox
		pg.TestCmdEditor[port_id].Flags |= nucular.EditNeverInsertMode
		pg.TestCmdEditor[port_id].Buffer = []rune("AT")
		pg.TestCmdEditor[port_id].Maxlen = 128
	}

	myBtnTab[Btn_CMD_Produce] = BtnDoTable{Btn_CMD_Produce, "produce", serialProduce}
	myBtnTab[Btn_CMD_CheckDo] = BtnDoTable{Btn_CMD_CheckDo, "checkdo", serialCheckDo}
	myBtnTab[Btn_CMD_Close] = BtnDoTable{Btn_CMD_Close, "close", serialClose}

	return pg
}

func (pg *portGroup) showUI(w *nucular.Window) {
	pg.showMenuBar(w)
	w.Row(5).Dynamic(1)

	for port_id := 0; port_id < gConfig.Serial.Serial_max; port_id++ {
		pg.showPortG(w, port_id)
	}

	w.Row(5).Dynamic(1)
	w.Row(30).Dynamic(4)
	if w.Button(label.T("RefreshPort"), false) {
		pg.btnRefreshPort(w)
	}

	if w.Button(label.T("LoadToken"), false) {
		pg.btnLoadToken(w)
	}

	if w.Button(label.T("ProduceAll"), false) {
		pg.btnHandleAll(w, Btn_CMD_Produce, true)
	}

	if w.Button(label.T("CheckDoAll"), false) {
		pg.btnHandleAll(w, Btn_CMD_CheckDo, true)
	}

	w.Row(30).Dynamic(2)
	if w.Button(label.T("CloseAll"), false) {
		pg.btnHandleAll(w, Btn_CMD_Close, false)
	}

	if w.Button(label.T("Quit"), false) {
		pg.btnExit(w)
	}
}

func (pg *portGroup) showMenuBar(w *nucular.Window) {
	w.Row(25).Static(400, 100, 100)
	clryellow := color.RGBA{0xff, 0xff, 0x00, 0xff}
	w.LabelColored(fmt.Sprintf("** %s  (%s) **", szBanner, szVersion), "CC", clryellow)

	w.MenubarBegin()
	if w := w.Menu(label.TA("Module", "RC"), 150, nil); w != nil {
		w.Row(25).Dynamic(1)
		newmodule := pg.Module
		if w.OptionText("SIM800C", newmodule == SIM800C) {
			newmodule = SIM800C
		}
		if w.OptionText("EC20", newmodule == EC20) {
			newmodule = EC20
		}
		if w.OptionText("EC20_AUTO", newmodule == EC20_AUTO) {
			newmodule = EC20_AUTO
		}

		if gConfig.Simfake == 1 {
			if w.OptionText("EC20_TC1", newmodule == EC20_TC1) {
				newmodule = EC20_TC1
			}
			if w.OptionText("EC20_TC3", newmodule == EC20_TC3) {
				newmodule = EC20_TC3
			}
		}

		if newmodule != pg.Module {
			pg.Module = newmodule
			pg.setModule(w, newmodule)
		}
	}

	if w := w.Menu(label.TA("Theme", "RC"), 150, nil); w != nil {
		w.Row(25).Dynamic(1)
		newtheme := pg.Theme
		if w.OptionText("Default Theme", newtheme == nstyle.DefaultTheme) {
			newtheme = nstyle.DefaultTheme
		}
		if w.OptionText("White Theme", newtheme == nstyle.WhiteTheme) {
			newtheme = nstyle.WhiteTheme
		}
		if w.OptionText("Red Theme", newtheme == nstyle.RedTheme) {
			newtheme = nstyle.RedTheme
		}
		if w.OptionText("Dark Theme", newtheme == nstyle.DarkTheme) {
			newtheme = nstyle.DarkTheme
		}
		if newtheme != pg.Theme {
			pg.Theme = newtheme
			w.Master().SetStyle(nstyle.FromTheme(pg.Theme, w.Master().Style().Scaling))
		}
	}

	w.MenubarEnd()
}

func (pg *portGroup) showPortG(w *nucular.Window, portid int) {
	w.Row(40).Dynamic(1)
	clrred := color.RGBA{0xff, 0x00, 0x00, 0xff}
	clrgreen := color.RGBA{0x00, 0xff, 0x00, 0xff}

	if sw := w.GroupBegin("Group Port", nucular.WindowNoScrollbar|nucular.WindowBorder); sw != nil {
		sw.Row(2).Dynamic(1)
		sw.Row(26).Static(85, 45, 70, 70, 70, 10, 70, 10, 70, 600)
		sw.CheckboxText(fmt.Sprintf("Port[%d]:", portid), &pg.Checkbox[portid])

		if serial_port[portid].strInfo == "" {
			sw.Label(string("(*)"), "LC")
		} else if serial_port[portid].strInfo == "XXX" {
			sw.LabelColored(fmt.Sprintf("(%s)", serial_port[portid].strInfo), "LC", clrred)
		} else {
			sw.LabelColored(fmt.Sprintf("(%s)", serial_port[portid].strInfo), "LC", clrgreen)
		}

		pg.CurrentPortId[portid] = sw.ComboSimple(pg.ComboPorts, pg.CurrentPortId[portid], 20)
		strCom := COM_RNAME_PREFIX + pg.ComboPorts[pg.CurrentPortId[portid]]

		if sw.Button(label.T("Open"), false) {
			pg.btnOpen(sw, portid, strCom)
		}

		if sw.Button(label.T("Produce"), false) {
			pg.btnProduce(sw, portid, strCom)
		}

		sw.Label(string(" "), "LC")
		if sw.Button(label.T("Close"), false) {
			pg.btnClose(sw, portid, strCom)
		}

		sw.Label(string(" "), "LC")
		if sw.Button(label.T("ATsend"), false) {
			pg.btnATSend(sw, portid, strCom)
		}

		pg.TestCmdEditor[portid].Edit(sw)

		sw.GroupEnd()
	}
}

func (pg *portGroup) openMessage(w *nucular.Window, message string) {
	var wf nucular.WindowFlags
	wf |= nucular.WindowBorder
	wf |= nucular.WindowMovable
	wf |= nucular.WindowNoScrollbar
	wf |= nucular.WindowClosable
	wf |= nucular.WindowTitle
	pg.Message = message
	w.Master().PopupOpen("Message", wf, rect.Rect{170, 100, 300, 190}, true, pg.openMessageBox)
}

func (pg *portGroup) openMessageBox(w *nucular.Window) {
	w.Row(30).Dynamic(1)
	w.Label(fmt.Sprintf("%s", pg.Message), "CC")
	w.Row(30).Dynamic(1)
	w.Row(30).Dynamic(1)
	if w.Button(label.T("OK"), false) {
		w.Close()
	}
}

//////////////////////////////////////////////////////////////////////////////////
// btn handle

func (pg *portGroup) btnOpen(w *nucular.Window, portid int, strCom string) {
	if serialOpen(portid, strCom) != 0 {
		msg := fmt.Sprintf("Filed to open the %s!", strCom)
		pg.openMessage(w, msg)
	}
}

func (pg *portGroup) btnClose(w *nucular.Window, portid int, strCom string) {
	serialClose(portid)
}

func (pg *portGroup) btnATSend(w *nucular.Window, portid int, strCom string) {
	if portIsOK(portid) == 0 {
		msg := fmt.Sprintf("Please open the port[%d]!", portid)
		pg.openMessage(w, msg)
	} else {
		strCmd := string(pg.TestCmdEditor[portid].Buffer)
		serialATsendCmd(portid, strCom, strCmd)
	}
}

func (pg *portGroup) btnProduce(w *nucular.Window, portid int, strCom string) {
	if portIsOK(portid) == 0 {
		msg := fmt.Sprintf("Please open the port[%d]!", portid)
		pg.openMessage(w, msg)
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

func (pg *portGroup) btnHandleAll(w *nucular.Window, oper int, check bool) {
	if check && (pg.checkBox(w) == false) {
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

func (pg *portGroup) btnExit(w *nucular.Window) {
	os.Exit(1)
}

func (pg *portGroup) btnLoadToken(w *nucular.Window) {
	loadTokenCfg(gConfig.Token.Cmcc_file, OPER_CN_MOBILE)
	loadTokenCfg(gConfig.Token.Uni_file, OPER_CN_UNICOM)
	loadTokenCfg(gConfig.Token.Tel_file, OPER_CN_TELECOM)
}

func (pg *portGroup) btnRefreshPort(w *nucular.Window) {
	pg.btnHandleAll(w, Btn_CMD_Close, false)

	for port_id := 0; port_id < gConfig.Serial.Serial_max; port_id++ {
		pg.CurrentPortId[port_id] = 0
	}

	ports_list = serialList()
	pg.ComboPorts = ports_list
	vlog.Info("Portlists: %v", ports_list)
}

//////////////////////////////////////////////////////////////////////////////////
// other handle

func (pg *portGroup) setTaskBtn(oper int, portid int, taskCH chan PortResult) {
	ret := myBtnTab[oper].BtnFunc(portid)
	///wg.Done()
	resp := PortResult{
		Portid: portid,
		Oper:   oper,
		Result: ret,
	}
	vlog.Info("Handle-Put result of %s port[%d]: %d", myBtnTab[resp.Oper].BtnStr, resp.Portid, resp.Result)
	taskCH <- resp
}

func (pg *portGroup) getTaskBtn(count int, taskCH chan PortResult) {
	for i := 0; i < count; i++ {
		resp := <-taskCH
		vlog.Info("Handle-Get result of %s port[%d]: %d", myBtnTab[resp.Oper].BtnStr, resp.Portid, resp.Result)
		if (resp.Oper == Btn_CMD_Produce) || (resp.Oper == Btn_CMD_CheckDo) {
			if resp.Result == 0 {
				serial_port[resp.Portid].strInfo = "OK"
			} else {
				serial_port[resp.Portid].strInfo = "XXX"
			}
		}
	}
	close(taskCH)
}

func (pg *portGroup) checkBox(w *nucular.Window) bool {
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
		pg.openMessage(w, "Please select(open) the ports!")
		return false
	} else if portlist != "" {
		msg := fmt.Sprintf("Please select(open) the ports: [%s]!", portlist)
		pg.openMessage(w, msg)
		return false
	}

	return true
}

func (pg *portGroup) setModule(w *nucular.Window, module Module_cfg) {
	pg.btnHandleAll(w, Btn_CMD_Close, false)
	module_reinit(pg.Module)
}
