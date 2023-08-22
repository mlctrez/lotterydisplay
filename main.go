package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"

	owm "github.com/briandowns/openweathermap"
	"github.com/kardianos/service"
	"github.com/mlctrez/lotterydisplay/display"
	"github.com/mlctrez/servicego"
	"github.com/robfig/cron/v3"
)

type LotteryDisplay struct {
	Cron      *cron.Cron
	DG        *display.Group
	TempMin   string
	TempMax   string
	Temp      string
	FeelsLike string
	Time      string
}

func (ldb *LotteryDisplay) updateBrightness() {
	ldb.DG.Intensity = 4
	now := time.Now()
	if now.Hour() > 7 && now.Hour() < 22 {
		ldb.DG.Intensity = 15
	}
}

func (ldb *LotteryDisplay) addSchedule(spec string, cmd func()) {
	_, err := ldb.Cron.AddFunc(spec, cmd)
	if err != nil {
		log.Fatal(err)
	}
}

func (ldb *LotteryDisplay) updateTemp() {
	var err error
	var forecast *owm.CurrentWeatherData
	var key = os.Getenv("API_KEY")
	var zip = os.Getenv("ZIP_CODE")

	if forecast, err = owm.NewCurrent("F", "en", key); err == nil {
		if err = forecast.CurrentByZipcode(zip, "US"); err == nil {
			ldb.Temp = fmt.Sprintf("%3.0f", forecast.Main.Temp)
			ldb.FeelsLike = fmt.Sprintf("%3.0f", forecast.Main.FeelsLike)
			ldb.TempMin = fmt.Sprintf("%3.0f", forecast.Main.TempMin)
			ldb.TempMax = fmt.Sprintf("%3.0f", forecast.Main.TempMax)
		}
	}
	if err != nil {
		log.Println(err)
	}
}

func (ldb *LotteryDisplay) updateTime() {
	now := time.Now()
	ldb.Time = now.Format("0304")
	//fmt.Println("time", ldb.Time)
}

func (ldb *LotteryDisplay) updateDisplay() {
	ldb.DG.Clear()
	for i := 0; i < 8; i++ {
		ldb.DG.SetDecimalPoint(i, false)
	}

	dpOn := (time.Now().Second()/5)%2 == 0
	if dpOn {
		ldb.DG.WriteString(ldb.Time + " " + ldb.TempMax)
	} else {
		ldb.DG.WriteString(ldb.Time + " " + ldb.Temp)
	}
	ldb.DG.SetDecimalPoint(3, dpOn)

	err := Transmit(os.Getenv("BROADCAST_ADDRESS"), ldb.DG.Packet())
	if err != nil {
		log.Println(err)
	}
}

type program struct {
	servicego.Defaults
	ldb *LotteryDisplay
}

func (p *program) Start(s service.Service) error {

	p.ldb = &LotteryDisplay{
		Cron: cron.New(cron.WithSeconds()),
		DG:   &display.Group{Intensity: 4},
		Temp: "--",
	}

	p.ldb.updateBrightness()
	p.ldb.updateTemp()
	p.ldb.updateTime()

	p.ldb.addSchedule("* * * * * *", p.ldb.updateDisplay)
	p.ldb.addSchedule("0 * * * * *", p.ldb.updateTime)
	p.ldb.addSchedule("0 */30 * * * *", p.ldb.updateTemp)
	p.ldb.addSchedule("0 * * * * *", p.ldb.updateBrightness)
	p.ldb.Cron.Start()

	return nil
}

func Transmit(address string, data []byte) (err error) {
	var udpAddr *net.UDPAddr
	var con *net.UDPConn

	defer func() {
		if con != nil {
			conErr := con.Close()
			if err == nil {
				err = conErr
			}
		}
	}()

	if udpAddr, err = net.ResolveUDPAddr("udp4", address); err == nil {
		if con, err = net.DialUDP("udp4", nil, udpAddr); err == nil {
			_, err = con.Write(data)
		}
	}
	return
}

func (p *program) Stop(s service.Service) error {
	p.ldb.Cron.Stop()
	return nil
}

func main() {
	readEnvironment()
	servicego.Run(&program{})
}

func readEnvironment() {
	if _, err := os.Stat(".environment"); os.IsNotExist(err) {
		return
	}
	if contents, err := os.ReadFile(".environment"); err != nil {
		log.Fatal("error opening .environment file")
	} else {
		scanner := bufio.NewScanner(bytes.NewBuffer(contents))
		for scanner.Scan() {
			parts := strings.SplitN(strings.TrimSpace(scanner.Text()), "=", 2)
			if len(parts) == 2 {
				err = os.Setenv(strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]))
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	}

}
