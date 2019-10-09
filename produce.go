/* VSIM Serial Product Toolkit
 * Author: daixijiang@gmail.com (2019)
 */

package main

import (
	"encoding/json"
	"fmt"
	"time"
	"vlog"
)

func serialProduce(portid int) int {
	ret := getDevInfo(portid)
	if ret != 0 {
		vlog.Info("!!!!Port[%d] => failed(%d) to getDevInfo ...", portid, ret)
		return ret
	}

	ret = getProToken(portid)
	if ret != 0 {
		vlog.Info("!!!!Port[%d] => failed(%d) to getProToken ...", portid, ret)
		return ret
	}

	ret = getServInfo(portid)
	if ret != 0 {
		vlog.Info("!!!!Port[%d] => failed(%d) to getServInfo ...", portid, ret)
		return ret
	}

	ret = getVsimDe(portid)
	if ret != 0 {
		vlog.Info("!!!!Port[%d] => failed(%d) to getVsimDe ...", portid, ret)
		return ret
	}

	ret = setVsimData(portid)
	if ret != 0 {
		vlog.Info("!!!!Port[%d] => failed(%d) to setVsimData ...", portid, ret)
		return ret
	}

	serial_port[portid].port_status = PORT_STATUS_OPEN
	return 0
}

func serialCheckDo(portid int) int {
	var result string
	vlog.Info("Port[%d] => check device produce ...", portid)

	//check1 ccid "AT+CCID"
	(*myProduce.Mod)[Module_CMD3_CCID].CmdFunc(
		Module_CMD3_CCID, portid,
		serial_port[portid].comPort,
		&result)

	if result != serial_port[portid].sim_pv1.Iccid {
		vlog.Info("!!!!check ccid err: %s %s", result, serial_port[portid].sim_pv1.Iccid)
	} else {
		vlog.Info("    check ccid ok: %s", result)
	}

	//TODO, check2 creg "AT+CREG?"
	(*myProduce.Mod)[Module_CMD3_CREG].CmdFunc(
		Module_CMD3_CREG, portid,
		serial_port[portid].comPort,
		&result)
	vlog.Info("    get creg: %s", result)
	return 0
}

func getDevInfo(portid int) int {
	vlog.Info("Port[%d] p(1.0)=> get device info ...", portid)
	time.Sleep(time.Duration(1) * time.Second)

	(*myProduce.Mod)[Module_CMD2_IMEI].CmdFunc(
		Module_CMD2_IMEI, portid,
		serial_port[portid].comPort,
		&serial_port[portid].dev_data.Imei)

	(*myProduce.Mod)[Module_CMD2_CHIPID].CmdFunc(
		Module_CMD2_CHIPID, portid,
		serial_port[portid].comPort,
		&serial_port[portid].dev_data.Chipid)

	(*myProduce.Mod)[Module_CMD2_VER].CmdFunc(
		Module_CMD2_VER, portid,
		serial_port[portid].comPort,
		&serial_port[portid].dev_data.Ver)

	return 0
}

func getProToken(portid int) int {
	vlog.Info("Port[%d] p(2.0)=> get token info ...", portid)
	serial_port[portid].dev_data.Token = getToken(serial_port[portid].dev_data.Imei)
	return 0
}

func getServInfo(portid int) int {
	var ret int
	vlog.Info("Port[%d] p(3.0)=> get server info ...", portid)
	time.Sleep(time.Duration(1) * time.Second)

	ver := myProduce.UrlVer
	if ver == SERVER_PLAIN_v0 {
		ret = getServInfo_pv1(portid)
	} else if ver == SERVER_Cipher_v1 {
		ret = getServInfo_cv1(portid, SERVER_Cipher_v1)
	} else if ver == SERVER_Cipher_v2 {
		ret = getServInfo_cv1(portid, SERVER_Cipher_v2)
	} else if ver == SERVER_Cipher_v3 {
		ret = getServInfo_cv1(portid, SERVER_Cipher_v3)
	} else if ver == SERVER_Cipher_v4 {
		ret = getServInfo_cv1(portid, SERVER_Cipher_v4)
	}

	serial_port[portid].strInfo = fmt.Sprintf("%s", "S")
	return ret
}

