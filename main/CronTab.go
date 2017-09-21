package main

import (
	"os"
	"strings"
	"bufio"
	"fmt"
	"io"
	"strconv"
	"time"
	"path/filepath"
)

const middle = "========="
const (
	DAY string = "day"
	HOUR string = "hour"
	MINUTES string = "minutes"
	SECOND string = "second"
)

type Config struct {
	Mymap  map[string]string
	strcet string
}

func (c *Config) InitConfig(path string) {
	c.Mymap = make(map[string]string)

	f, err := os.OpenFile(path, os.O_RDONLY, 0666)
	//f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	r := bufio.NewReader(f)
	for {
		b, _, err := r.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}

		s := strings.TrimSpace(string(b))
		//fmt.Println(s)
		if strings.Index(s, "#") == 0 {
			continue
		}

		n1 := strings.Index(s, "[")
		n2 := strings.LastIndex(s, "]")
		if n1 > -1 && n2 > -1 && n2 > n1 + 1 {
			c.strcet = strings.TrimSpace(s[n1 + 1 : n2])
			continue
		}

		if len(c.strcet) == 0 {
			continue
		}
		index := strings.Index(s, "=")
		if index < 0 {
			continue
		}

		frist := strings.TrimSpace(s[:index])
		if len(frist) == 0 {
			continue
		}
		second := strings.TrimSpace(s[index + 1:])

		pos := strings.Index(second, "\t#")
		if pos > -1 {
			second = second[0:pos]
		}

		pos = strings.Index(second, " #")
		if pos > -1 {
			second = second[0:pos]
		}

		pos = strings.Index(second, "\t//")
		if pos > -1 {
			second = second[0:pos]
		}

		pos = strings.Index(second, " //")
		if pos > -1 {
			second = second[0:pos]
		}

		if len(second) == 0 {
			continue
		}

		key := c.strcet + middle + frist
		c.Mymap[key] = strings.TrimSpace(second)
	}
}

func (c Config) Read(node, key string) string {
	key = node + middle + key
	v, found := c.Mymap[key]
	if !found {
		return ""
	}
	return v
}

/*func job(c chan time.Ticker){

}*/

func exeCronJob() {
	myConfig := new(Config)
	myConfig.InitConfig("conf.yaml")

	var interval int
	var unit string
	for k := range myConfig.Mymap {
		if strings.HasPrefix(k, "trigger=========") {
			if strings.Contains(k, "interval") {
				interval, _ = strconv.Atoi(myConfig.Mymap[k])
			} else if strings.Contains(k, "unit") {
				unit = myConfig.Mymap[k]
			}
		}
	}

	fmt.Println("the cron job will triggered per " + strconv.Itoa(interval) + " "+unit)

	var ch *time.Ticker
	if unit == HOUR {
		ch = time.NewTicker(time.Hour * time.Duration(interval))
	} else if unit == DAY {
		ch = time.NewTicker(time.Hour * 24 * time.Duration(interval))
	} else if unit == MINUTES {
		ch = time.NewTicker(time.Minute * time.Duration(interval))
	} else if unit == SECOND {
		ch = time.NewTicker(time.Second * time.Duration(interval))
	}

	//debugLog.SetPrefix("[Info]")
	//debugLog.SetFlags(debugLog.Flags() | log.LstdFlags)

	for {
		<-ch.C

		for k := range myConfig.Mymap {
			if strings.HasPrefix(k, "fp=========") {
				filepath.Walk(myConfig.Mymap[k],
					func(path string, f os.FileInfo, err error) error {
						if f == nil {
							return err
						}
						if f.IsDir() {
							fmt.Println("dir:", path)
							return nil
						}

						for k := range myConfig.Mymap {
							if strings.HasPrefix(k, "logs=========") {
								if strings.Contains(path, myConfig.Mymap[k]) {
									fmt.Printf("[%s] %v", time.Now().Format("2006-01-02 15:04:05"), "===========delete:============" + path)
									fmt.Println()
									err := os.Remove(path)               //删除
									if err != nil {
										fmt.Printf("%s", err)
									} else {
										fmt.Println(path + " removed done!")
									}
								}
							}
						}
						return nil
					})
			}
		}

	}
}

func main() {
	exeCronJob()
}
