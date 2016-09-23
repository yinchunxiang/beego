package beego

import (
	"encoding/json"
	"mime"
	"net/http"
	"path/filepath"

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
		if sessionConfig == "" {
			conf := map[string]interface{}{
				"cookieName":      BConfig.WebConfig.Session.SessionName,
				"gclifetime":      BConfig.WebConfig.Session.SessionGCMaxLifetime,
				"providerConfig":  filepath.ToSlash(BConfig.WebConfig.Session.SessionProviderConfig),
				"secure":          BConfig.Listen.EnableHTTPS,
				"enableSetCookie": BConfig.WebConfig.Session.SessionAutoSetCookie,
				"domain":          BConfig.WebConfig.Session.SessionDomain,
				"cookieLifeTime":  BConfig.WebConfig.Session.SessionCookieLifeTime,
			}
			confBytes, err := json.Marshal(conf)
			if err != nil {
				return err
			}
			sessionConfig = string(confBytes)
		}
		if GlobalSessions, err = session.NewManager(BConfig.WebConfig.Session.SessionProvider, sessionConfig); err != nil {
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
			Warn(err)
		}
		return err
	}
	return nil
}

func registerDocs() error {
	if BConfig.WebConfig.EnableDocs {
		Get("/docs", serverDocs)
		Get("/docs/*", serverDocs)
	}
	return nil
}

func registerAdmin() error {
	if BConfig.Listen.EnableAdmin {
		go beeAdminApp.Run()
	}
	return nil
}
