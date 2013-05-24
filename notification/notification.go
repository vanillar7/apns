package notification

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
)

type header struct {
	RequestType uint8
	Identifier  uint32
	Expiry      uint32
	TokenLength uint16
	Token       [32]byte
}

type Notification struct {
	Header  header
	Payload string
}

type Failure struct {
	FailureType uint8
	Status      uint8
	Identifier  uint32
}

func DeviceTokenAsBinary(token string) ([32]byte, error) {
	decoded, err := hex.DecodeString(token)
	b := [32]byte{}
	copy(b[:], decoded)
	return b, err
}

func MakeNotification(identifier int, token string, payload string) Notification {
	binaryToken, _ := DeviceTokenAsBinary(token)
	return Notification{Header: header{1, uint32(identifier), 0, 32, binaryToken}, Payload: payload}
}

func (n *Notification) Bytes() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.BigEndian, &n.Header); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.BigEndian, uint16(len(n.Payload))); err != nil {
		return nil, err
	}
	buf.WriteString(n.Payload)
	return buf.Bytes(), nil
}

func NotificationToBytes(n Notification, payload []byte) (*bytes.Buffer, error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.BigEndian, &n); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.BigEndian, uint16(len(payload))); err != nil {
		return nil, err
	}
	buf.Write(payload)
	return buf, nil
}

func FailureFromBytes(resp *bytes.Buffer) Failure {
	var f Failure
	binary.Read(resp, binary.BigEndian, &f)
	return f
}

func ApsPayload(payload string) ([]byte, error) {
	type tree map[string]interface{}
	jsonPayload := tree{"aps": tree{"alert": payload}}
	return json.Marshal(jsonPayload)
}
