package odbch

import (
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"gopkg.in/ini.v1"
)

// getDrivers returns the installed drivers from odbcinst.ini
func getDrivers() ([]string, error) {
	var ss []string
	locations := []string{"/usr/local/etc/odbcinst.ini", "/etc/odbcinst.ini"}
	iniFile := firstfile(locations)
	if iniFile == "" {
		return ss, errors.New("odbcinst.ini not found")
	}
	bb, err := afero.ReadFile(fs, iniFile)
	if err != nil {
		return ss, errors.Wrap(err, "readfile")
	}
	cfg, err := ini.Load(bb)
	if err != nil {
		return ss, errors.Wrap(err, "ini.load")
	}
	for _, sect := range cfg.Sections() {
		ss = append(ss, sect.Name())
	}
	return ss, nil
}

// InstalledDrivers returns the available SQL Server drivers on the computer
func InstalledDrivers() ([]string, error) {
	var drivers []string

	dd, err := getDrivers()
	if err != nil {
		return drivers, errors.Wrap(err, "getdrivers")
	}

	for _, v := range dd {
		for _, d := range orderedDrivers {
			if strings.EqualFold(d, v) {
				drivers = append(drivers, v)
			}
		}
	}
	return drivers, nil
}

// BestDriver returns the "best" driver installed on the machine
func BestDriver() (string, error) {
	drivers, err := getDrivers()
	if err != nil {
		return "", errors.Wrap(err, "getDrivers")
	}

	for _, d := range orderedDrivers {
		for _, v := range drivers {
			if strings.EqualFold(d, v) {
				return d, nil
			}
		}
	}
	return "", ErrNoDrivers
}

func ValidDriver(d string) error {
	drivers, err := InstalledDrivers()
	if err != nil {
		return errors.Wrap(err, "availabledrivers")
	}

	for _, v := range drivers {
		if v == d {
			return nil
		}
	}
	return ErrInvalidDriver
}

func firstfile(files []string) string {
	for _, f := range files {
		_, err := fs.Stat(f)
		if err == nil {
			return f
		}
	}
	return ""
}
