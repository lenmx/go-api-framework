package byteDecoder

import "testing"

func TestDecode(t *testing.T) {
	str := "test"

	_, err := Decode([]byte(str), GB18030)
	if err != nil {
		t.Error(err)
	}

	t.Log("Decode test pass")
}

func BenchmarkDecode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Decode([]byte("test"), GB18030)
	}
}