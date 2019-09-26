package main

/*
#include <string.h>
#include <stdlib.h>
#include <ctype.h>

#define	ICCID_STR_LEN					20
#define IMSI_STR_LEN					15
#define KI_STR_LEN						32
#define OPC_STR_LEN						32
#define IMSI_M_STR_LEN					15
#define UIMID_STR_LEN					8
#define HRDUPP_STR_LEN					25
#define	IMEI_STR_LEN					15
#define CHIPID_STR_LEN					32
#define ENC_DATA_192					192
#define ENC_DATA_64						64

enum operator{
	OPER_CN_MOBILE = 0,
	OPER_CN_UNICOM = 1,
	OPER_CN_TELECOM = 2,
	OPER_MAX,
} oper;

typedef struct {
	char iccid[ICCID_STR_LEN+1];
	char imsi[IMSI_STR_LEN+1];
	char ki[KI_STR_LEN+1];
	char opc[OPC_STR_LEN+1];
} SIM_DATA;

typedef struct {
	char imsi_m[IMSI_M_STR_LEN+1];
	char uim_id[UIMID_STR_LEN+1];
	char hrdupp[HRDUPP_STR_LEN+1];
} CDMA_DATA;

typedef struct {
	char imei[IMEI_STR_LEN+1];
	char chipID[CHIPID_STR_LEN+1];
	SIM_DATA vsimData[OPER_MAX];
	CDMA_DATA cdmaData;
} SRC_SIM_DATA;

typedef struct {
	char encData192[ENC_DATA_192+1];
	char encData64[ENC_DATA_64+1];
} ENC_SIM_DATA;

*/
import "C"

import (
	"fmt"
	"syscall"
	"unsafe"
)

const ICCID_STR_LEN = 20
const IMSI_STR_LEN = 15
const KI_STR_LEN = 32
const OPC_STR_LEN = 32
const IMSI_M_STR_LEN = 15
const UIMID_STR_LEN = 8
const HRDUPP_STR_LEN = 25
const IMEI_STR_LEN = 15
const CHIPID_STR_LEN = 32
const ENC_DATA_192 = 192
const ENC_DATA_64 = 64

const OPER_CN_MOBILE = 0
const OPER_CN_UNICOM = 1
const OPER_CN_TELECOM = 2
const OPER_MAX = 3

type SIM_DATA struct {
	Iccid string
	Imsi  string
	Ki    string
	Opc   string
}

type CDMA_DATA struct {
	Imsi_m string
	Uim_id string
	Hrdupp string
}

type SRC_SIM_DATA struct {
	Imei     string
	ChipID   string
	VsimData [OPER_MAX]SIM_DATA
	CdmaData CDMA_DATA
}

type ENC_SIM_DATA struct {
	EncData192 string
	EncData64  string
}

func Min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func IntPtr(n int) uintptr {
	return uintptr(n)
}

func StrPtr(s string) uintptr {
	return uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(s)))
}

