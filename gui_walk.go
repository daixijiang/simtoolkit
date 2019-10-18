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

type walkUI struct {
	MainWindow *walk.MainWindow
	portText   [SERIAL_PORT_MAX]*walk.LineEdit
	portCombo  [SERIAL_PORT_MAX]*walk.ComboBox
	modRadiuo  [MODULE_MAX]*walk.RadioButton

	Binfo *TBaseInfo
}

var defColor = walk.RGB(255, 255, 0)
var errColor = walk.RGB(255, 0, 0)
var okColor = walk.RGB(255, 255, 0)

func newWalkUI() (wui *walkUI) {
	newui := &walkUI{}
	newui.SelfInit()
	newui.Binfo = newTBaseInfo()
	newui.Binfo.wui = newui

	return newui
}

//////////////////////////////////////////////////////////////////////////////////

func (wui *walkUI) SelfInit() {}

func (wui *walkUI) RunUI() {
	wui.showUI()
}

func (wui *walkUI) showUI() {
	portWidget := make([]Widget, 0)
	portWidget = append(portWidget, wui.showBanner())
	for port_id := 0; port_id < gConfig.Serial.Serial_max; port_id++ {
		portWidget = append(portWidget, wui.showPortG(port_id))
	}
	portWidget = append(portWidget, wui.showBtnG())

	if err := (MainWindow{
		AssignTo: &wui.MainWindow,
		Title:    szTitle,
		//Icon: "./main.ico",
		MinSize: Size{600, 400},
		Font:    Font{Family: "", PointSize: 10},
		Layout:  VBox{MarginsZero: true},

		MenuItems: wui.showMenu(),
		Children:  portWidget,
	}.Create()); err != nil {
		vlog.Error("MainWindow create: %s", err)
	}

	wui.modRadiuo[wui.Binfo.Module].SetChecked(true)
	wui.MainWindow.Run()
}

func (wui *walkUI) showMenu() []MenuItem {
	return []MenuItem{
		Menu{
			Text: "&File",
			Items: []MenuItem{
				Separator{},
				Action{
					Text: "Exit",
					OnTriggered: func() {
						wui.Binfo.btnExit()
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
						walk.MsgBox(wui.MainWindow, aboutTitle, aboutMsg, walk.MsgBoxIconInformation)
					},
				},
			},
		},
	}
}

func (wui *walkUI) showBanner() Composite {
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
						AssignTo: &wui.modRadiuo[SIM800C],
						Name:     "SIM800C",
						Text:     "SIM800C",
						Value:    "0",
						OnClicked: func() {
							wui.Binfo.setModule(SIM800C)
						},
					},

					RadioButton{
						AssignTo: &wui.modRadiuo[EC20],
						Name:     "EC20",
						Text:     "EC20",
						Value:    "1",
						OnClicked: func() {
							wui.Binfo.setModule(EC20)
						},
					},
					RadioButton{
						AssignTo: &wui.modRadiuo[EC20_AUTO],
						Name:     "EC20_AUTO",
						Text:     "EC20_AUTO",
						Value:    "2",
						OnClicked: func() {
							wui.Binfo.setModule(EC20_AUTO)
						},
					},
				},
			},
		},
	}
}

