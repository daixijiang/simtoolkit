/* VSIM Serial Product Toolkit
 * Author: daixijiang@gmail.com (2019)
 */

package main

import (
	"fmt"
	"image/color"

	"github.com/aarzilli/nucular"
	"github.com/aarzilli/nucular/label"
	"github.com/aarzilli/nucular/rect"
	nstyle "github.com/aarzilli/nucular/style"
)

type nucularUI struct {
	MainWindow    *nucular.Window
	CmdEditor     [SERIAL_PORT_MAX]nucular.TextEditor
	CurrentPortId [SERIAL_PORT_MAX]int

	Binfo *TBaseInfo
}

func newNucularUI() (nui *nucularUI) {
	newui := &nucularUI{}
	newui.SelfInit()
	newui.Binfo = newTBaseInfo()
	newui.Binfo.nui = newui

	return newui
}

//////////////////////////////////////////////////////////////////////////////////

func (nui *nucularUI) SelfInit() {
	for port_id := 0; port_id < gConfig.Serial.Serial_max; port_id++ {
		nui.CurrentPortId[port_id] = 0
		nui.CmdEditor[port_id].Flags |= nucular.EditBox
		nui.CmdEditor[port_id].Buffer = []rune("AT")
		nui.CmdEditor[port_id].Maxlen = 128
	}
}

func (nui *nucularUI) RunUI() {
	wnd := nucular.NewMasterWindow(0, szTitle, nui.showUI)
	wnd.SetStyle(nstyle.FromTheme(nstyle.DarkTheme, gConfig.Scaling))
	wnd.Main()
}

func (nui *nucularUI) showUI(w *nucular.Window) {
	nui.MainWindow = w

	nui.showMenuBar(w)
	w.Row(5).Dynamic(1)

	for port_id := 0; port_id < gConfig.Serial.Serial_max; port_id++ {
		nui.showPortG(w, port_id)
	}

	w.Row(5).Dynamic(1)
	w.Row(30).Dynamic(3)
	if w.Button(label.T("ProduceAll"), false) {
		nui.Binfo.btnHandleAll(Btn_CMD_Produce, true)
	}

	if w.Button(label.T("CheckDoAll"), false) {
		nui.Binfo.btnHandleAll(Btn_CMD_CheckDo, true)
	}

	if w.Button(label.T("CloseAll"), false) {
		nui.Binfo.btnHandleAll(Btn_CMD_Close, false)
	}

	w.Row(30).Dynamic(3)
	if w.Button(label.T("RefreshPort"), false) {
		nui.Binfo.btnRefreshPort()
	}

	if w.Button(label.T("LoadToken"), false) {
		nui.Binfo.btnLoadToken()
	}

	if w.Button(label.T("Quit"), false) {
		nui.Binfo.btnExit()
	}
}

func (nui *nucularUI) showMenuBar(w *nucular.Window) {
	w.Row(25).Static(400, 100, 100)
	clryellow := color.RGBA{0xff, 0xff, 0x00, 0xff}
	w.LabelColored(fmt.Sprintf("** %s  (%s) **", szBanner, szVersion), "CC", clryellow)

	w.MenubarBegin()
	if w := w.Menu(label.TA("Module", "RC"), 150, nil); w != nil {
		w.Row(25).Dynamic(1)
		newmodule := nui.Binfo.Module
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

		if newmodule != nui.Binfo.Module {
			nui.Binfo.Module = newmodule
			nui.Binfo.setModule(newmodule)
		}
	}

	w.MenubarEnd()
}

func (nui *nucularUI) showPortG(w *nucular.Window, portid int) {
	w.Row(40).Dynamic(1)
	clrred := color.RGBA{0xff, 0x00, 0x00, 0xff}
	clrgreen := color.RGBA{0x00, 0xff, 0x00, 0xff}

	if sw := w.GroupBegin("Group Port", nucular.WindowNoScrollbar|nucular.WindowBorder); sw != nil {
		sw.Row(2).Dynamic(1)
		sw.Row(26).Static(85, 45, 70, 70, 70, 10, 70, 10, 70, 600)
		sw.CheckboxText(fmt.Sprintf("Port[%d]:", portid), &nui.Binfo.Checkbox[portid])

		if serial_port[portid].strInfo == "" {
			sw.Label(string("(*)"), "LC")
		} else if serial_port[portid].strInfo == "XXX" {
			sw.LabelColored(fmt.Sprintf("(%s)", serial_port[portid].strInfo), "LC", clrred)
		} else {
			sw.LabelColored(fmt.Sprintf("(%s)", serial_port[portid].strInfo), "LC", clrgreen)
		}

		nui.CurrentPortId[portid] = sw.ComboSimple(nui.Binfo.PortList, nui.CurrentPortId[portid], 20)
		strCom := COM_RNAME_PREFIX + nui.Binfo.PortList[nui.CurrentPortId[portid]]

		if sw.Button(label.T("Open"), false) {
			nui.Binfo.btnOpen(portid, strCom)
		}

		if sw.Button(label.T("Produce"), false) {
			nui.Binfo.btnProduce(portid, strCom)
		}

		sw.Label(string(" "), "LC")
		if sw.Button(label.T("Close"), false) {
			nui.Binfo.btnClose(portid, strCom)
		}

		sw.Label(string(" "), "LC")

		if sw.Button(label.T("ATsend"), false) {
			nui.Binfo.btnATSend(portid, strCom, string(nui.CmdEditor[portid].Buffer))
		}

		nui.CmdEditor[portid].Edit(sw)

		sw.GroupEnd()
	}
}

//////////////////////////////////////////////////////////////////////////////////

func (nui *nucularUI) openMessage(message string) {
	var wf nucular.WindowFlags
	wf |= nucular.WindowBorder
	wf |= nucular.WindowMovable
	wf |= nucular.WindowNoScrollbar
	wf |= nucular.WindowClosable
	wf |= nucular.WindowTitle
	nui.Binfo.Message = message
	nui.MainWindow.Master().PopupOpen("Message", wf, rect.Rect{170, 100, 300, 190}, true, nui.openMessageBox)
}

func (nui *nucularUI) openMessageBox(w *nucular.Window) {
	w.Row(30).Dynamic(1)
	w.Label(fmt.Sprintf("%s", nui.Binfo.Message), "CC")
	w.Row(30).Dynamic(1)
	w.Row(30).Dynamic(1)
	if w.Button(label.T("OK"), false) {
		w.Close()
	}
}

func (nui *nucularUI) doResult(resp PortResult) {
	if (resp.Oper == Btn_CMD_Produce) || (resp.Oper == Btn_CMD_CheckDo) {
		if resp.Result == 0 {
			serial_port[resp.Portid].strInfo = "OK"
		} else {
			serial_port[resp.Portid].strInfo = "XXX"
		}
	} else if resp.Oper == Btn_CMD_Close {
		serial_port[resp.Portid].strInfo = "*"
	}
}