func Lib_vsim_encrypt(reqsim SRC_SIM_DATA, encsim *ENC_SIM_DATA) {
	//processProfileData(SRC_SIM_DATA *profile, ENC_SIM_DATA *encProfile)
	lib := syscall.NewLazyDLL("simcrypt.dll")
	//fmt.Println("dll:", lib.Name)
	vsim_encrypt := lib.NewProc("processProfileData")

	// C struct set and transform
	var srcSim C.SRC_SIM_DATA
	var encSim C.ENC_SIM_DATA

	cs := C.CString(reqsim.Imei)
	C.strncpy(&srcSim.imei[0], cs, IMEI_STR_LEN)
	C.free(unsafe.Pointer(cs))

	cs = C.CString(reqsim.ChipID)
	C.strncpy(&srcSim.chipID[0], C.CString(reqsim.ChipID), CHIPID_STR_LEN)
	C.free(unsafe.Pointer(cs))

	cs = C.CString(reqsim.CdmaData.Imsi_m)
	C.strncpy(&srcSim.cdmaData.imsi_m[0], C.CString(reqsim.CdmaData.Imsi_m), IMSI_M_STR_LEN)
	C.free(unsafe.Pointer(cs))

	cs = C.CString(reqsim.CdmaData.Uim_id)
	C.strncpy(&srcSim.cdmaData.uim_id[0], C.CString(reqsim.CdmaData.Uim_id), UIMID_STR_LEN)
	C.free(unsafe.Pointer(cs))

	cs = C.CString(reqsim.CdmaData.Hrdupp)
	C.strncpy(&srcSim.cdmaData.hrdupp[0], C.CString(reqsim.CdmaData.Hrdupp), HRDUPP_STR_LEN)
	C.free(unsafe.Pointer(cs))

	for index := 0; index < OPER_MAX; index++ {
		cs = C.CString(reqsim.VsimData[index].Iccid)
		C.strncpy(&srcSim.vsimData[index].iccid[0], C.CString(reqsim.VsimData[index].Iccid), ICCID_STR_LEN)
		C.free(unsafe.Pointer(cs))

		cs = C.CString(reqsim.VsimData[index].Imsi)
		C.strncpy(&srcSim.vsimData[index].imsi[0], C.CString(reqsim.VsimData[index].Imsi), IMSI_STR_LEN)
		C.free(unsafe.Pointer(cs))

		cs = C.CString(reqsim.VsimData[index].Ki)
		C.strncpy(&srcSim.vsimData[index].ki[0], C.CString(reqsim.VsimData[index].Ki), KI_STR_LEN)
		C.free(unsafe.Pointer(cs))

		cs = C.CString(reqsim.VsimData[index].Opc)
		C.strncpy(&srcSim.vsimData[index].opc[0], C.CString(reqsim.VsimData[index].Opc), OPC_STR_LEN)
		C.free(unsafe.Pointer(cs))
	}

	// C Call DLL
	ret, _, err := vsim_encrypt.Call(uintptr(unsafe.Pointer(&srcSim)), uintptr(unsafe.Pointer(&encSim)))
	if err != nil {
		fmt.Printf("lib.dll运算结果为: %d\n", ret)
		(*encsim).EncData192 = C.GoString(&encSim.encData192[0])
		(*encsim).EncData64 = C.GoString(&encSim.encData64[0])

		/*
			fmt.Printf("imei: %s\n", C.GoString(&srcSim.imei[0]))
			fmt.Printf("chipID: %s\n", C.GoString(&srcSim.chipID[0]))
			fmt.Printf("imsi_m: %s\n", C.GoString(&srcSim.cdmaData.imsi_m[0]))
			fmt.Printf("uim_id: %s\n", C.GoString(&srcSim.cdmaData.uim_id[0]))
			fmt.Printf("hrdupp: %s\n", C.GoString(&srcSim.cdmaData.hrdupp[0]))
			fmt.Printf("iccid: %s\n", C.GoString(&srcSim.vsimData[0].iccid[0]))
			fmt.Printf("imsi: %s\n", C.GoString(&srcSim.vsimData[0].imsi[0]))
			fmt.Printf("ki: %s\n", C.GoString(&srcSim.vsimData[0].ki[0]))
			fmt.Printf("opc: %s\n", C.GoString(&srcSim.vsimData[0].opc[0]))

			fmt.Printf("\nEncData192:\n")
			for index := 0; index < ENC_DATA_192; index++ {
				ens := []byte((*encsim).EncData192)
				fmt.Printf("%02X ", ens[index])
			}
			fmt.Printf("\n")
			fmt.Printf("\nEncData64:\n")
			for index := 0; index < ENC_DATA_64; index++ {
				ens := []byte((*encsim).EncData64)
				fmt.Printf("%02X ", ens[index])
			}
			fmt.Printf("\n")
		*/
	}
}

