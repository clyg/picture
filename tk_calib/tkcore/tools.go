package tkcore

import (
	"bufio"
	"fmt"
	"os"
)

// 调用函数就等待
func WaitForKye() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("输入一个按键")
	reader.ReadByte()
}
