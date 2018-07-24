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
	Id int64

	TeaVersion      string  `var:"teaVersion"`
	RemoteAddr      string  `var:"remoteAddr"`
	RemotePort      string  `var:"remotePort"`
	RemoteUser      string  `var:"remoteUser"`
	RequestURI      string  `var:"requestURI"`
	RequestPath     string  `var:"requestPath"`
	RequestLength   int64   `var:"requestLength"`
	RequestTime     float64 `var:"requestTime"` // 从请求到所有响应数据发送到请求端所花时间，单位为带有小数点的秒，精确到纳秒，比如：0.000260081
	RequestMethod   string  `var:"requestMethod"`
	RequestFilename string  `var:"requestFilename"` // @TODO
	Scheme          string  `var:"requestScheme"`
	Proto           string  `var:"proto"`
	BytesSent       int64   `var:"bytesSent"`     // 响应的字节数
	BodyBytesSent   int64   `var:"bodyBytesSent"` // 响应的字节数（目前同BytesSent）
	Status          int     `var:"status"`        // 响应的状态码
	StatusMessage   string  `var:"statusMessage"` // 响应的信息
	TimeISO8601     string  `var:"timeISO8601"`   // ISO 8601格式的本地时间，比如 2018-07-16T23:52:24.839+08:00
	TimeLocal       string  `var:"timeLocal"`     // 本地时间，比如 17/Jul/2018:09:52:24 +0800
	Msec            string  `var:"msec"`          // 带有毫秒的时间，比如 1531756823.054
	Host            string  `var:"host"`
	Referer         string  `var:"referer"`
	UserAgent       string  `var:"userAgent"`
	Request         string  `var:"request"`
	ContentType     string  `var:"contentType"`
	Cookie          map[string]string              // Cookie cookie.name, cookie.sid
	Arg             map[string][]string            // arg_name, arg_id
	Args            string  `var:"args"`           // name=liu&age=20
	QueryString     string  `var:"queryString"`    // 同 Args
	Header          map[string][]string            // 请求的头部信息，支持header_*和http_*，header_content_type, header_expires, http_content_type, http_user_agent
	ServerName      string  `var:"serverName"`     // @TODO
	ServerPort      string  `var:"serverPort"`     // @TODO
	ServerProtocol  string  `var:"serverProtocol"` // @TODO

	// 代理相关
	BackendAddress string // 代理的后端的地址

	// 扩展
	Extend struct {
		File   AccessLogFile
		Client AccessLogClient
		Geo    AccessLogGeo
	}

	// 格式化的正则表达式
	formatReg *regexp.Regexp
	headerReg *regexp.Regexp
}

type AccessLogFile struct {
	MimeType  string // 类似于 image/jpeg
	Extension string // 扩展名，不带点（.）
	Charset   string // 字符集，统一大写
}

type AccessLogClient struct {
	OS      AccessLogClientOS
	Device  AccessLogClientDevice
	Browser AccessLogClientBrowser
}

type AccessLogClientOS struct {
	Family     string
	Major      string
	Minor      string
	Patch      string
	PatchMinor string
}

type AccessLogClientDevice struct {
	Family string
	Brand  string
	Model  string
}

type AccessLogClientBrowser struct {
	Family string
	Major  string
	Minor  string
	Patch  string
}

type AccessLogGeo struct {
	Country  string
	City     string
	Location AccessLogGeoLocation
}

type AccessLogGeoLocation struct {
	Latitude       float64
	Longitude      float64
	TimeZone       string
	AccuracyRadius uint16
	MetroCode      uint
}

