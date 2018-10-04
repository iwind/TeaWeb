package tealog

import (
	"strings"
	"github.com/iwind/TeaGo/utils/string"
	"github.com/iwind/TeaGo/logs"
	"path/filepath"
	"sync"
	"hash/crc32"
	"github.com/ua-parser/uap-go/uaparser"
	"github.com/oschwald/geoip2-golang"
	"net"
	"regexp"
	"github.com/iwind/TeaGo/Tea"
	"reflect"
	"fmt"
	"github.com/iwind/TeaWebCode/teautils"
	"github.com/pquerna/ffjson/ffjson"
)

var userAgentParserCache = &sync.Map{}
var userAgentParser *uaparser.Parser
var geoDB *geoip2.Reader
var accessLogVars = map[string]string{}

// 参考：http://nginx.org/en/docs/http/ngx_http_log_module.html#log_format
type AccessLog struct {
	ServerId   string `var:"serverId" bson:"serverId" json:"serverId"`       // 服务ID
	BackendId  string `var:"backendId" bson:"backendId" json:"backendId"`    // 后端服务ID
	LocationId string `var:"locationId" bson:"locationId" json:"locationId"` // 路径配置ID
	FastcgiId  string `var:"fastcgiId" bson:"fastcgiId" json:"fastcgiId"`    // Fastcgi配置ID
	RewriteId  string `var:"rewriteId" bson:"rewriteId" json:"rewriteId"`    // 重写规则ID

	TeaVersion      string  `var:"teaVersion" bson:"teaVersion" json:"teaVersion"` // TeaWeb版本
	RemoteAddr      string  `var:"remoteAddr" bson:"remoteAddr" json:"remoteAddr"` // 终端地址，通常是：ip:port
	RemotePort      int     `var:"remotePort" bson:"remotePort" json:"remotePort"` // 终端端口
	RemoteUser      string  `var:"remoteUser" bson:"remoteUser" json:"remoteUser"` // 终端用户，基于BasicAuth认证
	RequestURI      string  `var:"requestURI" bson:"requestURI" json:"requestURI"`
	RequestPath     string  `var:"requestPath" bson:"requestPath" json:"requestPath"`
	RequestLength   int64   `var:"requestLength" bson:"requestLength" json:"requestLength"`       // 请求内容长度
	RequestTime     float64 `var:"requestTime" bson:"requestTime" json:"requestTime"`             // 从请求到所有响应数据发送到请求端所花时间，单位为带有小数点的秒，精确到纳秒，比如：0.000260081
	RequestMethod   string  `var:"requestMethod" bson:"requestMethod" json:"requestMethod"`       // 请求方法
	RequestFilename string  `var:"requestFilename" bson:"requestFilename" json:"requestFilename"` // 请求的文件名，包含完整的路径
	Scheme          string  `var:"scheme" bson:"scheme" json:"scheme"`                            // 请求协议，http或者https
	Proto           string  `var:"proto" bson:"proto" json:"proto"`                               // 请求协议，比如HTTP/1.0, HTTP/1.1

	BytesSent     int64               `var:"bytesSent" bson:"bytesSent" json:"bytesSent"`             // 响应的字节数，目前同 bodyBytesSent
	BodyBytesSent int64               `var:"bodyBytesSent" bson:"bodyBytesSent" json:"bodyBytesSent"` // 响应的字节数
	Status        int                 `var:"status" bson:"status" json:"status"`                      // 响应的状态码
	StatusMessage string              `var:"statusMessage" bson:"statusMessage" json:"statusMessage"` // 响应的信息
	SentHeader    map[string][]string `var:"sentHeader" bson:"sentHeader" json:"sentHeader"`          // 响应的头信息

	TimeISO8601    string              `var:"timeISO8601" bson:"timeISO8601" json:"timeISO8601"` // ISO 8601格式的本地时间，比如 2018-07-16T23:52:24.839+08:00
	TimeLocal      string              `var:"timeLocal" bson:"timeLocal" json:"timeLocal"`       // 本地时间，比如 17/Jul/2018:09:52:24 +0800
	Msec           float64             `var:"msec" bson:"msec" json:"msec"`                      // 带有毫秒的时间，比如 1531756823.054
	Timestamp      int64               `var:"timestamp" bson:"timestamp" json:"timestamp"`       // unix时间戳，单位为秒
	Host           string              `var:"host" bson:"host" json:"host"`
	Referer        string              `var:"referer" bson:"referer" json:"referer"`
	UserAgent      string              `var:"userAgent" bson:"userAgent" json:"userAgent"`
	Request        string              `var:"request" bson:"request" json:"request"`                      // 请求的简要说明，格式类似于 GET /hello/world HTTP/1.1
	ContentType    string              `var:"contentType" bson:"contentType" json:"contentType"`          // 请求头部的Content-Type
	Cookie         map[string]string   `bson:"cookie" json:"cookie"`                                      // Cookie cookie.name, cookie.sid
	Arg            map[string][]string `bson:"arg" json:"arg"`                                            // arg_name, arg_id
	Args           string              `var:"args" bson:"args" json:"args"`                               // name=liu&age=20
	QueryString    string              `var:"queryString" bson:"queryString" json:"queryString"`          // 同 Args
	Header         map[string][]string `bson:"header" json:"header"`                                      // 请求的头部信息，支持header_*和http_*，header_content_type, header_expires, http_content_type, http_user_agent
	ServerName     string              `var:"serverName" bson:"serverName" json:"serverName"`             // 接收请求的服务器名
	ServerPort     int                 `var:"serverPort" bson:"serverPort" json:"serverPort"`             // 服务器端口
	ServerProtocol string              `var:"serverProtocol" bson:"serverProtocol" json:"serverProtocol"` // 服务器协议，类似于HTTP/1.0”

	// 代理相关
	BackendAddress string `var:"backendAddress" bson:"backendAddress" json:"backendAddress"` // 代理的后端的地址
	FastcgiAddress string `var:"fastcgiAddress" bson:"fastcgiAddress" json:"fastcgiAddress"` // Fastcgi后端地址

	// 扩展
	Extend struct {
		File   AccessLogFile   `bson:"file" json:"file"`
		Client AccessLogClient `bson:"client" json:"client"`
		Geo    AccessLogGeo    `bson:"geo" json:"geo"`
	} `bson:"extend" json:"extend"`

	// 格式化的正则表达式
	formatReg *regexp.Regexp
	headerReg *regexp.Regexp
}

