/**
 * Copyright (c) 2014 Deepin, Inc.
 *               2014 Xu FaSheng
 *
 * Author:      Xu FaSheng <fasheng.xu@gmail.com>
 * Maintainer:  Xu FaSheng <fasheng.xu@gmail.com>
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, see <http://www.gnu.org/licenses/>.
 **/

package network

import (
	"pkg.linuxdeepin.com/lib/dbus"
	"time"
)

type switchHandler struct {
	config *config

	NetworkingEnabled bool // airplane mode for NetworkManager
	WirelessEnabled   bool
	WwanEnabled       bool
	WiredEnabled      bool
	VpnEnabled        bool
}

const (
	autoConnectTimeout = 20
)

func newSwitchHandler(c *config) (sh *switchHandler) {
	sh = &switchHandler{config: c}
	sh.init()

	// connect global switch signals
	nmManager.NetworkingEnabled.ConnectChanged(func() {
		sh.setPropNetworkingEnabled()
	})
	nmManager.WirelessEnabled.ConnectChanged(func() {
		sh.setPropWirelessEnabled()
	})
	nmManager.WwanEnabled.ConnectChanged(func() {
		sh.setPropWwanEnabled()
	})

	return
}

func destroySwitchHandler(sh *switchHandler) {
}

func (sh *switchHandler) init() {
	// initialize global switches
	sh.initPropNetworkingEnabled()
	sh.initPropWirelessEnabled()
	sh.initPropWwanEnabled()

	// initialize virtual global switches
	sh.initPropWiredEnabled()
	sh.initPropVpnEnabled()
}

func (sh *switchHandler) setNetworkingEnabled(enabled bool) {
	nmSetNetworkingEnabled(enabled)
}
func (sh *switchHandler) setWirelessEnabled(enabled bool) {
	if sh.NetworkingEnabled {
		nmSetWirelessEnabled(enabled)
	} else {
		// if NetworkingEnabled is off, turn it on, and only keep
		// current global device switch alive
		sh.config.setLastGlobalSwithes(false)
		sh.config.setLastWirelessEnabled(true)
		sh.setNetworkingEnabled(true)
	}
}
func (sh *switchHandler) setWwanEnabled(enabled bool) {
	if sh.NetworkingEnabled {
		nmSetWwanEnabled(enabled)
	} else {
		// if NetworkingEnabled is off, turn it on, and only keep
		// current global device switch alive
		sh.config.setLastGlobalSwithes(false)
		sh.config.setLastWwanEnabled(true)
		sh.setNetworkingEnabled(true)
	}
}
func (sh *switchHandler) setWiredEnabled(enabled bool) {
	if sh.NetworkingEnabled {
		sh.setPropWiredEnabled(enabled)
	} else {
		// if NetworkingEnabled is off, turn it on, and only keep
		// current global device switch alive
		sh.config.setLastGlobalSwithes(false)
		sh.config.setLastWiredEnabled(true)
		sh.setNetworkingEnabled(true)
	}
}
func (sh *switchHandler) setVpnEnabled(enabled bool) {
	if sh.NetworkingEnabled {
		sh.setPropVpnEnabled(enabled)
	} else {
		// if NetworkingEnabled is off, turn it on, and only keep
		// current global device switch alive
		sh.config.setLastGlobalSwithes(false)
		sh.config.setLastVpnEnabled(true)
		sh.setNetworkingEnabled(true)
	}
}

