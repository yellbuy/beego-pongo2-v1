# beego-pongo2 v3 version

##download install
go get -u github.com/astaxie/beego
go get -u gopkg.in/flosch/pongo2.v3
go get -u github.com/yansuan/beego-pongo2

###Latest stable release: v3.0 (go get -u gopkg.in/flosch/pongo2.v3 / v3-branch)

##code
```go
package controllers

import (
    "github.com/astaxie/beego"
    "github.com/yansuan/beego-pongo2"
)

type MainController struct {
    beego.Controller
}

func (this *MainController) Get() {
    pongo2.Render(this.Ctx, "page.html", pongo2.Context{
        "name": "value"},
    })
}
```
