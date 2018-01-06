package main

import (
	"fmt"
	"github.com/godbus/dbus"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

const (
	NETBOOK_BRIGHTNESS_FILE = "/sys/class/backlight/intel_backlight/brightness"
	SUSPEND_FILE            = "/sys/power/state"
)

type foo string

func (f foo) GetDevice() (string, *dbus.Error) {
	fmt.Println(f)
	return string(f), nil
}

func (f foo) FooPlus(what string) (string, *dbus.Error) {
	r := string(f) + " plus < " + what + " >"
	fmt.Println(r)
	return r, nil
}

func (f foo) GetBrightness() (string, float64, *dbus.Error) {
	fmt.Println(f)

	contents, err := ioutil.ReadFile(NETBOOK_BRIGHTNESS_FILE)
	if err != nil {
		panic(err)
	}

	brightness, err := strconv.ParseFloat(strings.TrimSpace(string(contents)), 64)
	if err != nil {
		panic(err)
	}

	fmt.Println(brightness)

	return string(f), float64(brightness), nil
}

func main() {
	conn, err := dbus.SystemBus()
	if err != nil {
		panic(err)
	}
	reply, err := conn.RequestName("org.cneo.cneod",
		dbus.NameFlagDoNotQueue)
	if err != nil {
		panic(err)
	}
	if reply != dbus.RequestNameReplyPrimaryOwner {
		fmt.Fprintln(os.Stderr, "name already taken")
		os.Exit(1)
	}

	var f foo
	if _, err := os.Stat(NETBOOK_BRIGHTNESS_FILE); !os.IsNotExist(err) {
		f = foo("/org/cneo/Netbook")
	} else {
		f = foo("/org/cneo/GTA02")
	}

	conn.Export(f, "/", "org.cneo.System")
	fmt.Println("Listening on org.cneo.cneod ...")
	select {}
}
