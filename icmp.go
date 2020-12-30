/*
Ping
*/
package icmp_tunnel

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: ", os.Args[0], "host")
		os.Exit(1)
	}

	addr, err := net.ResolveIPAddr("ip", os.Args[1])
	if err != nil {
		fmt.Println("Resolution error", err.Error())
		os.Exit(1)
	}

	conn, err := net.DialIP("ip4:icmp", nil, addr)
	checkError(err)

	mb := []byte("hello, world!!!!!!!")
	body := &icmp.Echo{
		ID:   1,
		Seq:  1,
		Data: mb[:],
	}

	msg := &icmp.Message{
		Type: (ipv4.ICMPType)(8),
		Code: 0,
		Body: body,
	}

	bytes, err := msg.Marshal(nil)

	_, err = conn.Write(bytes)
	checkError(err)

	var buffer [1024]byte

	_, err = conn.Read(buffer[0:])
	checkError(err)

	fmt.Println("Got response")

	/*
			 response, err := icmp.ParseMessage(ipv4.ICMPTypeEcho.Protocol(), buffer[:])

			 if err != nil {
				fmt.Println("error occurred.")
		        os.Exit(0)
		    }*/

	echo, err := parseEcho(1, buffer[20:])

	fmt.Println(string(echo.Data))

	os.Exit(0)
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}

// parseEcho parses b as an ICMP echo request or reply message body.
func parseEcho(proto int, b []byte) (*icmp.Echo, error) {
	bodyLen := len(b)
	if bodyLen < 4 {
		return nil, io.EOF
	}
	p := &icmp.Echo{ID: int(binary.BigEndian.Uint16(b[:2])), Seq: int(binary.BigEndian.Uint16(b[2:4]))}
	if bodyLen > 4 {
		p.Data = make([]byte, bodyLen-4)
		copy(p.Data, b[4:])
	}
	return p, nil
}

func readFully(conn net.Conn) ([]byte, error) {
	defer conn.Close()

	result := bytes.NewBuffer(nil)
	var buf [512]byte
	for {
		n, err := conn.Read(buf[0:])
		result.Write(buf[0:n])
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
	}
	return result.Bytes(), nil
}
