/* VSIM Serial Product Toolkit
 * Author: daixijiang@gmail.com (2019)
 */

package main

import (
	"github.com/tarm/serial"
)

/* ModuleTable */
const (
	// prepare
	Module_CMD1_VER             = 10
	Module_CMD1_SOFTMODE        = 11
	Module_CMD1_REBOOT          = 12
	Module_PRE1_CMD_NETWORK_SET = 13
	Module_CMD1_URL_SET         = 14
	Module_CMD1_AUTOSWITCH_ON   = 15
	Module_CMD1_AUTOSWITCH_OFF  = 16
	Module_CMD1_CFG_BACKUP      = 17

	// produce
	Module_CMD2_IMEI   = 20
	Module_CMD2_VER    = 21
	Module_CMD2_CHIPID = 22
	Module_CMD2_SIM192 = 23
	Module_CMD2_SIM64  = 24

	// check
	Module_CMD3_CCID       = 30
	Module_CMD3_CIMI       = 31
	Module_CMD3_CREG       = 32
	Module_CMD3_CEREG      = 33
	Module_CMD3_COPS       = 34
	Module_CMD3_SWITCH_CM  = 35
	Module_CMD3_SWITCH_CU  = 36
	Module_CMD3_SWITCH_TEL = 37

	Module_TAB_AT_CMD_MAX = 40
)

type Module_cfg int

const (
	EC20 = Module_cfg(iota)
	SIM800C
	EC20_PT
	EC20_CT1
	EC20_CT3
)

const MODULE_TEST = true

type cmdHandler func(cmdid int, portid int, s *serial.Port, reply *string) int

type ModuleTable struct {
	CmdID   int
	CmdStr  string
	CmdFunc cmdHandler
}

type ProModule struct {
	Mod      *[Module_TAB_AT_CMD_MAX]ModuleTable
	UrlVer   int
	TestFlag int
}

var modEC20 [Module_TAB_AT_CMD_MAX]ModuleTable
var modSIM800C [Module_TAB_AT_CMD_MAX]ModuleTable
var myProduce *ProModule

func module_init() {
	module_ec20_init()
	module_sim800c_init()
	// add list of module init

	// default is sim800c
	myProduce = &ProModule{
		Mod:      &modSIM800C,
		UrlVer:   SERVER_Cipher,
		TestFlag: 0,
	}

}

func module_reinit(module Module_cfg) {
	if module == EC20 {
		myProduce = &ProModule{
			Mod:      &modEC20,
			UrlVer:   SERVER_PLAIN_v0,
			TestFlag: 0,
		}
	} else if module == SIM800C {
		myProduce = &ProModule{
			Mod:      &modSIM800C,
			UrlVer:   SERVER_Cipher,
			TestFlag: 0,
		}
	} else {
		// default is sim800c
		myProduce = &ProModule{
			Mod:      &modSIM800C,
			UrlVer:   SERVER_Cipher,
			TestFlag: 0,
		}

		/* test */
		if MODULE_TEST {
			if module == EC20_PT {
				myProduce = &ProModule{
					Mod:      &modEC20,
					UrlVer:   SERVER_PLAIN_v0,
					TestFlag: 1,
				}
			} else if module == EC20_CT1 {
				myProduce = &ProModule{
					Mod:      &modEC20,
					UrlVer:   SERVER_Cipher_v1,
					TestFlag: 1,
				}
			} else if module == EC20_CT3 {
				myProduce = &ProModule{
					Mod:      &modEC20,
					UrlVer:   SERVER_Cipher_v3,
					TestFlag: 1,
				}
			}
		}
		/* test end */
	}
}

/* ModuleTable end */
