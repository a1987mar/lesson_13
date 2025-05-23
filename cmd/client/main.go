package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {

	us := bufio.NewScanner(os.Stdin)
	uw := bufio.NewWriter(os.Stdout)

	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		panic(fmt.Errorf("error connecting: %w", err))
	}

	defer conn.Close()

	sr := bufio.NewReader(conn)
	sw := bufio.NewWriter(conn)

	for us.Scan() {
		msg := us.Text()

		// if

		//fmt.Printf("User Input: %s\n", msg)

		sw.WriteString(fmt.Sprintf("%s\n", msg))
		sw.Flush()

		resp, err := sr.ReadString('\n')
		if err != nil {
			panic(fmt.Errorf("error reading response: %w", err))
		}

		uw.WriteString(resp)
		uw.Flush()
	}
}
