package cmd

import (
	"github.com/spf13/cobra"
	"log"
	"net"
	model "tree-ftp/tree"
	constant "tree-ftp/util/global"
)

const minimumArgs = 2

var rootTree = model.Node{}

var (
	addressServer string
	port          int
	user          string
	password      string
	maxDepth      int
	rootCmd       = &cobra.Command{
		Use:   "tree-ftp",
		Short: "Display a tree-like output of the content of a ftp server ",
		Run: func(cmd *cobra.Command, args []string) {

			addr, err := GetIpFromURL()
			if err != nil {
				log.Fatalf(err.Error())
			}

			conn, err := net.DialTCP(constant.TcpString, nil, addr)

			if err != nil {
				log.Fatal(err.Error(), "are you sure your port is correct ?")
			}
			reply := make([]byte, constant.SizeAnswer)
			_, err = conn.Read(reply)
			if err != nil {
				log.Fatal(err.Error())
			}
			err = UserConn(user, password, conn)
			if err != nil {
				log.Fatalf(err.Error())
			}
			dataConn, err := GetDataConn(conn)
			if err != nil {
				log.Fatalf(err.Error())
			}
			rootTree.InitNode()
			rootTree.Filepath = "/"
			rootTree.Filename = "/"
			rootTree.IsDirectory = true
			rootTree.Depth = 0
			err = sendList(conn, dataConn, rootTree.Filepath, maxDepth, 1, &rootTree)
			if err != nil {
				log.Fatal(err)
			}
			err = dataConn.Close()
			if err != nil {
				log.Fatalf("Err while closing conn\n")
			}
			rootTree.DisplayTree()
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
	rootCmd.MarkFlagsRequiredTogether("addressServer", "port")
	rootCmd.MarkFlagsRequiredTogether("user", "password")
	err := rootCmd.MarkFlagRequired("addressServer")
	if err != nil {
		log.Fatalf(err.Error())
	}
	err = rootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
