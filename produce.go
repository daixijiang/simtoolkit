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
		vlog.Info("######## Port[%d] => failed(%d) to getDevInfo ...", portid, ret)
		return ret
	}

	ret = getProToken(portid)
	if ret != 0 {
		vlog.Info("######## Port[%d] => failed(%d) to getProToken ...", portid, ret)
		return ret
	}

	ret = getServInfo(portid)
	if ret != 0 {
		vlog.Info("######## Port[%d] => failed(%d) to getServInfo ...", portid, ret)
		return ret
	}

	ret = getVsimDe(portid)
	if ret != 0 {
		vlog.Info("######## Port[%d] => failed(%d) to getVsimDe ...", portid, ret)
		return ret
	}

	ret = setVsimData(portid)
	if ret != 0 {
		vlog.Info("######## Port[%d] => failed(%d) to setVsimData ...", portid, ret)
		return ret
	}

	serial_port[portid].port_status = PORT_STATUS_OPEN
	return 0
}

func serialCheckDo(portid int) int {
	var result string
	vlog.Info("Port[%d] => check device produce ...", portid)

	//check1 ccid "AT+CCID"
	if (*myProduce.Mod)[Module_CMD3_CCID].CmdFunc != nil {
		(*myProduce.Mod)[Module_CMD3_CCID].CmdFunc(
			Module_CMD3_CCID, portid,
			serial_port[portid].comPort,
			&result)
	}

	//TODO, check iccid of cmcc/uni/tel
	cur_iccid := serial_port[portid].devInfo.sim_src.VsimData[OPER_CN_MOBILE].Iccid
	if result != cur_iccid {
		vlog.Info("!!!!check ccid err: %s %s", result, cur_iccid)
	} else {
		vlog.Info("    check ccid ok: %s", result)
	}

	//TODO, check2 creg "AT+CREG?"
	if (*myProduce.Mod)[Module_CMD3_CREG].CmdFunc != nil {
		(*myProduce.Mod)[Module_CMD3_CREG].CmdFunc(
			Module_CMD3_CREG, portid,
			serial_port[portid].comPort,
			&result)
	}
	vlog.Info("    get creg: %s", result)
	return 0
}

func getDevInfo(portid int) int {
	vlog.Info("Port[%d] p(1.0)=> get device info ...", portid)
	time.Sleep(time.Duration(1) * time.Second)

	if (*myProduce.Mod)[Module_CMD2_IMEI].CmdFunc != nil {
		(*myProduce.Mod)[Module_CMD2_IMEI].CmdFunc(
			Module_CMD2_IMEI, portid,
			serial_port[portid].comPort,
			&serial_port[portid].devInfo.sim_src.Imei)
	}

	if (*myProduce.Mod)[Module_CMD2_CHIPID].CmdFunc != nil {
		(*myProduce.Mod)[Module_CMD2_CHIPID].CmdFunc(
			Module_CMD2_CHIPID, portid,
			serial_port[portid].comPort,
			&serial_port[portid].devInfo.sim_src.ChipID)
	}

	if (*myProduce.Mod)[Module_CMD2_VER].CmdFunc != nil {
		(*myProduce.Mod)[Module_CMD2_VER].CmdFunc(
			Module_CMD2_VER, portid,
			serial_port[portid].comPort,
			&serial_port[portid].devInfo.version)
	}

	return 0
}

func getProToken(portid int) int {
	vlog.Info("Port[%d] p(2.0)=> get token info ...", portid)

	for index := 0; index < OPER_MAX; index++ {
		serial_port[portid].devInfo.token[index] = getToken(serial_port[portid].devInfo.sim_src.Imei, index)
	}
	return 0
}

func getServInfo(portid int) int {
	var ret int
	vlog.Info("Port[%d] p(3.0)=> get server info ...", portid)
	time.Sleep(time.Duration(1) * time.Second)

	ver := myProduce.UrlVer
	if ver == SERVER_PLAIN_v0 {
		ret = getServInfo_pv1(portid)
	} else if ver == SERVER_Cipher {
		ret = getServInfo_cv1(portid, SERVER_Cipher)
	} else if ver == SERVER_Cipher_v1 {
		ret = getServInfo_cv1(portid, SERVER_Cipher_v1)
	} else if ver == SERVER_Cipher_v3 {
		ret = getServInfo_cv3(portid, SERVER_Cipher_v3)
	}

	serial_port[portid].strInfo = fmt.Sprintf("%s", "S")
	return ret
}

func getServInfo_pv1(portid int) int {
	var res []byte
	var req_data devReqPlainData
	var res_data devResPlainData

	ret := 403
	for index := 0; index < OPER_MAX; index++ {
		res_data = devResPlainData{}
		req_data = devReqPlainData{
			Imei:  serial_port[portid].devInfo.sim_src.Imei,
			Token: serial_port[portid].devInfo.token[index],
		}

		/* test */
		if myProduce.TestFlag == 1 {
			if index == OPER_CN_MOBILE {
				req_data.Imei = "867732034973305"
			} else if index == OPER_CN_TELECOM {
				req_data.Imei = "863412049788253"
			}
		}
		/* test end */

		if req_data.Token == "" {
			continue
		}

		reqSimServer(SERVER_PLAIN_v0, req_data, &res)
		err := json.Unmarshal(res, &res_data)
		if checkerr(err, 3, "Json parse server data") != 0 {
			continue
		}

		vlog.Info("    Get siminfo[%d](%d): %+v", index, res_data.Status, res_data)
		if res_data.Status == 200 {
			serial_port[portid].devInfo.sim_src.VsimData[index].Iccid = res_data.Iccid
			serial_port[portid].devInfo.sim_src.VsimData[index].Ki = res_data.Ki
			serial_port[portid].devInfo.sim_src.VsimData[index].Opc = res_data.Opc

			if index == OPER_CN_TELECOM {
				serial_port[portid].devInfo.sim_src.VsimData[index].Imsi = res_data.ImsiLTE
				serial_port[portid].devInfo.sim_src.CdmaData.Imsi_m = res_data.ImsiM
				serial_port[portid].devInfo.sim_src.CdmaData.Uim_id = res_data.Uimid
				serial_port[portid].devInfo.sim_src.CdmaData.Hrdupp = res_data.Hrpdupp
			} else {
				serial_port[portid].devInfo.sim_src.VsimData[index].Imsi = res_data.Imsi
			}

			ret = res_data.Status
		}
	}

	serial_port[portid].strInfo = fmt.Sprintf("%s", "S")

	if ret != 200 {
		return 403
	}

	return 0
}

