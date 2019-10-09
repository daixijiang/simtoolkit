/* VSIM Serial Product Toolkit
 * Author: daixijiang@gmail.com (2019)
 */

package main

import (
	"vlog"

	"github.com/aarzilli/nucular"
	nstyle "github.com/aarzilli/nucular/style"
)

const szTitle string = "VSIM Serial Product Toolkit"
const szBanner string = "Welcome use vsim toolkit!"
const szVersion string = "V2019.09.09"

const scaling = 1.3

var mytheme nstyle.Theme = nstyle.DarkTheme

func main() {
	log_init()
	module_init()
	token_init()
	serial_util_init()

	pg := newportGroup()
	wnd := nucular.NewMasterWindow(0, szTitle, pg.showUI)
	wnd.SetStyle(nstyle.FromTheme(mytheme, scaling))
	wnd.Main()
}

func logInit() {
	vlog.InitLog("file", "vsimtoolkit.log", "Debug", 7)
}
