package ftpconn

import (
	"bufio"
	"errors"
	"fmt"
	"strings"
	"syscall"
	"time"
	constantFTP "tree-ftp/util/ftp/codes"
	commandFTPConstant "tree-ftp/util/ftp/commands"
	constant "tree-ftp/util/global"
)

// SendUser send FTP User command with the given user
func (conn *FTPConn) SendUser(user string) error {
	stringToSend := commandFTPConstant.User + " " + user + "\n"
	_, err := conn.MainConn.Write([]byte(stringToSend))
	if err != nil {
		return err
	}

	reply := make([]byte, constant.SizeAnswer)
	_, err = conn.MainConn.Read(reply)
	if err != nil && !strings.Contains(string(reply), constantFTP.CodeUserOk) {
		return err
	}
	return nil
}

// SendPass send ftp Pass command with the given password
func (conn *FTPConn) SendPass(password string) error {
	stringToSend := commandFTPConstant.Pass + " " + password + "\n"

	_, err := conn.MainConn.Write([]byte(stringToSend))
	if err != nil {
		return err
	}
	return nil
}

// SendList send ftp List command and control the result with readerMainConn
func (conn *FTPConn) SendList(readerMainConn *bufio.Reader) error {
	stringToSend := commandFTPConstant.List + "\n"
	var err error = nil
	for i := 0; i < constant.MaxRetry; i++ {
		_, err := conn.MainConn.Write([]byte(stringToSend))
		if err != nil {
			if errors.Is(err, syscall.EPIPE) {
				return err
			}
			fmt.Printf("Err SendList %s, retrying in 10s\n", err.Error())
			time.Sleep(constant.TimeBeforeRetry * time.Second)
			continue
		}
		line, _, err := readerMainConn.ReadLine()
		if err != nil || !strings.Contains(string(line), constantFTP.CodeOkList) {
			var stringToPrint = ""
			if err != nil {
				stringToPrint = fmt.Sprintf("Err SendList %s, retrying in 10s\n", err.Error())
			} else {
				stringToPrint = fmt.Sprintf("Wrong code response : %s\n", string(line))
			}
			fmt.Printf(stringToPrint)
			time.Sleep(constant.TimeBeforeRetry * time.Second)
			continue
		} else {
			break
		}
	}
	if err != nil {
		fmt.Printf("CWD failed 3 times, aborting")
		return err
	}

	return nil
}

// SendCwd Send Ftp CWD command with the given filepath, and control the result on readerMainConn
func (conn *FTPConn) SendCwd(readerMainConn *bufio.Reader, filepath string) error {
	stringToSend := commandFTPConstant.Cwd + " " + filepath + "\n"
	var err error = nil
	for i := 0; i < constant.MaxRetry; i++ {
		_, err := conn.MainConn.Write([]byte(stringToSend))
		if err != nil {
			if errors.Is(err, syscall.EPIPE) {
				return err
			}
			fmt.Printf("Err SendCWD %s, retrying in 10s\n", err.Error())
			time.Sleep(constant.TimeBeforeRetry * time.Second)
			continue
		}
		line, _, err := readerMainConn.ReadLine()
		if err != nil || !strings.Contains(string(line), constantFTP.CodeCWDOk) {
			var stringToPrint = ""
			if err != nil {
				stringToPrint = fmt.Sprintf("Err SendCWD %s, retrying in 10s\n", err.Error())
			} else {
				stringToPrint = fmt.Sprintf("Wrong code response : %s\n", string(line))
			}
			fmt.Printf(stringToPrint)
			time.Sleep(constant.TimeBeforeRetry * time.Second)
			continue
		}
		break
	}
	if err != nil {
		fmt.Printf("CWD failed 3 times, aborting")
		return err
	}
	return nil
}

// SendPasv Send Ftp Pasv and control the result on readerMainConn
func (conn *FTPConn) SendPasv(readerMainConn *bufio.Reader) (string, error) {
	stringToSend := commandFTPConstant.Pasv + "\n"
	var err error = nil
	var line []byte
	for i := 0; i < constant.MaxRetry; i++ {
		_, err := conn.MainConn.Write([]byte(stringToSend))
		if err != nil {
			if errors.Is(err, syscall.EPIPE) {
				return "", err
			}
			fmt.Printf("Err PASV %s, retrying in 10s\n", err.Error())
			time.Sleep(constant.TimeBeforeRetry * time.Second)
			continue
		}
		line, _, err = readerMainConn.ReadLine()
		if err != nil || !strings.Contains(string(line), constantFTP.CodePasvOk) {
			var stringToPrint = ""
			if err != nil {
				stringToPrint = fmt.Sprintf("Err Pasv %s, retrying in 10s\n", err.Error())
			} else {
				stringToPrint = fmt.Sprintf("Wrong code response : %s\n", string(line))
			}
			fmt.Printf(stringToPrint)
			time.Sleep(constant.TimeBeforeRetry * time.Second)
			continue
		}
		break
	}
	if err != nil {
		fmt.Printf("CWD failed 3 times, aborting")
		return "", err
	}
	return string(line), nil
}
