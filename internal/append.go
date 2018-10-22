package internal

import "bytes"

func AppendStrToBuf(b *bytes.Buffer, str string) {
	if b.Len() == 0 {
		return
	}
	b.WriteString(str)
}
