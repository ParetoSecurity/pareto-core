package cmd

import (
	"github.com/godbus/dbus/v5"
)

func Notify(message string) error {
	conn, err := dbus.SessionBus()
	if err != nil {
		return err
	}
	defer conn.Close()

	obj := conn.Object("org.freedesktop.Notifications", "/org/freedesktop/Notifications")
	call := obj.Call("org.freedesktop.Notifications.Notify", 0,
		"ParetoSecurity",         // Application name
		uint32(0),                // Replace ID
		"",                       // Icon (empty for default)
		"Pareto Linux",           // Summary
		message,                  // Body
		[]string{},               // Actions
		map[string]interface{}{}, // Hints
		int32(5000))              // Timeout (5 seconds)

	return call.Err
}
