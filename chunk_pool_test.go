package sptp

import (
	"bytes"
	"testing"
	"time"
)

func TestChunkGroup(t *testing.T) {
	cg := NewChunkGroup(1, time.Now(), ModeChunked, 2)
	if cg.IsCompleted() {
		t.Fatal("early completed")
	}
	cg.Add([]byte{0x00}, 0)
	if cg.IsCompleted() {
		t.Fatal("early completed")
	}
	cg.Add([]byte{0x01}, 0)
	if cg.IsCompleted() {
		t.Fatal("early completed")
	}
	cg.Add([]byte{0x02}, 1)
	if !cg.IsCompleted() {
		t.Fatal("not completed")
	}
	if !bytes.Equal(cg.Bytes(), []byte{0x01, 0x02}) {
		t.Fatal("bad bytes")
	}
}

func TestChunkPool(t *testing.T) {
	cp := NewChunkPool(time.Second)
	var p []byte
	var err error
	p, err = cp.Consume(0x01, ModeChunked, 3, 0, []byte{0x01})
	if p != nil {
		t.Fatal("should be nil")
	}
	if err != nil {
		t.Fatal(err)
	}
	p, err = cp.Consume(0x01, ModeChunked, 3, 1, []byte{0x02, 0x03})
	if p != nil {
		t.Fatal("should be nil")
	}
	if err != nil {
		t.Fatal(err)
	}
	p, err = cp.Consume(0x01, ModeChunked|ModeGzipped, 3, 2, []byte{0x04, 0x05})
	if p != nil {
		t.Fatal("should be nil")
	}
	if err == nil {
		t.Fatal("should error")
	}
	p, err = cp.Consume(0x01, ModeChunked, 3, 2, []byte{0x04, 0x05})
	if p == nil {
		t.Fatal("should not be nil")
	}
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(p, []byte{0x01, 0x02, 0x03, 0x04, 0x05}) {
		t.Fatal("bad bytes")
	}
	p, err = cp.Consume(0x01, ModeChunked, 3, 0, []byte{0x01})
	time.Sleep(time.Second * 2)
	p, err = cp.Consume(0x02, ModeChunked, 3, 0, []byte{0x01})
	if cp.ChunkGroups[0x01] != nil {
		t.Fatal("not GCed")
	}
}
