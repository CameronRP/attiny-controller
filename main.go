package main

import (
	"errors"
	"log"
	"sync"
	"time"

	"github.com/TheCacophonyProject/window"
	arg "github.com/alexflint/go-arg"
	i2c "golang.org/x/exp/io/i2c"
)

const (
	// 3 was just a randomly chosen as the number for the attiny to return
	// to indicate its presence.
	magicReturn = 0x03

	// Check for the ATtiny for up to a minute.
	connectAttempts        = 20
	connectAttemptInterval = 3 * time.Second

	watchdogTimerAddress = 0x12
	sleepAddress         = 0x11

	// How long to wait before checking the recording window. This
	// gives time to do something with the device before it turns off.
	initialGracePeriod = 20 * time.Minute
)

var version = "<not set>"

var mu sync.Mutex // Protecting stayOnUntil

var stayOnUntil = time.Now()

type Args struct {
	ConfigFile string `arg:"-c,--config" help:"path to configuration file"`
	SkipWait   bool   `arg:"-s,--skip-wait" help:"will not wait for the date to update"`
}

func (Args) Version() string {
	return version
}

func procArgs() Args {
	var args Args
	args.ConfigFile = "/etc/cacophony-attiny.yaml"
	arg.MustParse(&args)
	return args
}

func main() {
	err := runMain()
	if err != nil {
		log.Fatal(err)
	}
}

func runMain() error {
	log.SetFlags(0) // Removes default timestamp flag

	args := procArgs()
	log.Printf("running version: %s", version)

	err := StartService()
	if err != nil {
		return err
	}
	log.Println("started service")

	attiny, err := i2c.Open(&i2c.Devfs{Dev: "/dev/i2c-1"}, 0x04)
	if err != nil {
		return err
	}
	defer attiny.Close()

	if !connected(attiny) {
		log.Println("attiny not present")
		return nil
	}
	log.Println("connected to attiny")

	go updateWatchdogTimer(attiny)

	if !args.SkipWait {
		log.Printf("waiting for %s before checking recording window", initialGracePeriod)
		time.Sleep(initialGracePeriod)
	}

	var lastloggedError string = ""

	for {
		conf, err := ParseAttinyConfigFile(args.ConfigFile)
		if err == nil {
			window := window.New(conf.PiWakeUpTime, conf.PiSleepTime)
			minutesUntilActive := int(window.Until().Minutes())
			log.Printf("minutes until active %d", minutesUntilActive)
			if shouldTurnOff(minutesUntilActive) {
				minutesUntilActive = minutesUntilActive - 2
				lb := byte(minutesUntilActive / 256)
				rb := byte(minutesUntilActive % 256)
				err = attiny.Write([]byte{sleepAddress, lb, rb})
				if err != nil {
					log.Fatal(err)
				}
			}
		} else if lastloggedError != err.Error() {
			lastloggedError = err.Error()
			log.Printf("attiny-config error: %s", err.Error())
		}
		time.Sleep(time.Minute * 5)
	}
}

func shouldTurnOff(minutesUntilActive int) bool {
	mu.Lock()
	defer mu.Unlock()
	if time.Now().Before(stayOnUntil) {
		return false
	}
	return minutesUntilActive > 15
}

// SetStayOnUntil will not trigger the pi to turn off through the attiny until the given time
func SetStayOnUntil(newTime time.Time) error {
	if time.Until(newTime) > 12*time.Hour {
		return errors.New("can not delay over 12 hours")
	}
	mu.Lock()
	stayOnUntil = newTime
	mu.Unlock()
	log.Println("staying on until", newTime.Format(time.UnixDate))
	return nil
}

func connected(attiny *i2c.Device) bool {
	for i := 0; i < connectAttempts; i++ {
		time.Sleep(connectAttemptInterval)

		buf := make([]byte, 1)
		attiny.Read(buf)
		if buf[0] == magicReturn {
			return true
		}
	}
	return false
}

func updateWatchdogTimer(attiny *i2c.Device) {
	log.Println("sending watchdog timer updates")
	for {
		err := attiny.Write([]byte{watchdogTimerAddress})
		if err != nil {
			log.Fatal(err)
		}
		time.Sleep(time.Minute)
	}
}
