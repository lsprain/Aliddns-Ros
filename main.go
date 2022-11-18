package main

import (
	middlewares "Aliddns-Ros/log-handler"
	"github.com/denverdino/aliyungo/dns"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func init() {

}

// ConfigInfo 定义域名相关配置信息
type ConfigInfo struct {
	AccessKeyID     string
	AccessKeySecret string
	DomainName      string
	RR              string
	IpAddr          string
}

func main() {
	// 同时将日志写入文件和控制台
	//f, _ := os.Create("Aliddns_Get.log")
	//gin.DefaultWriter = io.MultiWriter(f, os.Stdout)

	//gin.SetMode(gin.ReleaseMode)
	//gin.DisableConsoleColor()

	r := gin.Default()
	r.Use(middlewares.Logger())
	r.GET("/aliddns", AddUpdateAliddns)
	r.Run(":8800")
}

func AddUpdateAliddns(c *gin.Context) {

	// 读取获取配置信息
	conf := new(ConfigInfo)
	conf.AccessKeyID = c.Query("AccessKeyID")
	conf.AccessKeySecret = c.Query("AccessKeySecret")
	conf.DomainName = c.Query("DomainName")
	conf.RR = c.Query("RR")
	conf.IpAddr = c.Query("IpAddr")

	//Info.Print("当前路由公网IP：" + conf.IpAddr)
	//log.SetOutput()
	log.Println("当前路由公网IP：" + conf.IpAddr)
	log.Println("进行阿里云登录……")

	// 连接阿里云服务器，获取DNS信息
	client := dns.NewClient(conf.AccessKeyID, conf.AccessKeySecret)
	client.SetDebug(false)
	domainInfo := new(dns.DescribeDomainRecordsArgs)
	domainInfo.DomainName = conf.DomainName
	oldRecord, err := client.DescribeDomainRecords(domainInfo)
	if err != nil {
		log.Println("阿里云登录失败！请查看错误日志！", err)
		c.String(http.StatusOK, "loginerr")
		return
	}
	log.Println("阿里云登录成功！")
	log.Println("进行域名及IP比对……")

	var exsitRecordID string
	for _, record := range oldRecord.DomainRecords.Record {
		if record.DomainName == conf.DomainName && record.RR == conf.RR {
			if record.Value == conf.IpAddr {
				log.Println("当前配置解析地址与公网IP相同，不需要修改。")
				c.String(http.StatusOK, "same")
				return
			}
			exsitRecordID = record.RecordId
		}
	}

	if 0 < len(exsitRecordID) {
		// 有配置记录，则匹配配置文件，进行更新操作
		updateRecord := new(dns.UpdateDomainRecordArgs)
		updateRecord.RecordId = exsitRecordID
		updateRecord.RR = conf.RR
		updateRecord.Value = conf.IpAddr
		updateRecord.Type = dns.ARecord
		rsp := new(dns.UpdateDomainRecordResponse)
		rsp, err := client.UpdateDomainRecord(updateRecord)
		if nil != err {
			log.Println("修改解析地址信息失败!", err)
			c.String(http.StatusOK, "iperr")
		} else {
			log.Println("修改解析地址信息成功!", rsp)
			c.String(http.StatusOK, "ip")
		}
	} else {
		// 没有找到配置记录，那么就新增一个
		newRecord := new(dns.AddDomainRecordArgs)
		newRecord.DomainName = conf.DomainName
		newRecord.RR = conf.RR
		newRecord.Value = conf.IpAddr
		newRecord.Type = dns.ARecord
		rsp := new(dns.AddDomainRecordResponse)
		rsp, err = client.AddDomainRecord(newRecord)
		if nil != err {
			log.Println("添加新域名解析失败！", err)
			c.String(http.StatusOK, "domainerr")
		} else {
			log.Println("添加新域名解析成功！", rsp)
			c.String(http.StatusOK, "domain")
		}
	}
}
