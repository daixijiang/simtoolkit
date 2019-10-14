/* VSIM Serial Product Toolkit
 * Author: daixijiang@gmail.com (2019)
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

const (
	SERVER_PLAIN_v0  = 0
	SERVER_Cipher    = 1
	SERVER_Cipher_v1 = 2
	SERVER_Cipher_v3 = 3
	SERVER_MAX       = 4
)

var serverUrl [SERVER_MAX]string

type devReqCipherDataV1 struct {
	Ver    string `json:"version"`
	Imei   string `json:"imei"`
	Chipid string `json:"chipid"`
	Token  string `json:"token"`
}

type devReqCipherDataV3 struct {
	Ver    string   `json:"version"`
	Imei   string   `json:"imei"`
	Chipid string   `json:"chipid"`
	Token  []string `json:"token"`
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
	Imsi    string `json:"imsi"`
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

func server_init() {
	serverUrl[SERVER_PLAIN_v0] = gConfig.Server.Plain_url
	serverUrl[SERVER_Cipher] = gConfig.Server.Cipher_url
	serverUrl[SERVER_Cipher_v1] = gConfig.Server.Cipherv1_url
	serverUrl[SERVER_Cipher_v3] = gConfig.Server.Cipherv3_url
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
