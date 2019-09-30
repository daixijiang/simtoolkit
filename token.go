/* VSIM Serial Product Toolkit
 * Author: daixijiang@gmail.com (2019)
 */

package main

import (
	"bufio"
	"io"
	"os"
	"strings"
	"vlog"
)

const TOKEN_MAX = 100
const TOKEN_FILE_DAFAULT string = "./token.cfg"
const VLOG_FILE_DAFAULT string = "./vsimkit.log"
const VLOG_LEVEL_DAFAULT string = "info"
const VLOG_MAXDAY_DAFAULT = 7

type simToken struct {
	Token   string
	Imei    string
	Useflag int
}

type simTokenStat struct {
	total   int
	current int
	sToken  [TOKEN_MAX]simToken
}

var myTokenStat simTokenStat

func log_init() {
	vlog.InitLog("file", VLOG_FILE_DAFAULT, VLOG_LEVEL_DAFAULT, VLOG_MAXDAY_DAFAULT)
}

func token_init() {
	myTokenStat.total = 0
	myTokenStat.current = 0

	loadTokenCfg(TOKEN_FILE_DAFAULT)
}

func loadTokenCfg(filename string) {
	readTokenfile(filename)
	vlog.Info("Load token total %d", myTokenStat.total)
}

func readTokenfile(fileName string) int {
	vlog.Info("Load token file %s", fileName)
	file, err := os.OpenFile(fileName, os.O_RDWR, 0666)
	if err != nil {
		vlog.Error("Open file %s error: %s!", fileName, err)
		return -1
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		vlog.Error("Stat file %s error: %s!", fileName, err)
		return -2
	}

	var size = stat.Size()
	vlog.Info("file size=%d", size)

	buf := bufio.NewReader(file)
	for {
		line, err := buf.ReadString('\n')
		line = strings.TrimSpace(line)
		if err != nil {
			if err == io.EOF {
				vlog.Info("File read ok!")
				break
			} else {
				vlog.Error("Read file error: %s!", err)
				return 0
			}
		} else {
			//vlog.Info("%s", line)
			ldata := []byte(line)
			if ldata[0] != '#' {
				//vlog.Info("%s", ldata[0])
				addToken(line, myTokenStat.total)
			}
		}
	}

	return 0
}

func addToken(token string, index int) {
	if index < TOKEN_MAX {
		myTokenStat.sToken[index].Token = token
		myTokenStat.sToken[index].Useflag = 0
		myTokenStat.sToken[index].Imei = ""
		myTokenStat.total++
	}
}

func getToken(imei string) (token string) {
	// general token
	if myTokenStat.total == 1 {
		vlog.Info("    get token[%s] for imei[%s]",
			myTokenStat.sToken[0].Token, imei)
		return myTokenStat.sToken[0].Token
	}

	// fast lookup
	for index := myTokenStat.current; index < myTokenStat.total; index++ {
		if myTokenStat.sToken[index].Token != "" {
			if myTokenStat.sToken[index].Imei == imei {
				vlog.Info("    get token[%s] for imei[%s]",
					myTokenStat.sToken[index].Token, imei)
				return myTokenStat.sToken[index].Token
			} else if myTokenStat.sToken[index].Useflag == 0 {
				vlog.Info("    get token[%s] for imei[%s]",
					myTokenStat.sToken[index].Token, imei)
				myTokenStat.current = index + 1
				return myTokenStat.sToken[index].Token
			}
		}
	}

	//again lookup
	for index := 0; index < myTokenStat.total; index++ {
		if myTokenStat.sToken[index].Token != "" {
			if myTokenStat.sToken[index].Imei == imei {
				vlog.Info("    get token[%s] for imei[%s]",
					myTokenStat.sToken[index].Token, imei)
				return myTokenStat.sToken[index].Token
			} else if myTokenStat.sToken[index].Useflag == 0 {
				vlog.Info("    get token[%s] for imei[%s]",
					myTokenStat.sToken[index].Token, imei)
				myTokenStat.current = index + 1
				return myTokenStat.sToken[index].Token
			}
		}
	}

	vlog.Info("    get no token for imei[%s]", imei)
	return ""
}