func getServInfo_cv1(portid int, version int) int {
	var res []byte
	var res_data devResCipherData

	req_data := devReqCipherDataV1{
		Ver:    serial_port[portid].devInfo.version,
		Imei:   serial_port[portid].devInfo.sim_src.Imei,
		Chipid: serial_port[portid].devInfo.sim_src.ChipID,
		Token:  serial_port[portid].devInfo.token[OPER_CN_MOBILE],
		//TODO? only Token of CMCC
	}

	/* test */
	if myProduce.TestFlag == 1 {
		req_data.Ver = "CosVer_1.1.4"
		req_data.Imei = "867732034973305"
		req_data.Chipid = "3934363531303236320A3A373B3C3A3B"
	}
	/* test end */

	reqSimServer(version, req_data, &res)
	err := json.Unmarshal(res, &res_data)
	if checkerr(err, 3, "Json parse server data") != 0 {
		return 3
	}

	vlog.Info("    Get siminfo(%d): %+v", res_data.Status, res_data)
	serial_port[portid].strInfo = fmt.Sprintf("%s", "S")

	if res_data.Status == 200 {
		serial_port[portid].devInfo.servde = res_data.De
	} else {
		return res_data.Status
	}

	return 0
}

func getServInfo_cv3(portid int, version int) int {
	var res []byte
	var res_data devResCipherData

	req_data := devReqCipherDataV3{
		Ver:    serial_port[portid].devInfo.version,
		Imei:   serial_port[portid].devInfo.sim_src.Imei,
		Chipid: serial_port[portid].devInfo.sim_src.ChipID,
	}

	token_json := make([]string, 0)
	for index := 0; index < OPER_MAX; index++ {
		token_json = append(token_json, serial_port[portid].devInfo.token[index])
	}
	req_data.Token = token_json

	/* test */
	if myProduce.TestFlag == 1 {
		req_data.Ver = "CosVer_1.1.4"
		req_data.Imei = "867732034973305"
		req_data.Chipid = "3934363531303236320A3A373B3C3A3B"
	}

	reqSimServer(version, req_data, &res)
	err := json.Unmarshal(res, &res_data)
	if checkerr(err, 3, "Json parse server data") != 0 {
		return 3
	}

	vlog.Info("    Get siminfo(%d): %+v", res_data.Status, res_data)
	serial_port[portid].strInfo = fmt.Sprintf("%s", "S")

	if res_data.Status == 200 {
		serial_port[portid].devInfo.servde = res_data.De
	} else {
		return res_data.Status
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
		ret = getVsimDe_cv(portid)
	}

	vlog.Info("    crypto de192: %s", serial_port[portid].devInfo.sim_ens.EncData192)
	vlog.Info("    crypto de64: %s", serial_port[portid].devInfo.sim_ens.EncData64)

	return ret
}

func getVsimDe_pv(portid int) int {
	return Lib_vsim_encrypt(serial_port[portid].devInfo.sim_src, &serial_port[portid].devInfo.sim_ens)
}

func getVsimDe_cv(portid int) int {
	encdata := []byte(serial_port[portid].devInfo.servde)
	delen := len(serial_port[portid].devInfo.servde)
	if delen >= ENC_DATA_192 {
		serial_port[portid].devInfo.sim_ens.EncData192 = string(encdata[0 : ENC_DATA_192-1])
		serial_port[portid].devInfo.sim_ens.EncData64 = string(encdata[ENC_DATA_192:delen])
	} else {
		serial_port[portid].devInfo.sim_ens.EncData64 = serial_port[portid].devInfo.servde
	}

	return 0
}

func setVsimData(portid int) int {
	var result string
	vlog.Info("Port[%d] p(5.0)=> do producing on ...", portid)

	ret := 0

	if ((*myProduce.Mod)[Module_CMD2_SIM192].CmdFunc != nil) &&
		(serial_port[portid].devInfo.sim_ens.EncData192 != "") {
		ret = (*myProduce.Mod)[Module_CMD2_SIM192].CmdFunc(
			Module_CMD2_SIM192, portid,
			serial_port[portid].comPort,
			&result)

		vlog.Info("    set de192 %s", result)
	}

	if ((*myProduce.Mod)[Module_CMD2_SIM64].CmdFunc != nil) &&
		(serial_port[portid].devInfo.sim_ens.EncData64 != "") {
		ret = (*myProduce.Mod)[Module_CMD2_SIM64].CmdFunc(
			Module_CMD2_SIM64, portid,
			serial_port[portid].comPort,
			&result)

		vlog.Info("    set de64 %s", result)
	}

	if ret == 0 {
		vlog.Info("    Port[%d]: failed to set sim data on", portid)
		return 170
	}

	vlog.Info("Port[%d] p(5.1)=> do producing ok!", portid)
	return 0
}
