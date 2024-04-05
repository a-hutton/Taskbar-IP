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

var Addresses []string
var ipIcons []*systray.Systray
var ipMenuItems []*glSystray.MenuItem

func main() {
	Addresses = getAddressesRegex()
	err := clipboard.Init()
	if err != nil {
		println("Error initialising clipboard")
	}
	glSystray.Run(func() {
		createControlIcon()
		createIpTrayIcons(Addresses, 0)
	}, func() {
		println("On Exit")
	})

}

func createIpTrayIcons(addresses []string, index int) {
	allAddressesHint := strings.Join(addresses, "\n")
	address := addresses[index]
	numbers := strings.Split(address, ".")
	slices.Reverse(numbers)
	ipIcons = make([]*systray.Systray, 0)

	var tray *systray.Systray
	for i, number := range numbers {
		tray = systray.New("", "localhost", 6320+i)
		ipIcons = append(ipIcons, tray)
		tray.OnClick(func() {
			clipboard.Write(clipboard.FmtText, []byte(address))
		})
		err := tray.Show(fmt.Sprintf("icons/icon-%s.ico", number), allAddressesHint)
		if err != nil {
			panic(err)
		}
	}

	go func() {
		if err := tray.Run(); err != nil {
			println(err.Error())
		}
	}()
}

func removeIpIcons() {
	for _, trayIcon := range ipIcons {
		if err := trayIcon.Stop(); err != nil {
			println(err.Error())
		}
	}
}

func createControlIcon() {
	glSystray.SetIcon(icon.Data)
	glSystray.SetTitle("Refresh Address")
	glSystray.SetTooltip("Do refreshing")
	refreshBtn := glSystray.AddMenuItem("Refresh", "Reload to refresh IPs")
	addIpIcons()
	go func() {
		for {
			<-refreshBtn.ClickedCh
			println("Refresh clicked")
			Addresses = getAddressesRegex()
			removeIpIcons()
			removeControlIpsList()
			addIpIcons()
			createIpTrayIcons(Addresses, 0)
		}
	}()

}

func addIpIcons() {
	ipMenuItems = make([]*glSystray.MenuItem, 0)
	for i, address := range Addresses {
		btn := glSystray.AddMenuItem(fmt.Sprintf("Change to %s", address), "")
		ipMenuItems = append(ipMenuItems, btn)
		go func(button *glSystray.MenuItem, index int) {
			for {
				<-button.ClickedCh
				fmt.Printf("swtiching to %d\n", index)
				removeIpIcons()
				createIpTrayIcons(Addresses, index)
			}

		}(btn, i)
	}
}

func removeControlIpsList() {
	for _, item := range ipMenuItems {
		item.Disable()
		item.Hide()
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
