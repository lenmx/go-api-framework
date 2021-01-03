package byteDecoder

import "golang.org/x/text/encoding/simplifiedchinese"

type Charset string

const (
	UTF8    = Charset("UTF-8")
	GB18030 = Charset("GB18030")
	GBK     = Charset("GBK")
)

func Decode(source []byte, targetCharset Charset) (result []byte, err error) {
	switch targetCharset {
	case GB18030:
		result, err = simplifiedchinese.GB18030.NewDecoder().Bytes(source)
	case GBK:
		result, err = simplifiedchinese.GBK.NewDecoder().Bytes(source)
	case UTF8:
		fallthrough
	default:
		result = source
	}
	return
}
