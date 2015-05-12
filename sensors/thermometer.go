package sensors

import (
	"io/ioutil"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const (
	sensorBaseDirectory = "/sys/bus/w1/devices/"
)

var (
	sensorPaths = []string{}
)

func init() {
	command := exec.Command("modprobe", "w1-gpio")
	command.Run()
	command = exec.Command("modprobe", "w1-therm")
	command.Run()

	dirs, err := ioutil.ReadDir(sensorBaseDirectory)
	if err != nil {
		panic(err)
	}

	for _, dir := range dirs {
		sensorPaths = append(sensorPaths, sensorBaseDirectory+dir.Name()+"/w1_slave")
	}
}

func getThermometerReading(sensor string) (float64, float64) {

	stringContents := ""

	for goodReading := false; ; {
		contents, err := ioutil.ReadFile(sensor)
		if err != nil {
			panic(err)
		}

		stringContents = string(contents[:])
		goodReading = strings.Contains(stringContents, "YES")

		if goodReading {
			break
		}
		time.Sleep(time.Millisecond * 200)
	}

	tempAvailable := strings.Contains(stringContents, "t=")

	if tempAvailable {
		panic("Could not find 't=' inside sensor file.")
	}

	reading := strings.SplitAfter(stringContents, "t=")
	raw, err := strconv.ParseFloat(reading[1], 64)
	if err != nil {
		panic(err)
	}
	celsius := raw / 1000.0
	fahrenheit := celsius*9.0/5.0 + 32.0

	return celsius, fahrenheit
}

func GetThermometerReadings() []float64 {
	temperatures := []float64{0, 0}

	// for _, sensor := range sensorPaths {
	// 	_, f := getThermometerReading(sensor)
	// 	temperatures = append(temperatures, f)
	// }

	return temperatures
}
