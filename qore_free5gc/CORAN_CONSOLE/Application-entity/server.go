package cmd

import (
	"os"
	"path/filepath"
	"runtime/debug"

	
	"github.com/urfave/cli"
    "github.com/coranlabs/CORAN_LIB_UTIL/version"
	"github.com/coranlabs/CORAN_CONSOLE/backend/factory"
	"github.com/coranlabs/CORAN_CONSOLE/backend/logger"
	"github.com/coranlabs/CORAN_CONSOLE/backend/console_service"
	logger_util "github.com/coranlabs/CORAN_LIB_UTIL/logger"
)

var WEBUI *console_service.WebuiApp

func main() {
	defer func() {
		if p := recover(); p != nil {
			// Print stack for panic to log. Fatalf() will let program exit.
			logger.MainLog.Fatalf("panic: %v\n%s", p, string(debug.Stack()))
		}
	}()
	app := cli.NewApp()
	app.Name = "coran_console"
	app.Usage = "Coran Web Console"
	app.Action = Action
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config, c",
			Usage: "Load configuration from `FILE`",
		},
		cli.StringSliceFlag{
			Name:  "log, l",
			Usage: "Output NF log to `FILE`",
		},
	}
	if err := app.Run(os.Args); err != nil {
		logger.MainLog.Errorf("console Run error: %v\n", err)
	}
}

func Action() error {
	

	logger.MainLog.Infoln("WEBUI version: ", version.GetVersion())

	cfg, err := factory.ReadConfig("./config/CORAN_CONSOLE.yaml")
	if err != nil {
		return err
	}
	factory.WebuiConfig = cfg

	webui, err := console_service.NewApp(cfg)
	if err != nil {
		return err
	}
	WEBUI = webui

	webui.Start()

	return nil
}

func initLogFile(logNfPath []string) (string, error) {
	logTlsKeyPath := ""

	for _, path := range logNfPath {
		if err := logger_util.LogFileHook(logger.Log, path); err != nil {
			return "", err
		}
		if logTlsKeyPath != "" {
			continue
		}

		nfDir, _ := filepath.Split(path)
		tmpDir := filepath.Join(nfDir, "key")
		if err := os.MkdirAll(tmpDir, 0o775); err != nil {
			logger.InitLog.Errorf("Make directory %s failed: %+v", tmpDir, err)

			return "", err
		}

		_, name := filepath.Split(factory.WebuiDefaultTLSKeyLogPath)
		logTlsKeyPath = filepath.Join(tmpDir, name)
	}

	return logTlsKeyPath, nil
}
