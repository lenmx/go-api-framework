package byteDecoder

import "golang.org/x/text/encoding/simplifiedchinese"

type Charset string

const (
	UTF8    = Charset("UTF-8")
	GB18030 = Charset("GB18030")
	GBK     = Charset("GBK")
)

func Decode(source []byte, targetCharset Charset) (result []byte) {
	switch targetCharset {
	case GB18030:
		result, _ = simplifiedchinese.GB18030.NewDecoder().Bytes(source)
	case GBK:
		result, _ = simplifiedchinese.GBK.NewDecoder().Bytes(source)
	case UTF8:
		fallthrough
	default:
		result = source
	}
	return
}
