/* VSIM Serial Product Toolkit
 * Author: daixijiang@gmail.com (2019)
 */

package main

import (
	"github.com/tarm/serial"
	"vlog"
)

/* ModuleTable */
const (
	// prepare
	Module_CMD1_SYSVER         = 10
	Module_CMD1_SOFTMODE       = 11
	Module_CMD1_RESET0         = 12
	Module_CMD1_RESET1         = 13
	Module_PRE1_SET_NETWORK0   = 14
	Module_PRE1_SET_NETWORK1   = 15
	Module_CMD1_SET_SERVURL    = 16
	Module_CMD1_AUTOSWITCH_ON  = 17
	Module_CMD1_AUTOSWITCH_OFF = 18
	Module_CMD1_BACKUP_CONFIG  = 19

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
	EC20_PT
	EC20_CT1
	EC20_CT3
)

const MODULE_TEST = false

type cmdHandler func(cmdid int, cmdstr string, portid int, s *serial.Port, reply *string) int

type ModCmdTable struct {
	CmdID   int
	CmdStr  string
	CmdFunc cmdHandler
}

type ModuleProduce struct {
	Type     Module_cfg
	ModCmd   *[Module_TAB_AT_CMD_MAX]ModCmdTable
	UrlVer   int
	TestFlag int
}

var modCmd_EC20 *[Module_TAB_AT_CMD_MAX]ModCmdTable
var modCmd_SIM800C *[Module_TAB_AT_CMD_MAX]ModCmdTable
var thisModule *ModuleProduce

func (mp *ModuleProduce) DoComCMD(cmdid int, portid int, result *string) int {
	if cmdid > Module_TAB_AT_CMD_MAX {
		return -1
	}

	funcdo := (mp.ModCmd)[cmdid].CmdFunc
	if funcdo == nil {
		return 0
	}

	if (mp.ModCmd)[cmdid].CmdStr == "" {
		return 0
	}

	vlog.Debug("DoComCMD: cmdid %d, cmd %s, func %p", cmdid, (mp.ModCmd)[cmdid].CmdStr, funcdo)
	return funcdo(cmdid, (mp.ModCmd)[cmdid].CmdStr, portid, serial_port[portid].comPort, result)
}

func module_init() {
	modCmd_EC20 = module_ec20_init()
	modCmd_SIM800C = module_sim800c_init()
	// add list of module init

	// default is sim800c
	thisModule = &ModuleProduce{
		Type:     SIM800C,
		ModCmd:   modCmd_SIM800C,
		UrlVer:   SERVER_Cipher,
		TestFlag: 0,
	}

}

func module_reinit(module Module_cfg) {
	if module == SIM800C {
		thisModule = &ModuleProduce{
			Type:     module,
			ModCmd:   modCmd_SIM800C,
			UrlVer:   SERVER_Cipher,
			TestFlag: 0,
		}
	} else if module == EC20 {
		thisModule = &ModuleProduce{
			Type:     module,
			ModCmd:   modCmd_EC20,
			UrlVer:   SERVER_PLAIN_v0,
			TestFlag: 0,
		}
	} else if module == EC20_AUTO {
		thisModule = &ModuleProduce{
			Type:     module,
			ModCmd:   modCmd_EC20,
			UrlVer:   SERVER_PLAIN_v0,
			TestFlag: 0,
		}
	} else {
		// default is sim800c
		thisModule = &ModuleProduce{
			Type:     module,
			ModCmd:   modCmd_SIM800C,
			UrlVer:   SERVER_Cipher,
			TestFlag: 0,
		}

		/* test */
		if MODULE_TEST {
			if module == EC20_PT {
				thisModule = &ModuleProduce{
					Type:     module,
					ModCmd:   modCmd_EC20,
					UrlVer:   SERVER_PLAIN_v0,
					TestFlag: 1,
				}
			} else if module == EC20_CT1 {
				thisModule = &ModuleProduce{
					Type:     module,
					ModCmd:   modCmd_EC20,
					UrlVer:   SERVER_Cipher_v1,
					TestFlag: 1,
				}
			} else if module == EC20_CT3 {
				thisModule = &ModuleProduce{
					Type:     module,
					ModCmd:   modCmd_EC20,
					UrlVer:   SERVER_Cipher_v3,
					TestFlag: 1,
				}
			}
		}
		/* test end */
	}
}

/* ModuleTable end */
