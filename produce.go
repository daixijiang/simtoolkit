/* VSIM Serial Product Toolkit
 * Author: daixijiang@gmail.com (2019)
 */

package main

import (
	"encoding/json"
	"time"
	"vlog"
)

func (mp *ModuleProduce) GoProduce(portid int) int {
	var result string

	//prepare
	ret := mp.PreProduce(portid)
	if ret != 0 {
		vlog.Info("######## Port[%d] => failed(%d) to prepare ...", portid, ret)
		return ret
	}

	//produce
	ret = mp.getDevInfo(portid)
	if ret != 0 {
		vlog.Info("######## Port[%d] => failed(%d) to getDevInfo ...", portid, ret)
		return ret
	}

	ret = mp.getProToken(portid)
	if ret != 0 {
		vlog.Info("######## Port[%d] => failed(%d) to getProToken ...", portid, ret)
		return ret
	}

	ret = mp.getServInfo(portid)
	if ret != 0 {
		vlog.Info("######## Port[%d] => failed(%d) to getServInfo ...", portid, ret)
		return ret
	}

	ret = mp.getVsimDe(portid)
	if ret != 0 {
		vlog.Info("######## Port[%d] => failed(%d) to getVsimDe ...", portid, ret)
		return ret
	}

	ret = mp.setVsimData(portid)
	if ret != 0 {
		vlog.Info("######## Port[%d] => failed(%d) to setVsimData ...", portid, ret)
		return ret
	}

	//hot-reset module
	if gConfig.Produce.Hot_reset_timeout > 0 {
		ret = mp.DoComCMD(Module_CMD1_RESET0, portid, &result)
		vlog.Info("Port[%d] reset0 module: %s", portid, result)
		ret = mp.DoComCMD(Module_CMD1_RESET1, portid, &result)
		vlog.Info("Port[%d] reset1 module: %s", portid, result)
		//hot-reset wait
		time.Sleep(time.Duration(gConfig.Produce.Hot_reset_timeout) * time.Second)
	}

	//check
	ret = mp.GoCheck(portid)
	if ret != 0 {
		vlog.Info("######## Port[%d] => failed(%d) to checkdo ...", portid, ret)
		return ret
	}

	serial_port[portid].port_status = PORT_STATUS_OPEN
	return 0
}

func (mp *ModuleProduce) PreProduce(portid int) int {
	var result string

	vlog.Info("Port[%d] => prepare device produce ...", portid)

	//check imei
	ret := mp.DoComCMD(Module_CMD2_IMEI, portid, &result)
	if ret <= 0 {
		vlog.Info("######## Port[%d] get imei error: %s", portid, result)
		return -1
	}
	vlog.Info("Port[%d] get imei: %s", portid, result)

	//check system version
	ret = mp.DoComCMD(Module_CMD1_SYSVER, portid, &result)
	vlog.Info("Port[%d] get system version: %s", portid, result)

	//set softsimmode
	ret = mp.DoComCMD(Module_CMD1_SOFTMODE, portid, &result)
	vlog.Info("Port[%d] set softsimmode: %s", portid, result)

	//cold-reset module
	if gConfig.Produce.Cold_reset_timeout > 0 {
		ret = mp.DoComCMD(Module_CMD1_RESET2, portid, &result)
		vlog.Info("Port[%d] reset2 module: %s", portid, result)
		//cold-reset wait
		time.Sleep(time.Duration(gConfig.Produce.Cold_reset_timeout) * time.Second)
	}

	if mp.Type >= EC20 && mp.Type <= EC20_TC3 {
		//set network
		ret = mp.DoComCMD(Module_PRE1_SET_NETWORK0, portid, &result)
		vlog.Info("Port[%d] set network0: %s", portid, result)

		ret = mp.DoComCMD(Module_PRE1_SET_NETWORK1, portid, &result)
		vlog.Info("Port[%d] set network1: %s", portid, result)

		//set server url
		ret = mp.DoComCMD(Module_CMD1_SET_SERVURL, portid, &result)
		vlog.Info("Port[%d] set server url: %s", portid, result)

		if mp.Type == EC20_AUTO {
			//set autoswitch-on
			ret = mp.DoComCMD(Module_CMD1_AUTOSWITCH_ON, portid, &result)
			vlog.Info("Port[%d] set autoswitch-on: %s", portid, result)
		} else {
			//set autoswitch-off
			ret = mp.DoComCMD(Module_CMD1_AUTOSWITCH_OFF, portid, &result)
			vlog.Info("Port[%d] set autoswitch-off: %s", portid, result)
		}

		//backup config
		ret = mp.DoComCMD(Module_CMD1_BACKUP_CONFIG, portid, &result)
		vlog.Info("Port[%d] backup config: %s", portid, result)
	}

	return 0
}

