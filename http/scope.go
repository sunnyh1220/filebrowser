/**
 * @Create on: 2023/9/12
 * @Author: sunnyh
 * @File:  scope
 * @Desc:
 */

package http

import (
	"fmt"
	"github.com/spf13/afero"
	"net/http"
	"path/filepath"
)

func scopeInner(fn handleFunc, w http.ResponseWriter, r *http.Request, d *data) (int, error) {
	if d.user.Username != "admin" {
		// 当请求头包含 Arena-Scope 参数时，该参数值将被视为 scope 值，并用于作用域及权限的校验。
		arenaScope := r.Header.Get("Arena-Scope")
		if arenaScope == "Namespaced" || arenaScope == "" {
			// query参数必须包含 namespace，adtype， adname， adversion 四个参数
			namespace := r.URL.Query().Get("namespace")
			adtype := r.URL.Query().Get("adtype")
			adname := r.URL.Query().Get("adname")
			adversion := r.URL.Query().Get("adversion")
			if namespace == "" || adtype == "" || adname == "" || adversion == "" {
				return http.StatusBadRequest, nil
			}

			// todo 鉴权

			// 鉴权通过后，修改user.Fs
			d.user.Fs = afero.NewBasePathFs(afero.NewOsFs(), filepath.Join(namespace, adtype, adname, adversion))

		}
	}
	return fn(w, r, d)
}

/*
example:

var resourceGetHandler = withScope(func(w http.ResponseWriter, r *http.Request, d *data) (int, error) {
*/
func withScope(fn handleFunc) handleFunc {
	return withUser(func(w http.ResponseWriter, r *http.Request, d *data) (int, error) {
		if d.user.Username != "admin" {
			// 当请求头包含 Arena-Scope 参数时，该参数值将被视为 scope 值，并用于作用域及权限的校验。
			arenaScope := r.Header.Get("Arena-Scope")
			if arenaScope == "Namespaced" || arenaScope == "" {
				// query参数必须包含 namespace，adtype， adname， adversion 四个参数
				namespace := r.URL.Query().Get("namespace")
				adtype := r.URL.Query().Get("adtype")
				adname := r.URL.Query().Get("adname")
				adversion := r.URL.Query().Get("adversion")
				if namespace == "" || adtype == "" || adname == "" || adversion == "" {
					return http.StatusBadRequest, nil
				}

				//r.WithContext(context.WithValue(r.Context(), "scopeMeta", scope.Metadata{
				//	Namespace:        namespace,
				//	ArenaDataType:    adtype,
				//	ArenaDataName:    adname,
				//	ArenaDataVersion: adversion,
				//}))

				// todo 鉴权

				// 鉴权通过后，修改user.Fs
				d.user.Fs = afero.NewBasePathFs(afero.NewOsFs(), filepath.Join(namespace, adtype, adname, adversion))

			}
		}
		return fn(w, r, d)
	})
}

/*
call before withUser
example:

	func withUser(fn handleFunc) handleFunc {
		return beforeWithUser(func(w http.ResponseWriter, r *http.Request, d *data) (int, error) {
			...
			return fn(w, r, d)
		})
	}

or:

beforeWithUser(resourceGetHandler)
*/
func beforeWithUser(fn handleFunc) handleFunc {
	return func(w http.ResponseWriter, r *http.Request, d *data) (int, error) {
		fmt.Println("beforeWithUser start")

		return fn(w, r, d)
	}
}
