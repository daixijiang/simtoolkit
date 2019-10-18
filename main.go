/* VSIM Serial Product Toolkit
 * Author: daixijiang@gmail.com (2019)
 */

package main

import (
	"vlog"
)

const szTitle string = "VSIM Serial Product Toolkit"
const szBanner string = "Welcome use vsim toolkit!"
const szVersion string = "V2019.09.09"

func main() {
	config_init()
	log_init()

	vlog.Info("version %s", szVersion)

	module_init()
	token_init()
	server_init()
	serial_util_init()

	if gConfig.UIstyle == "walk" {
		newui := newWalkUI()
		newui.RunUI()
	} else if gConfig.UIstyle == "nucular" {
		newui := newNucularUI()
		newui.RunUI()
	} else {
		newui := newNucularUI()
		newui.RunUI()
	}
}
