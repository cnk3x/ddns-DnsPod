package main

import (
	"fmt"
	"os"
	"time"

	"github.com/kardianos/service"
)

type program struct{}

var (
	stop    = false
	records = make([]*Record, 0)
)

func (p *program) Start(s service.Service) error {
	if service.Interactive() {
		fmt.Println("动态域名解析(DnsPod)")
		fmt.Println("用法:\n ddns command")
		fmt.Println(" command:")
		fmt.Println("\tinstall:安装服务")
		fmt.Println("\tuninstall:卸载服务")
		fmt.Println("\tstart:启动服务")
		fmt.Println("\tstop:停止服务")
		os.Exit(0)
	} else {
		Info("正在以服务方式运行...")
		go p.run()
	}
	return nil
}

func (p *program) Stop(s service.Service) error {
	Info("服务正在停止...")
	stop = true
	return nil
}

func (p *program) run() {
	Info("服务已启动")

	var err error

	for _, item := range config.Records {
		records, err = FixAdd(records, item.Domain, item.Host, item.Token)

		if err != nil {
			Error(err.Error())
		}
	}

	for {
		time.Sleep(time.Second * 5)
		ip, err := GetIp()
		if err != nil {
			Error(`获取Ip错误:` + err.Error())
			continue
		}
		if len(ip) == 0 {
			continue
		}
		for _, record := range records {
			if stop {
				InfoR("已停止", record)
				return
			}

			InfoR("ip:"+record.RecordValue+" -> "+ip, record)

			if record.RecordValue == ip {
				continue
			}

			msg, err := DoUpdate(record.DomainId, record.RecordId, record.RecordName, ip, record.LoginKey)
			if err != nil {
				ErrorR(err.Error(), record)
			} else {
				record.RecordValue = ip
				InfoR(msg, record)
			}
		}

		if stop {
			Info("服务已停止")
			return
		}
	}
}

func main() {
	svcConfig := &service.Config{
		Name:        "HaoDDns",
		DisplayName: "Hao DDns",
		Description: "动态DNS控制",
	}

	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		Error(err.Error())
	}
	if len(os.Args) > 1 {
		fmt.Println(os.Args[1] + " ...")
		err = service.Control(s, os.Args[1])
		if err != nil {
			fmt.Println(os.Args[1] + " fatel.")
			Error(err.Error())
		} else {
			fmt.Println(os.Args[1] + " complated.")
		}
		return
	}

	err = s.Run()
	if err != nil {
		Error(err.Error())
	}
}
