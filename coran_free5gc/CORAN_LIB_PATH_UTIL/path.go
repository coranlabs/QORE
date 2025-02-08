package path_util

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"

	"github.com/coranlabs/CORAN_LIB_PATH_UTIL/logger"
)

// CoranlabsPath ...
/*
 * Author: Roger Chu aka Sasuke
 *
 * This package is used to locate the root directory of gocoranlabs project
 * Compatible with Windows and Linux
 *
 * Please import "github.com/coranlabs/CORAN_LIB_PATH_UTIL"
 *
 * Return value:
 * A string value of the relative path between the working directory and the root directory of the gocoranlabs project
 *
 * Usage:
 * path_util.CoranlabsPath("your file location starting with gocoranlabs")
 *
 * Example:
 * path_util.CoranlabsPath("coranlabs/abcdef/abcdef.pem")
 */
func CoranlabsPath(path string) string {
	rootCode := strings.Split(path, "/")[0]
	cleanPath := filepath.Clean(path)
	targetFilePath := cleanPath[len(rootCode)+1:]

	var pwd string
	if pwdTmp, err := os.Getwd(); err != nil {
		logger.PathLog.Errorln(err)
	} else {
		pwd = pwdTmp
	}
	currentPath := filepath.Clean(pwd)

	// Module mode
	target := ""
	if returnPath, ok := FindModuleRoot(currentPath, rootCode); ok {
		target = returnPath + filepath.Clean("/"+targetFilePath)
	}

	// Non-module mode
	if target == "" {
		binPathDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			logger.PathLog.Errorln(err)
		}

		rootPath := ""
		if strings.Contains(currentPath, rootCode) {
			if returnPath, ok := FindRoot(currentPath, rootCode, targetFilePath); ok {
				rootPath = returnPath
			} else if returnPath, ok := FindRoot(currentPath, rootCode, "lib"); ok {
				rootPath = returnPath
			}
		}
		if strings.Contains(binPathDir, rootCode) {
			if returnPath, ok := FindRoot(binPathDir, rootCode, targetFilePath); ok {
				rootPath = returnPath
			} else if returnPath, ok := FindRoot(binPathDir, rootCode, "lib"); ok {
				rootPath = returnPath
			}
		}

		if rootPath != "" {
			target = rootPath + cleanPath
		} else {
			binPathDirParent := GetParentDirectory(binPathDir)
			binPathDirParentWithTargetFilePath := binPathDirParent + filepath.Clean("/"+targetFilePath)
			target = binPathDirParentWithTargetFilePath
		}
	}

	location, err := filepath.Rel(currentPath, target)
	if err != nil {
		logger.PathLog.Errorln(err)
	}

	return location
}

func Exists(fpath string) bool {
	_, err := os.Stat(fpath)
	return !os.IsNotExist(err)
}

func FindRoot(path string, rootCode string, objName string) (string, bool) {
	rootPath := path
	loc := strings.LastIndex(rootPath, rootCode)
	for loc != -1 {
		rootPath = rootPath[:loc+len(rootCode)]
		if Exists(rootPath + filepath.Clean("/"+objName)) {
			return rootPath[:loc], true
		}
		rootPath = rootPath[:loc]
		loc = strings.LastIndex(rootPath, rootCode)
	}
	return "", false
}

func FindModuleRoot(path string, rootCode string) (string, bool) {
	moduleFilePath := path + filepath.Clean("/go.mod")
	if Exists(moduleFilePath) {
		var file *os.File
		if fileTmp, err := os.Open(moduleFilePath); err != nil {
			logger.PathLog.Fatalf("Cannot open %s: %+v", moduleFilePath, err)
		} else {
			file = fileTmp
		}
		defer func() {
			if err := file.Close(); err != nil {
				logger.PathLog.Warnf("File %s cannot close: %v", moduleFilePath, err)
			}
		}()

		reader := bufio.NewReader(file)
		moduleDeclearation, _, err := reader.ReadLine()
		if err != nil {
			logger.PathLog.Warnf("Read Line failed: %+v", err)
		}
		if string(moduleDeclearation) == "module "+rootCode {
			return path, true
		}
	}

	abs, err := filepath.Abs(path + string(filepath.Separator) + "..")
	if err != nil || abs == filepath.Clean("/") {
		return "", false
	}

	return FindModuleRoot(abs, rootCode)
}

func GetParentDirectory(dirctory string) string {
	return filepath.Dir(dirctory)
}