func (wui *walkUI) showPortG(portid int) Composite {
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
						wui.Binfo.Checkbox[portid] = false
					} else {
						wui.Binfo.Checkbox[portid] = true
					}
				},
			},
			LineEdit{
				AssignTo: &wui.portText[portid],
				MaxSize:  Size{40, 0},
				Text:     "*",
				ReadOnly: true,
			},
			ComboBox{
				AssignTo:     &wui.portCombo[portid],
				Model:        wui.Binfo.PortList,
				CurrentIndex: 0,
			},
			PushButton{
				Text: "Open",
				OnClicked: func() {
					comId := wui.portCombo[portid].CurrentIndex()
					if comId < 0 {
						comId = 0
					}
					strCom := COM_RNAME_PREFIX + wui.Binfo.PortList[comId]
					wui.Binfo.btnOpen(portid, strCom)
					statusTE.SetText(fmt.Sprintf("Open port %s", strCom))

					wui.portText[portid].SetTextColor(okColor)
					wui.portText[portid].SetText(fmt.Sprintf("%s", serial_port[portid].strInfo))
				},
			},
			PushButton{
				Text: "Produce",
				OnClicked: func() {
					comId := wui.portCombo[portid].CurrentIndex()
					if comId < 0 {
						comId = 0
					}
					strCom := COM_RNAME_PREFIX + wui.Binfo.PortList[comId]
					wui.Binfo.btnProduce(portid, strCom)
					statusTE.SetText(fmt.Sprintf("Produce port %s", strCom))
				},
			},
			PushButton{
				Text: "Close",
				OnClicked: func() {
					comId := wui.portCombo[portid].CurrentIndex()
					if comId < 0 {
						comId = 0
					}
					strCom := COM_RNAME_PREFIX + wui.Binfo.PortList[comId]
					wui.Binfo.btnClose(portid, strCom)
					statusTE.SetText(fmt.Sprintf("Close port %s", strCom))

					wui.portText[portid].SetTextColor(okColor)
					wui.portText[portid].SetText(fmt.Sprintf("%s", serial_port[portid].strInfo))
				},
			},
			PushButton{
				Text: "ATsend",
				OnClicked: func() {
					comId := wui.portCombo[portid].CurrentIndex()
					if comId < 0 {
						comId = 0
					}
					strCom := COM_RNAME_PREFIX + wui.Binfo.PortList[comId]
					resp := wui.Binfo.btnATSend(portid, strCom, cmdLE.Text())
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

func (wui *walkUI) showBtnG() Composite {
	return Composite{
		Layout: Grid{Columns: 3, Spacing: 10},
		Children: []Widget{
			PushButton{
				Text: "ProduceAll",
				OnClicked: func() {
					fmt.Printf("ProduceAll: %v\n", wui.Binfo.Checkbox)
					wui.Binfo.btnHandleAll(Btn_CMD_Produce, true)
				},
			},
			PushButton{
				Text: "CheckDoAll",
				OnClicked: func() {
					fmt.Printf("CheckDoAll: %v\n", wui.Binfo.Checkbox)
					wui.Binfo.btnHandleAll(Btn_CMD_CheckDo, true)
				},
			},
			PushButton{
				Text: "CloseAll",
				OnClicked: func() {
					fmt.Printf("CloseAll: %v\n", wui.Binfo.Checkbox)
					wui.Binfo.btnHandleAll(Btn_CMD_Close, false)
				},
			},
			PushButton{
				Text: "RefreshPort",
				OnClicked: func() {
					fmt.Printf("RefreshPort: %v\n", wui.Binfo.Checkbox)
					wui.Binfo.btnRefreshPort()
				},
			},
			PushButton{
				Text: "LoadToken",
				OnClicked: func() {
					fmt.Printf("LoadToken: %v\n", wui.Binfo.Checkbox)
					wui.Binfo.btnLoadToken()
				},
			},
			PushButton{
				Text: "Quit",
				OnClicked: func() {
					fmt.Printf("Quit: %v\n", wui.Binfo.Checkbox)
					wui.Binfo.btnExit()
				},
			},
		},
	}
}

//////////////////////////////////////////////////////////////////////////////////

func (wui *walkUI) openMessage(message string) {
	walk.MsgBox(wui.MainWindow, "Message", message, walk.MsgBoxIconInformation)
}

func (wui *walkUI) doResult(resp PortResult) {
	if (resp.Oper == Btn_CMD_Produce) || (resp.Oper == Btn_CMD_CheckDo) {
		if resp.Result == 0 {
			serial_port[resp.Portid].strInfo = "OK"
			wui.portText[resp.Portid].SetTextColor(okColor)
			wui.portText[resp.Portid].SetText(fmt.Sprintf("%s", serial_port[resp.Portid].strInfo))
		} else {
			serial_port[resp.Portid].strInfo = "XXX"
			wui.portText[resp.Portid].SetTextColor(errColor)
			wui.portText[resp.Portid].SetText(fmt.Sprintf("%s", serial_port[resp.Portid].strInfo))
		}
	} else if resp.Oper == Btn_CMD_Close {
		serial_port[resp.Portid].strInfo = "*"
		wui.portText[resp.Portid].SetTextColor(defColor)
		wui.portText[resp.Portid].SetText(fmt.Sprintf("%s", serial_port[resp.Portid].strInfo))
	}
}