func test_vsim_encrypt() {
	var encsim ENC_SIM_DATA

	//test data
	imeiStr := []byte{0x38, 0x36, 0x32, 0x31, 0x30, 0x37, 0x30, 0x34, 0x32, 0x38, 0x38, 0x30, 0x38, 0x32, 0x33}
	chipidStr := []byte{0x33, 0x31, 0x33, 0x33, 0x33, 0x32, 0x33, 0x33, 0x33, 0x35, 0x33, 0x30, 0x33, 0x33, 0x33, 0x30, 0x33, 0x39, 0x33, 0x37, 0x30, 0x41, 0x33, 0x32, 0x33, 0x36, 0x33, 0x37, 0x33, 0x41, 0x33, 0x45}
	iccidStr := []byte{0x38, 0x39, 0x38, 0x36, 0x30, 0x32, 0x42, 0x32, 0x32, 0x31, 0x31, 0x36, 0x43, 0x30, 0x30, 0x30, 0x39, 0x38, 0x36, 0x38}
	imsiStr := []byte{0x34, 0x36, 0x30, 0x30, 0x34, 0x30, 0x32, 0x34, 0x30, 0x32, 0x31, 0x35, 0x33, 0x36, 0x38}
	kiStr := []byte{0x32, 0x31, 0x32, 0x31, 0x33, 0x46, 0x46, 0x36, 0x33, 0x36, 0x38, 0x37, 0x32, 0x30, 0x35, 0x46, 0x46, 0x34, 0x36, 0x43, 0x42, 0x31, 0x37, 0x32, 0x46, 0x34, 0x44, 0x35, 0x41, 0x39, 0x34, 0x44}
	opcStr := []byte{0x43, 0x42, 0x31, 0x35, 0x30, 0x32, 0x44, 0x32, 0x38, 0x35, 0x44, 0x31, 0x45, 0x34, 0x39, 0x46, 0x44, 0x35, 0x39, 0x41, 0x44, 0x45, 0x35, 0x39, 0x44, 0x30, 0x36, 0x33, 0x34, 0x34, 0x45, 0x35}
	imsimStr := []byte{0x34, 0x36, 0x30, 0x30, 0x33, 0x38, 0x37, 0x34, 0x31, 0x35, 0x38, 0x37, 0x32, 0x37, 0x38}
	uimidStr := []byte{0x38, 0x30, 0x32, 0x41, 0x44, 0x41, 0x37, 0x44}
	hrduppStr := []byte{0x34, 0x36, 0x30, 0x30, 0x33, 0x38, 0x37, 0x34, 0x31, 0x35, 0x38, 0x37, 0x32, 0x37, 0x38, 0x40, 0x6D, 0x79, 0x63, 0x64, 0x6D, 0x61, 0x2E, 0x63, 0x6E}

	srcsim := SRC_SIM_DATA{
		Imei:   string(imeiStr[:]),
		ChipID: string(chipidStr[:]),
		CdmaData: CDMA_DATA{
			Imsi_m: string(imsimStr[:]),
			Uim_id: string(uimidStr[:]),
			Hrdupp: string(hrduppStr[:]),
		},
	}

	srcsim.VsimData[OPER_CN_MOBILE] = SIM_DATA{
		Iccid: string(iccidStr[:]),
		Imsi:  string(imsiStr[:]),
		Ki:    string(kiStr[:]),
		Opc:   string(opcStr[:]),
	}

	Lib_vsim_encrypt(srcsim, &encsim)
	fmt.Printf("\nEncData192:\n")
	for index := 0; index < ENC_DATA_192; index++ {
		ens := []byte(encsim.EncData192)
		fmt.Printf("%02X ", ens[index])
	}
	fmt.Printf("\nEncData64:\n")
	for index := 0; index < ENC_DATA_64; index++ {
		ens := []byte(encsim.EncData64)
		fmt.Printf("%02X ", ens[index])
	}
	fmt.Printf("\n")
}
