package cmd

import (
	"log"
	"net"

	"github.com/spf13/cobra"
)

const minimumArgs = 2
const TcpString = "tcp"

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

			conn, err := net.DialTCP(TcpString, nil, addr)

			if err != nil {
				log.Fatal(err.Error(), "are you sure your port is correct ?")
			}
			reply := make([]byte, 1024)
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

			err = sendList(conn, dataConn, "/", maxDepth, 0)
			if err != nil {
				log.Fatal(err)
			}
			err = dataConn.Close()
			if err != nil {
				log.Fatalf("wtf bro\n")
			}
			tree()
			return
		},
	}
)

func Execute() {
	rootCmd.Flags().StringVar(&addressServer, "addressServer", "", "Address to server")
	rootCmd.Flags().IntVar(&port, "port", 21, "Port to access ftp server")
	rootCmd.Flags().StringVar(&user, "user", "anonymous", "User for connexion")
	rootCmd.Flags().StringVar(&password, "password", "anonymous", "Password for connexion")
	rootCmd.Flags().IntVar(&maxDepth, "maxDepth", -1, "Max depths of tree")
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