func init() {
	var err error
	userAgentParser, err = uaparser.New(Tea.Root + Tea.Ds + "resources" + Tea.Ds + "regexes.yaml")
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

func (log *AccessLog) Format(format string) string {
	if log.formatReg == nil {
		log.formatReg = regexp.MustCompile("\\${[\\w.]+}")
	}

	if log.headerReg == nil {
		log.headerReg = regexp.MustCompile("([A-Z])")
	}

	refValue := reflect.ValueOf(*log)

	// 处理变量${varName}
	format = log.formatReg.ReplaceAllStringFunc(format, func(s string) string {
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
			values, found := log.Arg[varName]
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
			value, found := log.Cookie[varName]
			if found {
				return value
			}
			return ""
		}

		// http
		if strings.HasPrefix(varName, "http.") {
			varName = varName[5:]
			values, found := log.Header[varName]
			if found {
				if len(values) > 0 {
					return values[0]
				}
			} else {
				varName = strings.TrimPrefix(log.headerReg.ReplaceAllString(varName, "-${1}"), "-")
				values, found := log.Header[varName]
				if found && len(values) > 0 {
					return values[0]
				}
			}

			return ""
		}

		// header
		if strings.HasPrefix(varName, "header.") {
			varName = varName[7:]
			values, found := log.Header[varName]
			if found {
				if len(values) > 0 {
					return values[0]
				}
			} else {
				varName = strings.TrimPrefix(log.headerReg.ReplaceAllString(varName, "-${1}"), "-")
				values, found := log.Header[varName]
				if found && len(values) > 0 {
					return values[0]
				}
			}

			return ""
		}

		// extend
		if strings.HasPrefix(varName, "extend.") {
			value := teautils.Get(log.Extend, strings.Split(varName[7:], "."))
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

func (log *AccessLog) parse() {
	log.parseMime()
	log.parseExtension()
	log.parseUserAgent()
	log.parseGeoIP()
}

func (log *AccessLog) parseMime() {
	semicolonIndex := strings.Index(log.ContentType, ";")
	if semicolonIndex == -1 {
		log.Extend.File.MimeType = log.ContentType
		log.Extend.File.Charset = ""
		return
	}

	log.Extend.File.MimeType = log.ContentType[:semicolonIndex]
	reg, err := stringutil.RegexpCompile("(?i)charset\\s*=\\s*([\\w-]+)")
	if err != nil {
		logs.Error(err)
	} else {
		match := reg.FindStringSubmatch(log.ContentType)
		if len(match) > 0 {
			log.Extend.File.Charset = strings.ToUpper(match[1])
		} else {
			log.Extend.File.Charset = ""
		}
	}
}

func (log *AccessLog) parseExtension() {
	ext := filepath.Ext(log.RequestPath)
	if len(ext) == 0 {
		log.Extend.File.Extension = ""
	} else {
		log.Extend.File.Extension = strings.ToLower(ext[1:])
	}
}

func (log *AccessLog) parseUserAgent() {
	// MDN上的参考：https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/User-Agent
	// 浏览器集合测试：http://www.browserscope.org/
	// 多种变成语言识别：https://github.com/ua-parser/uap-php

	userAgent := log.UserAgent
	crc := crc32.ChecksumIEEE([]byte(userAgent))
	item, found := userAgentParserCache.Load(crc)
	if found {
		log.Extend.Client = item.(AccessLogClient)
		return
	}

	client := userAgentParser.Parse(log.UserAgent)
	if client != nil {
		log.Extend.Client = AccessLogClient{
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

		userAgentParserCache.Store(crc, log.Extend.Client)
	}
}

func (log *AccessLog) parseGeoIP() {
	log.RemoteAddr = "111.197.170.6" //@TODO

	// 参考：https://dev.maxmind.com/geoip/geoip2/geolite2/
	ip := net.ParseIP(log.RemoteAddr)
	record, err := geoDB.City(ip)
	if err != nil {
		logs.Error(err)
		return
	}

	log.Extend.Geo.Location.AccuracyRadius = record.Location.AccuracyRadius
	log.Extend.Geo.Location.MetroCode = record.Location.MetroCode
	log.Extend.Geo.Location.TimeZone = record.Location.TimeZone
	log.Extend.Geo.Location.Latitude = record.Location.Latitude
	log.Extend.Geo.Location.Longitude = record.Location.Longitude

	if len(record.Country.Names) > 0 {
		name, found := record.Country.Names["zh-CN"]
		if found {
			log.Extend.Geo.Country = name
		}
	}

	if len(record.City.Names) > 0 {
		name, found := record.City.Names["zh-CN"]
		if found {
			log.Extend.Geo.City = name
		}
	}
}
