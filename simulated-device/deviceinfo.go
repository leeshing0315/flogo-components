package main

import "sync"

var electricalCommunicationFrequency string

var noElectricityCommunicationFrequency string

var noElectricityAcquisitionFrequency string

var settingIpAddress string

var settingIpPort string

var forcedNotToSleep string

var rwMutex sync.RWMutex

func setDeviceInfo(arg0, arg1, arg2, arg3, arg4, arg5 string) {
	rwMutex.Lock()
	electricalCommunicationFrequency = arg0
	noElectricityCommunicationFrequency = arg1
	noElectricityAcquisitionFrequency = arg2
	settingIpAddress = arg3
	settingIpPort = arg4
	forcedNotToSleep = arg5
	rwMutex.Unlock()
}

func getDeviceInfo() (string, string, string, string, string, string) {
	rwMutex.RLock()
	defer rwMutex.RUnlock()
	return electricalCommunicationFrequency, noElectricityCommunicationFrequency, noElectricityAcquisitionFrequency, settingIpAddress, settingIpPort, forcedNotToSleep
}