func getServInfo_pv1(portid int) int {
	var res []byte

	dev_data := devReqPlainData{
		Imei:  serial_port[portid].dev_data.Imei,
		Token: serial_port[portid].dev_data.Token,
	}

	/* test */
	if myProduce.TestFlag == 1 {
		dev_data.Imei = "863412049788253"
		dev_data.Token = "YR0NI-259CE-R3JI5-01DJN-ENY2Z"
	}

	reqSimServer(SERVER_PLAIN_v0, dev_data, &res)
	err := json.Unmarshal(res, &serial_port[portid].sim_pv1)
	if checkerr(err, 3, "Json parse server data") != 0 {
		return 3
	}

	vlog.Info("    Get siminfo(%d): %+v", serial_port[portid].sim_pv1.Status, serial_port[portid].sim_pv1)
	serial_port[portid].strInfo = fmt.Sprintf("%s", "S")

	if serial_port[portid].sim_pv1.Status != 200 {
		return serial_port[portid].sim_pv1.Status
	}

	return 0
}

func getServInfo_cv1(portid int, version int) int {
	var res []byte

	dev_data := devReqData{
		Ver:    serial_port[portid].dev_data.Ver,
		Imei:   serial_port[portid].dev_data.Imei,
		Chipid: serial_port[portid].dev_data.Chipid,
		Token:  serial_port[portid].dev_data.Token,
	}

	/* test */
	if myProduce.TestFlag == 1 {
		dev_data.Ver = "CosVer_1.1.4"
		dev_data.Imei = "867732034973305"
		dev_data.Chipid = "3934363531303236320A3A373B3C3A3B"
		dev_data.Token = "WPAFE-7O2T3-SPEX9-DUWBJ"
	}

	reqSimServer(version, dev_data, &res)
	err := json.Unmarshal(res, &serial_port[portid].sim_cv1)
	if checkerr(err, 3, "Json parse server data") != 0 {
		return 3
	}

	vlog.Info("    Get siminfo(%d): %+v", serial_port[portid].sim_cv1.Status, serial_port[portid].sim_cv1)
	serial_port[portid].strInfo = fmt.Sprintf("%s", "S")

	if serial_port[portid].sim_cv1.Status != 200 {
		return serial_port[portid].sim_cv1.Status
	}

	return 0
}

func getVsimDe(portid int) int {
	var ret int
	vlog.Info("Port[%d] p(4.0)=> get vsim de ...", portid)

	ver := myProduce.UrlVer
	if ver == SERVER_PLAIN_v0 {
		ret = getVsimDe_pv(portid)
	} else {
		ret = getVsimDe_cv1(portid)
	}

	vlog.Info("    crypto de192: %s", serial_port[portid].sim_ens.EncData192)
	vlog.Info("    crypto de64: %s", serial_port[portid].sim_ens.EncData64)

	return ret
}

func getVsimDe_pv(portid int) int {
	srcsim := SRC_SIM_DATA{
		Imei:   serial_port[portid].dev_data.Imei,
		ChipID: serial_port[portid].dev_data.Chipid,
		CdmaData: CDMA_DATA{
			Imsi_m: serial_port[portid].sim_pv1.ImsiM,
			Uim_id: serial_port[portid].sim_pv1.Uimid,
			Hrdupp: serial_port[portid].sim_pv1.Hrpdupp,
		},
	}

	srcsim.VsimData[OPER_CN_MOBILE] = SIM_DATA{
		Iccid: serial_port[portid].sim_pv1.Iccid,
		Imsi:  serial_port[portid].sim_pv1.ImsiLTE,
		Ki:    serial_port[portid].sim_pv1.Ki,
		Opc:   serial_port[portid].sim_pv1.Opc,
	}

	return Lib_vsim_encrypt(srcsim, &serial_port[portid].sim_ens)
}

func getVsimDe_cv1(portid int) int {
	encdata := []byte(serial_port[portid].sim_cv1.De)
	delen := len(serial_port[portid].sim_cv1.De)
	if delen > ENC_DATA_192 {
		serial_port[portid].sim_ens.EncData192 = string(encdata[0 : ENC_DATA_192-1])
		serial_port[portid].sim_ens.EncData64 = string(encdata[ENC_DATA_192:delen])
	} else {
		serial_port[portid].sim_ens.EncData192 = serial_port[portid].sim_cv1.De
	}

	return 0
}

func setVsimData(portid int) int {
	var result string
	vlog.Info("Port[%d] p(5.0)=> do producing on ...", portid)

	(*myProduce.Mod)[Module_CMD2_SIM192].CmdFunc(
		Module_CMD2_SIM192, portid,
		serial_port[portid].comPort,
		&result)
	vlog.Info("    set de192 %s", result)

	(*myProduce.Mod)[Module_CMD2_SIM64].CmdFunc(
		Module_CMD2_SIM64, portid,
		serial_port[portid].comPort,
		&result)
	vlog.Info("    set de64 %s", result)

	vlog.Info("Port[%d] p(5.1)=> do producing ok!", portid)
	return 0
}
