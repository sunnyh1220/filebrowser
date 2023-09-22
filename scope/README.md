# Scope support
- admin用户不能使用scope功能;
- 其他用户必须使用scope功能;
- scope 类型可通过请求头参数 Arena-Scope 指定,当取值为空时，默认为Namespaced;
- 支持的scope值： Namespaced； 示例：`Arena-Scope: Namespaced`。


## Namespaced
当使用namespaced时，必须在设置query参数：namespace=<namespace>&adtype=<adtype>&adname=<adname>&adversion=<adversion>;
示例：`namespace=janedoe&adtype=dataset&adname=animal&adversion=default`

-  namespace：命名空间
-  adtype：数据类型
-  adname：数据名称
-  adversion：数据版本


### test
reverse proxy by caddy

ref:
- https://caddyserver.com/docs/caddyfile/directives/reverse_proxy#headers
- https://caddyserver.com/docs/caddyfile/directives/rewrite

#### caddy cli
```powershell
# .\filebrowser.exe config set --auth.method=noauth

go build -o filebrowser.exe .\main.go
.\filebrowser.exe config set --auth.method=proxy --auth.header=X-Arena-User
.\filebrowser.exe -a 0.0.0.0


# ref: https://caddyserver.com/docs/quick-starts/reverse-proxy
./caddy.exe reverse-proxy -H "X-Arena-User: jane" -H "Arena-Scope: Namespaced" --from :9090 --to :8080
./caddy.exe reverse-proxy -H "X-Arena-User: admin" --from :9999 --to :8080

http://localhost:9090/api/resources/?namespace=janedoe&adtype=dataset&adname=animal&adversion=default
```

#### Caddyfile
Caddyfile
```
{
	admin 127.0.0.1:8088
}

:9090 {
	reverse_proxy {
		to :8080
		header_up X-Arena-User "jane"
		header_up Arena-Scope "Namespaced"
	}

	rewrite /api/* ?{query}&namespace=janedoe&adtype=dataset&adname=animal&adversion=default
}

```

```powershell
.\caddy.exe run
```