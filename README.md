# beego-pongo2 v4 version

##download and install
go get -u github.com/astaxie/beego

go get -u gopkg.in/flosch/pongo2.v4

go get -u github.com/flosch/beego-pongo2

#####Latest stable release: v4.0 (go get -u gopkg.in/flosch/pongo2.v4 / v4-branch)

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
