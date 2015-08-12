package main

import (
	"os"
	"path"
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/lib/glib-2.0"
	dutils "pkg.deepin.io/lib/utils"
)

const (
	AppMimeTerminal = "application/x-terminal"
)

const (
	dbusDest = "com.deepin.api.Mime"
	dbusPath = "/com/deepin/api/Manager"
	dbusIFC  = "com.deepin.api.Manager"
)

type Manager struct {
	media *Media
}

func NewManager() *Manager {
	m := new(Manager)
	m.initConfigData()
	return m
}

func (m *Manager) initConfigData() {
	if dutils.IsFileExist(path.Join(glib.GetUserConfigDir(),
		"mimeapps.list")) {
		return
	}

	err := m.doInitConfigData()
	if err != nil {
		logger.Warning("Init mime config file failed", err)
	}
}

func (m *Manager) doInitConfigData() error {
	os.Remove(path.Join(glib.GetUserConfigDir(), "mimeapps.list"))

	var data = "data.json"
	switch os.Getenv("LANGUAGE") {
	case "zh_CN", "zh_TW", "zh_HK":
		data = "data-zh_CN.json"
	}
	return genMimeAppsFile(
		findFilePath(path.Join("dde-api", "mime", data)))
}

// Reset reset mimes default app
func (m *Manager) Reset() {
	err := m.doInitConfigData()
	if err != nil {
		logger.Warning("Init mime config file failed", err)
	}
	resetTerminal()
}

// GetDefaultApp get the default app id for the special mime
// ty: the special mime
// ret0: the default app info
// ret1: error message
func (m *Manager) GetDefaultApp(ty string) (string, error) {
	var (
		info *AppInfo
		err  error
	)
	if ty == AppMimeTerminal {
		info, err = getDefaultTerminal()
	} else {
		info, err = GetAppInfo(ty)
	}
	if err != nil {
		return "", err
	}

	return marshal(info)
}

// SetDefaultApp set the default app for the special mime
// ty: the special mime
// deskId: the default app desktop id
// ret0: error message
func (m *Manager) SetDefaultApp(ty string, deskId string) error {
	if ty == AppMimeTerminal {
		return setDefaultTerminal(deskId)
	}
	return SetAppInfo(ty, deskId)
}

// ListApps list the apps that supported the special mime
// ty: the special mime
// ret0: the app infos
func (m *Manager) ListApps(ty string) string {
	var infos AppInfos
	if ty == AppMimeTerminal {
		infos = getTerminalInfos()
	} else {
		infos = GetAppInfos(ty)
	}

	content, _ := marshal(infos)
	return content
}

func (m *Manager) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		Dest:       dbusDest,
		ObjectPath: dbusPath,
		Interface:  dbusIFC,
	}
}