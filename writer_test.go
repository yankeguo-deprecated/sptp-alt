package sptp

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
	"testing"
)

func TestWriter_WritePlain(t *testing.T) {
	var err error

	w := &RecordedRW{}
	pw := NewWriterWithOptions(w, WriterOptions{ChunkThreshold: sampleThreshold})

	var n int
	if n, err = pw.Write(sampleLG); err != nil {
		t.Fatal(err)
	}
	if n != 9 {
		t.Fatal("length return mismatch")
	}
	if len(w.data) != 3 {
		t.Fatal("length mismatch")
	}
	if w.data[0][0] != Magic || w.data[1][0] != Magic || w.data[2][0] != Magic {
		t.Fatal("missing magic")
	}
	if w.data[0][1] != ModeChunked || w.data[1][1] != ModeChunked || w.data[2][1] != ModeChunked {
		t.Fatal("missing mode chunked")
	}
	id := w.data[0][2:10]
	if !bytes.Equal(id, w.data[1][2:10]) {
		t.Fatal("line 2, id not equal")
	}
	if !bytes.Equal(id, w.data[2][2:10]) {
		t.Fatal("line 3, id not equal")
	}
	if !bytes.Equal([]byte{0x03, 0x00}, w.data[0][10:12]) {
		t.Fatal("line 1, bad count/index")
	}
	if !bytes.Equal([]byte{0x03, 0x01}, w.data[1][10:12]) {
		t.Fatal("line 2, bad count/index")
	}
	if !bytes.Equal([]byte{0x03, 0x02}, w.data[2][10:12]) {
		t.Fatal("line 3, bad count/index")
	}
	if !bytes.Equal([]byte{0x01, 0x02, 0x03, 0x04}, w.data[0][12:]) {
		t.Fatal("line 1, bad payload")
	}
	if !bytes.Equal([]byte{0x05, 0x06, 0x07, 0x08}, w.data[1][12:]) {
		t.Fatal("line 2, bad payload")
	}
	if !bytes.Equal([]byte{0x09}, w.data[2][12:]) {
		t.Fatal("line 3, bad payload")
	}
}

func TestWriter_WriteCompression(t *testing.T) {
	var err error

	w := &RecordedRW{}
	pw := NewWriterWithOptions(w, WriterOptions{ChunkThreshold: sampleThreshold, GzipLevel: gzip.BestCompression})

	var n int
	if n, err = pw.Write(sampleLG); err != nil {
		t.Fatal(err)
	}
	if n != 9 {
		t.Fatal("length return mismatch")
	}
	out := make([]byte, 0)
	for i, d := range w.data {
		if d[0] != Magic {
			t.Fatal("bad magic", i+1)
		}
		if d[1] != ModeChunked|ModeGzipped {
			t.Fatal("bad mode", i+1)
		}
		out = append(out, d[12:]...)
	}

	var gr *gzip.Reader
	if gr, err = gzip.NewReader(bytes.NewReader(out)); err != nil {
		t.Fatal(err)
	}
	var res []byte
	if res, err = ioutil.ReadAll(gr); err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(res, sampleLG) {
		t.Logf("% 02x", res)
		t.Fatal("failed to resume")
	}
}
