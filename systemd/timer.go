package systemd

func IsTimerEnabled() bool {
	return isEnabled("paretosecurity-user.timer") && isEnabled("paretosecurity-user.service")
}

func EnableTimer() error {
	if err := enable("paretosecurity-user.timer"); err != nil {
		return err
	}
	return enable("paretosecurity-user.service")
}

func DisableTimer() error {
	if err := disable("paretosecurity-user.timer"); err != nil {
		return err
	}
	return disable("paretosecurity-user.service")
}
