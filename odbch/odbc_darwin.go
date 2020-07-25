package odbch

import "github.com/pkg/errors"

func getDrivers() ([]string, error) {
	return []string{}, errors.New("no linux support")
}

func InstalledDrivers() ([]string, error) {
	return []string{}, errors.New("no linux support")
}

func BestDriver() (string, error) {
	return "", errors.New("no linux support")
}

func ValidDriver(d string) error {
	return errors.New("no linux support")
}
