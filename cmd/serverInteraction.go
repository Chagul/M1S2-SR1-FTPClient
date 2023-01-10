package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
)

var retryDone = 0

type fileList struct {
	files       []fileList
	directories []fileList
	directory   bool
	path        string
}

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

	_, err = conn.Write([]byte(stringToSend))
	fmt.Printf("Sending %s\n", stringToSend)
	if err != nil {
		log.Fatalf(err.Error())
	}
	reply = make([]byte, 1024)
	_, err = conn.Read(reply)
	if err != nil {
		log.Fatalf(err.Error())
	}
	println("response :", string(reply))
	/*if strings.Contains(string(reply), "530") && retryDone < 3 {
		fmt.Printf("Retry conn, type user : ")
		fmt.Scan(&user)
		fmt.Printf("Retry conn, type password : ")
		fmt.Scan(&password)
		retryDone = retryDone + 1
		return UserConn(user, password, conn)
	}*/
	retryDone = 0
	return nil
}

func GetDataConn(conn *net.TCPConn) (*net.TCPConn, error) {
	_, err := conn.Write([]byte("PASV\n"))
	fmt.Println("Sending PASV")
	if err != nil {
		return nil, err
	}
	reader := bufio.NewReader(io.Reader(conn))
	line, _, err := reader.ReadLine()
	lineString := string(line)
	responseForIPAndPort := lineString[strings.Index(lineString, "(")+1 : strings.LastIndex(lineString, ")")]
	arrayResponseForIPAndPort := strings.Split(responseForIPAndPort, ",")
	ipAddr := strings.Join(arrayResponseForIPAndPort[0:4], ".")
	numberToMultiply, err := strconv.Atoi(arrayResponseForIPAndPort[4])
	if err != nil {
		return nil, err
	}
	numberToAdd, err := strconv.Atoi(arrayResponseForIPAndPort[5])
	if err != nil {
		return nil, err
	}
	port := numberToMultiply*256 + numberToAdd

	ip := &net.TCPAddr{
		IP:   net.ParseIP(ipAddr),
		Port: port,
	}
	portStr := strconv.Itoa(port)
	fmt.Printf("Adress : %s, port: %s\n", ipAddr, portStr)
	connData, err := net.DialTCP(TCP_STRING, nil, ip)
	if err != nil {
		return nil, err
	}
	return connData, nil
}
func constructStringToSend(cmd string, stringToAppend string) (string, error) {
	switch cmd {
	case "USER":
		return "USER " + stringToAppend + "\n", nil
	case "PASS":
		return "PASS " + stringToAppend + "\n", nil
	case "LIST":
		return "LIST\n", nil
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

func sendList(mainConn *net.TCPConn, dataConn *net.TCPConn, filelist fileList) (fileList, error) {
	req, err := constructStringToSend("LIST", "")
	if err != nil {
		log.Fatalf(err.Error())
	}
	fmt.Printf("sending %s", req)
	_, err = mainConn.Write([]byte(req))
	if err != nil {
		log.Fatalf(err.Error())
	}
	reply := make([]byte, 1024)
	mainConn.Read(reply)
	println("From conn", string(reply))
	reply = make([]byte, 1024)
	dataConn.Read(reply)
	println("From dataConn", string(reply))
	path := "pouet"
	directories, files, pathBase := parseAnswerList(string(reply), path)
	fileListThisLevel := &fileList{
		files:       files,
		directories: directories,
		directory:   false,
		path:        pathBase,
	}
	return *fileListThisLevel, nil
}

func parseAnswerList(answer string, pathBase string) ([]fileList, []fileList, string) {
	return nil, nil, ""
}
