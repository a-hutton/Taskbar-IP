package main

import (
	"fmt"
	glSystray "github.com/getlantern/systray"
	"github.com/xilp/systray"
	"golang.design/x/clipboard"
	"ipIcons/icons"
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
		createIpTrayIcons(Addresses, -1)
	}, func() {
		println("On Exit")
	})

}

// createIpTrayIcons displays one of the user's current IP addresses, specified by index.
// if index is negative, it will use the user's preference, or 0 if it isn't connected
func createIpTrayIcons(addresses []string, index int) {
	if index < 0 {
		config := ReadConfig()
		if slices.Contains(Addresses, config.PreferredIp) {
			index = slices.Index(Addresses, config.PreferredIp)
		} else {
			index = 0
		}
	}

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
	glSystray.SetIcon(icons.More)
	glSystray.SetTitle("Refresh Address")
	glSystray.SetTooltip("Do refreshing")
	refreshBtn := glSystray.AddMenuItem("Refresh", "Reload to refresh IPs")
	quitBtn := glSystray.AddMenuItem("Quit", "Quit application and remove from taskbar")
	addIpIcons()
	go func() {
		for {
			select {
			case <-refreshBtn.ClickedCh:
				println("Refresh clicked")
				Addresses = getAddressesRegex()
				removeIpIcons()
				removeControlIpsList()
				addIpIcons()
				createIpTrayIcons(Addresses, -1)
			case <-quitBtn.ClickedCh:
				glSystray.Quit()
			}

		}
	}()

}

func addIpIcons() {
	ipMenuItems = make([]*glSystray.MenuItem, 0)
	for i, address := range Addresses {
		addressButton := glSystray.AddMenuItem(address, fmt.Sprintf("Options for %s", address))
		ipMenuItems = append(ipMenuItems, addressButton)
		go func(button *glSystray.MenuItem, index int) {
			for {
				<-button.ClickedCh
				fmt.Printf("swtiching to %d and saving it as default\n", index)
				config := ReadConfig()
				config.PreferredIp = Addresses[index]
				config.Save()
				removeIpIcons()
				createIpTrayIcons(Addresses, index)
			}

		}(addressButton, i)
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
