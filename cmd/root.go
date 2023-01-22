package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"net"
	"os"
	"tree-ftp/ftpconn"
	model "tree-ftp/tree"
	constant "tree-ftp/util/global"
)

var rootTree = model.Node{}
var (
	addressServer string
	port          int
	user          string
	password      string
	maxDepth      int
	directoryOnly bool
	fullPath      bool
	toJson        bool
	jsonFile      string
	rootCmd       = &cobra.Command{
		Use:   "tree-ftp",
		Short: "Display a tree-like output of the content of a ftp server ",
		Run: func(cmd *cobra.Command, args []string) {

			Addr, err := ftpconn.GetIpFromURL(port, addressServer)
			if err != nil {
				log.Fatalf(err.Error())
			}

			conn, err := net.DialTCP(constant.TcpString, nil, Addr)
			ftpConn := ftpconn.FTPConn{MainConn: conn}
			if err != nil {
				log.Fatal(err.Error(), "are you sure your port is correct ?")
			}
			reply := make([]byte, constant.SizeAnswer)
			_, err = conn.Read(reply)
			if err != nil {
				log.Fatal(err.Error())
			}
			err = ftpConn.UserConn(user, password)
			if err != nil {
				log.Fatalf(err.Error())
			}
			dataConn, err := ftpConn.GetDataConn()
			if err != nil {
				log.Fatalf(err.Error())
			}
			rootTree.InitNode()
			rootTree.Filepath = "/"
			rootTree.Filename = "/"
			rootTree.IsDirectory = true
			rootTree.Depth = 0
			err, lastDirVisited := ftpConn.ListFtpFiles(dataConn, rootTree.Filepath, maxDepth, 1, &rootTree)
			if err != nil {
				fmt.Printf("%s while trying to visit %s\n Seems like the connection to server is lost \naborting\n", err, lastDirVisited)
				return
			}
			err = dataConn.MainConn.Close()
			if err != nil {
				log.Fatalf("Err while closing conn\n")
			}
			if toJson {
				marshal, err := json.Marshal(rootTree)
				if err != nil {
					log.Fatalf("Unable to marshal\n")
				}
				file, err := os.Create(jsonFile)
				if err != nil {
					log.Fatalf("Unable to create %s\n", jsonFile)
				}
				var outJson bytes.Buffer
				err = json.Indent(&outJson, marshal, "", "	")
				_, err = file.Write(outJson.Bytes())
				if err != nil {
					log.Fatalf("Unable to write to %s\n", jsonFile)
				}
				file.Close()
			} else {
				rootTree.DisplayTree(fullPath, directoryOnly)
			}
			return
		},
	}
)

func Execute() {
	rootCmd.Flags().StringVar(&addressServer, "addressServer", "", "Address to server")
	rootCmd.Flags().IntVar(&port, "port", constant.DefaultPortTCP, "Port to access ftp server")
	rootCmd.Flags().StringVar(&user, "user", "anonymous", "User for connexion")
	rootCmd.Flags().StringVar(&password, "password", "anonymous", "Password for connexion")
	rootCmd.Flags().IntVar(&maxDepth, "maxDepth", constant.DefaultMaxDepth, "Max depths of tree")
	rootCmd.Flags().BoolVar(&fullPath, "fullPath", false, "Display fullpath of files")
	rootCmd.Flags().BoolVar(&directoryOnly, "directoryOnly", false, "Display directories only")

	rootCmd.Flags().StringVar(&jsonFile, "jsonFile", "", "Path for json file")
	rootCmd.Flags().BoolVar(&toJson, "toJson", false, "Output is directed in a file as a json")

	err := rootCmd.MarkFlagRequired("addressServer")
	if err != nil {
		log.Fatalf(err.Error())
	}
	rootCmd.MarkFlagsRequiredTogether("user", "password")
	err = rootCmd.MarkFlagRequired("addressServer")
	if err != nil {
		log.Fatalf(err.Error())
	}
	rootCmd.MarkFlagsRequiredTogether("toJson", "jsonFile")
	err = rootCmd.MarkFlagRequired("addressServer")
	if err != nil {
		log.Fatalf(err.Error())
	}
	err = rootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
