package rmbridge

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"github.com/melbahja/goph"
	"golang.org/x/crypto/ssh"
)

var client *goph.Client

func askIsHostTrusted(host string, key ssh.PublicKey) bool {

	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("Unknown Host: %s \nFingerprint: %s \n", host, ssh.FingerprintSHA256(key))
	fmt.Print("Would you like to add it? type yes or no: ")

	a, err := reader.ReadString('\n')

	if err != nil {
		log.Fatal(err)
	}

	return strings.ToLower(strings.TrimSpace(a)) == "yes"
}
func verifyHost(host string, remote net.Addr, key ssh.PublicKey) error {

	//
	// If you want to connect to new hosts.
	// here your should check new connections public keys
	// if the key not trusted you shuld return an error
	//

	// hostFound: is host in known hosts file.
	// err: error if key not in known hosts file OR host in known hosts file but key changed!
	hostFound, err := goph.CheckKnownHost(host, remote, key, "")

	// Host in known hosts but key mismatch!
	// Maybe because of MAN IN THE MIDDLE ATTACK!
	if hostFound && err != nil {

		return err
	}

	// handshake because public key already exists.
	if hostFound && err == nil {

		return nil
	}

	// Ask user to check if he trust the host public key.
	if askIsHostTrusted(host, key) == false {

		// Make sure to return error on non trusted keys.
		return errors.New("you typed no, aborted")
	}

	// Add the new host to known hosts file.
	return goph.AddKnownHost(host, remote, key, "")
}

// Connect will try to establish a connection to the reMarkable
func connect(passphrase string, remote string) bool {

	user := "root"
	port := uint(22)
	var err error
	client, err = goph.NewConn(&goph.Config{
		User:     user,
		Addr:     remote,
		Port:     port,
		Auth:     goph.Password(passphrase),
		Callback: verifyHost,
	})
	if err != nil {
		log.Println("could not create hostkeycallback function: ", err)
		return false
	}
	log.Println("successfully connected at " + remote)
	return true
}

// ConnectUSB will try to establish a connection to the reMarkable via USB
func ConnectUSB(passphrase string) bool {
	return connect(passphrase, "10.11.99.1")
}

// ConnectWifi will try to establish a connection to the reMarkable via Wifi at a given ip address
func ConnectWifi(passphrase string, ipAdress string) bool {
	return connect(passphrase, ipAdress)
}
