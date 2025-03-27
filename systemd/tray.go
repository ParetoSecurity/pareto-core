package systemd

func IsTrayIconEnabled() bool {
	return isEnabled("paretosecurity-trayicon.service")
}

func EnableTrayIcon() error {
	return enable("paretosecurity-trayicon.service")
}

func DisableTrayIcon() error {
	return disable("paretosecurity-trayicon.service")
}
