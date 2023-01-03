package cmd

import (
	"errors"
	"fmt"
	"log"
	"net"
)

func UserConn(user string, pwd string, conn *net.TCPConn) error {
	stringToSend, err := constructStringToSend("USER", user)
	if err != nil {
		return err
	}

	_, err = conn.Write([]byte(stringToSend))
	if err != nil {
		log.Fatalf(err.Error())
	}
	fmt.Printf("Sending %s", stringToSend)

	reply := make([]byte, 1024)
	_, err = conn.Read(reply)
	if err != nil {
		log.Fatalf(err.Error())
	}
	println("response :", string(reply))

	stringToSend, err = constructStringToSend("PASS", pwd)
	if err != nil {
		return err
	}

	_, err = conn.Write([]byte("PASS anonymous\n"))
	fmt.Printf("Sending PASS anonymous\n")
	if err != nil {
		log.Fatalf(err.Error())
	}
	reply = make([]byte, 1024)
	_, err = conn.Read(reply)
	if err != nil {
		log.Fatalf(err.Error())
	}
	println("response :", string(reply))

	return nil
}

func constructStringToSend(cmd string, stringToAppend string) (string, error) {
	switch cmd {
	case "USER":
		return "USER " + stringToAppend + "\n", nil
	case "PASS":
		return "PASS " + stringToAppend + "\n", nil
	}
	return "", errors.New("command" + cmd + "not found/supported")
}

func GetIpFromURL(url string) (*net.TCPAddr, error) {
	ip, err := net.LookupIP(addressServer)
	if err != nil {
		return nil, err
	}
	fmt.Printf("IP adress found for %s : %s\n", addressServer, ip[0].String())

	addr := &net.TCPAddr{
		IP:   ip[0],
		Port: port,
	}
	return addr, nil
}
