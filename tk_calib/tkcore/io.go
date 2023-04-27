package tkcore

import (
	"mytools/fio"
	"strconv"
)

func (tk *TKData) WriteACC(path string) error {
	writeData := make([][]string, 0)
	writeData = append(writeData, []string{"时间戳", "ACC_X", "ACC_Y", "ACC_Z"})
	timeData := tk.convTimestamp()
	accData := tk.convAccData()
	for i := range timeData {
		writeData = append(writeData, []string{timeData[i], accData[0][i], accData[1][i], accData[2][i]})
	}
	err := fio.WriteCSV(path, writeData)
	if err != nil {
		return err
	}
	return nil
}

func (tk *TKData) WriteGyro(path string) error {
	writeData := make([][]string, 0)
	writeData = append(writeData, []string{"时间戳", "Gyro_X", "Gyro_Y", "Gyro_Z"})
	timeData := tk.convTimestamp()
	gyroData := tk.convGyroData()
	for i := range timeData {
		writeData = append(writeData, []string{timeData[i], gyroData[0][i], gyroData[1][i], gyroData[2][i]})
	}
	err := fio.WriteCSV(path, writeData)
	if err != nil {
		return err
	}
	return nil
}

func (tk *TKData) Write(path string) error {
	writeData := make([][]string, 0)
	writeData = append(writeData, []string{"时间戳", "ACC_X", "ACC_Y", "ACC_Z", "Gyro_X", "Gyro_Y", "Gyro_Z"})
	timeData := tk.convTimestamp()
	accData := tk.convAccData()
	gyroData := tk.convGyroData()
	for i := range timeData {
		writeData = append(writeData, []string{
			timeData[i],
			accData[0][i],
			accData[1][i],
			accData[2][i],
			gyroData[0][i],
			gyroData[1][i],
			gyroData[2][i],
		})
	}
	err := fio.WriteCSV(path, writeData)
	if err != nil {
		return err
	}
	return nil
}

func (tk *TKData) convTimestamp() []string {
	result := make([]string, 0)
	for _, v := range tk.time {
		t := strconv.FormatInt(v, 10)
		result = append(result, t)
	}
	return result
}

func (tk *TKData) convAccData() [3][]string {
	var result [3][]string = [3][]string{}
	for i := 0; i < 3; i++ {
		for j := 0; j < tk.Size(); j++ {
			a := strconv.FormatInt(tk.acc[i][j], 10)
			result[i] = append(result[i], a)
		}
	}
	return result
}

func (tk *TKData) convGyroData() [3][]string {
	var result [3][]string = [3][]string{}
	for i := 0; i < 3; i++ {
		for j := 0; j < tk.Size(); j++ {
			a := strconv.FormatInt(tk.gyro[i][j], 10)
			result[i] = append(result[i], a)
		}
	}
	return result
}
