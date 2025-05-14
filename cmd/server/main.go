package main

import (
	"bufio"
	"fmt"
	"les12/internal/commands"
	"les12/internal/documentstore"
	"net"
	"strings"

	"github.com/bytedance/sonic"
)

func main() {
	l, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(map[string]string{"error listening": err.Error()})
	}
	s := documentstore.NewStore()
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		fmt.Println("Start server")
		go handleConnection(conn, s)
	}
}

func handleConnection(conn net.Conn, s *documentstore.Store) {
	defer conn.Close()
	scanner := bufio.NewScanner(conn)
	w := bufio.NewWriter(conn)
	for scanner.Scan() {
		msg := scanner.Text()
		elems := strings.SplitN(msg, " ", 2)
		if len(elems) != 2 {
			w.WriteString("invalid command\n")
			w.Flush()
			continue
		}
		var resp string
		var err error
		switch elems[0] {
		case commands.GetCommandName:
			resp, err = GetCollect(s, elems[1])
		case commands.PutCommandName:
			resp, err = PutCollection(s, elems[1])
		case commands.DeleteCommandName:
			resp = DeleteCollect(s, elems[1])
		default:
			w.WriteString("invalid command\n")
		}
		if err != nil {
			w.WriteString(fmt.Sprintf("error: %v\n", err))
			w.Flush()
			continue
		}

		fmt.Printf("%s\n", msg)
		w.WriteString(fmt.Sprintf("your request has been processed:%s\n", resp))
		w.Flush()
	}
}

func PutCollection(s *documentstore.Store, raw string) (string, error) {
	var nc commands.NewCollection
	err := sonic.Unmarshal([]byte(raw), &nc)
	if err != nil {
		return "", err
	}
	nameMap := strings.ToLower(nc.Name)
	c, err := s.GetCollection(nameMap)
	if err != nil {
		c, err = s.CreateCollection(nameMap, nc.ID)
		if err != nil {
			return "", err
		}
	}
	doc := documentstore.Document{
		Fields: map[string]documentstore.DocumentField{
			nc.ID:   {Type: documentstore.DocumentFieldTypeString, Value: nc.Doc.Id},
			nc.Name: {Type: documentstore.DocumentFieldTypeString, Value: nc.Doc.Name},
		},
	}
	c.Put(doc)

	return "data added to store", nil
}

func GetCollect(s *documentstore.Store, raw string) (string, error) {
	var nc commands.NewCollection
	err := sonic.Unmarshal([]byte(raw), &nc)
	if err != nil {
		return "", err
	}
	nameMap := strings.ToLower(nc.Name)
	gDoc, err := s.GetCollection(nameMap)
	if err != nil {
		return "", err
	}
	colDTO := gDoc.ToDto()
	res, err := sonic.Marshal(colDTO)
	if err != nil {
		return "", err
	}

	return string(res), nil
}

func DeleteCollect(s *documentstore.Store, raw string) string {
	var nc commands.NewCollection
	err := sonic.Unmarshal([]byte(raw), &nc)
	if err != nil {
		return "іnvalid data format"
	}
	nameMap := strings.ToLower(nc.Name)
	if !s.DeleteCollection(nameMap) {
		return "collection not deleted, check the data you have entered"
	}
	return "deleted"
}
