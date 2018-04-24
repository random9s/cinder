package response

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"compress/lzw"
)

//Encoder function compresses response
type Encoder func(*Response) ([]byte, error)

//Response compressions
const (
	GZIP   = "gzip"
	LZW    = "compress"
	FLATE  = "deflate"
	BROTLI = "br"
)

func encoder(contentEncoding string) Encoder {
	var enc Encoder

	switch contentEncoding {
	case GZIP:
		enc = gzipcompress
	case LZW:
		enc = lzwcompress
	case FLATE:
		enc = deflate
	default:
		enc = identity
	}

	return enc
}

func gzipcompress(r *Response) ([]byte, error) {
	r.header.Set("Content-Encoding", "gzip")

	var buff bytes.Buffer
	writer := gzip.NewWriter(&buff)
	_, err := writer.Write(r.Bytes())
	writer.Close()

	return buff.Bytes(), err
}

func lzwcompress(r *Response) ([]byte, error) {
	r.header.Set("Content-Encoding", "compress")

	var buff bytes.Buffer
	writer := lzw.NewWriter(&buff, lzw.MSB, 8)
	_, err := writer.Write(r.Bytes())
	writer.Close()

	return buff.Bytes(), err
}

func deflate(r *Response) ([]byte, error) {
	r.header.Set("Content-Encoding", "deflate")

	var buff bytes.Buffer
	writer, err := flate.NewWriter(&buff, -1)
	if err != nil {
		return nil, err
	}

	_, err = writer.Write(r.Bytes())
	writer.Close()

	return buff.Bytes(), err
}

func identity(r *Response) ([]byte, error) {
	r.header.Set("Content-Encoding", "identity")
	return r.Bytes(), nil
}
