package tcpReadWrite

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)
//4294967295
const MAX_DATA_SIZE = 4967295

//Read - reads bytes from reader. First reads length of data, than data. Return number of bytes that have been read. If error, return 0 and error.
func Read(rd io.Reader, buffer []byte) (int, error) {
	nBytes := make([]byte, 4)

	//Read first 4 bytes, that containing the length of incoming data
	_, err := rd.Read(nBytes)
	if err != nil {
		return 0, err
	}

	//Length of incoming data
	lenData := binary.BigEndian.Uint32(nBytes)
	if lenData == 0 {
		return 0, io.EOF
	}

	if lenData > MAX_DATA_SIZE {
		return 0, fmt.Errorf("Data size is too large. Max size: %d", MAX_DATA_SIZE)
	}

	if cap(buffer) < int(lenData) {
		return 0, fmt.Errorf("Message to large: %d", lenData)
	}

	//Double the len of slice if cumming data is larger then slice
	if len(buffer) < int(lenData) {
		buffer = buffer[:lenData*2]
	}

	reqLen := 0
	for reqLen < int(lenData) {
		tempLen, err := rd.Read(buffer[reqLen:lenData])
		reqLen += tempLen
		if err != nil {
			if err != io.EOF {
				return 0, fmt.Errorf("Error reading: %s", err.Error())
			} else {
				return 0, errors.New("EOF before reading all requested bytes")
			}
		}
	}
	return reqLen, nil
}

//Write - writes length of data + data into Writer. Return length of written data + length of data. If error return 0 and error.
func Write(wr io.Writer, data []byte) (int, error) {
	if len(data) > MAX_DATA_SIZE {
		return 0, fmt.Errorf("Data size is too large. Max size: %d", MAX_DATA_SIZE)
	}

	if len(data) == 0 {
		return 0, nil
	}

	lenMsg := make([]byte, 4)

	//Create 4 bytes that contains length of data
	binary.BigEndian.PutUint32(lenMsg, uint32(len(data)))

	//Combined data will be look like 0005(message bytes...) - means length 5
	lenMsg = append(lenMsg, data...)

	n, err := wr.Write(lenMsg)
	if err != nil {
		return 0, err
	}

	if n != len(lenMsg) {
		return n, errors.New(fmt.Sprintf("Lost some bytes got %d, sended %d",len(lenMsg), n))
	}
	return n, nil
}
