package packer

import (
	"botify/main/logger"
	"encoding/json"
)

func Decentralize(data []byte) (request Request, err error) {
	var header [2]byte
	var length [4]byte

	var intHEADER int
	var intLENGTH int

	// get the header
	header[0] = data[0]
	header[1] = data[1]

	// get the length of the data
	length[0] = data[2]
	length[1] = data[3]
	length[2] = data[4]
	length[3] = data[5]

	// convert the header and the length to int
	intHEADER = int(header[0])*256 + int(header[1])
	intLENGTH = int(length[0])*256*256*256 + int(length[1])*256*256 + int(length[2])*256 + int(length[3])

	// get the data
	var jsonData []byte
	for i := 0; i < intLENGTH; i++ {
		jsonData = append(jsonData, data[i+6])
	}

	// unmarshal the data
	err = json.Unmarshal(jsonData, &request.Data)
	if err != nil {
		logger.Error(err)
		return
	}
	request.Header = intHEADER

	return
}

// centralize a response
func Centralize(response Response) (data []byte, err error) {
	// convert the header and the length to byte
	var header [2]byte
	var length [4]byte
	var intLENGTH int

	header[0] = byte(response.Header / 256)
	header[1] = byte(response.Header % 256)

	// marshal the data
	var jsonData []byte
	jsonData, err = json.Marshal(response.Data)
	if err != nil {
		return
	}

	// get the length of the data
	intLENGTH = len(jsonData)

	length[0] = byte(intLENGTH / 256 / 256 / 256)
	length[1] = byte(intLENGTH / 256 / 256 % 256)
	length[2] = byte(intLENGTH / 256 % 256)
	length[3] = byte(intLENGTH % 256)

	// append the header and the length to the data
	data = append(data, header[0])
	data = append(data, header[1])
	data = append(data, length[0])
	data = append(data, length[1])
	data = append(data, length[2])
	data = append(data, length[3])
	data = append(data, jsonData...)

	return
}
