/* VSIM Serial Product Toolkit
 * Author: daixijiang@gmail.com (2019)
 */

package main

import (
	"fmt"
	"github.com/tarm/serial"
	"vlog"
)

/* ModuleTable */
const (
	// prepare
	Module_CMD1_SYSVER         = 1
	Module_CMD1_SOFTMODE       = 2
	Module_CMD1_RESET0         = 3
	Module_CMD1_RESET1         = 4
	Module_CMD1_RESET2         = 5
	Module_PRE1_SET_NETWORK0   = 6
	Module_PRE1_SET_NETWORK1   = 7
	Module_CMD1_SET_SERVURL    = 8
	Module_CMD1_AUTOSWITCH_ON  = 9
	Module_CMD1_AUTOSWITCH_OFF = 10
	Module_CMD1_BACKUP_CONFIG  = 11

	// produce
	Module_CMD2_IMEI   = 20
	Module_CMD2_COSVER = 21
	Module_CMD2_CHIPID = 22
	Module_CMD2_SIM192 = 23
	Module_CMD2_SIM64  = 24

	// check
	Module_CMD3_CCID       = 30
	Module_CMD3_CIMI       = 31
	Module_CMD3_CREG       = 32
	Module_CMD3_CEREG      = 33
	Module_CMD3_COPS       = 34
	Module_CMD3_SWITCH_TEL = 35
	Module_CMD3_SWITCH_CU  = 36
	Module_CMD3_SWITCH_CM  = 37

	Module_TAB_AT_CMD_MAX = 40
)

type Module_cfg int

const (
	SIM800C = Module_cfg(iota)
	EC20
	EC20_AUTO
	EC20_TP
	EC20_TC1
	EC20_TC3
)

type cmdHandler func(cmdid int, cmdstr string, portid int, s *serial.Port, reply *string) int

type ModCmdTable struct {
	CmdID   int
	CmdStr  string
	CmdFunc cmdHandler
}

type ModuleProduce struct {
	Type   Module_cfg
	ModCmd *[Module_TAB_AT_CMD_MAX]ModCmdTable
	UrlVer int
}

var modCmd_EC20 [Module_TAB_AT_CMD_MAX]ModCmdTable
var modCmd_SIM800C [Module_TAB_AT_CMD_MAX]ModCmdTable
var thisModule *ModuleProduce

func (mp *ModuleProduce) DoComCMD(cmdid int, portid int, result *string) int {
	if cmdid > Module_TAB_AT_CMD_MAX {
		return -1
	}

	funcdo := mp.ModCmd[cmdid].CmdFunc
	if funcdo == nil {
		return 0
	}

	if mp.ModCmd[cmdid].CmdStr == "" {
		return 0
	}

	vlog.Debug("DoComCMD: cmdid %d, cmd %s, func %p", cmdid, mp.ModCmd[cmdid].CmdStr, funcdo)
	return funcdo(cmdid, mp.ModCmd[cmdid].CmdStr, portid, serial_port[portid].comPort, result)
}

func module_init() {
	module_ec20_init(&modCmd_EC20)
	module_sim800c_init(&modCmd_SIM800C)
	// add list of module init

	// default is sim800c
	module := SIM800C
	if gConfig.Module == "sim800c" {
		module = SIM800C
	} else if gConfig.Module == "ec20" {
		module = EC20
	} else if gConfig.Module == "ec20_auto" {
		module = EC20_AUTO
	}
	module_reinit(module)
}

func module_get() Module_cfg {
	return thisModule.Type
}

func module_reinit(module Module_cfg) {
	fmt.Printf("module[%d] %s, simfake %d\n", module, gConfig.Module, gConfig.Simfake)
	if module == SIM800C {
		thisModule = &ModuleProduce{
			Type:   module,
			ModCmd: &modCmd_SIM800C,
			UrlVer: SERVER_Cipher,
		}
	} else if module == EC20 {
		thisModule = &ModuleProduce{
			Type:   module,
			ModCmd: &modCmd_EC20,
			UrlVer: SERVER_PLAIN_v0,
		}
	} else if module == EC20_AUTO {
		thisModule = &ModuleProduce{
			Type:   module,
			ModCmd: &modCmd_EC20,
			UrlVer: SERVER_PLAIN_v0,
		}
	} else {
		// default is sim800c
		thisModule = &ModuleProduce{
			Type:   module,
			ModCmd: &modCmd_SIM800C,
			UrlVer: SERVER_Cipher,
		}

		/* test */
		if gConfig.Simfake == 1 {
			if module == EC20_TC1 {
				thisModule = &ModuleProduce{
					Type:   module,
					ModCmd: &modCmd_EC20,
					UrlVer: SERVER_Cipher_v1,
				}
			} else if module == EC20_TC3 {
				thisModule = &ModuleProduce{
					Type:   module,
					ModCmd: &modCmd_EC20,
					UrlVer: SERVER_Cipher_v3,
				}
			}
		}
		/* test end */
	}
}

/* ModuleTable end */
