// Code generated by 'go generate'; DO NOT EDIT.

package windriver

import (
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

var _ unsafe.Pointer

// Do the interface allocations only once for common
// Errno values.
const (
	errnoERROR_IO_PENDING = 997
)

var (
	errERROR_IO_PENDING error = syscall.Errno(errnoERROR_IO_PENDING)
)

// errnoErr returns common boxed Errno values, to prevent
// allocations at runtime.
func errnoErr(e syscall.Errno) error {
	switch e {
	case 0:
		return nil
	case errnoERROR_IO_PENDING:
		return errERROR_IO_PENDING
	}
	// TODO: add more here, after collecting data on the common
	// error values see on Windows. (perhaps when running
	// all.bat?)
	return e
}

var (
	modkernel32 = windows.NewLazySystemDLL("kernel32.dll")
	moduser32   = windows.NewLazySystemDLL("user32.dll")
	modgdi32    = windows.NewLazySystemDLL("gdi32.dll")
	modshell32  = windows.NewLazySystemDLL("shell32.dll")

	procGetModuleHandleW         = modkernel32.NewProc("GetModuleHandleW")
	procGlobalLock               = modkernel32.NewProc("GlobalLock")
	procGlobalUnlock             = modkernel32.NewProc("GlobalUnlock")
	procGlobalAlloc              = modkernel32.NewProc("GlobalAlloc")
	procGetConsoleWindow         = modkernel32.NewProc("GetConsoleWindow")
	procGetCurrentProcessId      = modkernel32.NewProc("GetCurrentProcessId")
	procLoadImageW               = moduser32.NewProc("LoadImageW")
	procLoadCursorW              = moduser32.NewProc("LoadCursorW")
	procRegisterClassExW         = moduser32.NewProc("RegisterClassExW")
	procCreateWindowExW          = moduser32.NewProc("CreateWindowExW")
	procPostMessageW             = moduser32.NewProc("PostMessageW")
	procGetMessageW              = moduser32.NewProc("GetMessageW")
	procTranslateAccelerator     = moduser32.NewProc("TranslateAccelerator")
	procTranslateMessage         = moduser32.NewProc("TranslateMessage")
	procDispatchMessageW         = moduser32.NewProc("DispatchMessageW")
	procDefWindowProcW           = moduser32.NewProc("DefWindowProcW")
	procGetWindowRect            = moduser32.NewProc("GetWindowRect")
	procSetCursor                = moduser32.NewProc("SetCursor")
	procDestroyWindow            = moduser32.NewProc("DestroyWindow")
	procPostQuitMessage          = moduser32.NewProc("PostQuitMessage")
	procGetCursorPos             = moduser32.NewProc("GetCursorPos")
	procValidateRect             = moduser32.NewProc("ValidateRect")
	procInvalidateRect           = moduser32.NewProc("InvalidateRect")
	procBeginPaint               = moduser32.NewProc("BeginPaint")
	procEndPaint                 = moduser32.NewProc("EndPaint")
	procUpdateWindow             = moduser32.NewProc("UpdateWindow")
	procRedrawWindow             = moduser32.NewProc("RedrawWindow")
	procShowWindow               = moduser32.NewProc("ShowWindow")
	procShowWindowAsync          = moduser32.NewProc("ShowWindowAsync")
	procGetDC                    = moduser32.NewProc("GetDC")
	procReleaseDC                = moduser32.NewProc("ReleaseDC")
	procMapVirtualKeyW           = moduser32.NewProc("MapVirtualKeyW")
	procToUnicode                = moduser32.NewProc("ToUnicode")
	procGetKeyboardState         = moduser32.NewProc("GetKeyboardState")
	procGetKeyState              = moduser32.NewProc("GetKeyState")
	procSetCursorPos             = moduser32.NewProc("SetCursorPos")
	procMapWindowPoints          = moduser32.NewProc("MapWindowPoints")
	procClientToScreen           = moduser32.NewProc("ClientToScreen")
	procOpenClipboard            = moduser32.NewProc("OpenClipboard")
	procCloseClipboard           = moduser32.NewProc("CloseClipboard")
	procSetClipboardData         = moduser32.NewProc("SetClipboardData")
	procGetClipboardData         = moduser32.NewProc("GetClipboardData")
	procEmptyClipboard           = moduser32.NewProc("EmptyClipboard")
	procGetWindowThreadProcessId = moduser32.NewProc("GetWindowThreadProcessId")
	procSetWindowTextW           = moduser32.NewProc("SetWindowTextW")
	procSelectObject             = modgdi32.NewProc("SelectObject")
	procCreateBitmap             = modgdi32.NewProc("CreateBitmap")
	procCreateCompatibleBitmap   = modgdi32.NewProc("CreateCompatibleBitmap")
	procDeleteObject             = modgdi32.NewProc("DeleteObject")
	procCreateCompatibleDC       = modgdi32.NewProc("CreateCompatibleDC")
	procDeleteDC                 = modgdi32.NewProc("DeleteDC")
	procBitBlt                   = modgdi32.NewProc("BitBlt")
	procSetPixel                 = modgdi32.NewProc("SetPixel")
	procCreateBitmapIndirect     = modgdi32.NewProc("CreateBitmapIndirect")
	procGetObject                = modgdi32.NewProc("GetObject")
	procCreateDIBSection         = modgdi32.NewProc("CreateDIBSection")
	procDragAcceptFiles          = modshell32.NewProc("DragAcceptFiles")
	procDragQueryPoint           = modshell32.NewProc("DragQueryPoint")
	procDragQueryFileW           = modshell32.NewProc("DragQueryFileW")
	procDragFinish               = modshell32.NewProc("DragFinish")
)

func _GetModuleHandleW(name *uint16) (modH windows.Handle, err error) {
	r0, _, e1 := syscall.Syscall(procGetModuleHandleW.Addr(), 1, uintptr(unsafe.Pointer(name)), 0, 0)
	modH = windows.Handle(r0)
	if modH == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func _GlobalLock(h windows.Handle) (ptr uintptr, err error) {
	r0, _, e1 := syscall.Syscall(procGlobalLock.Addr(), 1, uintptr(h), 0, 0)
	ptr = uintptr(r0)
	if ptr == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func _GlobalUnlock(h windows.Handle) (ok bool) {
	r0, _, _ := syscall.Syscall(procGlobalUnlock.Addr(), 1, uintptr(h), 0, 0)
	ok = r0 != 0
	return
}

func _GlobalAlloc(uFlags uint32, dwBytes uintptr) (h windows.Handle, err error) {
	r0, _, e1 := syscall.Syscall(procGlobalAlloc.Addr(), 2, uintptr(uFlags), uintptr(dwBytes), 0)
	h = windows.Handle(r0)
	if h == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func _GetConsoleWindow() (cH windows.Handle) {
	r0, _, _ := syscall.Syscall(procGetConsoleWindow.Addr(), 0, 0, 0, 0)
	cH = windows.Handle(r0)
	return
}

func _GetCurrentProcessId() (pid uint32) {
	r0, _, _ := syscall.Syscall(procGetCurrentProcessId.Addr(), 0, 0, 0, 0)
	pid = uint32(r0)
	return
}

func _LoadImageW(hInstance windows.Handle, name uintptr, typ uint32, cx int32, cy int32, fuLoad uint32) (imgH windows.Handle, err error) {
	r0, _, e1 := syscall.Syscall6(procLoadImageW.Addr(), 6, uintptr(hInstance), uintptr(name), uintptr(typ), uintptr(cx), uintptr(cy), uintptr(fuLoad))
	imgH = windows.Handle(r0)
	if imgH == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func _LoadCursorW(hInstance windows.Handle, name uint32) (cursorH windows.Handle, err error) {
	r0, _, e1 := syscall.Syscall(procLoadCursorW.Addr(), 2, uintptr(hInstance), uintptr(name), 0)
	cursorH = windows.Handle(r0)
	if cursorH == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func _RegisterClassExW(wcx *_WndClassExW) (atom uint16, err error) {
	r0, _, e1 := syscall.Syscall(procRegisterClassExW.Addr(), 1, uintptr(unsafe.Pointer(wcx)), 0, 0)
	atom = uint16(r0)
	if atom == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func _CreateWindowExW(dwExStyle uint32, lpClassName *uint16, lpWindowName *uint16, dwStyle int32, x int32, y int32, nWidth int32, nHeight int32, hWndParent windows.Handle, hMenu windows.Handle, hInstance windows.Handle, lpParam uintptr) (wndH windows.Handle, err error) {
	r0, _, e1 := syscall.Syscall12(procCreateWindowExW.Addr(), 12, uintptr(dwExStyle), uintptr(unsafe.Pointer(lpClassName)), uintptr(unsafe.Pointer(lpWindowName)), uintptr(dwStyle), uintptr(x), uintptr(y), uintptr(nWidth), uintptr(nHeight), uintptr(hWndParent), uintptr(hMenu), uintptr(hInstance), uintptr(lpParam))
	wndH = windows.Handle(r0)
	if wndH == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func _PostMessageW(hwnd windows.Handle, msg uint32, wParam uintptr, lParam uintptr) (ok bool) {
	r0, _, _ := syscall.Syscall6(procPostMessageW.Addr(), 4, uintptr(hwnd), uintptr(msg), uintptr(wParam), uintptr(lParam), 0, 0)
	ok = r0 != 0
	return
}

func _GetMessageW(msg *_Msg, hwnd windows.Handle, msgFilterMin uint32, msgFilterMax uint32) (res int32, err error) {
	r0, _, e1 := syscall.Syscall6(procGetMessageW.Addr(), 4, uintptr(unsafe.Pointer(msg)), uintptr(hwnd), uintptr(msgFilterMin), uintptr(msgFilterMax), 0, 0)
	res = int32(r0)
	if res == -1 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func _TranslateAccelerator(hwnd windows.Handle, hAccTable windows.Handle, msg *_Msg) (ok bool) {
	r0, _, _ := syscall.Syscall(procTranslateAccelerator.Addr(), 3, uintptr(hwnd), uintptr(hAccTable), uintptr(unsafe.Pointer(msg)))
	ok = r0 != 0
	return
}

func _TranslateMessage(msg *_Msg) (translated bool) {
	r0, _, _ := syscall.Syscall(procTranslateMessage.Addr(), 1, uintptr(unsafe.Pointer(msg)), 0, 0)
	translated = r0 != 0
	return
}

func _DispatchMessageW(msg *_Msg) (res int32) {
	r0, _, _ := syscall.Syscall(procDispatchMessageW.Addr(), 1, uintptr(unsafe.Pointer(msg)), 0, 0)
	res = int32(r0)
	return
}

func _DefWindowProcW(hwnd windows.Handle, msg uint32, wparam uintptr, lparam uintptr) (ret uintptr) {
	r0, _, _ := syscall.Syscall6(procDefWindowProcW.Addr(), 4, uintptr(hwnd), uintptr(msg), uintptr(wparam), uintptr(lparam), 0, 0)
	ret = uintptr(r0)
	return
}

func _GetWindowRect(hwnd windows.Handle, r *_Rect) (ok bool) {
	r0, _, _ := syscall.Syscall(procGetWindowRect.Addr(), 2, uintptr(hwnd), uintptr(unsafe.Pointer(r)), 0)
	ok = r0 != 0
	return
}

func _SetCursor(cursorH windows.Handle) (prevCursorH windows.Handle) {
	r0, _, _ := syscall.Syscall(procSetCursor.Addr(), 1, uintptr(cursorH), 0, 0)
	prevCursorH = windows.Handle(r0)
	return
}

func _DestroyWindow(hwnd windows.Handle) (ok bool) {
	r0, _, _ := syscall.Syscall(procDestroyWindow.Addr(), 1, uintptr(hwnd), 0, 0)
	ok = r0 != 0
	return
}

func _PostQuitMessage(exitCode int32) {
	syscall.Syscall(procPostQuitMessage.Addr(), 1, uintptr(exitCode), 0, 0)
	return
}

func _GetCursorPos(p *_Point) (ok bool) {
	r0, _, _ := syscall.Syscall(procGetCursorPos.Addr(), 1, uintptr(unsafe.Pointer(p)), 0, 0)
	ok = r0 != 0
	return
}

func _ValidateRect(hwnd windows.Handle, r *_Rect) (ok bool) {
	r0, _, _ := syscall.Syscall(procValidateRect.Addr(), 2, uintptr(hwnd), uintptr(unsafe.Pointer(r)), 0)
	ok = r0 != 0
	return
}

func _InvalidateRect(hwnd windows.Handle, r *_Rect, erase bool) (ok bool) {
	var _p0 uint32
	if erase {
		_p0 = 1
	} else {
		_p0 = 0
	}
	r0, _, _ := syscall.Syscall(procInvalidateRect.Addr(), 3, uintptr(hwnd), uintptr(unsafe.Pointer(r)), uintptr(_p0))
	ok = r0 != 0
	return
}

func _BeginPaint(hwnd windows.Handle, paint *_Paint) (dcH windows.Handle, err error) {
	r0, _, e1 := syscall.Syscall(procBeginPaint.Addr(), 2, uintptr(hwnd), uintptr(unsafe.Pointer(paint)), 0)
	dcH = windows.Handle(r0)
	if dcH == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func _EndPaint(hwnd windows.Handle, paint *_Paint) (ok bool) {
	r0, _, _ := syscall.Syscall(procEndPaint.Addr(), 2, uintptr(hwnd), uintptr(unsafe.Pointer(paint)), 0)
	ok = r0 != 0
	return
}

func _UpdateWindow(hwnd windows.Handle) (ok bool) {
	r0, _, _ := syscall.Syscall(procUpdateWindow.Addr(), 1, uintptr(hwnd), 0, 0)
	ok = r0 != 0
	return
}

func _RedrawWindow(hwnd windows.Handle, r *_Rect, region windows.Handle, flags uint) (ok bool) {
	r0, _, _ := syscall.Syscall6(procRedrawWindow.Addr(), 4, uintptr(hwnd), uintptr(unsafe.Pointer(r)), uintptr(region), uintptr(flags), 0, 0)
	ok = r0 != 0
	return
}

func _ShowWindow(hwnd windows.Handle, nCmdShow int) (ok bool) {
	r0, _, _ := syscall.Syscall(procShowWindow.Addr(), 2, uintptr(hwnd), uintptr(nCmdShow), 0)
	ok = r0 != 0
	return
}

func _ShowWindowAsync(hwnd windows.Handle, nCmdShow int) (ok bool) {
	r0, _, _ := syscall.Syscall(procShowWindowAsync.Addr(), 2, uintptr(hwnd), uintptr(nCmdShow), 0)
	ok = r0 != 0
	return
}

func _GetDC(hwnd windows.Handle) (dcH windows.Handle, err error) {
	r0, _, e1 := syscall.Syscall(procGetDC.Addr(), 1, uintptr(hwnd), 0, 0)
	dcH = windows.Handle(r0)
	if dcH == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func _ReleaseDC(hwnd windows.Handle, dc windows.Handle) (ok bool) {
	r0, _, _ := syscall.Syscall(procReleaseDC.Addr(), 2, uintptr(hwnd), uintptr(dc), 0)
	ok = r0 != 0
	return
}

func _MapVirtualKeyW(uCode uint32, uMapType uint32) (code uint32) {
	r0, _, _ := syscall.Syscall(procMapVirtualKeyW.Addr(), 2, uintptr(uCode), uintptr(uMapType), 0)
	code = uint32(r0)
	return
}

func _ToUnicode(wVirtKey uint32, wScanCode uint32, lpKeyState *[256]byte, pwszBuff *uint16, cchBuff int32, wFlags uint32) (code int32) {
	r0, _, _ := syscall.Syscall6(procToUnicode.Addr(), 6, uintptr(wVirtKey), uintptr(wScanCode), uintptr(unsafe.Pointer(lpKeyState)), uintptr(unsafe.Pointer(pwszBuff)), uintptr(cchBuff), uintptr(wFlags))
	code = int32(r0)
	return
}

func _GetKeyboardState(state *[256]byte) (ok bool) {
	r0, _, _ := syscall.Syscall(procGetKeyboardState.Addr(), 1, uintptr(unsafe.Pointer(state)), 0, 0)
	ok = r0 != 0
	return
}

func _GetKeyState(vkey int32) (state uint16) {
	r0, _, _ := syscall.Syscall(procGetKeyState.Addr(), 1, uintptr(vkey), 0, 0)
	state = uint16(r0)
	return
}

func _SetCursorPos(x int32, y int32) (ok bool) {
	r0, _, _ := syscall.Syscall(procSetCursorPos.Addr(), 2, uintptr(x), uintptr(y), 0)
	ok = r0 != 0
	return
}

func _MapWindowPoints(hwndFrom windows.Handle, hwndTo windows.Handle, lpPoints *_Point, cPoints uint32) (res int32) {
	r0, _, _ := syscall.Syscall6(procMapWindowPoints.Addr(), 4, uintptr(hwndFrom), uintptr(hwndTo), uintptr(unsafe.Pointer(lpPoints)), uintptr(cPoints), 0, 0)
	res = int32(r0)
	return
}

func _ClientToScreen(hwnd windows.Handle, lpPoint *_Point) (ok bool) {
	r0, _, _ := syscall.Syscall(procClientToScreen.Addr(), 2, uintptr(hwnd), uintptr(unsafe.Pointer(lpPoint)), 0)
	ok = r0 != 0
	return
}

func _OpenClipboard(hWndNewOwner windows.Handle) (ok bool) {
	r0, _, _ := syscall.Syscall(procOpenClipboard.Addr(), 1, uintptr(hWndNewOwner), 0, 0)
	ok = r0 != 0
	return
}

func _CloseClipboard() (ok bool) {
	r0, _, _ := syscall.Syscall(procCloseClipboard.Addr(), 0, 0, 0, 0)
	ok = r0 != 0
	return
}

func _SetClipboardData(uFormat uint32, h windows.Handle) (dataH windows.Handle, err error) {
	r0, _, e1 := syscall.Syscall(procSetClipboardData.Addr(), 2, uintptr(uFormat), uintptr(h), 0)
	dataH = windows.Handle(r0)
	if dataH == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func _GetClipboardData(uFormat uint32) (dataH windows.Handle, err error) {
	r0, _, e1 := syscall.Syscall(procGetClipboardData.Addr(), 1, uintptr(uFormat), 0, 0)
	dataH = windows.Handle(r0)
	if dataH == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func _EmptyClipboard() (ok bool) {
	r0, _, _ := syscall.Syscall(procEmptyClipboard.Addr(), 0, 0, 0, 0)
	ok = r0 != 0
	return
}

func _GetWindowThreadProcessId(hwnd windows.Handle, pid *uint32) (threadId uint32) {
	r0, _, _ := syscall.Syscall(procGetWindowThreadProcessId.Addr(), 2, uintptr(hwnd), uintptr(unsafe.Pointer(pid)), 0)
	threadId = uint32(r0)
	return
}

func _SetWindowTextW(hwnd windows.Handle, lpString *uint16) (res bool) {
	r0, _, _ := syscall.Syscall(procSetWindowTextW.Addr(), 2, uintptr(hwnd), uintptr(unsafe.Pointer(lpString)), 0)
	res = r0 != 0
	return
}

func _SelectObject(hdc windows.Handle, obj windows.Handle) (prevObjH windows.Handle, err error) {
	r0, _, e1 := syscall.Syscall(procSelectObject.Addr(), 2, uintptr(hdc), uintptr(obj), 0)
	prevObjH = windows.Handle(r0)
	if prevObjH == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func _CreateBitmap(w int32, h int32, planes uint32, bitCount uint32, bits uintptr) (bmH windows.Handle, err error) {
	r0, _, e1 := syscall.Syscall6(procCreateBitmap.Addr(), 5, uintptr(w), uintptr(h), uintptr(planes), uintptr(bitCount), uintptr(bits), 0)
	bmH = windows.Handle(r0)
	if bmH == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func _CreateCompatibleBitmap(hdc windows.Handle, w int32, h int32) (bmH windows.Handle, err error) {
	r0, _, e1 := syscall.Syscall(procCreateCompatibleBitmap.Addr(), 3, uintptr(hdc), uintptr(w), uintptr(h))
	bmH = windows.Handle(r0)
	if bmH == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func _DeleteObject(obj windows.Handle) (ok bool) {
	r0, _, _ := syscall.Syscall(procDeleteObject.Addr(), 1, uintptr(obj), 0, 0)
	ok = r0 != 0
	return
}

func _CreateCompatibleDC(hdc windows.Handle) (dcH windows.Handle, err error) {
	r0, _, e1 := syscall.Syscall(procCreateCompatibleDC.Addr(), 1, uintptr(hdc), 0, 0)
	dcH = windows.Handle(r0)
	if dcH == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func _DeleteDC(dc windows.Handle) (ok bool) {
	r0, _, _ := syscall.Syscall(procDeleteDC.Addr(), 1, uintptr(dc), 0, 0)
	ok = r0 != 0
	return
}

func _BitBlt(hdc windows.Handle, x int32, y int32, w int32, h int32, hdcSrc windows.Handle, x2 int32, y2 int32, rOp uint32) (ok bool) {
	r0, _, _ := syscall.Syscall9(procBitBlt.Addr(), 9, uintptr(hdc), uintptr(x), uintptr(y), uintptr(w), uintptr(h), uintptr(hdcSrc), uintptr(x2), uintptr(y2), uintptr(rOp))
	ok = r0 != 0
	return
}

func _SetPixel(hdc windows.Handle, x int, y int, c _ColorRef) (colorSet int32, err error) {
	r0, _, e1 := syscall.Syscall6(procSetPixel.Addr(), 4, uintptr(hdc), uintptr(x), uintptr(y), uintptr(c), 0, 0)
	colorSet = int32(r0)
	if colorSet == -1 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func _CreateBitmapIndirect(bm *_Bitmap) (bmH windows.Handle, err error) {
	r0, _, e1 := syscall.Syscall(procCreateBitmapIndirect.Addr(), 1, uintptr(unsafe.Pointer(bm)), 0, 0)
	bmH = windows.Handle(r0)
	if bmH == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func _GetObject(h windows.Handle, c int32, v uintptr) (n int) {
	r0, _, _ := syscall.Syscall(procGetObject.Addr(), 3, uintptr(h), uintptr(c), uintptr(v))
	n = int(r0)
	return
}

func _CreateDIBSection(dc windows.Handle, bmi *_BitmapInfo, usage uint32, bits **byte, section windows.Handle, offset uint32) (bmH windows.Handle, err error) {
	r0, _, e1 := syscall.Syscall6(procCreateDIBSection.Addr(), 6, uintptr(dc), uintptr(unsafe.Pointer(bmi)), uintptr(usage), uintptr(unsafe.Pointer(bits)), uintptr(section), uintptr(offset))
	bmH = windows.Handle(r0)
	if bmH == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func _DragAcceptFiles(hwnd windows.Handle, fAccept bool) {
	var _p0 uint32
	if fAccept {
		_p0 = 1
	} else {
		_p0 = 0
	}
	syscall.Syscall(procDragAcceptFiles.Addr(), 2, uintptr(hwnd), uintptr(_p0), 0)
	return
}

func _DragQueryPoint(hDrop uintptr, ppt *_Point) (res bool) {
	r0, _, _ := syscall.Syscall(procDragQueryPoint.Addr(), 2, uintptr(hDrop), uintptr(unsafe.Pointer(ppt)), 0)
	res = r0 != 0
	return
}

func _DragQueryFileW(hDrop uintptr, iFile uint32, lpszFile *uint16, cch uint32) (res uint32) {
	r0, _, _ := syscall.Syscall6(procDragQueryFileW.Addr(), 4, uintptr(hDrop), uintptr(iFile), uintptr(unsafe.Pointer(lpszFile)), uintptr(cch), 0, 0)
	res = uint32(r0)
	return
}

func _DragFinish(hDrop uintptr) {
	syscall.Syscall(procDragFinish.Addr(), 1, uintptr(hDrop), 0, 0)
	return
}