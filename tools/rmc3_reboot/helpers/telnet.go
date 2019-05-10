package helpers

import (
	"errors"
	"fmt"
	"net"
	"strings"
	"time"
	//"github.com/byuoitav/common/log"
)

const (
	CARRIAGE_RETURN           = 0x0D
	LINE_FEED                 = 0x0A
	SPACE                     = 0x20
	DELAY_BETWEEN_CONNECTIONS = time.Second * 10
)

// GetConnection makes and returns a tcp connection with the given address
func GetConnection(address string, readWelcome bool) (*net.TCPConn, error) {

	addr, err := net.ResolveTCPAddr("tcp", address+":23") //checks for a valid address at port 5000 and returns the address
	if err != nil {
		return nil, err
	}

	conn, err := net.DialTCP("tcp", nil, addr) //connects to the given address
	if err != nil {
		return nil, err
	}

	if readWelcome {
		_, err := readUntil(CARRIAGE_RETURN, conn, 3) // skips welcome message
		if err != nil {
			return conn, err
		}
	}

	return conn, err
}

// SendCommand writes the given command to a device over the given tcp connection, returning the result
func SendCommand(conn *net.TCPConn, command string) (resp string, err error) {

	resp, err = writeCommand(conn, command)
	if err != nil {
		return "", err
	}

	return resp, nil
}

func writeCommand(conn *net.TCPConn, command string) (string, error) {

	command = strings.Replace(command, " ", string(SPACE), -1)
	command += string(CARRIAGE_RETURN) + string(LINE_FEED)
	conn.Write([]byte(command))

	// get response
	resp, err := readUntil(LINE_FEED, conn, 5)
	if err != nil {
		return "", err
	}
	return string(resp), nil
}

func readUntil(delimeter byte, conn *net.TCPConn, timeoutInSeconds int) ([]byte, error) {

	conn.SetReadDeadline(time.Now().Add(time.Duration(int64(timeoutInSeconds)) * time.Second))

	buffer := make([]byte, 128)
	message := []byte{}

	for !charInBuffer(delimeter, buffer) { //Loops while delimiter is not found
		_, err := conn.Read(buffer)
		if err != nil {
			err = errors.New(fmt.Sprintf("Error reading response: %s", err.Error()))
			return message, err
		}

		message = append(message, buffer...) // appends buffer to message
	}

	return removeNil(message), nil
}

// Returns true if the delimiter is found
func charInBuffer(toCheck byte, buffer []byte) bool {

	for _, b := range buffer {
		if toCheck == b {
			return true
		}
	}

	return false
}

func removeNil(b []byte) (ret []byte) {

	for _, c := range b {
		switch c {
		case '\x00':
			break
		default:
			ret = append(ret, c)
		}
	}
	return ret
}