type AccessLogFile struct {
	MimeType  string `bson:"mimeType" json:"mimeType"`   // 类似于 image/jpeg
	Extension string `bson:"extension" json:"extension"` // 扩展名，不带点（.）
	Charset   string `bson:"charset" json:"charset"`     // 字符集，统一大写
}

type AccessLogClient struct {
	OS      AccessLogClientOS      `bson:"os" json:"os"`
	Device  AccessLogClientDevice  `bson:"device" json:"device"`
	Browser AccessLogClientBrowser `bson:"browser" json:"browser"`
}

type AccessLogClientOS struct {
	Family     string `bson:"family" json:"family"`
	Major      string `bson:"major" json:"major"`
	Minor      string `bson:"minor" json:"minor"`
	Patch      string `bson:"patch" json:"patch"`
	PatchMinor string `bson:"patchMinor" json:"patchMinor"`
}

type AccessLogClientDevice struct {
	Family string `bson:"family" json:"family"`
	Brand  string `bson:"brand" json:"brand"`
	Model  string `bson:"model" json:"model"`
}

type AccessLogClientBrowser struct {
	Family string `bson:"family" json:"family"`
	Major  string `bson:"major" json:"major"`
	Minor  string `bson:"minor" json:"minor"`
	Patch  string `bson:"patch" json:"patch"`
}

