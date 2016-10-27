package beego

import (
	"encoding/json"
	"mime"
	"net/http"
	"path/filepath"

	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/session"
)

//
func registerMime() error {
	for k, v := range mimemaps {
		mime.AddExtensionType(k, v)
	}
	return nil
}

// register default error http handlers, 404,401,403,500 and 503.
func registerDefaultErrorHandler() error {
	m := map[string]func(http.ResponseWriter, *http.Request){
		"401": unauthorized,
		"402": paymentRequired,
		"403": forbidden,
		"404": notFound,
		"405": methodNotAllowed,
		"500": internalServerError,
		"501": notImplemented,
		"502": badGateway,
		"503": serviceUnavailable,
		"504": gatewayTimeout,
	}
	for e, h := range m {
		if _, ok := ErrorMaps[e]; !ok {
			ErrorHandler(e, h)
		}
	}
	return nil
}

func registerSession() error {
	if BConfig.WebConfig.Session.SessionOn {
		var err error
		sessionConfig := AppConfig.String("sessionConfig")
		conf := new(session.ManagerConfig)
		if sessionConfig == "" {
			conf.CookieName = BConfig.WebConfig.Session.SessionName
			conf.EnableSetCookie = BConfig.WebConfig.Session.SessionAutoSetCookie
			conf.Gclifetime = BConfig.WebConfig.Session.SessionGCMaxLifetime
			conf.Secure = BConfig.Listen.EnableHTTPS
			conf.CookieLifeTime = BConfig.WebConfig.Session.SessionCookieLifeTime
			conf.ProviderConfig = filepath.ToSlash(BConfig.WebConfig.Session.SessionProviderConfig)
			conf.Domain = BConfig.WebConfig.Session.SessionDomain
			conf.EnableSidInHttpHeader = BConfig.WebConfig.Session.EnableSidInHttpHeader
			conf.SessionNameInHttpHeader = BConfig.WebConfig.Session.SessionNameInHttpHeader
			conf.EnableSidInUrlQuery = BConfig.WebConfig.Session.EnableSidInUrlQuery
		} else {
			if err = json.Unmarshal([]byte(sessionConfig), conf); err != nil {
				return err
			}
		}
		if GlobalSessions, err = session.NewManager(BConfig.WebConfig.Session.SessionProvider, conf); err != nil {
			return err
		}
		///怎么让这个GC终止？
		///要是main退出了，这个还是继续？
		///经测试，在程序退出的时候，这个也就会退出了
		///TODO: 查查Go规范的说明
		go GlobalSessions.GC()
	}
	return nil
}

func registerTemplate() error {
	if err := BuildTemplate(BConfig.WebConfig.ViewsPath); err != nil {
		if BConfig.RunMode == DEV {
			logs.Warn(err)
		}
		return err
	}
	return nil
}

func registerAdmin() error {
	if BConfig.Listen.EnableAdmin {
		go beeAdminApp.Run()
	}
	return nil
}

func registerGzip() error {
	if BConfig.EnableGzip {
		context.InitGzip(
			AppConfig.DefaultInt("gzipMinLength", -1),
			AppConfig.DefaultInt("gzipCompressLevel", -1),
			AppConfig.DefaultStrings("includedMethods", []string{"GET"}),
		)
	}
	return nil
}
