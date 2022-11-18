package middlewares

import (
	"fmt"
	"github.com/gin-gonic/gin"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

func Logger() gin.HandlerFunc {
	log := logrus.New()

	// 设置输出文件
	filePath := "logs/"
	fileName := "Aliddns-API"
	// 打开指定处的文件，并指定权限为：可读可写，可创建
	file := filePath + fileName
	//log.SetOutput(io.MultiWriter(os.Stdout, "))
	log.Out = os.Stdout

	// 设置日志级别。低于 Debug 级别的 Trace 将不会被打印
	log.SetLevel(logrus.DebugLevel)

	// 设置日志切割 rotatelogs
	writer, _ := rotatelogs.New(
		file+"%Y%m%d.log",
		//日志最大保存时间
		rotatelogs.WithMaxAge(7*24*time.Hour),
		//设置日志切割时间间隔(1天)(隔多久分割一次)
		rotatelogs.WithRotationTime(24*time.Hour),
	)
	//log.SetOutput(writer)

	// lfshook 决定哪些日志级别可用日志分割
	writeMap := lfshook.WriterMap{
		logrus.PanicLevel: writer,
		logrus.FatalLevel: writer,
		logrus.ErrorLevel: writer,
		logrus.WarnLevel:  writer,
		logrus.InfoLevel:  writer,
		logrus.DebugLevel: writer,
	}

	// 配置 lfshook
	hook := lfshook.NewHook(writeMap, &logrus.JSONFormatter{
		// 设置日期格式
		TimestampFormat: "2006-01-02 15:04:05",
	})

	//为 logrus 实例添加自定义 hook
	log.AddHook(hook)
	return func(c *gin.Context) {
		// 一.配置所需的 Fields
		startTime := time.Now()
		c.Next()
		spendTime := time.Since(startTime).Milliseconds()
		ST := fmt.Sprintf("%d ms", spendTime) // 1.API 调用耗时
		//hostName, err := os.Hostname()        // 2.主机名
		//if err != nil {
		//	hostName = "unknown"
		//}
		statusCode := c.Writer.Status() // 3.状态码
		clientIP := c.ClientIP()        // 4.请求客户端的 IP
		//userAgent := c.Request.UserAgent() // 5.用户代理，通常是某个浏览器。dev环境下是apipost
		dataSize := c.Writer.Size() // 6.响应报文 body 的字节长度
		if dataSize < 0 {
			dataSize = 0
		}
		method := c.Request.Method   // 7.请求方法
		path := c.Request.RequestURI // 8.请求 URL

		// 二.从标准记录器创建一个条目，并向其中添加多个字段(隐式添加 log 本身的时间戳,信息等 fields )
		entry := log.WithFields(logrus.Fields{
			//"HostName":  hostName,
			"Status":    statusCode,
			"SpendTime": ST,
			"IP":        clientIP,
			//"UserAgent": userAgent,
			"Method": method,
			//"DataSize":  dataSize,
			"Path": path,
		})

		// Errors 保存了使用当前context的所有中间件/handler 所产生的全部错误信息。
		// 源码注释： Errors is a list of errors attached to all the handlers/middlewares who used this context.
		// 三.将系统内部的错误 log 出去
		if len(c.Errors) > 0 {
			log.Error(c.Errors.ByType(gin.ErrorTypePrivate).String())
		}

		// 四.根据状态码决定打印 log 的等级
		if statusCode >= 500 {
			entry.Error()
		} else if statusCode >= 400 {
			entry.Warn()
		} else {
			entry.Info()
		}
	}
}
