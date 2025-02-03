package codec

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io"
)

// Compressor 定义消息压缩接口
type Compressor interface {
	// Compress 压缩数据
	Compress(data []byte) ([]byte, error)
	// Decompress 解压数据
	Decompress(data []byte) ([]byte, error)
}

// Encoder 定义消息编码接口
type Encoder interface {
	// Encode 编码数据
	Encode(v interface{}) ([]byte, error)
	// Decode 解码数据
	Decode(data []byte, v interface{}) error
}

// GzipCompressor 实现基于Gzip的压缩器
type GzipCompressor struct{}

// NewGzipCompressor 创建新的Gzip压缩器
func NewGzipCompressor() *GzipCompressor {
	return &GzipCompressor{}
}

// Compress 使用Gzip压缩数据
func (g *GzipCompressor) Compress(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	writer := gzip.NewWriter(&buf)

	_, err := writer.Write(data)
	if err != nil {
		return nil, err
	}

	if err := writer.Close(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// Decompress 使用Gzip解压数据
func (g *GzipCompressor) Decompress(data []byte) ([]byte, error) {
	reader, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	return io.ReadAll(reader)
}

// JSONEncoder 实现基于JSON的编码器
type JSONEncoder struct{}

// NewJSONEncoder 创建新的JSON编码器
func NewJSONEncoder() *JSONEncoder {
	return &JSONEncoder{}
}

// Encode 使用JSON编码数据
func (j *JSONEncoder) Encode(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

// Decode 使用JSON解码数据
func (j *JSONEncoder) Decode(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}