func (mp *ModuleProduce) GoCheck(portid int) int {
	var result string

	vlog.Info("Port[%d] => check device produce ...", portid)
	time.Sleep(time.Duration(gConfig.Produce.Common_timeout) * time.Second)

	//check ccid "AT+CCID"
	mp.DoComCMD(Module_CMD3_CCID, portid, &result)
	vlog.Info("Port[%d] get ccid: %s", portid, result)

	//check cimi "AT+CIMI"
	mp.DoComCMD(Module_CMD3_CIMI, portid, &result)
	vlog.Info("Port[%d] get cimi: %s", portid, result)

	//check creg "AT+CREG?"
	mp.DoComCMD(Module_CMD3_CREG, portid, &result)
	vlog.Info("Port[%d] get creg: %s", portid, result)

	//check creg "AT+CEREG?"
	mp.DoComCMD(Module_CMD3_CEREG, portid, &result)
	vlog.Info("Port[%d] get cereg: %s", portid, result)

	//check creg "AT+COPS?"
	mp.DoComCMD(Module_CMD3_COPS, portid, &result)
	vlog.Info("Port[%d] get cops: %s", portid, result)

	if mp.Type >= EC20 && mp.Type <= EC20_TC3 {
		//check multi-oper, switch to tel
		mp.DoComCMD(Module_CMD3_SWITCH_TEL, portid, &result)
		vlog.Info("Port[%d] set switch tel: %s", portid, result)
		//creg wait
		time.Sleep(time.Duration(gConfig.Produce.Creg_timeout) * time.Second)

		mp.DoComCMD(Module_CMD3_CCID, portid, &result)
		vlog.Info("Port[%d] get ccid[tel]: %s", portid, result)

		//check multi-oper, switch to uni
		mp.DoComCMD(Module_CMD3_SWITCH_CU, portid, &result)
		vlog.Info("Port[%d] set switch uni: %s", portid, result)
		//creg wait
		time.Sleep(time.Duration(gConfig.Produce.Creg_timeout) * time.Second)

		mp.DoComCMD(Module_CMD3_CCID, portid, &result)
		vlog.Info("Port[%d] get ccid[uni]: %s", portid, result)

		//check multi-oper, switch to cmcc
		mp.DoComCMD(Module_CMD3_SWITCH_CM, portid, &result)
		vlog.Info("Port[%d] set switch cmcc: %s", portid, result)
		//creg wait
		time.Sleep(time.Duration(gConfig.Produce.Creg_timeout) * time.Second)

		mp.DoComCMD(Module_CMD3_CCID, portid, &result)
		vlog.Info("Port[%d] get ccid[cmcc]: %s", portid, result)
	}

	return 0
}

func (mp *ModuleProduce) getDevInfo(portid int) int {
	vlog.Info("Port[%d] p(1.0)=> get device info ...", portid)
	time.Sleep(time.Duration(gConfig.Produce.Common_timeout) * time.Second)

	mp.DoComCMD(Module_CMD2_IMEI, portid, &serial_port[portid].devInfo.sim_src.Imei)
	mp.DoComCMD(Module_CMD2_CHIPID, portid, &serial_port[portid].devInfo.sim_src.ChipID)
	mp.DoComCMD(Module_CMD2_COSVER, portid, &serial_port[portid].devInfo.version)

	return 0
}

