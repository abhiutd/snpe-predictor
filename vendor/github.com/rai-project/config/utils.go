package config

import "os"

func modeInfo() (isDebug bool, isVerbose bool) {
	if IsDebug {
		isDebug = IsDebug
	} else {
		isDebug = false
		if v := os.Getenv("DEBUG"); v == "1" || v == "TRUE" {
			isDebug = true
		}
	}
	if IsVerbose {
		isVerbose = IsVerbose
	} else {
		isVerbose = false
		if v := os.Getenv("VERBOSE"); v == "1" || v == "TRUE" {
			isVerbose = true
		}
		IsDebug = isDebug
		IsVerbose = isVerbose
	}
	return
}
