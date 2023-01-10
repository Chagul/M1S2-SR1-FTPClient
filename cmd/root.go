package cmd

import (
	"fmt"
	"log"
	"net"

	"github.com/spf13/cobra"
)

const minimumArgs = 2
const TCP_STRING = "tcp"

var (
	addressServer string
	port          int
	user          string
	password      string
	rootCmd       = &cobra.Command{
		Use:   "tree-ftp",
		Short: "Display a tree-like output of the content of a ftp server ",
		Run: func(cmd *cobra.Command, args []string) {

			addr, err := GetIpFromURL(addressServer)
			if err != nil {
				log.Fatalf(err.Error())
			}

			conn, err := net.DialTCP(TCP_STRING, nil, addr)
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
			_, err = conn.Write([]byte("PWD\n"))
			fmt.Printf("Sending PWD\n")
			if err != nil {
				log.Fatalf(err.Error())
			}
			reply = make([]byte, 1024)
			_, err = conn.Read(reply)
			if err != nil {
				log.Fatalf(err.Error())
			}
			println(string(reply))
			dataConn, err := GetDataConn(conn)
			if err != nil {
				log.Fatalf(err.Error())
			}

			err = sendList(conn, dataConn)
			if err != nil {
				log.Fatal(err)
			}
			return
		},
	}
)

func Execute() {
	rootCmd.Flags().StringVar(&addressServer, "addressServer", "", "Address to server")
	rootCmd.Flags().IntVar(&port, "port", 21, "Port to access ftp server")
	rootCmd.Flags().StringVar(&user, "user", "anonymous", "User for connexion")
	rootCmd.Flags().StringVar(&password, "password", "anonymous", "Password for connexion")
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
