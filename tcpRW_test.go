package tcpReadWrite

import (
	"bytes"
	"encoding/binary"
	"io"
	"log"
	"math/rand"
	"strings"
	"testing"
	"time"
)

func TestReadLimit(t *testing.T) {
	buf := bytes.NewBuffer(make([]byte, 0))

	data := make([]byte, MAX_DATA_SIZE+1)

	rand.Read(data)

	lenMsg := make([]byte, 4)

	//Create 4 bytes that contains length of data
	binary.BigEndian.PutUint32(lenMsg, uint32(len(data)))

	_, err := buf.Write(lenMsg)
	if err != nil {
		t.Fatal(err)
	}

	n, err := buf.Write(data)
	if err != nil {
		t.Fatal(err)
	}

	if n != len(data) {
		t.Fatalf("Lost some bytes got %d, sended %d", len(data), n)
	}

	readData := make([]byte, len(data))

	n, err = Read(buf, readData)
	if err != nil {
		if strings.Contains(err.Error(), "Data size is too large") {
			return
		}
		t.Fatal(err)
	}
}

func TestWriteLimit(t *testing.T) {
	buf := bytes.NewBuffer(make([]byte, 0))
	data := make([]byte, MAX_DATA_SIZE+1)
	rand.Read(data)

	_, err := Write(buf, data)
	if err != nil {
		if strings.Contains(err.Error(), "Data size is too large") {
			return
		}
		t.Fatal(err)
	}
	t.Fail()
}

func TestReadZero(t *testing.T) {
	buf := bytes.NewBuffer(make([]byte, 0))
	readData := make([]byte, 1)

	_, err := Read(buf, readData)
	if err != nil {
		if err == io.EOF {
			return
		}
		t.Fatal(err)
	}
}

func TestWriteZero(t *testing.T) {
	buf := bytes.NewBuffer(make([]byte, 0))
	data := []byte("")

	writtenN, err := Write(buf, data)
	if err != nil {
		t.Fatal(err)
	}

	if writtenN == 0 {
		return
	} else {
		t.Fail()
	}
}

func TestRead(t *testing.T) {
	buf := bytes.NewBuffer(make([]byte, 0))

	data := make([]byte, 16)

	rand.Read(data)

	lenMsg := make([]byte, 4)

	//Create 4 bytes that contains length of data
	binary.BigEndian.PutUint32(lenMsg, uint32(len(data)))

	_, err := buf.Write(lenMsg)
	if err != nil {
		t.Fatal(err)
	}

	n, err := buf.Write(data)
	if err != nil {
		t.Fatal(err)
	}

	if n != len(data) {
		t.Fatalf("Lost some bytes got %d, sended %d", len(data), n)
	}

	readData := make([]byte, len(data))

	n, err = Read(buf, readData)
	if err != nil {
		t.Fatal(err)
	}

	if n != len(data) {
		t.Fatalf("Lost some bytes got %d, sended %d", len(data), n)
	}

	if !bytes.Equal([]byte(data), readData[:n]) {
		t.Fatal("Read not full data")
	}
}

func TestWrite(t *testing.T) {
	buf := bytes.NewBuffer(make([]byte, 0))
	data := make([]byte, 16)
	rand.Read(data)

	actualN := len(data)

	writtenN, err := Write(buf, data)
	if err != nil {
		t.Fatal(err)
	}

	if writtenN-4 != actualN {
		t.Fatalf("Write not all bytes, have: %d, written: %d", actualN, writtenN-4 )
	}

	if !bytes.Equal(data, buf.Bytes()[4:writtenN]) {
		t.Fatal("Written bytes and actual are not equal")
	}
}

func TestRun(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	rand.Seed(time.Now().Unix())
	t.Run("Read", TestRead)
	t.Run("Write", TestWrite)
}