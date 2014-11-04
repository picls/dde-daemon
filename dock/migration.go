package dock

// this file should be rewrite

//#cgo pkg-config: glib-2.0 gio-unix-2.0 gtk+-3.0
//#include <stdlib.h>
// char* guess_app_id(long s_pid, const char* instance_name, const char* wmname, const char* wmclass, const char* icon_name);
// char* get_exec(int pid);
// char* get_exe(int pid);
//char* icon_name_to_path(const char* name, int size);
// void init_deepin();
// char* get_data_uri_by_path(const char* path);
import "C"
import "strings"
import "unsafe"

func find_app_id(pid uint, instanceName, wmName, wmClass, iconName string) string {
	iName := C.CString(instanceName)
	wName := C.CString(wmName)
	wClass := C.CString(wmClass)
	icon := C.CString(iconName)
	id := C.guess_app_id(C.long(pid), iName, wName, wClass, icon)
	defer func() {
		C.free(unsafe.Pointer(iName))
		C.free(unsafe.Pointer(wName))
		C.free(unsafe.Pointer(wClass))
		C.free(unsafe.Pointer(icon))
		C.free(unsafe.Pointer(id))
	}()
	return strings.ToLower(C.GoString(id))
}

func find_exec_by_pid(pid uint) string {
	e := C.GoString(C.get_exec(C.int(pid)))
	if e != "" {
		return e
	}
	return C.GoString(C.get_exe(C.int(pid)))
}

func get_theme_icon(name string, size int) string {
	iconName := C.CString(name)
	defer func() {
		C.free(unsafe.Pointer(iconName))
	}()
	return C.GoString(C.icon_name_to_path(iconName, C.int(size)))
}

func initDeepin() {
	C.init_deepin()
}

func xpm_to_dataurl(icon string) string {
	iconName := C.CString(icon)
	defer func() {
		C.free(unsafe.Pointer(iconName))
	}()
	return C.GoString(C.get_data_uri_by_path(iconName))

}
