package CORAN_LIB_LOGGER_CONF

import (
	"log"
	"os"
	"strconv"

	path_util "github.com/coranlabs/CORAN_LIB_PATH_UTIL"
)

var CoranlabsLogDir string = path_util.CoranlabsPath("coranlabs/log") + "/"
var LibLogDir string = CoranlabsLogDir + "lib/"
var NfLogDir string = CoranlabsLogDir + "nf/"

var CoranlabsLogFile string = CoranlabsLogDir + "coranlabs.log"

func init() {
	if err := os.MkdirAll(LibLogDir, 0775); err != nil {
		log.Printf("Mkdir %s failed: %+v", LibLogDir, err)
	}
	if err := os.MkdirAll(NfLogDir, 0775); err != nil {
		log.Printf("Mkdir %s failed: %+v", NfLogDir, err)
	}

	// Create log file or if it already exist, check if user can access it
	f, fileOpenErr := os.OpenFile(CoranlabsLogFile, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if fileOpenErr != nil {
		// user cannot access it.
		log.Printf("Cannot Open %s\n", CoranlabsLogFile)
	} else {
		// user can access it
		if err := f.Close(); err != nil {
			log.Printf("File %s cannot been closed\n", CoranlabsLogFile)
		}
	}

	sudoUID, errUID := strconv.Atoi(os.Getenv("SUDO_UID"))
	sudoGID, errGID := strconv.Atoi(os.Getenv("SUDO_GID"))

	if errUID == nil && errGID == nil {
		// if using sudo to run the program, errUID will be nil and sudoUID will get the uid who run sudo
		// else errUID will not be nil and sudoUID will be nil
		// If user using sudo to run the program and create log file, log will own by root,
		// here we change own to user so user can view and reuse the file
		if err := os.Chown(CoranlabsLogDir, sudoUID, sudoGID); err != nil {
			log.Printf("Dir %s chown to %d:%d error: %v\n", CoranlabsLogDir, sudoUID, sudoGID, err)
		}
		if err := os.Chown(LibLogDir, sudoUID, sudoGID); err != nil {
			log.Printf("Dir %s chown to %d:%d error: %v\n", LibLogDir, sudoUID, sudoGID, err)
		}
		if err := os.Chown(NfLogDir, sudoUID, sudoGID); err != nil {
			log.Printf("Dir %s chown to %d:%d error: %v\n", NfLogDir, sudoUID, sudoGID, err)
		}

		if fileOpenErr == nil {
			if err := os.Chown(CoranlabsLogFile, sudoUID, sudoGID); err != nil {
				log.Printf("File %s chown to %d:%d error: %v\n", CoranlabsLogFile, sudoUID, sudoGID, err)
			}
		}
	}
}