func (sh *switchHandler) initPropNetworkingEnabled() {
	sh.NetworkingEnabled = nmManager.NetworkingEnabled.Get()
	if !sh.NetworkingEnabled {
		sh.doTurnOffGlobalDeviceSwitches()
	}
	sh.NetworkingEnabled = nmManager.NetworkingEnabled.Get()
	manager.setPropNetworkingEnabled(sh.NetworkingEnabled)
}
func (sh *switchHandler) setPropNetworkingEnabled() {
	if sh.NetworkingEnabled == nmManager.NetworkingEnabled.Get() {
		return
	}
	sh.NetworkingEnabled = nmManager.NetworkingEnabled.Get()
	// setup global device switches
	if sh.NetworkingEnabled {
		sh.restoreGlobalDeviceSwitches()
	} else {
		sh.saveAndTurnOffGlobalDeviceSwitches()
	}
	sh.NetworkingEnabled = nmManager.NetworkingEnabled.Get()
	manager.setPropNetworkingEnabled(sh.NetworkingEnabled)
}
func (sh *switchHandler) restoreGlobalDeviceSwitches() {
	nmSetWirelessEnabled(sh.config.getLastWirelessEnabled())
	nmSetWwanEnabled(sh.config.getLastWwanEnabled())
	sh.setPropWiredEnabled(sh.config.getLastWiredEnabled())
	sh.setPropVpnEnabled(sh.config.getLastVpnEnabled())
}
func (sh *switchHandler) saveAndTurnOffGlobalDeviceSwitches() {
	sh.config.setLastWirelessEnabled(sh.WirelessEnabled)
	sh.config.setLastWwanEnabled(sh.WwanEnabled)
	sh.config.setLastWiredEnabled(sh.WiredEnabled)
	sh.config.setLastVpnEnabled(sh.VpnEnabled)
	sh.doTurnOffGlobalDeviceSwitches()
}
func (sh *switchHandler) doTurnOffGlobalDeviceSwitches() {
	nmSetWirelessEnabled(false)
	nmSetWwanEnabled(false)
	sh.setPropWiredEnabled(false)
	sh.setPropVpnEnabled(false)
}

func (sh *switchHandler) initPropWirelessEnabled() {
	sh.WirelessEnabled = nmManager.WirelessEnabled.Get()
	for _, devPath := range nmGetDevicesByType(NM_DEVICE_TYPE_WIFI) {
		if sh.WirelessEnabled {
			sh.doEnableDevice(devPath, sh.config.getDeviceEnabled(devPath))
		} else {
			sh.doEnableDevice(devPath, false)
		}
	}
	sh.WirelessEnabled = nmManager.WirelessEnabled.Get()
	manager.setPropWirelessEnabled(sh.WirelessEnabled)
}
func (sh *switchHandler) setPropWirelessEnabled() {
	if sh.WirelessEnabled == nmManager.WirelessEnabled.Get() {
		return
	}
	sh.WirelessEnabled = nmManager.WirelessEnabled.Get()
	logger.Debug("setPropWirelessEnabled", sh.WirelessEnabled)
	// setup wireless devices switches
	for _, devPath := range nmGetDevicesByType(NM_DEVICE_TYPE_WIFI) {
		if sh.WirelessEnabled {
			sh.restoreDeviceState(devPath)
		} else {
			sh.saveAndDisconnectDevice(devPath)
		}
	}
	sh.WirelessEnabled = nmManager.WirelessEnabled.Get()
	manager.setPropWirelessEnabled(sh.WirelessEnabled)
}

func (sh *switchHandler) initPropWwanEnabled() {
	sh.WwanEnabled = nmManager.WwanEnabled.Get()
	for _, devPath := range nmGetDevicesByType(NM_DEVICE_TYPE_MODEM) {
		if sh.WwanEnabled {
			sh.doEnableDevice(devPath, sh.config.getDeviceEnabled(devPath))
		} else {
			sh.doEnableDevice(devPath, false)
		}
	}
	sh.WwanEnabled = nmManager.WwanEnabled.Get()
	manager.setPropWwanEnabled(sh.WwanEnabled)
}
func (sh *switchHandler) setPropWwanEnabled() {
	if sh.WwanEnabled == nmManager.WwanEnabled.Get() {
		return
	}
	sh.WwanEnabled = nmManager.WwanEnabled.Get()
	// setup modem devices switches
	for _, devPath := range nmGetDevicesByType(NM_DEVICE_TYPE_MODEM) {
		if sh.WwanEnabled {
			sh.restoreDeviceState(devPath)
		} else {
			sh.saveAndDisconnectDevice(devPath)
		}
	}
	sh.WwanEnabled = nmManager.WwanEnabled.Get()
	manager.setPropWwanEnabled(sh.WwanEnabled)
}

