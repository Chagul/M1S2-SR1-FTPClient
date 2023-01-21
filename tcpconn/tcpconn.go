package tcpconn

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
	tree2 "tree-ftp/tree"
	constantFTP "tree-ftp/util/ftp"
	constant "tree-ftp/util/global"
)

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

	reply := make([]byte, constant.SizeAnswer)
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

	reply = make([]byte, constant.SizeAnswer)
	_, err = conn.Read(reply)
	if err != nil {
		log.Fatalf(err.Error())
	}
	fmt.Println("fin userconn")
	return nil
}

// GetDataConn Create a new data connection from mainConn, that send PASV/**
func GetDataConn(conn *net.TCPConn) (*net.TCPConn, error) {
	_, err := conn.Write([]byte("PASV\n"))
	if err != nil {
		return nil, err
	}

	reader := bufio.NewReader(io.Reader(conn))
	line, _, err := reader.ReadLine()
	if err != nil && !strings.Contains(string(line), constantFTP.CodePasvOk) {
		log.Fatalf(err.Error())
	}

	lineString := string(line)
	err, ipAddrDataConn, portDataConn := getIPAndPortFromResponse(lineString)
	if err != nil {
		return nil, err
	}

	ip := &net.TCPAddr{
		IP:   net.ParseIP(ipAddrDataConn),
		Port: portDataConn,
	}
	connData, err := net.DialTCP(constant.TcpString, nil, ip)
	if err != nil {
		return nil, err
	}

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

func GetIpFromURL(port int, addressServer string) (*net.TCPAddr, error) {
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

// SendList send the command list to the mainConn, and CWD to all directories returned to recursively call sendList with them**/
func SendList(mainConn *net.TCPConn, dataConn *net.TCPConn, base string, maxDepth int, currentDepth int, currentNode *tree2.Node) error {
	fmt.Printf("SendList %s \n", base)
	req, err := constructStringToSend("LIST", "")
	if err != nil {
		log.Fatalf(err.Error())
	}

	_, err = mainConn.Write([]byte(req))
	if err != nil {
		log.Fatalf(err.Error())
	}

	readerMainConn := bufio.NewReader(io.Reader(mainConn))
	line, _, err := readerMainConn.ReadLine()

	if err != nil || !strings.Contains(string(line), constantFTP.CodeOkList) {
		log.Fatalf("LIST RETURN ERROR")
	}

	readerDataConn := bufio.NewReader(io.Reader(dataConn))
	lines := getListLines(readerDataConn)
	children := parseAnswerList(lines, base, currentDepth)
	line, _, err = readerMainConn.ReadLine()
	if (err != nil && err != io.EOF) || !strings.Contains(string(line), constantFTP.CodeComingList) {
		log.Fatalf("Fin liste %s", err.Error())
	}
	if currentDepth == maxDepth {
		currentNode.AddChildren(children)
		return nil
	}
	for _, child := range children {
		if child.IsDirectory {
			dataConn, err = GetDataConn(mainConn)
			if err != nil {
				log.Fatalf("rip : %s", err.Error())
			}

			req, err = constructStringToSend("CWD", child.Filepath)
			if err != nil {
				log.Fatalf(err.Error())
			}
			_, err = mainConn.Write([]byte(req))
			if err != nil {
				log.Fatalf(err.Error())
			}

			line, _, err = readerMainConn.ReadLine()
			if err != nil || !strings.Contains(string(line), constantFTP.CodeCWDOk) {
				fmt.Printf("Err while readline")
			}

			err = SendList(mainConn, dataConn, child.Filepath, maxDepth, currentDepth+1, child)
			if err != nil {
				log.Fatalf("rip : %s", err.Error())
			}

			err = dataConn.Close()
		}
	}
	currentNode.AddChildren(children)
	return nil
}

// parseAnswerList parse the asnwer from the list command, add absolute path to the global array of struct paths
func parseAnswerList(lines []string, base string, depth int) []*tree2.Node {
	children := make([]*tree2.Node, 0)
	for _, val := range lines {
		currentNode := &tree2.Node{}
		currentNode.Depth = depth
		if val[0] == 'd' {
			currentNode.IsDirectory = true
			currentNode.Filepath = base + val[strings.LastIndex(val, " ")+1:] + "/"
		} else {
			currentNode.IsDirectory = false
			currentNode.Filepath = base + val[strings.LastIndex(val, " ")+1:]
		}
		currentNode.Filename = val[strings.LastIndex(val, " ")+1:]
		children = append(children, currentNode)
	}
	return children

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

// getListLines read and parse the lines returned by the reader, return a string array of the line
func getListLines(readerDataConn *bufio.Reader) []string {
	lines := make([]string, 0)
	line, _, err := readerDataConn.ReadLine()
	for err == nil {
		lines = append(lines, string(line))
		line, _, err = readerDataConn.ReadLine()
	}
	if err != nil && err != io.EOF {
		log.Fatalf(err.Error())
	}
	return lines
}