type AccessLogGeo struct {
	Region   string               `bson:"region" json:"region"`
	State    string               `bson:"state" json:"state"`
	City     string               `bson:"city" json:"city"`
	Location AccessLogGeoLocation `bson:"location" json:"location"`
}

type AccessLogGeoLocation struct {
	Latitude       float64 `bson:"latitude" json:"latitude"`
	Longitude      float64 `bson:"longitude" json:"longitude"`
	TimeZone       string  `bson:"timeZone" json:"timeZone"`
	AccuracyRadius uint16  `bson:"accuracyRadius" json:"accuracyRadius"`
	MetroCode      uint    `bson:"metroCode" json:"metroCode"`
}

func init() {
	var err error
	userAgentParser, err = uaparser.New(Tea.Root + Tea.DS + "resources" + Tea.DS + "regexes.yaml")
	if err != nil {
		logs.Error(err)
	}

	geoDB, err = geoip2.Open(Tea.Root + "/resources/GeoLite2-City/GeoLite2-City.mmdb")
	if err != nil {
		logs.Error(err)
	}

	// 变量
	reflectType := reflect.TypeOf(AccessLog{})
	countField := reflectType.NumField()
	for i := 0; i < countField; i ++ {
		field := reflectType.Field(i)
		value := field.Tag.Get("var")
		if len(value) > 0 {
			accessLogVars[value] = field.Name
		}
	}
}

func (this *AccessLog) SentContentType() string {
	contentType, ok := this.SentHeader["Content-Type"]
	if ok && len(contentType) > 0 {
		return contentType[0]
	}
	return ""
}

func (this *AccessLog) Format(format string) string {
	if this.formatReg == nil {
		this.formatReg = regexp.MustCompile("\\${[\\w.]+}")
	}

	if this.headerReg == nil {
		this.headerReg = regexp.MustCompile("([A-Z])")
	}

	refValue := reflect.ValueOf(*this)

	// 处理变量${varName}
	format = this.formatReg.ReplaceAllStringFunc(format, func(s string) string {
		varName := s[2 : len(s)-1]

		fieldName, found := accessLogVars[varName]
		if found {
			field := refValue.FieldByName(fieldName)
			if field.IsValid() {
				if field.Kind() == reflect.String {
					return field.String()
				} else {
					return fmt.Sprintf("%#v", field.Interface())
				}
			}

			return ""
		}

		// arg
		if strings.HasPrefix(varName, "arg.") {
			varName = varName[4:]
			values, found := this.Arg[varName]
			if found {
				countValues := len(values)
				if countValues == 1 {
					return values[0]
				} else if countValues > 1 {
					return "[" + strings.Join(values, ",") + "]"
				}
			}
			return ""
		}

		// cookie
		if strings.HasPrefix(varName, "cookie.") {
			varName = varName[7:]
			value, found := this.Cookie[varName]
			if found {
				return value
			}
			return ""
		}

		// http
		if strings.HasPrefix(varName, "http.") {
			varName = varName[5:]
			values, found := this.Header[varName]
			if found {
				if len(values) > 0 {
					return values[0]
				}
			} else {
				varName = strings.TrimPrefix(this.headerReg.ReplaceAllString(varName, "-${1}"), "-")
				values, found := this.Header[varName]
				if found && len(values) > 0 {
					return values[0]
				}
			}

			return ""
		}

		// header
		if strings.HasPrefix(varName, "header.") {
			varName = varName[7:]
			values, found := this.Header[varName]
			if found {
				if len(values) > 0 {
					return values[0]
				}
			} else {
				varName = strings.TrimPrefix(this.headerReg.ReplaceAllString(varName, "-${1}"), "-")
				values, found := this.Header[varName]
				if found && len(values) > 0 {
					return values[0]
				}
			}

			return ""
		}

		// extend
		if strings.HasPrefix(varName, "extend.") {
			value := teautils.Get(this.Extend, strings.Split(varName[7:], "."))
			jsonValue, err := ffjson.Marshal(value)
			if err != nil {
				logs.Error(err)
			} else {
				return string(jsonValue)
			}
		}

		return s
	})

	return format
}

