package sptp

import "io"

var (
	sampleLG        = []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09}
	sampleSM        = []byte{0x01, 0x02}
	sampleThreshold = 4
)

type RecordedRW struct {
	i    int
	data [][]byte
}

func (w *RecordedRW) Write(p []byte) (int, error) {
	if w.data == nil {
		w.data = [][]byte{}
	}
	d := make([]byte, len(p), len(p))
	copy(d, p)
	w.data = append(w.data, d)
	return len(p), nil
}

func (w *RecordedRW) Read(p []byte) (int, error) {
	if w.i >= len(w.data) {
		return 0, io.EOF
	}
	d := w.data[w.i]
	copy(p, d)
	w.i++
	return len(d), nil
}
