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
	NETBOOK_BRIGHTNESS_FILE     = "/sys/class/backlight/intel_backlight/brightness"
	NETBOOK_MAX_BRIGHTNESS_FILE = "/sys/class/backlight/intel_backlight/max_brightness"
	SUSPEND_FILE                = "/sys/power/state"
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

func GetMaxBrightness() float64 {
	contents, err := ioutil.ReadFile(NETBOOK_MAX_BRIGHTNESS_FILE)
	if err != nil {
		panic(err)
	}

	max_brightness, err := strconv.ParseFloat(strings.TrimSpace(string(contents)), 64)
	if err != nil {
		panic(err)
	}

	return float64(max_brightness)
}

func (f foo) GetBrightness() (string, int, *dbus.Error) {
	fmt.Println(f)

	contents, err := ioutil.ReadFile(NETBOOK_BRIGHTNESS_FILE)
	if err != nil {
		panic(err)
	}

	brightness, err := strconv.ParseFloat(strings.TrimSpace(string(contents)), 64)
	if err != nil {
		panic(err)
	}

	fmt.Println(brightness, GetMaxBrightness())

	return string(f), int(brightness * 100.0 / GetMaxBrightness()), nil
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
