package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"sort"
	"strconv"
	"strings"
	"time"
)

var firstPASV = true
var paths = make([]file, 0)

type file struct {
	path      string
	directory bool
}

// UserConn Init TCP conn with given user and pwd, if both are not precised, anonymous is the default/**
func UserConn(user string, pwd string, conn *net.TCPConn) error {
	fmt.Println("userconn")
	stringToSend, err := constructStringToSend("USER", user)
	if err != nil {
		return err
	}

	_, err = conn.Write([]byte(stringToSend))
	if err != nil {
		log.Fatalf(err.Error())
	}

	reply := make([]byte, 1024)
	_, err = conn.Read(reply)
	if err != nil {
		log.Fatalf(err.Error())
	}

	stringToSend, err = constructStringToSend("PASS", pwd)
	if err != nil {
		return err
	}

	_, err = conn.Write([]byte(stringToSend))
	if err != nil {
		log.Fatalf(err.Error())
	}

	reply = make([]byte, 1024)
	_, err = conn.Read(reply)
	if err != nil {
		log.Fatalf(err.Error())
	}
	fmt.Println("fin userconn")
	return nil
}

// GetDataConn Create a new data connection from mainConn, that send PASV/**
func GetDataConn(conn *net.TCPConn) (*net.TCPConn, error) {
	fmt.Printf("dataconn\n")
	_, err := conn.Write([]byte("PASV\n"))
	if err != nil {
		return nil, err
	}

	reader := bufio.NewReader(io.Reader(conn))
	line, _, err := reader.ReadLine()
	if err != nil {
		fmt.Printf("POUET1 ; %s", err)
	}

	lineString := string(line)
	if !firstPASV {
		line, _, err = reader.ReadLine()
		if err != nil {
			fmt.Printf("POUET2 ; %s", err)
		}
		lineString = string(line)
	}

	err, ipAddrDataConn, portDataConn := getIPAndPortFromResponse(lineString)
	if err != nil {
		return nil, err
	}

	ip := &net.TCPAddr{
		IP:   net.ParseIP(ipAddrDataConn),
		Port: portDataConn,
	}
	connData, err := net.DialTCP(TcpString, nil, ip)
	if err != nil {
		return nil, err
	}

	firstPASV = false
	return connData, nil
}

// constructStringToSend Construct the command that will be sent to the ftp server**/
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

func GetIpFromURL() (*net.TCPAddr, error) {
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

// sendList  send the command list to the mainConn, and CWD to all directories returned to recursively call sendList with them**/
func sendList(mainConn *net.TCPConn, dataConn *net.TCPConn, base string) error {
	fmt.Printf("sendList\n")
	req, err := constructStringToSend("LIST", "")
	if err != nil {
		log.Fatalf(err.Error())
	}

	_, err = mainConn.Write([]byte(req))
	if err != nil {
		log.Fatalf(err.Error())
	}

	reply := make([]byte, 1024)

	//to ignore ok
	_, err = mainConn.Read(reply)
	if err != nil {
		log.Fatalf("%s", err.Error())
	}
	reply = make([]byte, 1024)
	_, err = dataConn.Read(reply)
	if err != nil && err != io.EOF {
		log.Fatalf("EOF ?%s\n", err.Error())
	}
	pathss := parseAnswerList(string(reply), base)
	for _, val := range pathss {
		if val.directory {
			dataConn, err = GetDataConn(mainConn)
			if err != nil {
				log.Fatalf("rip : %s", err.Error())
			}

			req, err = constructStringToSend("CWD", val.path)
			if err != nil {
				log.Fatalf(err.Error())
			}

			_, err = mainConn.Write([]byte(req))
			if err != nil {
				log.Fatalf(err.Error())
			}

			readerMainConn := bufio.NewReader(io.Reader(mainConn))
			line, _, err := readerMainConn.ReadLine()
			lineString := string(line)
			fmt.Println(lineString)
			if err != nil {
				fmt.Printf("Err while readline")
			}
			err = sendList(mainConn, dataConn, val.path)
			if err != nil {
				log.Fatalf("rip : %s", err.Error())
			}

			err = dataConn.Close()
			if err != nil {
				log.Fatalf("pouet %s", err.Error())
			}
		}
	}
	for _, vals := range pathss {
		paths = append(paths, vals)
	}
	return nil
}

// parseAnswerList parse the asnwer from the list command, add absolute path to the global array of struct paths
func parseAnswerList(answer string, base string) []file {
	pathss := make([]file, 0)
	lines := make([]string, 0)
	var j = 0
	for i := 0; i < len(answer); i++ {
		if answer[i] == '\n' {
			lines = append(lines, answer[j:i-1])
			j = i + 1
		}
	}
	for _, val := range lines {
		currentFile := file{}
		if val[0] == 'd' {
			currentFile.directory = true
			currentFile.path = base + val[strings.LastIndex(val, " ")+1:] + "/"
		} else {
			currentFile.directory = false
			currentFile.path = base + val[strings.LastIndex(val, " ")+1:]
		}
		pathss = append(pathss, currentFile)
	}
	return pathss

}

// tree Construct and print the tree-output from path in paths **/
func tree() {
	sort.Slice(paths, func(i, j int) bool {
		return paths[i].path < paths[j].path
	})
	for _, val := range paths {
		fmt.Printf("%s\n", val.path)
	}
	time.Sleep(time.Second * 10)
	//depth := 0
	space := "    "
	trail := "---"
	//branch := "│   "
	tee := "├── "
	//last :=   "└── "
	parent := paths[0].path
	parent = strings.TrimLeft(parent, "/")
	fmt.Printf("%s\n\t%s", parent, tee)
	//depth := 0
	for _, val := range paths {
		pathss := strings.Split(val.path, "/")
		fmt.Println(pathss)
		for i := range pathss {
			fmt.Printf(pathss[i] + "\n")

			if pathss[i] == parent {
				fmt.Printf("%s%s", space, trail)
			} else {
				if val.directory {
					fmt.Printf("%s%s", tee, pathss[i])
				} else {
					fmt.Printf("%s%s", trail, pathss[i])
				}
				break
			}
		}
	}
	/*val = strings.TrimPrefix(val, "/")
	currentParent := val[0:strings.Index(val, "/")]
	cutVal := strings.TrimLeft(val, val[0:strings.Index(val, "/")+1])
	for currentParent == parent {
		fmt.Printf("%s", tee)
		currentParent = cutVal[0:strings.Index(val, "/")]
		cutVal = strings.TrimLeft(cutVal, cutVal[0:strings.Index(val, "/")+1])
	}
	fmt.Printf("%s %s\n", tee, cutVal)*/
}

/*
* getIPAndPortFromResponse parse the reply from PASV request to calculate and return the IP address and the port*
 */
func getIPAndPortFromResponse(reply string) (error, string, int) {
	responseForIPAndPort := reply[strings.Index(reply, "(")+1 : strings.LastIndex(reply, ")")]
	arrayResponseForIPAndPort := strings.Split(responseForIPAndPort, ",")
	ipAddr := strings.Join(arrayResponseForIPAndPort[0:4], ".")
	numberToMultiply, err := strconv.Atoi(arrayResponseForIPAndPort[4])
	if err != nil {
		return err, "", 0
	}
	numberToAdd, err := strconv.Atoi(arrayResponseForIPAndPort[5])
	if err != nil {
		return err, "", 0
	}
	port := numberToMultiply*256 + numberToAdd
	return nil, ipAddr, port
}
