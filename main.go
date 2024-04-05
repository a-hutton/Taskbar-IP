package main

import (
	"fmt"
	glSystray "github.com/getlantern/systray"
	"github.com/getlantern/systray/example/icon"
	"github.com/xilp/systray"
	"golang.design/x/clipboard"
	"os/exec"
	"regexp"
	"slices"
	"strings"
)

var addresses []string

func main() {
	addresses = getAddressesRegex()
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

	go func() {
		if err = tray.Run(); err != nil {
			println(err.Error())
		}
	}()
	glSystray.Run(glOnReady, func() {})

}

func glOnReady() {
	glSystray.SetIcon(icon.Data)
	glSystray.SetTitle("Refresh Address")
	glSystray.SetTooltip("Do refreshing")
	refreshBtn := glSystray.AddMenuItem("Refresh", "Reload to refresh IPs")
	go func() {
		for {
			<-refreshBtn.ClickedCh
			println("Refresh clicked")
		}
	}()
	for _, address := range addresses {
		btn := glSystray.AddMenuItem(fmt.Sprintf("Change to %s", address), "")
		go func(button *glSystray.MenuItem, address string) {
			for {
				<-button.ClickedCh
				fmt.Printf("swtiching to %s\n", address)
			}

		}(btn, address)
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
	foundAddresses := make([]string, 0)
	for _, match := range ipMatches {
		if len(match) > 0 {
			index := re.SubexpIndex("address")
			newAddress := match[index]
			foundAddresses = append(foundAddresses, newAddress)
			fmt.Printf("Address:%s\n", match[index])
		}
	}
	return foundAddresses
}
