package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"net"
	"strings"
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

			_, err = conn.Write([]byte("PWD\n"))
			fmt.Printf("Sending PWD\n")
			if err != nil {
				log.Fatalf(err.Error())
			}
			reply := make([]byte, 1024)
			_, err = conn.Read(reply)
			if err != nil {
				log.Fatalf(err.Error())
			}
			println("response :", string(reply))

			err = UserConn("anonymous", "anonymous", conn)
			if err != nil {
				log.Fatalf(err.Error())
			}

			_, err = conn.Write([]byte("FEAT\n"))
			fmt.Printf("Sending FEAT\n")
			if err != nil {
				log.Fatalf(err.Error())
			}
			reply = make([]byte, 1024)
			_, err = conn.Read(reply)
			for strings.Contains(string(reply), "211") {
				if err != nil {
					log.Fatalf(err.Error())
				}
				println(string(reply))
				_, err = conn.Read(reply)
			}
			println(string(reply))

			return
		},
	}
)

func Execute() {
	rootCmd.Flags().StringVar(&addressServer, "addressServer", "", "Address to server")
	rootCmd.Flags().IntVar(&port, "port", 21, "Port to access ftp server")

	rootCmd.MarkFlagsRequiredTogether("addressServer", "port")
	err := rootCmd.MarkFlagRequired("addressServer")
	if err != nil {
		log.Fatalf(err.Error())
	}
	err = rootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
