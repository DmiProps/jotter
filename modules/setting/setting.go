package setting

import (
	// System packages.
	"encoding/json"
	"os"
	"strings"

	// Jotter packages.
	"github.com/dmiprops/jotter/modules/auth"
	"github.com/dmiprops/jotter/modules/log"
)

// Application constants.
const (
	DefaultPort         string = "3030"
	Localhost           string = "127.0.0.1"
	SettingsDir         string = "/etc/jotter.d"
	StoredSettingsFile  string = "jotter.json"
	CurrentSettingsFile string = "session.json"
)

// Application settings.
var (
	AppVer   string
	Protocol string

	StoredAdminSettings  storedAdminSettingsType
	CurrentAdminSettings currentAdminSettingsType
)

type storedAdminSettingsType struct {
	Password string `json:"password"`
	Address  string `json:"address"`
	Database string `json:"database"`
}

type currentAdminSettingsType struct {
	Address  string `json:"address"`
	Database string `json:"database"`
}

func init() {
	err := InitSettings()
	if err != nil {
		log.Fatal(err.Error())
	}
}

// InitSettings initializas administrative service options.
func InitSettings() error {
	_, err := os.Stat(SettingsDir)
	if err == nil {
		_, err := os.Stat(SettingsDir + "/" + StoredSettingsFile)
		if err == nil {
			return ReadStoredAdminSettings()
		}
		return createDefaultAdminSettings()
	}
	err = os.Mkdir(SettingsDir, 0777)
	if err == nil {
		return createDefaultAdminSettings()
	}
	return err
}

func createDefaultAdminSettings() error {
	path := SettingsDir + "/" + StoredSettingsFile

	// Create default settings struct.
	hash, _ := auth.HashPassword("jotter")
	adminSettings := struct {
		Version       string                  `json:"version"`
		Configuration storedAdminSettingsType `json:"configuration"`
	}{
		Version: AppVer,
		Configuration: storedAdminSettingsType{
			Password: hash,
			Address:  ":" + DefaultPort,
		},
	}

	// Apply default settings.
	StoredAdminSettings.Password = adminSettings.Configuration.Password
	StoredAdminSettings.Address = adminSettings.Configuration.Address
	StoredAdminSettings.Database = adminSettings.Configuration.Database

	// Store defailt settings.
	f, err := os.Create(path)
	if err == nil {
		defer f.Close()

		encoder := json.NewEncoder(f)
		return encoder.Encode(&adminSettings)
	}
	return err
}

// ReadStoredAdminSettings read stored administrative settings.
func ReadStoredAdminSettings() error {
	path := SettingsDir + "/" + StoredSettingsFile

	adminSettings := struct {
		Version       string                  `json:"version"`
		Configuration storedAdminSettingsType `json:"configuration"`
	}{}

	f, err := os.Open(path)
	if err == nil {
		decoder := json.NewDecoder(f)
		err := decoder.Decode(&adminSettings)
		if err == nil {
			StoredAdminSettings.Password = adminSettings.Configuration.Password
			StoredAdminSettings.Address = adminSettings.Configuration.Address
			StoredAdminSettings.Database = adminSettings.Configuration.Database
		}
		return err
	}
	return err
}

// SaveStoredAdminSettings write administrative settings.
func SaveStoredAdminSettings() error {
	path := SettingsDir + "/" + StoredSettingsFile

	adminSettings := struct {
		Version       string                  `json:"version"`
		Configuration storedAdminSettingsType `json:"configuration"`
	}{
		Version: AppVer,
		Configuration: storedAdminSettingsType{
			Password: StoredAdminSettings.Password,
			Address:  StoredAdminSettings.Address,
			Database: StoredAdminSettings.Database,
		},
	}

	// Save stored administrative settings.
	f, err := os.Create(path)
	if err == nil {
		defer f.Close()

		encoder := json.NewEncoder(f)
		return encoder.Encode(&adminSettings)
	}
	return err
}

// ReadCurrentAdminSettings read current administrative settings.
func ReadCurrentAdminSettings() error {
	path := SettingsDir + "/" + CurrentSettingsFile

	adminSettings := struct {
		Version       string                   `json:"version"`
		Configuration currentAdminSettingsType `json:"configuration"`
	}{}

	f, err := os.Open(path)
	if err == nil {
		decoder := json.NewDecoder(f)
		err := decoder.Decode(&adminSettings)
		if err == nil {
			CurrentAdminSettings.Address = adminSettings.Configuration.Address
			CurrentAdminSettings.Database = adminSettings.Configuration.Database
		}
		return err
	}
	return err
}

// SaveCurrentAdminSettings write administrative settings.
func SaveCurrentAdminSettings() error {
	path := SettingsDir + "/" + CurrentSettingsFile

	adminSettings := struct {
		Version       string                   `json:"version"`
		Configuration currentAdminSettingsType `json:"configuration"`
	}{
		Version: AppVer,
		Configuration: currentAdminSettingsType{
			Address:  StoredAdminSettings.Address,
			Database: StoredAdminSettings.Database,
		},
	}

	// Save current administrative settings.
	f, err := os.Create(path)
	if err == nil {
		defer f.Close()

		encoder := json.NewEncoder(f)
		return encoder.Encode(&adminSettings)
	}
	return err
}

// ConnectionStringWithoutPassword returns connection string to database without password.
func ConnectionStringWithoutPassword(database string) string {
	i1 := strings.Index(database, ":")
	if i1 < 0 {
		return database
	}
	i2 := strings.Index(database, "@")
	return string([]rune(database)[:i1]) + string([]rune(database)[i2+1:])
}
