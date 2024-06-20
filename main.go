package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/netip"
	"os"
	"strconv"
	"strings"

	prompt "github.com/c-bata/go-prompt"
	_ "github.com/mattn/go-sqlite3"
)

func lookUpIpAddress(ip string) ([]string, error) {
	db, err := sql.Open("sqlite3", "./db/cliproxy.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	query := `
	SELECT d.device_name
	FROM devices d
	JOIN ip_addresses ip ON d.device_id = ip.device_id
	WHERE ip.ip_address = ?;
	`
	rows, err := db.Query(query, ip)
	if err != nil {
		return nil, err
	}
	var devs []string
	for rows.Next() {
		var dev string
		err = rows.Scan(&dev)
		if err != nil {
			return nil, err
		}
		devs = append(devs, dev)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return devs, nil
}
func lookupHostName(hostname string) (string, error) {
	var ip string

	return ip, nil
}

func addDevice(hostname, mgmtip string) error {
	db, err := sql.Open("sqlite3", "./db/cliproxy.db")
	if err != nil {
		return err
	}
	defer db.Close()

	insertDevice := `INSERT INTO devices (device_name, management_ip) VALUES (?, ?);`
	_, err = db.Exec(insertDevice, hostname, mgmtip)
	if err != nil {
		return err
	}
	return nil
}
func checkIsIP(input string) bool {
	_, err := netip.ParseAddr(input)
	if err != nil {
		return false
	}
	return true
}
func connect(arg string) error {
	if checkIsIP(arg) {
		res, err := lookUpIpAddress(arg)
		if err != nil {
			return err
		}
		fmt.Printf("Multiple Devices found with IP address %s.\n Please select from the list below:\n", arg)
		for i, dev := range res {
			fmt.Printf("%d:  %s\n", i, dev)
		}
		var input string
		_, err = fmt.Scan(&input)
		if err != nil {
			return err
		}
		selection, err := strconv.Atoi(input)
		if err != nil {
			return fmt.Errorf("Invalid Selection\n")
		}
		if selection < 0 || selection > len(res)-1 {
			return fmt.Errorf("Invalid Selection")
		}
		fmt.Printf("Connecting to %s\n", res[selection])
	}
	return nil
}
func executor(t string) {
	cmd := strings.Split(t, " ")
	if cmd[0] == "connect" {

		//err := sshclient.Connect(ssh.ClientConfig{}, cmd[1])
		err := connect(cmd[1])
		if err != nil {
			log.Printf("Error connecting to %s\nError: %v", cmd[1], err)
		}
	}
	if cmd[0] == "add" && cmd[1] == "device" {
		if len(cmd) < 4 {
			log.Printf("Incorrect number of arguments. Please use \"add device <hostname> <mgmt ip>\"\n")
		}
		err := addDevice(cmd[2], cmd[3])
		if err != nil {
			log.Printf("Error adding device to Database:\n %v", err)
		}
	}
	if cmd[0] == "lookup" {
		res, err := lookUpIpAddress(cmd[1])
		if err != nil {
			log.Print("Database error during lookup:\n", err)
		}
		for _, dev := range res {
			fmt.Printf("%s: %s\n", cmd[1], dev)
		}
	}
	if t == "exit" {
		os.Exit(0)
	}
}

func completer(t prompt.Document) []prompt.Suggest {
	return []prompt.Suggest{
		{Text: "connect"},
		{Text: "exit"},
		{Text: "add device"},
		{Text: "lookup"},
	}
}

func main() {

	p := prompt.New(
		executor,
		completer,
	)
	p.Run()
}
