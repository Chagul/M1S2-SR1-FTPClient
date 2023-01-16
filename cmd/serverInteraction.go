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
	"sort"
)

var retryDone = 0
var firstPASV = true
var paths = make([]file,0)

type file struct {
	path string
	directory   bool
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
	if(!firstPASV){
		line, _, err = reader.ReadLine()
		lineString = string(line)
	}
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
	firstPASV = false;
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
	case "CWD":
		return "CWD " + stringToAppend + "\n", nil
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

func sendList(mainConn *net.TCPConn, dataConn *net.TCPConn, base string) error {
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
	defer dataConn.Close()
	//println("From dataConn", string(reply))
	pathss := parseAnswerList(string(reply), base)
	for _, val := range pathss{
		if(val.directory){
			dataConn, err := GetDataConn(mainConn)
			if(err != nil){
				log.Fatalf("rip : %s", err.Error())
			}
			req, err = constructStringToSend("CWD", val.path)
			if err != nil {
				log.Fatalf(err.Error())
			}
			fmt.Printf("sending %s\n", req)
			_, err = mainConn.Write([]byte(req))
			if err != nil {
				log.Fatalf(err.Error())
			}
			reader := bufio.NewReader(io.Reader(mainConn))
			line, _, err := reader.ReadLine()
			lineString := string(line)
			fmt.Println(lineString)
			//go func() {
			err = sendList(mainConn, dataConn, val.path)
			if err != nil{
				log.Fatalf("rip : %s", err.Error())
			}
			dataConn.Close()
			//}()
		}
	}
	for _,vals := range pathss{
		paths = append(paths, vals)
	}
	return nil
}

func parseAnswerList(answer string, base string) []file{
	pathss := make([]file, 0)
	lines := make([]string, 0) 
	var j = 0
	for i := 0; i < len(answer); i++{
		if(answer[i] == '\n'){
			lines = append(lines,answer[j:i-1])
			j = i+1
		}
	}
	for _, val := range lines{
		currentFile := file{}
		if(val[0] == 'd'){
			currentFile.directory = true
			currentFile.path = base + val[strings.LastIndex(val, " ")+1:] + "/"
		}else{
			currentFile.directory = false
			currentFile.path = base + val[strings.LastIndex(val, " ")+1:]
		}
	pathss = append(pathss, currentFile)
	}
	return pathss

}

func tree() {
	filePaths := make([]string, 0)
	for _, file := range paths{
		filePaths = append(filePaths, file.path)
	}
	sort.Strings(filePaths)
	//depth := 0
	parent := ""
	for _,val := range filePaths{
		val = strings.TrimPrefix(val, "/")
		currentParent = [0:strings.Index(val, "/")-1]
		if parent == ""{
			parent = currentParent
		}
		val = strings.TrimLeft(val, val[0:strings.Index(val, "/")])
		fmt.Println(val);
	}
}
