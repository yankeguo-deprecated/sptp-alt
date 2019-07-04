package sptp

import (
	"bytes"
	"compress/gzip"
	"io"
	"testing"
	"time"
)

func TestReceiverNoCompressNoChunk(t *testing.T) {
	rw := &RecordedRW{}
	w := NewWriterWithOptions(rw, WriterOptions{ChunkThreshold: sampleThreshold})
	r := NewReceiverWithOptions(rw, ReceiverOptions{ChunkTimeout: time.Second / 2})

	var buf []byte
	var err error
	var n int

	if n, err = w.Write(sampleSM); err != nil {
		t.Fatal(err)
	}
	if n != len(sampleSM) {
		t.Fatal("w bad n")
	}
	if buf, err = r.Receive(); err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(buf, sampleSM) {
		t.Fatal("bad result")
	}
}

func TestReceiverNoCompressChunked(t *testing.T) {
	rw := &RecordedRW{}
	w := NewWriterWithOptions(rw, WriterOptions{ChunkThreshold: sampleThreshold})
	r := NewReceiverWithOptions(rw, ReceiverOptions{ChunkTimeout: time.Second / 2})

	var buf []byte
	var err error
	var n int

	if n, err = w.Write(sampleLG); err != nil {
		t.Fatal(err)
	}
	if n != len(sampleLG) {
		t.Fatal("w bad n")
	}

	var i int
	for {
		i++
		if i > 3 {
			t.Fatal("too much loop")
		}
		if buf, err = r.Receive(); err != nil {
			if err != io.EOF {
				t.Fatal(err)
			}
		}
		if buf == nil {
			continue
		}
		if !bytes.Equal(buf, sampleLG) {
			t.Fatal("bad result")
		}
		break
	}
}

func TestReceiverCompressNoChunk(t *testing.T) {
	rw := &RecordedRW{}
	w := NewWriterWithOptions(rw, WriterOptions{GzipLevel: gzip.BestCompression})
	r := NewReceiverWithOptions(rw, ReceiverOptions{ChunkTimeout: time.Second / 2})

	var buf []byte
	var err error
	var n int

	if n, err = w.Write(sampleSM); err != nil {
		t.Fatal(err)
	}
	if n != len(sampleSM) {
		t.Fatal("w bad n")
	}
	if buf, err = r.Receive(); err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(buf, sampleSM) {
		t.Fatal("bad result")
	}
}

func TestReceiverCompressChunked(t *testing.T) {
	rw := &RecordedRW{}
	w := NewWriterWithOptions(rw, WriterOptions{ChunkThreshold: sampleThreshold, GzipLevel: gzip.BestCompression})
	r := NewReceiverWithOptions(rw, ReceiverOptions{ChunkTimeout: time.Second / 2})

	var buf []byte
	var err error
	var n int

	if n, err = w.Write(sampleLG); err != nil {
		t.Fatal(err)
	}
	if n != len(sampleLG) {
		t.Fatal("w bad n")
	}

	var i int
	for {
		i++
		if i > 100 {
			t.Fatal("too much loop")
		}
		if buf, err = r.Receive(); err != nil {
			if err != io.EOF {
				t.Fatal(err)
			}
		}
		if buf == nil {
			continue
		}
		if !bytes.Equal(buf, sampleLG) {
			t.Fatal("bad result")
		}
		break
	}
}
