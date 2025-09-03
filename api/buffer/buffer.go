package buffer

import (
	"bytes"
	"encoding/json"
	"io"
)

// Buffer 连续缓冲区结构
// 用于管理多个字节切片的指针，实现连续数据存储
// 内部维护一个[]*[]byte切片，每个元素指向一个字节切片
type Buffer []*[]byte

// bufferCall 函数类型，实现了io.Reader接口
// 该类型用于创建一个可以连续读取多个缓冲区的读取器
type bufferCall func(p []byte) (n int, err error)

// Read 实现io.Reader接口的Read方法
// 参数:
//
//	p - 要读取数据的字节切片
//
// 返回:
//
//	n - 实际读取的字节数
//	err - 读取过程中遇到的错误
func (call bufferCall) Read(p []byte) (n int, err error) { return call(p) }

// ReadJson 将数据v以JSON格式写入缓冲区的第i个位置
// 参数:
//
//	i - 缓冲区索引(从0开始)
//	v - 需要序列化的数据(任意类型)
//
// 返回:
//
//	error - 序列化过程中可能发生的错误
//
// 注意:
//  1. 如果索引i超出范围，函数会直接返回nil
//  2. 使用json.NewEncoder进行序列化
//  3. 会覆盖缓冲区i位置原有的数据
func (buf Buffer) ReadJson(i int, v any) error {
	if i < len(buf) {
		buffer := bytes.NewBuffer((*buf[i])[:0])
		err := json.NewEncoder(buffer).Encode(v)
		if err != nil {
			return err
		}
		*buf[i] = buffer.Bytes()
	}
	return nil
}

// Get 获取一个io.Reader用于连续读取缓冲区内容
// 返回:
//
//	io.Reader - 可用于连续读取所有缓冲区内容的读取器
//
// 功能说明:
//  1. 返回的读取器会按顺序读取每个缓冲区的内容
//  2. 当所有缓冲区读取完毕时返回io.EOF
//  3. 内部使用bufferCall类型实现连续读取
func (buf *Buffer) Get() io.Reader {
	var i, x, y, z int
	return (bufferCall)(func(p []byte) (int, error) {
		// 输入判断
		if len(p) == 0 || len(*buf) == i {
			return 0, io.EOF
		}
		// 读取当前内容到结束
		for x = 0; ; {
			// 读取
			ptr := (*buf)[i]
			z = copy(p[x:], (*ptr)[y:])
			// 记录读取的长度
			x, y = x+z, y+z
			if x >= len(p) {
				// 缓冲区用完
				break
			}
			if y >= len(*ptr) {
				// 当前切片读取完成
				i++   // 下一个切片
				y = 0 // 新读取位置
			}
			if i >= len(*buf) {
				break // 内容读取完成
			}
		}
		return x, nil
	})
}

// AddDyte 添加数据
func (buf *Buffer) AddDyte(data []byte) {
	*buf = append(*buf, &data)
}

// AddString 添加数据
func (buf *Buffer) AddString(data string) {
	ptr := new([]byte)
	*ptr = []byte(data)
	*buf = append(*buf, ptr)
}

// AddPtr 添加数据
func (buf *Buffer) AddPtr(ptr *[]byte) {
	*buf = append(*buf, ptr)
}