func (this *AccessLog) Parse() {
	this.parseMime()
	this.parseExtension()
	this.parseUserAgent()
	this.parseGeoIP()
}

func (this *AccessLog) parseMime() {
	semicolonIndex := strings.Index(this.ContentType, ";")
	if semicolonIndex == -1 {
		this.Extend.File.MimeType = this.ContentType
		this.Extend.File.Charset = ""
		return
	}

	this.Extend.File.MimeType = this.ContentType[:semicolonIndex]
	reg, err := stringutil.RegexpCompile("(?i)charset\\s*=\\s*([\\w-]+)")
	if err != nil {
		logs.Error(err)
	} else {
		match := reg.FindStringSubmatch(this.ContentType)
		if len(match) > 0 {
			this.Extend.File.Charset = strings.ToUpper(match[1])
		} else {
			this.Extend.File.Charset = ""
		}
	}
}

func (this *AccessLog) parseExtension() {
	ext := filepath.Ext(this.RequestPath)
	if len(ext) == 0 {
		this.Extend.File.Extension = ""
	} else {
		this.Extend.File.Extension = strings.ToLower(ext[1:])
	}
}

func (this *AccessLog) parseUserAgent() {
	// MDN上的参考：https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/User-Agent
	// 浏览器集合测试：http://www.browserscope.org/
	// 多种变成语言识别：https://github.com/ua-parser/uap-php

	userAgent := this.UserAgent
	crc := crc32.ChecksumIEEE([]byte(userAgent))
	item, found := userAgentParserCache.Load(crc)
	if found {
		this.Extend.Client = item.(AccessLogClient)
		return
	}

	client := userAgentParser.Parse(this.UserAgent)
	if client != nil {
		this.Extend.Client = AccessLogClient{
			OS: AccessLogClientOS{
				Family:     client.Os.Family,
				Major:      client.Os.Major,
				Minor:      client.Os.Minor,
				Patch:      client.Os.Patch,
				PatchMinor: client.Os.PatchMinor,
			},
			Device: AccessLogClientDevice{
				Family: client.Device.Family,
				Brand:  client.Device.Brand,
				Model:  client.Device.Model,
			},
			Browser: AccessLogClientBrowser{
				Family: client.UserAgent.Family,
				Major:  client.UserAgent.Major,
				Minor:  client.UserAgent.Minor,
				Patch:  client.UserAgent.Patch,
			},
		}

		userAgentParserCache.Store(crc, this.Extend.Client)
	}
}

func (this *AccessLog) parseGeoIP() {
	if geoDB == nil {
		return
	}

	// 参考：https://dev.maxmind.com/geoip/geoip2/geolite2/
	ip := net.ParseIP(this.RemoteAddr)
	record, err := geoDB.City(ip)
	if err != nil {
		logs.Error(err)
		return
	}

	this.Extend.Geo.Location.AccuracyRadius = record.Location.AccuracyRadius
	this.Extend.Geo.Location.MetroCode = record.Location.MetroCode
	this.Extend.Geo.Location.TimeZone = record.Location.TimeZone
	this.Extend.Geo.Location.Latitude = record.Location.Latitude
	this.Extend.Geo.Location.Longitude = record.Location.Longitude

	if len(record.Country.Names) > 0 {
		name, found := record.Country.Names["zh-CN"]
		if found {
			this.Extend.Geo.Region = name
		}
	}

	if len(record.Subdivisions) > 0 && len(record.Subdivisions[0].Names) > 0 {
		name, found := record.Subdivisions[0].Names["zh-CN"]
		if found {
			this.Extend.Geo.State = name
		}
	}

	if len(record.City.Names) > 0 {
		name, found := record.City.Names["zh-CN"]
		if found {
			this.Extend.Geo.City = name
		}
	}
}
