# 一、部署方式
支持服务器部署及Docker容器部署两种方式
---
## 1.1、服务器部署
    *自行编译并运行；   
    *使用已编译好的Release版本；  
   ~~信得过的也可以直接用我自己部署的服务：tts.lsprain.tech:3000，免部署;~~
## 1.2、Docker容器部署
Docker镜像：[Docker Hub](https://hub.docker.com/r/lsprain/aliddns).  
---
# 二、RouterOS7.x 脚本代码
ROS路由脚本如下：  
```
#xxxx处替换为需要解析的域名AccessKeyID  
:local AccessKeyID "xxxx"  
#xxxx处替换为需要解析的域名AccessKeySecret  
:local AccessKeySecret "xxxx"   
#xxxx处替换为需要解析的域名，如baidu.com  
:local DomainName "lsprain.xxxx"  
#xxxx处替换为需要解析的子域名，如home.baidu.com只需要填home即可   
:local RR "home"   
#xxxx处替换为网络出口名称，如pppoe-out1  
:local pppoe "pppoe-out1"   

:local IpAddr [/ip address get [/ip address find interface=$pppoe] address]  
:set IpAddr [:pick $IpAddr 0 ([len $IpAddr] -3)]  
:log warning "当前公网IP地址：$IpAddr"  

:local aliddns "http://服务地址:8800/aliddns?AccessKeyID=$AccessKeyID&AccessKeySecret=$AccessKeySecret&RR=$RR&DomainName=$DomainName&IpAddr=$IpAddr"  

:local result [/tool fetch url=("$aliddns") mode=http http-method=get as-value output=user];  
#:log warning $result  

:if ($result->"status" = "finished") do={  

:if ($result->"data" = "loginerr") do={  
:log warning "阿里云登录失败！!";  
}  
:if ($result->"data" = "iperr") do={  
:log warning "修改解析地址信息失败!";  
}  
:if ($result->"data" = "ip") do={  
:log warning "修改解析地址信息成功!";  
}  
:if ($result->"data" = "domainerr") do={  
:log warning "添加新域名解析失败!";  
}  
:if ($result->"data" = "domain") do={  
:log warning "添加新域名解析成功!";  
}  
:if ($result->"data" = "same") do={  
:log warning "当前配置解析地址与公网IP相同，不需要修改!";  
}  
:if ($result->"data" = "ip") do={  
:log warning "更新IP信息成功!";  
:log warning "$IpAddr!";  
}  
:if ($result->"data" = "domain") do={  
:log warning "增加域名信息成功!";  
}  
}  
}  
```
# 三、其它方式
method：```post```   
url：```http://服务地址:8800/aliddns?AccessKeyID=XXXXXX&AccessKeySecret=XXXXXX&RR=XX&DomainName=XXX&IpAddr=XXX```

