package sensors

import (
	"io/ioutil"
	"os/exec"
	"strings"
	"time"
	"strconv"
)

const (
	sensorBaseDirectory = "/sys/bus/w1/devices/"
)

var (
	sensorPaths = []string{}
)

/*
os.system('modprobe w1-gpio')
os.system('modprobe w1-therm')
 
base_dir = '/sys/bus/w1/devices/'
device_folder = glob.glob(base_dir + '28*')[0]
device_file = device_folder + '/w1_slave'
 
def read_temp_raw():
    f = open(device_file, 'r')
    lines = f.readlines()
    f.close()
    return lines
 
def read_temp():
    lines = read_temp_raw()
    while lines[0].strip()[-3:] != 'YES':
        time.sleep(0.2)
        lines = read_temp_raw()
    equals_pos = lines[1].find('t=')
    if equals_pos != -1:
        temp_string = lines[1][equals_pos+2:]
        temp_c = float(temp_string) / 1000.0
        temp_f = temp_c * 9.0 / 5.0 + 32.0
        return temp_c, temp_f
	
while True:
	print(read_temp())	
	time.sleep(1)
*/


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
		sensorPaths = append(sensorPaths, sensorBaseDirectory + dir.Name() + "/w1_slave")
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
	fahrenheit := celsius * 9.0 / 5.0 + 32.0
	
	return celsius, fahrenheit
}

func GetThermometerReadings() []float64 {
	temperatures := []float64{0,0}

	// for _, sensor := range sensorPaths {
	// 	_, f := getThermometerReading(sensor)
	// 	temperatures = append(temperatures, f)
	// }

	return temperatures
}
