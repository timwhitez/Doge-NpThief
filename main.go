package main

import (
	"fmt"
	"strings"
	"syscall"
	"unsafe"
)

func main(){
	const WM_GETTEXT = 0x000D
	var lent = 30000

	pFindWindowExA := syscall.NewLazyDLL("user32.dll").NewProc("FindWindowExA")
	notep_, _ := syscall.BytePtrFromString("Notepad")
	notep := uintptr(unsafe.Pointer(notep_))

	noteHwnd,_,_ := pFindWindowExA.Call(0,0,notep,0)

	if noteHwnd!= 0{
		fmt.Printf("[+] Found Window %x\n",noteHwnd)

		editp_, _ := syscall.BytePtrFromString("Edit")
		editp := uintptr(unsafe.Pointer(editp_))

		editHwnd,_,_ := pFindWindowExA.Call(noteHwnd,0,editp,0)
		if editHwnd != 0{

			buff := make([]byte,lent+1)
			pSendMessageA := syscall.NewLazyDLL("user32.dll").NewProc("SendMessageA")
			pSendMessageA.Call(editHwnd,WM_GETTEXT,uintptr(lent),uintptr(unsafe.Pointer(&buff[0])))

			fmt.Printf("[+] Content: \n%s\n\n",strings.Replace(string(buff), "\x00", "", -1))

		}

	}else{
		fmt.Println("[-] Cannot Found Window")
	}

}