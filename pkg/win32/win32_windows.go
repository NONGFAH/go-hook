//go:build windows
// +build windows

package win32

import (
	"fmt"
	"syscall"
	"unsafe"

	"github.com/nongfah/go-hook/pkg/types"
)

var (
	modUser32, _               = syscall.LoadDLL("user32.dll")
	procCallNextHookEx, _      = modUser32.FindProc("CallNextHookEx")
	procSetWindowsHookExW, _   = modUser32.FindProc("SetWindowsHookExW")
	procGetMessageW, _         = modUser32.FindProc("GetMessageW")
	procTranslateMessage, _    = modUser32.FindProc("TranslateMessage")
	procDispatchMessageW, _    = modUser32.FindProc("DispatchMessageW")
	procGetModuleHandleW, _    = modUser32.FindProc("GetModuleHandleW")
	procUnhookWindowsHookEx, _ = modUser32.FindProc("UnhookWindowsHookEx")
	procSendInput, _           = modUser32.FindProc("SendInput")
	procMapVirtualKeyW, _      = modUser32.FindProc("MapVirtualKeyW")
)

func CallNextHookEx(hhk uintptr, code int32, wParam, lParam uintptr) uintptr {
	r, _, _ := procCallNextHookEx.Call(hhk, uintptr(code), wParam, lParam)

	return r
}

func SetWindowsHookEx(idHook types.Hook, lpfn, hmod uintptr, dwThreadId uint32) uintptr {
	r, _, _ := procSetWindowsHookExW.Call(uintptr(idHook), lpfn, hmod, uintptr(dwThreadId))

	return r
}

func UnhookWindowsHookEx(hhk uintptr) bool {
	r, _, _ := procUnhookWindowsHookEx.Call(hhk)

	if r == 0 {
		return false
	}

	return true
}

func GetMessage(lpMsg **types.MSG, hWnd uintptr, wMsgFilterMin, wMsgFilterMax uint32) int32 {
	r, _, _ := procGetMessageW.Call(
		uintptr(unsafe.Pointer(lpMsg)),
		hWnd,
		uintptr(wMsgFilterMin),
		uintptr(wMsgFilterMin))

	return int32(r)
}

func TranslateMessage(lpMsg **types.MSG) int32 {
	r, _, _ := procTranslateMessage.Call(uintptr(unsafe.Pointer(lpMsg)))

	return int32(r)
}

func DispatchMessage(lpMsg **types.MSG) int32 {
	r, _, _ := procDispatchMessageW.Call(uintptr(unsafe.Pointer(lpMsg)))

	return int32(r)
}

func GetModuleHandle(lpModuleName uintptr) uintptr {
	r, _, _ := procSetWindowsHookExW.Call(lpModuleName)

	return r
}

func SendInput(numInputs uint, unsafePointerToVal unsafe.Pointer, inputStructSizeBytes uintptr) error {
	numSent, _, err := procSendInput.Call(
		uintptr(numInputs),
		uintptr(unsafePointerToVal),
		inputStructSizeBytes)
	if uint(numSent) == numInputs {
		return nil
	} else if err != nil {
		return err
	}

	return fmt.Errorf("failed to send input, unknown errror")
}

func MapVirtualKey(scanCode uint32) (types.VKCode, error) {
	r1, _, _ := procMapVirtualKeyW.Call(uintptr(scanCode), 3)
	if r1 == 0 {
		return 0, fmt.Errorf("failed to map scanCode to vkCode, unknown errror")
	}
	return types.VKCode(r1), nil
}
