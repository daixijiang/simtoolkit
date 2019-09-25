/* VSIM Serial Product Toolkit
 * Author: daixijiang@gmail.com (2019)
 *
 * check: curl -X POST -k -i 'https://rdp.showmac.cn/api/v1/profile/get' --data '{"imei":"868575021892064","chipid":"20171026050559A399032A3416886391","token":"KD1MQ-BWJGR-8ZB29-9D59J"}'
 */

package main

import (
	"crypto/tls"
	"encoding/json"
	"vlog"
	"bytes"
	"io/ioutil"
	"net/http"
)

type devReqData struct {
	Imei	string  `json:"imei"`
	Chipid  string  `json:"chipid"`
	Token	string  `json:"token"`
}

type devResData struct {
	Status	int     `json:"status"`
	Imei	string  `json:"imei"`
	Iccid	string  `json:"iccid"`
	De	string  `json:"de"`
}

func check(err error, code int) int {
	if err != nil {
        	vlog.Error("Json parse error %d", code)
        	return code
	}

	return 0
}

const serverUrl string = "https://rdp.showmac.cn/api/v1/profile/get"

func reqSimServer(req devReqData, res devResData) {
	var resp *http.Response
	var message devResData
	var data []byte

	serverData, err := json.Marshal(req)
	if (check(err, 1) != 0) { return }

	vlog.Info("%s", string(serverData))
	req_new := bytes.NewBuffer(serverData)

	tr := &http.Transport{
		TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
		DisableCompression: true,
	}

	client := &http.Client{Transport: tr}
	request, _ := http.NewRequest("POST", serverUrl, req_new)
	request.Header.Set("Content-type", "application/json")


	resp, err = client.Do(request)
	if (check(err, 2) != 0) { return }

	if data, err = ioutil.ReadAll(resp.Body); err == nil {
		vlog.Info("%s", data)
	}

	err = json.Unmarshal(data, &message)
	if (check(err, 3) != 0) { return }

	vlog.Info("%+v", message)
	res.Status = message.Status
	res.Imei   = message.Imei
	res.Iccid  = message.Iccid
	res.De     = message.De
}

/*
func main() {
	var dev_data devReqData
	var sim_data devResData

	dev_data.Imei = "868575021892064"
	dev_data.Chipid = "20171026050559A399032A3416886391"
	dev_data.Token = "KD1MQ-BWJGR-8ZB29-9D59J"

	reqSimServer(dev_data, sim_data)
}
*/
