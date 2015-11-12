package soundeffect

import (
	"pkg.deepin.io/dde/daemon/soundplayer"
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/lib/dbus/property"
	"pkg.deepin.io/lib/gio-2.0"
)

const (
	soundEffectSchema = "com.deepin.dde.sound-effect"

	dbusDest = "com.deepin.daemon.SoundEffect"
	dbusPath = "/com/deepin/daemon/SoundEffect"
	dbusIFC  = dbusDest
)

type Manager struct {
	Login         *property.GSettingsBoolProperty `access:"readwrite"`
	Shutdown      *property.GSettingsBoolProperty `access:"readwrite"`
	Logout        *property.GSettingsBoolProperty `access:"readwrite"`
	Wakeup        *property.GSettingsBoolProperty `access:"readwrite"`
	Notification  *property.GSettingsBoolProperty `access:"readwrite"`
	UnableOperate *property.GSettingsBoolProperty `access:"readwrite"`
	EmptyTrash    *property.GSettingsBoolProperty `access:"readwrite"`
	VolumeChange  *property.GSettingsBoolProperty `access:"readwrite"`
	BatteryLow    *property.GSettingsBoolProperty `access:"readwrite"`
	PowerPlug     *property.GSettingsBoolProperty `access:"readwrite"`
	PowerUnplug   *property.GSettingsBoolProperty `access:"readwrite"`
	DevicePlug    *property.GSettingsBoolProperty `access:"readwrite"`
	DeviceUnplug  *property.GSettingsBoolProperty `access:"readwrite"`
	IconToDesktop *property.GSettingsBoolProperty `access:"readwrite"`
	Screenshot    *property.GSettingsBoolProperty `access:"readwrite"`

	setting *gio.Settings
}

func NewManager() *Manager {
	var m = new(Manager)

	m.setting = gio.NewSettings(soundEffectSchema)
	m.Login = property.NewGSettingsBoolProperty(
		m, "Login",
		m.setting, soundplayer.KeyLogin)
	m.Shutdown = property.NewGSettingsBoolProperty(
		m, "Shutdown",
		m.setting, soundplayer.KeyShutdown)
	m.Logout = property.NewGSettingsBoolProperty(
		m, "Logout",
		m.setting, soundplayer.KeyLogout)
	m.Wakeup = property.NewGSettingsBoolProperty(
		m, "Wakeup",
		m.setting, soundplayer.KeyWakeup)
	m.Notification = property.NewGSettingsBoolProperty(
		m, "Notification",
		m.setting, soundplayer.KeyNotification)
	m.UnableOperate = property.NewGSettingsBoolProperty(
		m, "UnableOperate",
		m.setting, soundplayer.KeyUnableOperate)
	m.EmptyTrash = property.NewGSettingsBoolProperty(
		m, "EmptyTrash",
		m.setting, soundplayer.KeyEmptyTrash)
	m.VolumeChange = property.NewGSettingsBoolProperty(
		m, "VolumeChange",
		m.setting, soundplayer.KeyVolumeChange)
	m.BatteryLow = property.NewGSettingsBoolProperty(
		m, "BatteryLow",
		m.setting, soundplayer.KeyBatteryLow)
	m.PowerPlug = property.NewGSettingsBoolProperty(
		m, "PowerPlug",
		m.setting, soundplayer.KeyPowerPlug)
	m.PowerUnplug = property.NewGSettingsBoolProperty(
		m, "PowerUnplug",
		m.setting, soundplayer.KeyPowerUnplug)
	m.DevicePlug = property.NewGSettingsBoolProperty(
		m, "DevicePlug",
		m.setting, soundplayer.KeyDevicePlug)
	m.DeviceUnplug = property.NewGSettingsBoolProperty(
		m, "DeviceUnplug",
		m.setting, soundplayer.KeyDeviceUnplug)
	m.IconToDesktop = property.NewGSettingsBoolProperty(
		m, "IconToDesktop",
		m.setting, soundplayer.KeyIconToDesktop)
	m.Screenshot = property.NewGSettingsBoolProperty(
		m, "Screenshot",
		m.setting, soundplayer.KeyScreenshot)

	return m
}

func (*Manager) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		Dest:       dbusDest,
		ObjectPath: dbusPath,
		Interface:  dbusIFC,
	}
}