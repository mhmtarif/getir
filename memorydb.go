package main

var memoryData map[string]string

func initMemoryDb() {
	memoryData = make(map[string]string, 15)
}

func getMemoryData(key string) string {
	return memoryData[key]
}

func setMemoryData(key, value string) {
	memoryData[key] = value
}
