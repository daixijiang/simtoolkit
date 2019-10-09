/* VSIM Serial Product Toolkit
 * Author: daixijiang@gmail.com (2019)
 */

package main

import (
	"fmt"
	"image/color"
	"os"
	"sync"
	"vlog"

	"github.com/aarzilli/nucular"
	"github.com/aarzilli/nucular/label"
	"github.com/aarzilli/nucular/rect"
	nstyle "github.com/aarzilli/nucular/style"
)

type portGroup struct {
	Theme         nstyle.Theme
	Message       string
	Checkbox      [SERAIL_PORT_MAX]bool
	TestCmdEditor [SERAIL_PORT_MAX]nucular.TextEditor
	CheckValues   [SERAIL_PORT_MAX]int
	CurrentPortId [SERAIL_PORT_MAX]int
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

var wg sync.WaitGroup

func newportGroup() (pg *portGroup) {
	pg = &portGroup{}

	pg.Message = ""
	pg.Theme = mytheme
	pg.ComboPorts = ports_list

	for port_id := 0; port_id < SERAIL_PORT_MAX; port_id++ {
		pg.TestCmdEditor[port_id].Flags = nucular.EditSelectable
		pg.TestCmdEditor[port_id].Buffer = []rune("AT")
		pg.TestCmdEditor[port_id].Maxlen = 64
	}

	myBtnTab[Btn_CMD_Produce] = BtnDoTable{Btn_CMD_Produce, "produce", serialProduce}
	myBtnTab[Btn_CMD_CheckDo] = BtnDoTable{Btn_CMD_CheckDo, "checkdo", serialCheckDo}
	myBtnTab[Btn_CMD_Close] = BtnDoTable{Btn_CMD_Close, "close", serialClose}

	return pg
}

func (pg *portGroup) showUI(w *nucular.Window) {
	pg.showMenuBar(w)
	w.Row(5).Dynamic(1)

	for port_id := 0; port_id < SERAIL_PORT_MAX; port_id++ {
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
		os.Exit(1)
	}
}

func (pg *portGroup) showMenuBar(w *nucular.Window) {
	w.Row(25).Static(450, 150)
	clryellow := color.RGBA{0xff, 0xff, 0x00, 0xff}
	w.LabelColored(fmt.Sprintf("** %s  (%s) **", szBanner, szVersion), "CC", clryellow)

	w.MenubarBegin()
	if w := w.Menu(label.TA("Module", "RC"), 150, nil); w != nil {
		w.Row(25).Dynamic(1)
		curMod := myProduce.Mod
		if w.OptionText("EC20", myProduce.Mod == &modEC20) {
			myProduce.Mod = &modEC20
			myProduce.UrlVer = SERVER_PLAIN_v0
		}
		if w.OptionText("SIM800C", myProduce.Mod == &modSIM800C) {
			myProduce.Mod = &modSIM800C
			myProduce.UrlVer = SERVER_Cipher_v1
		}
		if curMod != myProduce.Mod {
			pg.btnHandleAll(w, Btn_CMD_Close, false)
		}
	}

	w.MenubarEnd()
}

func (pg *portGroup) showPortG(w *nucular.Window, portid int) {
	w.Row(40).Dynamic(1)
	if sw := w.GroupBegin("Group Port", nucular.WindowNoScrollbar|nucular.WindowBorder); sw != nil {
		sw.Row(4).Dynamic(1)
		sw.Row(25).Static(90, 30, 70, 70, 135, 70, 70, 70)
		sw.CheckboxText(fmt.Sprintf("Port[%d]: ", portid), &pg.Checkbox[portid])

		if serial_port[portid].strInfo == "" {
			sw.Label(string("(*)"), "LC")
		} else {
			sw.Label(string("("+serial_port[portid].strInfo+")"), "LC")
		}

		if portid >= SERAIL_PORT_MAX {
			vlog.Info("showPortG %d", portid)
			return
		}
		pg.CurrentPortId[portid] = sw.ComboSimple(pg.ComboPorts, pg.CurrentPortId[portid], 25)
		strCom := COM_RNAME_PREFIX + pg.ComboPorts[pg.CurrentPortId[portid]]

		if sw.Button(label.T("Open"), false) {
			if serialOpen(portid, strCom) != 0 {
				msg := fmt.Sprintf("Filed to open the %s!", strCom)
				pg.openMessage(w, msg)
			}
		}

		pg.TestCmdEditor[portid].Edit(sw)
		strCmd := string(pg.TestCmdEditor[portid].Buffer)
		if sw.Button(label.T("ATsend"), false) {
			if portIsOK(portid) == 0 {
				msg := fmt.Sprintf("Please open the port[%d]!", portid)
				pg.openMessage(w, msg)
			} else {
				serialATsendCmd(portid, strCom, strCmd)
			}
		}

		if sw.Button(label.T("Produce"), false) {
			if portIsOK(portid) == 0 {
				msg := fmt.Sprintf("Please open the port[%d]!", portid)
				pg.openMessage(w, msg)
			} else {
				pg.btnProduce(sw, portid, strCom)
			}
		}

		if sw.Button(label.T("Close"), false) {
			serialClose(portid)
		}

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

func (pg *portGroup) btnProduce(w *nucular.Window, portid int, strCom string) {
	if portIsOK(portid) == 0 {
		msg := fmt.Sprintf("Please open the port[%d]!", portid)
		pg.openMessage(w, msg)
	} else {
		vlog.Info("Port[%d] => start produce %s", portid, strCom)

		if serial_port[portid].port_status != PORT_STATUS_PRODUCE {
			serial_port[portid].port_status = PORT_STATUS_PRODUCE
			wg.Add(1)
			go pg.taskBtnHandle(Btn_CMD_Produce, portid)
			wg.Wait()
		}
		serial_port[portid].port_status = PORT_STATUS_OPEN
		serial_port[portid].strInfo = fmt.Sprintf("%s", "P")
	}
}

func (pg *portGroup) btnHandleAll(w *nucular.Window, oper int, check bool) {
	if check {
		cntCheck := 0
		portlist := ""

		for port_id := 0; port_id < SERAIL_PORT_MAX; port_id++ {
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
		} else if check && (portlist != "") {
			msg := fmt.Sprintf("Please select(open) the ports: [%s]!", portlist)
			pg.openMessage(w, msg)
		}

		return
	}

	vlog.Info("start %s all", myBtnTab[oper].BtnStr)
	for port_id := 0; port_id < SERAIL_PORT_MAX; port_id++ {
		if (pg.Checkbox[port_id] || !check) && (portIsOK(port_id) != 0) {
			wg.Add(1)
			go pg.taskBtnHandle(oper, port_id)
		}
	}
	wg.Wait()
}

func (pg *portGroup) taskBtnHandle(oper int, portid int) {
	myBtnTab[oper].BtnFunc(portid)
	wg.Done()
}

func (pg *portGroup) btnLoadToken(w *nucular.Window) {
	loadTokenCfg("./token.cfg")
}

func (pg *portGroup) btnRefreshPort(w *nucular.Window) {
	pg.btnHandleAll(w, Btn_CMD_Close, false)
	ports_list = serialList()
	pg.ComboPorts = ports_list
	vlog.Info("Portlists: %v", ports_list)
}
