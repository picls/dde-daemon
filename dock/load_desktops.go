package dock

import (
	"dbus/com/deepin/daemon/dock"
)

func loadAll() []string {
	DOCKED_APP_MANAGER, err := dock.NewDockedAppManager(
		"com.deepin.daemon.Dock",
		"/dde/dock/DockedAppManager",
	)
	if err != nil {
		logger.Warning("get DockedAppManager failed", err)
		return make([]string, 0)
	}

	l, _ := DOCKED_APP_MANAGER.DockedAppList()
	return l
}
