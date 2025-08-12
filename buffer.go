package catai

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
	var i, n int // i-当前缓冲区索引, n-当前缓冲区已读取位置
	return (bufferCall)(func(p []byte) (int, error) {
		// 检查缓冲区是否为空
		if len(*buf) == 0 {
			return 0, io.EOF
		}
		// 从当前缓冲区拷贝数据
		ptr := (*buf)[i]
		l := copy(p, (*ptr)[n:]) // 拷贝数据到p
		n += l                   // 更新已读取位置

		// 检查是否读完当前缓冲区
		if n >= len((*ptr)) {
			n = 0 // 重置位置
			i++   // 移动到下一个缓冲区
			// 检查是否所有缓冲区都已读完
			if i >= len(*buf) {
				return l, io.EOF
			}
		}
		return l, nil // 返回实际读取的字节数
	})
}
