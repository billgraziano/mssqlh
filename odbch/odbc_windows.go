package odbch

import (
	"sort"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/sys/windows/registry"
)

// getDrivers returns the ODBC drivers from the Windows registery
func getDrivers() ([]string, error) {

	k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\ODBC\ODBCINST.INI\ODBC Drivers`, registry.QUERY_VALUE)
	if err != nil {
		return nil, errors.Wrap(err, "openkey")
	}
	defer k.Close()

	s, err := k.ReadValueNames(0)
	if err != nil {
		return nil, errors.Wrap(err, "readvaluenames")
	}

	sort.Strings(s)

	return s, nil
}

// InstalledDrivers returns the available SQL Server drivers on the computer
func InstalledDrivers() ([]string, error) {
	var drivers []string

	d, err := getDrivers()
	if err != nil {
		return drivers, errors.Wrap(err, "getdrivers")
	}

	for _, v := range d {
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

// ValidDriver tests if a string is a valid SQL Server Driver on this machine
func ValidDriver(d string) error {
	d = strings.TrimSpace(d)
	drivers, err := InstalledDrivers()
	if err != nil {
		return errors.Wrap(err, "availabledrivers")
	}

	for _, v := range drivers {
		if strings.EqualFold(v, d) {
			return nil
		}
	}
	return ErrInvalidDriver
}
