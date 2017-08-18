package main

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/enicho/go-wbgw2/src/config"

	"log"
)

var cfg config.Config

const INTERVAL_PERIOD time.Duration = 15 * time.Minute

const MINUTE_TO_TICK int = 00

func init() {
	//please put discord guild id
	cfg = config.InitConfig()
}

type jobTicker struct {
	t *time.Timer
}

func getNextTickDuration() time.Duration {
	now := time.Now().UTC()
	nextTick := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), MINUTE_TO_TICK, 0, 0, time.UTC)
	for nextTick.Before(now) {
		nextTick = nextTick.Add(INTERVAL_PERIOD)
	}
	log.Println(nextTick)
	return nextTick.Sub(time.Now().UTC())
}

func NewJobTicker() jobTicker {
	fmt.Println("running ticker")
	return jobTicker{time.NewTimer(getNextTickDuration())}
}

func (jt jobTicker) updateJobTicker() {
	now := time.Now().UTC()

	//hard
	hardHrInt := now.Hour() % cfg.Schedules.TimeInfo.HardInterval
	hardHr := "0" + strconv.Itoa(hardHrInt)
	hardMn := "00"

	if cfg.Schedules.HardBosses[hardHr+hardMn] != nil {
		// log.Println(cfg.Schedules.HardBosses[hardHr+hardMn].Name)
		notifier(cfg.Schedules.HardBosses[hardHr+hardMn].Name, cfg.Schedules.HardBosses[hardHr+hardMn].Location)
	}

	//medium
	midHrInt := now.Hour() % cfg.Schedules.TimeInfo.MidInterval
	midHr := "0" + strconv.Itoa(midHrInt)
	midMn := now.Minute()
	var midMnStr string
	midMnStr = strconv.Itoa(midMn)
	if midMn < 10 {
		midMnStr = "0" + midMnStr
	}

	if cfg.Schedules.MidBosses[midHr+midMnStr] != nil {
		// log.Println(cfg.Schedules.MidBosses[midHr+midMnStr].Name)
		notifier(cfg.Schedules.MidBosses[midHr+midMnStr].Name, cfg.Schedules.MidBosses[midHr+midMnStr].Location)
	}

	//easy
	ezHrInt := now.Hour() % cfg.Schedules.TimeInfo.EasyInterval
	ezHr := "0" + strconv.Itoa(ezHrInt)
	ezMn := now.Minute()
	var ezMnStr string
	ezMnStr = strconv.Itoa(ezMn)
	if ezMn < 10 {
		ezMnStr = "0" + ezMnStr
	}

	if cfg.Schedules.EasyBosses[ezHr+ezMnStr] != nil {
		// log.Println(cfg.Schedules.EasyBosses[ezHr+ezMnStr].Name)
		notifier(cfg.Schedules.EasyBosses[ezHr+ezMnStr].Name, cfg.Schedules.EasyBosses[ezHr+ezMnStr].Location)
	}

	jt.t.Reset(getNextTickDuration())
}

func notifier(bossName string, location string) {
	fullURL := cfg.Credentials.CredentialCfg.DiscordURL + cfg.Credentials.CredentialCfg.GuildID + "/" + cfg.Credentials.CredentialCfg.WebhookToken
	log.Println(fullURL)

	clientRfd := &http.Client{}
	formData := url.Values{}
	formData.Add("username", "World Boss Notifier - GW2")
	formData.Add("avatar_url", "https://dviw3bl0enbyw.cloudfront.net/uploads/forum_attachment/file/27353/GuildWars2tile.png")
	formData.Add("content", "Currently Active - **"+bossName+"**! Waypoint (paste in GW2 chat): **"+location+"**")
	reqRfd, _ := http.NewRequest("POST", fullURL, bytes.NewBufferString(formData.Encode()))
	reqRfd.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, errReq := clientRfd.Do(reqRfd)
	if errReq != nil {
		log.Println("[tmp] ", errReq)
	}

	defer resp.Body.Close()
}

func main() {
	jt := NewJobTicker()
	for {
		<-jt.t.C
		// fmt.Println(time.Now(), "- just ticked")
		jt.updateJobTicker()
	}
}