func (sh *switchHandler) initPropWiredEnabled() {
	sh.WiredEnabled = sh.config.getWiredEnabled()
	for _, devPath := range nmGetDevicesByType(NM_DEVICE_TYPE_ETHERNET) {
		if sh.WiredEnabled {
			sh.doEnableDevice(devPath, sh.config.getDeviceEnabled(devPath))
		} else {
			sh.doEnableDevice(devPath, false)
		}
	}
	manager.setPropWiredEnabled(sh.WiredEnabled)
}
func (sh *switchHandler) setPropWiredEnabled(enabled bool) {
	if sh.config.WiredEnabled == enabled {
		return
	}
	logger.Debug("setPropWiredEnabled", enabled)
	sh.WiredEnabled = enabled
	sh.config.setWiredEnabled(enabled)
	// setup wired devices switches
	for _, devPath := range nmGetDevicesByType(NM_DEVICE_TYPE_ETHERNET) {
		if enabled {
			sh.restoreDeviceState(devPath)
		} else {
			sh.saveAndDisconnectDevice(devPath)
		}
	}
	manager.setPropWiredEnabled(sh.WiredEnabled)
}

func (sh *switchHandler) initPropVpnEnabled() {
	sh.doEnableVpn(sh.config.getVpnEnabled())
}
func (sh *switchHandler) setPropVpnEnabled(enabled bool) {
	if sh.config.getVpnEnabled() == enabled {
		return
	}
	sh.doEnableVpn(enabled)
}
func (sh *switchHandler) doEnableVpn(enabled bool) {
	sh.VpnEnabled = enabled
	sh.config.setVpnEnabled(enabled)
	// setup vpn connections
	for _, uuid := range nmGetConnectionUuidsByType(NM_SETTING_VPN_SETTING_NAME) {
		if enabled {
			sh.restoreVpnConnectionState(uuid)
			sh.turnOffVpnSwitchIfNeed(autoConnectTimeout * 2)
		} else {
			sh.deactivateVpnConnection(uuid)
		}
	}
	manager.setPropVpnEnabled(sh.VpnEnabled)
}

func (sh *switchHandler) turnOffVpnSwitchIfNeed(timeout int) {
	// turn off vpn switch if all connections disconnected
	go func() {
		time.Sleep(time.Duration(timeout) * time.Second)
		vpnConnected := false
		for _, apath := range nmGetActiveConnections() {
			if nmGetActiveConnectionVpn(apath) {
				vpnConnected = true
				break
			}
		}
		if !vpnConnected {
			sh.setPropVpnEnabled(false)
		}
	}()
}

func (sh *switchHandler) initDeviceState(devPath dbus.ObjectPath) (err error) {
	err = sh.doEnableDevice(devPath, sh.config.getDeviceEnabled(devPath))
	return
}
func (sh *switchHandler) restoreDeviceState(devPath dbus.ObjectPath) (err error) {
	sh.config.restoreDeviceState(devPath)
	err = sh.doEnableDevice(devPath, sh.config.getDeviceEnabled(devPath))
	return
}
func (sh *switchHandler) saveAndDisconnectDevice(devPath dbus.ObjectPath) (err error) {
	sh.config.saveDeviceState(devPath)
	err = sh.doEnableDevice(devPath, false)
	return
}

func (sh *switchHandler) isDeviceEnabled(devPath dbus.ObjectPath) (enabled bool) {
	return sh.config.getDeviceEnabled(devPath)
}