func (mp *ModuleProduce) getProToken(portid int) int {
	vlog.Info("Port[%d] p(2.0)=> get token info ...", portid)

	for index := 0; index < OPER_MAX; index++ {
		serial_port[portid].devInfo.token[index] = getToken(serial_port[portid].devInfo.sim_src.Imei, index)
	}
	return 0
}

func (mp *ModuleProduce) getServInfo(portid int) int {
	var ret int
	vlog.Info("Port[%d] p(3.0)=> get server info ...", portid)
	time.Sleep(time.Duration(gConfig.Produce.Common_timeout) * time.Second)

	ver := mp.UrlVer
	if ver == SERVER_PLAIN_v0 {
		ret = mp.getServInfo_pv1(portid)
	} else if ver == SERVER_Cipher {
		ret = mp.getServInfo_cv1(portid, SERVER_Cipher)
	} else if ver == SERVER_Cipher_v1 {
		ret = mp.getServInfo_cv1(portid, SERVER_Cipher_v1)
	} else if ver == SERVER_Cipher_v3 {
		ret = mp.getServInfo_cv3(portid, SERVER_Cipher_v3)
	}

	return ret
}

func (mp *ModuleProduce) getServInfo_pv1(portid int) int {
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
		if gConfig.Simfake == 1 {
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

	if ret != 200 {
		return 403
	}

	return 0
}

func (mp *ModuleProduce) getServInfo_cv1(portid int, version int) int {
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
	if gConfig.Simfake == 1 {
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

	if res_data.Status == 200 {
		serial_port[portid].devInfo.servde = res_data.De
	} else {
		return res_data.Status
	}

	return 0
}

func (mp *ModuleProduce) getServInfo_cv3(portid int, version int) int {
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
	if gConfig.Simfake == 1 {
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

	if res_data.Status == 200 {
		serial_port[portid].devInfo.servde = res_data.De
	} else {
		return res_data.Status
	}

	return 0
}

func (mp *ModuleProduce) getVsimDe(portid int) int {
	var ret int
	vlog.Info("Port[%d] p(4.0)=> get vsim de ...", portid)

	ver := mp.UrlVer
	if ver == SERVER_PLAIN_v0 {
		ret = mp.getVsimDe_pv(portid)
	} else {
		ret = mp.getVsimDe_cv(portid)
	}

	vlog.Info("    crypto de192: %s", serial_port[portid].devInfo.sim_ens.EncData192)
	vlog.Info("    crypto de64: %s", serial_port[portid].devInfo.sim_ens.EncData64)

	return ret
}

func (mp *ModuleProduce) getVsimDe_pv(portid int) int {
	return Lib_vsim_encrypt(serial_port[portid].devInfo.sim_src, &serial_port[portid].devInfo.sim_ens)
}

func (mp *ModuleProduce) getVsimDe_cv(portid int) int {
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

func (mp *ModuleProduce) setVsimData(portid int) int {
	var result string
	vlog.Info("Port[%d] p(5.0)=> do producing on ...", portid)

	ret := 0
	if serial_port[portid].devInfo.sim_ens.EncData192 != "" {
		ret = mp.DoComCMD(Module_CMD2_SIM192, portid, &result)
		vlog.Info("    set de192 %s", result)
	}

	if serial_port[portid].devInfo.sim_ens.EncData64 != "" {
		ret = mp.DoComCMD(Module_CMD2_SIM64, portid, &result)
		vlog.Info("    set de64 %s", result)
	}

	if ret == 0 {
		vlog.Info("    Port[%d]: failed to set sim data on", portid)
		return 170
	}

	vlog.Info("Port[%d] p(5.1)=> do producing ok!", portid)
	return 0
}
