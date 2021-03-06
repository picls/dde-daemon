/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package sessionwatcher

import (
	libdisplay "dbus/com/deepin/daemon/display"
	"dbus/org/freedesktop/login1"
	"pkg.deepin.io/lib/dbus"
	"sync"
)

const (
	login1Dest         = "org.freedesktop.login1"
	login1Path         = "/org/freedesktop/login1"
	displayDBusDest    = "com.deepin.daemon.Display"
	displayDBusObjPath = "/com/deepin/daemon/Display"
)

type Manager struct {
	display       *libdisplay.Display
	loginManager  *login1.Manager
	sessionLocker sync.Mutex
	sessions      map[string]*login1.Session
	IsActive      bool
}

func newManager() (*Manager, error) {
	manager := &Manager{
		sessions: make(map[string]*login1.Session),
	}
	var err error
	manager.loginManager, err = login1.NewManager(login1Dest, login1Path)
	if err != nil {
		logger.Warning("New login1 manager failed:", err)
		return nil, err
	}

	manager.display, err = libdisplay.NewDisplay(displayDBusDest, displayDBusObjPath)
	if err != nil {
		logger.Warning(err)
		return nil, err
	}

	return manager, nil
}

func (m *Manager) destroy() {
	if m.sessions != nil {
		m.destroySessions()
	}

	if m.display != nil {
		libdisplay.DestroyDisplay(m.display)
		m.display = nil
	}

	if m.loginManager != nil {
		login1.DestroyManager(m.loginManager)
		m.loginManager = nil
	}
}

func (m *Manager) initUserSessions() {
	list, err := m.loginManager.ListSessions()
	if err != nil {
		logger.Warning("List sessions failed:", err)
		return
	}

	for _, v := range list {
		// v info: (id, uid, username, seat id, session path)
		if len(v) != 5 {
			logger.Warning("Invalid session info:", v)
			continue
		}

		id, ok := v[0].(string)
		if !ok {
			continue
		}

		p, ok := v[4].(dbus.ObjectPath)
		if !ok {
			continue
		}

		m.addSession(id, p)
	}

	m.loginManager.ConnectSessionNew(func(id string, path dbus.ObjectPath) {
		logger.Debug("Session added:", id, path)
		m.addSession(id, path)
	})

	m.loginManager.ConnectSessionRemoved(func(id string, path dbus.ObjectPath) {
		logger.Debug("Session removed:", id, path)
		m.deleteSession(id, path)
	})
}

func (m *Manager) destroySessions() {
	m.sessionLocker.Lock()
	for _, s := range m.sessions {
		login1.DestroySession(s)
		s = nil
	}
	m.sessions = nil
	m.sessionLocker.Unlock()
}

func (m *Manager) addSession(id string, path dbus.ObjectPath) {
	uid, session := newLoginSession(path)
	if session == nil {
		return
	}

	logger.Debug("Add session:", id, path, uid)
	if !isCurrentUser(uid) {
		logger.Debug("Not the current user session:", id, path, uid)
		return
	}

	m.sessionLocker.Lock()
	m.sessions[id] = session
	m.sessionLocker.Unlock()

	session.Active.ConnectChanged(func() {
		if session == nil {
			return
		}
		m.handleSessionChanged()
	})
	m.handleSessionChanged()
}

func (m *Manager) deleteSession(id string, path dbus.ObjectPath) {
	m.sessionLocker.Lock()
	session, ok := m.sessions[id]
	if !ok {
		m.sessionLocker.Unlock()
		return
	}

	logger.Debug("Delete session:", id, path)
	login1.DestroySession(session)
	session = nil
	delete(m.sessions, id)
	m.sessionLocker.Unlock()
	m.handleSessionChanged()
}

func (m *Manager) handleSessionChanged() {
	if len(m.sessions) == 0 {
		return
	}

	isActive := m.checkIsActive()
	changed := m.setIsActive(isActive)
	if !changed {
		return
	}

	if isActive {
		logger.Debug("[handleSessionChanged] Resume pulse")
		suspendPulseSinks(0)
		suspendPulseSources(0)

		logger.Debug("[handleSessionChanged] Reset Brightness")
		m.display.ResetChanges()
	} else {
		logger.Debug("[handleSessionChanged] Suspend pulse")
		suspendPulseSinks(1)
		suspendPulseSources(1)
	}
}

// return is changed?
func (m *Manager) setIsActive(val bool) bool {
	if m.IsActive != val {
		m.IsActive = val
		logger.Debug("[setIsActive] IsActive changed:", val)
		dbus.NotifyChange(m, "IsActive")
		return true
	}
	return false
}

func (m *Manager) checkIsActive() bool {
	var active bool = false
	m.sessionLocker.Lock()
	for _, session := range m.sessions {
		isSessionActive := session.Active.Get()
		logger.Debugf("[checkIsActive] session path: %q,isActive: %v", session.Path, isSessionActive)
		if isSessionActive {
			active = true
			break
		}
	}
	m.sessionLocker.Unlock()

	logger.Debugf("[checkIsActive] result user: %v isActive: %v", curUid, active)
	return active
}