func (sh *switchHandler) enableDevice(devPath dbus.ObjectPath, enabled bool) (err error) {
	if nmGetDeviceType(devPath) == NM_DEVICE_TYPE_WIFI {
		if !nmGetWirelessHardwareEnabled() {
			notifyWirelessHardSwitchOff()
			return
		}
	}
	return sh.doEnableDevice(devPath, enabled)
}
func (sh *switchHandler) doEnableDevice(devPath dbus.ObjectPath, enabled bool) (err error) {
	if enabled && sh.trunOnGlobalDeviceSwitchIfNeed(devPath) {
		return
	}
	devConfig, err := sh.config.getDeviceConfigForPath(devPath)
	if err != nil {
		return
	}
	logger.Debugf("doEnableDevice %s %v %#v", devPath, enabled, devConfig)

	sh.config.setDeviceEnabled(devPath, enabled)
	if enabled {
		// active last connection if device is disconnected
		if len(devConfig.LastConnectionUuid) > 0 {
			activeUuid, _ := nmGetDeviceActiveConnectionUuid(devPath)
			if devConfig.LastConnectionUuid != activeUuid {
				nmRunOnceUntilDeviceAvailable(devPath, func() {
					manager.ActivateConnection(devConfig.LastConnectionUuid, devPath)
				})
			}
		}
		uuids := nmGetConnectionUuidsForAutoConnect(devPath, devConfig.LastConnectionUuid)
		if len(uuids) > 0 {
			// TODO:
		}
	} else {
		err = manager.doDisconnectDevice(devPath)
	}
	return
}

func (sh *switchHandler) restoreVpnConnectionState(uuid string) (err error) {
	vpnConfig, err := sh.config.getVpnConfig(uuid)
	if err != nil {
		return
	}
	if vpnConfig.lastActivated || vpnConfig.AutoConnect {
		sh.activateVpnConnection(uuid)
	} else {
		err = manager.DeactivateConnection(uuid)
	}
	return
}
func (sh *switchHandler) activateVpnConnection(uuid string) {
	if _, err := nmGetActiveConnectionByUuid(uuid); err == nil {
		// connection already activated
		return
	}
	nmRunOnceUtilNetworkAvailable(func() {
		manager.ActivateConnection(uuid, "/")
	})
}
func (sh *switchHandler) deactivateVpnConnection(uuid string) (err error) {
	vpnConfig, err := sh.config.getVpnConfig(uuid)
	if err != nil {
		return
	}
	vpnConfig.lastActivated = vpnConfig.activated
	err = manager.DeactivateConnection(uuid)
	sh.config.save()
	return
}

func (sh *switchHandler) trunOnGlobalDeviceSwitchIfNeed(devPath dbus.ObjectPath) (need bool) {
	// if global device switch is off, turn it on, and only keep
	// current device alive
	need = (sh.generalGetGlobalDeviceEnabled(devPath) == false)
	if !need {
		return
	}
	sh.config.setAllDeviceLastEnabled(false)
	sh.config.setDeviceLastEnabled(devPath, true)
	sh.generalSetGlobalDeviceEnabled(devPath, true)
	return
}

func (sh *switchHandler) generalGetGlobalDeviceEnabled(devPath dbus.ObjectPath) (enabled bool) {
	switch devType := nmGetDeviceType(devPath); devType {
	case NM_DEVICE_TYPE_ETHERNET:
		enabled = sh.WiredEnabled
	case NM_DEVICE_TYPE_WIFI:
		enabled = sh.WirelessEnabled
	case NM_DEVICE_TYPE_MODEM:
		enabled = sh.WwanEnabled
	}
	return
}
func (sh *switchHandler) generalSetGlobalDeviceEnabled(devPath dbus.ObjectPath, enabled bool) {
	switch devType := nmGetDeviceType(devPath); devType {
	case NM_DEVICE_TYPE_ETHERNET:
		sh.setWiredEnabled(enabled)
	case NM_DEVICE_TYPE_WIFI:
		sh.setWirelessEnabled(enabled)
	case NM_DEVICE_TYPE_MODEM:
		sh.setWwanEnabled(enabled)
	}
}