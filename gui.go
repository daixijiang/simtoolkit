/* VSIM Serial Product Toolkit
 * Author: daixijiang@gmail.com (2019)
 */

package main

import (
	"fmt"
	"vlog"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

type portGroup struct {
	MainWindow *walk.MainWindow
	portText   [SERIAL_PORT_MAX]*walk.LineEdit
	portCombo  [SERIAL_PORT_MAX]*walk.ComboBox
	modRadiuo  [MODULE_MAX]*walk.RadioButton

	Module   Module_cfg
	Message  string
	PortList []string
	Checkbox [SERIAL_PORT_MAX]bool
}

var defColor = walk.RGB(255, 255, 0)
var errColor = walk.RGB(255, 0, 0)
var okColor = walk.RGB(255, 255, 0)

//////////////////////////////////////////////////////////////////////////////////

func (pg *portGroup) showUI() {
	portWidget := make([]Widget, 0)
	portWidget = append(portWidget, pg.showBanner())
	for port_id := 0; port_id < gConfig.Serial.Serial_max; port_id++ {
		portWidget = append(portWidget, pg.showPortG(port_id))
	}
	portWidget = append(portWidget, pg.showBtnG())

	if err := (MainWindow{
		AssignTo: &pg.MainWindow,
		Title:    szTitle,
		//Icon: "./main.ico",
		MinSize: Size{600, 400},
		Font:    Font{Family: "", PointSize: 10},
		Layout:  VBox{MarginsZero: true},

		MenuItems: pg.showMenu(),
		Children:  portWidget,
	}.Create()); err != nil {
		vlog.Error("MainWindow create: %s", err)
	}

	pg.modRadiuo[pg.Module].SetChecked(true)
	pg.MainWindow.Run()
}

func (pg *portGroup) showMenu() []MenuItem {
	return []MenuItem{
		Menu{
			Text: "&File",
			Items: []MenuItem{
				Separator{},
				Action{
					Text: "Exit",
					OnTriggered: func() {
						pg.btnExit()
					},
				},
			},
		},
		Menu{
			Text: "&Help",
			Items: []MenuItem{
				Action{
					Text: "About",
					OnTriggered: func() {
						aboutTitle := "About"
						aboutMsg := fmt.Sprintf("%s\r\n(%s)", szBanner, szVersion)
						walk.MsgBox(pg.MainWindow, aboutTitle, aboutMsg, walk.MsgBoxIconInformation)
					},
				},
			},
		},
	}
}

func (pg *portGroup) showBanner() Composite {
	return Composite{
		Layout: Grid{Columns: 4, Spacing: 5},
		Children: []Widget{
			Label{
				Text: "Module: ",
				//Background: SolidColorBrush{255,0,0},
			},
			RadioButtonGroup{
				Buttons: []RadioButton{
					RadioButton{
						AssignTo: &pg.modRadiuo[SIM800C],
						Name:     "SIM800C",
						Text:     "SIM800C",
						Value:    "0",
						OnClicked: func() {
							pg.setModule(SIM800C)
						},
					},

					RadioButton{
						AssignTo: &pg.modRadiuo[EC20],
						Name:     "EC20",
						Text:     "EC20",
						Value:    "1",
						OnClicked: func() {
							pg.setModule(EC20)
						},
					},
					RadioButton{
						AssignTo: &pg.modRadiuo[EC20_AUTO],
						Name:     "EC20_AUTO",
						Text:     "EC20_AUTO",
						Value:    "2",
						OnClicked: func() {
							pg.setModule(EC20_AUTO)
						},
					},
				},
			},
		},
	}
}

func (pg *portGroup) showPortG(portid int) Composite {
	var checkCKB *walk.CheckBox
	var statusTE *walk.LineEdit
	var cmdLE *walk.LineEdit

	return Composite{
		Layout:  HBox{},
		MaxSize: Size{0, 30},
		Children: []Widget{
			CheckBox{
				AssignTo: &checkCKB,
				MaxSize:  Size{60, 0},
				Text:     fmt.Sprintf("Port[%d]", portid),
				OnClicked: func() {
					if checkCKB.CheckState() == 0 {
						pg.Checkbox[portid] = false
					} else {
						pg.Checkbox[portid] = true
					}
				},
			},
			LineEdit{
				AssignTo: &pg.portText[portid],
				MaxSize:  Size{40, 0},
				Text:     "*",
				ReadOnly: true,
			},
			ComboBox{
				AssignTo:     &pg.portCombo[portid],
				Model:        pg.PortList,
				CurrentIndex: 0,
			},
			PushButton{
				Text: "Open",
				OnClicked: func() {
					comId := pg.portCombo[portid].CurrentIndex()
					if comId < 0 {
						comId = 0
					}
					strCom := COM_RNAME_PREFIX + pg.PortList[comId]
					pg.btnOpen(portid, strCom)
					statusTE.SetText(fmt.Sprintf("Open port %s", strCom))

					pg.portText[portid].SetTextColor(okColor)
					pg.portText[portid].SetText(fmt.Sprintf("%s", serial_port[portid].strInfo))
				},
			},
			PushButton{
				Text: "Produce",
				OnClicked: func() {
					comId := pg.portCombo[portid].CurrentIndex()
					if comId < 0 {
						comId = 0
					}
					strCom := COM_RNAME_PREFIX + pg.PortList[comId]
					pg.btnProduce(portid, strCom)
					statusTE.SetText(fmt.Sprintf("Produce port %s", strCom))
				},
			},
			PushButton{
				Text: "Close",
				OnClicked: func() {
					comId := pg.portCombo[portid].CurrentIndex()
					if comId < 0 {
						comId = 0
					}
					strCom := COM_RNAME_PREFIX + pg.PortList[comId]
					pg.btnClose(portid, strCom)
					statusTE.SetText(fmt.Sprintf("Close port %s", strCom))

					pg.portText[portid].SetTextColor(okColor)
					pg.portText[portid].SetText(fmt.Sprintf("%s", serial_port[portid].strInfo))
				},
			},
			PushButton{
				Text: "ATsend",
				OnClicked: func() {
					comId := pg.portCombo[portid].CurrentIndex()
					if comId < 0 {
						comId = 0
					}
					strCom := COM_RNAME_PREFIX + pg.PortList[comId]
					resp := pg.btnATSend(portid, strCom, cmdLE.Text())
					statusTE.SetText(fmt.Sprintf("%s", resp))
				},
			},
			LineEdit{
				AssignTo: &cmdLE,
				Text:     "AT",
			},
			LineEdit{
				AssignTo: &statusTE,
				ReadOnly: true,
			},
		},
	}
}

func (pg *portGroup) showBtnG() Composite {
	return Composite{
		Layout: Grid{Columns: 3, Spacing: 10},
		Children: []Widget{
			PushButton{
				Text: "ProduceAll",
				OnClicked: func() {
					fmt.Printf("ProduceAll: %v\n", pg.Checkbox)
					pg.btnHandleAll(Btn_CMD_Produce, true)
				},
			},
			PushButton{
				Text: "CheckDoAll",
				OnClicked: func() {
					fmt.Printf("CheckDoAll: %v\n", pg.Checkbox)
					pg.btnHandleAll(Btn_CMD_CheckDo, true)
				},
			},
			PushButton{
				Text: "CloseAll",
				OnClicked: func() {
					fmt.Printf("CloseAll: %v\n", pg.Checkbox)
					pg.btnHandleAll(Btn_CMD_Close, false)
				},
			},
			PushButton{
				Text: "RefreshPort",
				OnClicked: func() {
					fmt.Printf("RefreshPort: %v\n", pg.Checkbox)
					pg.btnRefreshPort()
				},
			},
			PushButton{
				Text: "LoadToken",
				OnClicked: func() {
					fmt.Printf("LoadToken: %v\n", pg.Checkbox)
					pg.btnLoadToken()
				},
			},
			PushButton{
				Text: "Quit",
				OnClicked: func() {
					fmt.Printf("Quit: %v\n", pg.Checkbox)
					pg.btnExit()
				},
			},
		},
	}
}

//////////////////////////////////////////////////////////////////////////////////

func (pg *portGroup) openMessage(message string) {
	walk.MsgBox(pg.MainWindow, "Message", message, walk.MsgBoxIconInformation)
}

func (pg *portGroup) doResult(resp PortResult) {
	if (resp.Oper == Btn_CMD_Produce) || (resp.Oper == Btn_CMD_CheckDo) {
		if resp.Result == 0 {
			serial_port[resp.Portid].strInfo = "OK"
			pg.portText[resp.Portid].SetTextColor(okColor)
			pg.portText[resp.Portid].SetText(fmt.Sprintf("%s", serial_port[resp.Portid].strInfo))
		} else {
			serial_port[resp.Portid].strInfo = "XXX"
			pg.portText[resp.Portid].SetTextColor(errColor)
			pg.portText[resp.Portid].SetText(fmt.Sprintf("%s", serial_port[resp.Portid].strInfo))
		}
	} else if resp.Oper == Btn_CMD_Close {
		serial_port[resp.Portid].strInfo = "*"
		pg.portText[resp.Portid].SetTextColor(defColor)
		pg.portText[resp.Portid].SetText(fmt.Sprintf("%s", serial_port[resp.Portid].strInfo))
	}
}
