/* VSIM Serial Product Toolkit
 * Author: daixijiang@gmail.com (2019)
 *
 * check: curl -X POST -k -i 'https://rdp.showmac.cn/api/v1/profile/get' --data '{"imei":"868575021892064","chipid":"20171026050559A399032A3416886391","token":"KD1MQ-BWJGR-8ZB29-9D59J"}'
 */

package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"vlog"
)

const SERVER_PLAIN_v0 = 0
const SERVER_Cipher_v1 = 1
const SERVER_Cipher_v2 = 2
const SERVER_Cipher_v3 = 3
const SERVER_Cipher_v4 = 4
const SERVER_MAX = 5

// plaintext
const serverPlainUrl string = "https://rdp.showmac.cn/api/v1/profile/clear/get"

// ciphertext
const serverCipherUrlV1 string = "https://rdp.showmac.cn/api/v1/profile/get"

var serverUrl = [SERVER_MAX]string{serverPlainUrl, serverCipherUrlV1}

type devReqData struct {
	Ver    string `json:"version"`
	Imei   string `json:"imei"`
	Chipid string `json:"chipid"`
	Token  string `json:"token"`
}

type devReqCipherData struct {
	Imei   string `json:"imei"`
	Chipid string `json:"chipid"`
	Token  string `json:"token"`
}

type devResCipherData struct {
	Status int    `json:"status"`
	Imei   string `json:"imei"`
	Iccid  string `json:"iccid"`
	De     string `json:"de"`
}

type devReqPlainData struct {
	Imei  string `json:"imei"`
	Token string `json:"token"`
}

type devResPlainData struct {
	Status  int    `json:"status"`
	Imei    string `json:"imei"`
	Iccid   string `json:"iccid"`
	Ki      string `json:"ki"`
	Opc     string `json:"opc"`
	ImsiLTE string `json:"imsiLTE"`
	ImsiM   string `json:"imsiM"`
	Uimid   string `json:"uimID"`
	Hrpdupp string `json:"hrpdupp"`
}

func checkerr(err error, code int, opername string) int {
	if err != nil {
		vlog.Error("    %s error %d", opername, code)
		return code
	}

	return 0
}

func reqSimServer(version int, req interface{}, resp_data *[]byte) int {
	var resp *http.Response

	serverData, err := json.Marshal(req)
	if checkerr(err, 1, "Json parse sim post-data") != 0 {
		return 1
	}

	vlog.Info("    request: %s", string(serverData))
	req_new := bytes.NewBuffer(serverData)

	tr := &http.Transport{
		TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
		DisableCompression: true,
	}

	client := &http.Client{Transport: tr}
	request, _ := http.NewRequest("POST", serverUrl[version], req_new)
	request.Header.Set("Content-type", "application/json")

	resp, err = client.Do(request)
	if checkerr(err, 2, "Json parse https request") != 0 {
		return 2
	}

	if *resp_data, err = ioutil.ReadAll(resp.Body); err == nil {
		vlog.Info("    response: %s", *resp_data)
	}

	if checkerr(err, 3, "Json parse https response") != 0 {
		return 3
	}

	return 0
}

func test_server_main() {
	var dev_data devReqPlainData
	var sim_data devResPlainData
	var res []byte

	dev_data.Imei = "863412049788253"
	dev_data.Token = "YR0NI-259CE-R3JI5-01DJN-ENY2Z"

	reqSimServer(SERVER_PLAIN_v0, dev_data, &res)
	err := json.Unmarshal(res, &sim_data)
	if checkerr(err, 3, "Json parse https response") != 0 {
		return
	}

	vlog.Info("%+v", sim_data)
}
