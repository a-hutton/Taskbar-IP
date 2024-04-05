package main

import (
	"fmt"
	"github.com/xilp/systray"
	"golang.design/x/clipboard"
	"os/exec"
	"regexp"
	"slices"
	"strings"
)

func main() {
	addresses := getAddressesRegex()
	allAddressesHint := strings.Join(addresses, "\n")
	address := addresses[0]
	numbers := strings.Split(address, ".")
	slices.Reverse(numbers)

	err := clipboard.Init()

	var tray *systray.Systray
	for i, number := range numbers {
		tray = systray.New("", "localhost", 6320+i)
		tray.OnClick(func() {
			clipboard.Write(clipboard.FmtText, []byte(address))
		})
		err := tray.Show(fmt.Sprintf("icons/icon-%s.ico", number), allAddressesHint)
		if err != nil {
			panic(err)
		}
	}

	if err = tray.Run(); err != nil {
		println(err.Error())
	}
}

func getAddressesRegex() []string {
	cmd := exec.Command("ipconfig")
	output, err := cmd.Output()
	if err != nil {
		println("Error executing command")
	}
	re := regexp.MustCompile(`Address[\s.]+:\s(?P<address>\d{1,3}.\d{1,3}\.\d{1,3}\.\d{1,3})`)
	ipMatches := re.FindAllStringSubmatch(string(output), -1)
	addresses := make([]string, 0)
	for _, match := range ipMatches {
		if len(match) > 0 {
			index := re.SubexpIndex("address")
			newAddress := match[index]
			addresses = append(addresses, newAddress)
			fmt.Printf("Address:%s\n", match[index])
		}
	}
	return addresses
}
