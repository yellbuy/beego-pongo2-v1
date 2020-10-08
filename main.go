// A small library that lets you use Pongo2 with Beego
//
// When Render is called, it will populate the render context with Beego's flash messages.
// You can also use {% urlfor "MyController.Action" ":key" "value" %} in your templates, and
// it'll work just like `urlfor` would with `html/template`. It takes one controller argument and
// zero or more key/value pairs to fill the URL.
//
package pongo2

import (
	"fmt"
	"net/url"
	"os"
	"strings"
	"sync"

	"path"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	p2 "github.com/flosch/pongo2"
)

const (
	// 默认模板路径
	templateDir = "templates"
	// 插件模板路径
	pluginDir = "plugins"
)

//var templateDir = beego.BConfig.WebConfig.ViewsPath

type Context map[string]interface{}

var fileModifyTime = make(map[string]int64)
var templates = map[string]*p2.Template{}
var mutex = &sync.RWMutex{}

var devMode bool

// Render takes a Beego context, template name and a Context (map[string]interface{}).
// The template is parsed and cached, and gets executed into beegoCtx's ResponseWriter.
//
// Templates are looked up in `templates/` instead of Beego's default `views/` so that
// Beego doesn't attempt to load and parse our templates with `html/template`.
func Render(beegoCtx *context.Context, basePath, tmpl string, ctx Context) error {
	var curModifyTime int64
	//获取文件修改时间
	if basePath != "" {
		tmpl = path.Join(pluginDir, basePath, templateDir, tmpl)
		curModifyTime = getFileModTime(tmpl)
	}
	if curModifyTime==0{
		tmpl = path.Join(templateDir, tmpl)
		curModifyTime = getFileModTime(tmpl)
	}

	if !devMode {
		if curModifyTime > 0 {
			// 文件有效，需比对文件修改时间
			mutex.RLock()
			modifyTime, ok := fileModifyTime[tmpl]
			if !ok || modifyTime != curModifyTime {
				fileModifyTime[tmpl] = curModifyTime
				//时间匹配不上，清除缓存，以便重新加载
				p2.DefaultSet.CleanCache(tmpl)
				//fmt.Println("缓存失效")
			}
			mutex.RUnlock()
		}
	}

	template, err := p2.FromCache(tmpl)
	if err != nil {
		panic(err)
	}

	var pCtx p2.Context
	if ctx == nil {
		pCtx = p2.Context{}
	} else {
		pCtx = p2.Context(ctx)
	}

	if xsrf, ok := beegoCtx.GetSecureCookie(beego.BConfig.WebConfig.XSRFKey, "_xsrf"); ok {
		pCtx["_xsrf"] = xsrf
	}

	// Only override "flash" if it hasn't already been set in Context
	if _, ok := ctx["flash"]; !ok {
		if ctx == nil {
			ctx = Context{}
		}
		ctx["flash"] = readFlash(beegoCtx)
	}

	return template.ExecuteWriter(pCtx, beegoCtx.ResponseWriter)
}
// func pathExists(path string) bool {
// 	_, err := os.Stat(path)
// 	if err == nil {
// 		return true
// 	}
// 	if os.IsNotExist(err) {
// 		return false
// 	}
// 	fmt.Println(err)
// 	return false
// }

// Same as Render() but returns a string
func RenderString(basePath, tmpl string, ctx Context) (string, error) {
	var curModifyTime int64
	//获取文件修改时间
	if basePath != "" {
		tmpl = path.Join(pluginDir, basePath, templateDir, tmpl)
		curModifyTime = getFileModTime(tmpl)
	}
	if curModifyTime==0{
		tmpl = path.Join(templateDir, tmpl)
		curModifyTime = getFileModTime(tmpl)
	}
	if !devMode {
		if curModifyTime > 0 {
			// 文件有效
			mutex.RLock()
			modifyTime, ok := fileModifyTime[tmpl]
			if !ok || modifyTime != curModifyTime {
				fileModifyTime[tmpl] = curModifyTime
				//时间匹配不上，清除缓存，以便重新加载
				p2.DefaultSet.CleanCache(tmpl)
			}
			mutex.RUnlock()
		}
	}

	template, err := p2.FromCache(tmpl)
	if err != nil {
		panic(err)
	}

	var pCtx p2.Context
	if ctx == nil {
		pCtx = p2.Context{}
	} else {
		pCtx = p2.Context(ctx)
	}

	// str, _ := template.Execute(pCtx)
	// return str
	return template.Execute(pCtx)
}

// readFlash is similar to beego.ReadFromRequest except that it takes a *context.Context instead
// of a *beego.Controller, and returns a map[string]string directly instead of a Beego.FlashData
// (which only has a Data field anyway).
func readFlash(ctx *context.Context) map[string]string {
	data := map[string]string{}
	if cookie, err := ctx.Request.Cookie(beego.BConfig.WebConfig.FlashName); err == nil {
		v, _ := url.QueryUnescape(cookie.Value)
		vals := strings.Split(v, "\x00")
		for _, v := range vals {
			if len(v) > 0 {
				kv := strings.Split(v, "\x23"+beego.BConfig.WebConfig.FlashSeparator+"\x23")
				if len(kv) == 2 {
					data[kv[0]] = kv[1]
				}
			}
		}
		// read one time then delete it
		ctx.SetCookie(beego.BConfig.WebConfig.FlashName, "", -1, "/")
	}
	return data
}

//func SetHtmlEncryptKey(key []byte) {
//	p2.DefaultSet.HtmlEncryptKey = key
//}
// 获取文件修改时间
func getFileModTime(path string) int64 {
	f, err := os.Open(path)
	if err != nil {
		fmt.Println("open file error")
		return 0
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		fmt.Println("stat fileinfo error")
		return 0
	}

	return fi.ModTime().UnixNano()
}

func init() {
	devMode = beego.AppConfig.String("runmode") == "dev"
	p2.DefaultSet.Debug = devMode
	beego.BConfig.WebConfig.AutoRender = false
}
