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

var myTokenStat [OPER_MAX]simTokenStat

func token_init() {
	for index := 0; index < OPER_MAX; index++ {
		myTokenStat[index].total = 0
		myTokenStat[index].current = 0
	}

	loadTokenCfg(gConfig.Token.Cmcc_file, OPER_CN_MOBILE)
	loadTokenCfg(gConfig.Token.Uni_file, OPER_CN_UNICOM)
	loadTokenCfg(gConfig.Token.Tel_file, OPER_CN_TELECOM)
}

func loadTokenCfg(filename string, oper int) {
	myTokenStat[oper].total = 0
	myTokenStat[oper].current = 0

	readTokenfile(filename, oper)
	vlog.Info("Load token[%d] total %d", oper, myTokenStat[oper].total)
}

func readTokenfile(fileName string, oper int) int {
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
			vlog.Debug("%s", line)
			ldata := []byte(line)
			if ldata[0] != '#' {
				vlog.Debug("%s", ldata[0])
				addToken(line, oper)
			}
		}
	}

	return 0
}

func addToken(token string, oper int) {
	index := myTokenStat[oper].total
	if index < TOKEN_MAX {
		myTokenStat[oper].sToken[index].Token = token
		myTokenStat[oper].sToken[index].Useflag = 0
		myTokenStat[oper].sToken[index].Imei = ""
		myTokenStat[oper].total++
	}
}

func getToken(imei string, oper int) (token string) {
	// general token
	if myTokenStat[oper].total == 1 {
		vlog.Info("    get token[%d][%s] for imei[%s]",
			oper, myTokenStat[oper].sToken[0].Token, imei)
		return myTokenStat[oper].sToken[0].Token
	}

	// fast lookup
	for index := myTokenStat[oper].current; index < myTokenStat[oper].total; index++ {
		if myTokenStat[oper].sToken[index].Token != "" {
			if myTokenStat[oper].sToken[index].Imei == imei {
				vlog.Info("    get token[%d][%s] for imei[%s]",
					oper, myTokenStat[oper].sToken[index].Token, imei)
				return myTokenStat[oper].sToken[index].Token
			} else if myTokenStat[oper].sToken[index].Useflag == 0 {
				vlog.Info("    get token[%d][%s] for imei[%s]",
					oper, myTokenStat[oper].sToken[index].Token, imei)
				myTokenStat[oper].current = index + 1
				return myTokenStat[oper].sToken[index].Token
			}
		}
	}

	//again lookup
	for index := 0; index < myTokenStat[oper].total; index++ {
		if myTokenStat[oper].sToken[index].Token != "" {
			if myTokenStat[oper].sToken[index].Imei == imei {
				vlog.Info("    get token[%d][%s] for imei[%s]",
					oper, myTokenStat[oper].sToken[index].Token, imei)
				return myTokenStat[oper].sToken[index].Token
			} else if myTokenStat[oper].sToken[index].Useflag == 0 {
				vlog.Info("    get token[%d][%s] for imei[%s]",
					oper, myTokenStat[oper].sToken[index].Token, imei)
				myTokenStat[oper].current = index + 1
				return myTokenStat[oper].sToken[index].Token
			}
		}
	}

	vlog.Info("    get no token[%d] for imei[%s]", oper, imei)
	return ""
}
