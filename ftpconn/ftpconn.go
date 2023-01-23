package ftpconn

import (
	"bufio"
	"errors"
	"fmt"
	"golang.org/x/term"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
	tree2 "tree-ftp/tree"
	constantFTP "tree-ftp/util/ftp/codes"
	constant "tree-ftp/util/global"
)

type FTPConn struct {
	MainConn *net.TCPConn
}

var retryLogin = 0

// UserConn Init TCP conn with given user and pwd, if both are not precised, anonymous is the default/**
func (conn *FTPConn) UserConn(user string, pwd string) error {
	fmt.Println("User connexion")

	err := conn.SendUser(user)
	if err != nil {
		return err
	}

	err = conn.SendPass(pwd)
	if err != nil {
		return err
	}

	reply := make([]byte, constant.SizeAnswer)
	_, err = conn.MainConn.Read(reply)

	if err != nil || !strings.Contains(string(reply), constantFTP.CodeLoginOk) {
		if strings.Contains(string(reply), constantFTP.CodeLoginNotOk) {
			if retryLogin == constant.MaxRetry {
				log.Fatalf("too many retry for login/password ")
			}
			fmt.Println("Wrong password/login !")
			fmt.Println("Enter your login")
			reader := bufio.NewReader(os.Stdin)
			login, _ := reader.ReadString('\n')
			login = strings.Replace(login, "\n", "", -1)
			fmt.Println("Enter your password")
			password, _ := term.ReadPassword(0)
			retryLogin++
			return conn.UserConn(login, string(password))

		}
		return err
	}
	fmt.Println("User connexion successful")
	return nil
}

// GetDataConn Create a new data connection from mainConn, that send PASV/**
func (conn *FTPConn) GetDataConn() (*FTPConn, error) {
	readerMainConn := bufio.NewReader(conn.MainConn)
	line, err := conn.SendPasv(readerMainConn)
	if err != nil {
		return nil, err
	}

	err, ipAddrDataConn, portDataConn := getIPAndPortFromResponse(line)
	if err != nil {
		return nil, err
	}

	ip := &net.TCPAddr{
		IP:   net.ParseIP(ipAddrDataConn),
		Port: portDataConn,
	}
	connData, err := net.DialTCP(constant.TcpString, nil, ip)
	if err != nil {
		for i := 0; i < constant.MaxRetry; i++ {
			time.Sleep(constant.TimeBeforeRetry * time.Second)
			fmt.Printf("Failed to dial %s retrying in 10secondes\n", ip.IP)
			connData, err = net.DialTCP(constant.TcpString, nil, ip)

		}
		return nil, err
	}

	return &FTPConn{MainConn: connData}, nil
}

// GetIpFromURL return the found IP for the given addressServer and port
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

// ListFtpFiles send the command list to the mainConn, and CWD to all directories returned to recursively call ListFtpFiles with them**/
func (conn *FTPConn) ListFtpFiles(dataConn *FTPConn, base string, maxDepth int, currentDepth int, currentNode *tree2.Node) (error, string) {
	readerMainConn := bufio.NewReader(conn.MainConn)
	err := conn.SendCwd(readerMainConn, base)
	if err != nil {
		return err, currentNode.Filepath
	}
	fmt.Printf("SendList %s \n", base)
	err = conn.SendList(readerMainConn)
	if err != nil {
		return err, base
	}

	readerDataConn := bufio.NewReader(io.Reader(dataConn.MainConn))
	lines, err := getListLines(readerDataConn)
	if err != nil {
		return err, base
	}
	children := parseAnswerList(lines, base, currentDepth)
	line, _, err := readerMainConn.ReadLine()
	if (err != nil && err != io.EOF) || !strings.Contains(string(line), constantFTP.CodeComingList) {
		return err, base
	}
	if currentDepth == maxDepth {
		currentNode.AddChildren(children)
		return nil, base
	}
	for _, child := range children {
		if child.IsDirectory {
			dataConn, err = conn.GetDataConn()
			if err != nil {
				return err, currentNode.Filepath
			}

			err := conn.SendCwd(readerMainConn, child.Filepath)
			if err != nil {
				return err, currentNode.Filepath
			}
			err, lastDirVisited := conn.ListFtpFiles(dataConn, child.Filepath, maxDepth, currentDepth+1, child)
			if err != nil {
				return err, lastDirVisited
			}

			err = dataConn.MainConn.Close()
			if err != nil {
				return err, currentNode.Filepath
			}
		}
	}
	currentNode.AddChildren(children)
	return err, ""
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
	if len(reply) == 0 {
		return errors.New("reply bizaroide %s"), "", 0
	}
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
func getListLines(readerDataConn *bufio.Reader) ([]string, error) {
	lines := make([]string, 0)
	line, _, err := readerDataConn.ReadLine()
	for err == nil {
		lines = append(lines, string(line))
		line, _, err = readerDataConn.ReadLine()
	}
	if err != nil && err != io.EOF {
		return nil, err
	}
	return lines, nil
}
