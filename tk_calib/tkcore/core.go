package tkcore

import (
	"fmt"
	"math/rand"
)

// 时间戳
type TimeStramp []int64

const (
	Microsecond int = 1
	Millisecond int = 1000 * Microsecond
	Second      int = 1000 * Millisecond
)

// ACC数据
type AccData [3][]int64

// Gyro数据
type GyroData [3][]int64

// 代表tk的数据
type TKData struct {
	time TimeStramp
	acc  AccData
	gyro GyroData
}

func NewCustomData(t TimeStramp, a AccData, g GyroData) *TKData {
	tk := &TKData{t, a, g}
	if tk.dataVerify() {
		return tk
	}
	return &TKData{}
}

// 获取tk数据的长度
func (tk *TKData) Size() int {
	return len(tk.time)
}

func (tk *TKData) GetAccDataAll() AccData {
	return tk.acc
}

func (tk *TKData) GetTKRange(startIndex int, endIndex int) (*TKData, error) {
	size := tk.Size()
	result := &TKData{}
	if startIndex < 0 || endIndex > size-1 {
		return nil, fmt.Errorf("初始下标或结束下标选择错误")
	}
	if startIndex > endIndex {
		return nil, fmt.Errorf("初始下标大于结束下标")
	}
	t := tk.time[startIndex : endIndex+1]
	a := AccData{}
	a[0] = tk.GetAccDataAll()[0][startIndex : endIndex+1]
	a[1] = tk.GetAccDataAll()[1][startIndex : endIndex+1]
	a[2] = tk.GetAccDataAll()[2][startIndex : endIndex+1]
	g := GyroData{}
	g[0] = tk.GetGyroDataAll()[0][startIndex : endIndex+1]
	g[1] = tk.GetGyroDataAll()[1][startIndex : endIndex+1]
	g[2] = tk.GetGyroDataAll()[2][startIndex : endIndex+1]
	result.time = t
	result.acc = a
	result.gyro = g
	return result, nil
}

func (tk *TKData) GetAccDataRange(startIndex int, endIndex int) (AccData, error) {
	size := tk.Size() - 1
	var result AccData = AccData{}
	if startIndex < 0 || endIndex < 0 {
		return result, fmt.Errorf("初始下标或结束下标小于0")
	}
	if startIndex > size || endIndex > size {
		return result, fmt.Errorf("初始下标或结束下标超过范围")
	}
	all := tk.GetAccDataAll()
	result[0] = all[0][startIndex : endIndex+1]
	result[1] = all[1][startIndex : endIndex+1]
	result[2] = all[2][startIndex : endIndex+1]
	return result, nil
}

func (tk *TKData) GetTimestampForIndex(index int) int64 {
	return tk.time[index]
}

// 根据时间戳范围来获取tk数据
func (tk *TKData) GetTimeRangeData(startTime int64, offset int64) (*TKData, error) {
	startIndex := tk.findTimeStamp(startTime)
	endIndex := tk.findTimeStamp(startTime + offset)
	data, err := tk.GetTKRange(startIndex, endIndex)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (tk *TKData) findTimeStamp(t int64) int {
	index := -1
	for i, v := range tk.time {
		if v >= t {
			return i
		}
	}
	return index
}

func (tk *TKData) GetGyroDataAll() GyroData {
	return tk.gyro
}

// GetGyroDataRange 获取一定范围内的Gyor数据
func (tk *TKData) GetGyroDataRange(startIndex int, endIndex int) (GyroData, error) {
	size := tk.Size()
	var result GyroData = GyroData{}
	if startIndex < 0 || endIndex < 0 {
		return result, fmt.Errorf("初始下标或结束下标小于0")
	}
	if startIndex > size || endIndex > size {
		return result, fmt.Errorf("初始下标或结束下标超过范围")
	}
	all := tk.GetGyroDataAll()
	result[0] = all[0][startIndex : endIndex+1]
	result[1] = all[1][startIndex : endIndex+1]
	result[2] = all[2][startIndex : endIndex+1]
	return result, nil
}

func (tk *TKData) GetTimeAll() TimeStramp {
	return tk.time
}

func (tk *TKData) dataVerify() bool {
	tSize := len(tk.time)
	ax := tk.acc[0]
	ay := tk.acc[1]
	az := tk.acc[2]
	if tSize != len(ax) || tSize != len(ay) || tSize != len(az) {
		return false
	}
	gx := tk.gyro[0]
	gy := tk.gyro[1]
	gz := tk.gyro[2]
	if tSize != len(gx) || tSize != len(gy) || tSize != len(gz) {
		return false
	}
	return true
}

func MakeTestDataAndWrite() *TKData {
	t := make(TimeStramp, 0)

	var a AccData = AccData{}
	var g GyroData = GyroData{}
	for i := 0; i < 100; i++ {
		t = append(t, int64(i))
	}
	for i := 0; i < 3; i++ {
		for j := 0; j < 100; j++ {
			a[i] = append(a[i], rand.Int63())
			g[i] = append(g[i], rand.Int63())
		}
	}
	tkData := NewCustomData(t, a, g)
	tkData.Write(".\\TestData.csv")
	return tkData
}
