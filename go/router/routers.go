/*
 * forum
 *
 * Тестовое задание для реализации проекта \"Форумы\" на курсе по базам данных в Технопарке Mail.ru (https://park.mail.ru).
 *
 * API version: 0.1.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package router

import (
	"fmt"
	"github.com/gorilla/mux"
	"go-server-server-generated/go"
	fdelivery "go-server-server-generated/src/forum/delivery"
	pdelivery "go-server-server-generated/src/post/delivery"
	"go-server-server-generated/src/service/delivery"
	tdelivery "go-server-server-generated/src/thread/delivery"
	ud "go-server-server-generated/src/user/delivery"
	"go-server-server-generated/src/utills"
	"net/http"
	"strings"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func NewRouter() *mux.Router {
	utills.CreateConnection()
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = swagger.Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/api/",
		Index,
	},

	Route{
		"Clear",
		strings.ToUpper("Post"),
		"/api/service/clear",
		delivery.ClearAll,
	},

	Route{
		"ForumCreate",
		strings.ToUpper("Post"),
		"/api/forum/create",
		fdelivery.ForumCreate,
	},

	Route{
		"ForumGetOne",
		strings.ToUpper("Get"),
		"/api/forum/{slug}/details",
		fdelivery.ForumGetOne,
	},

	Route{
		"ForumGetThreads",
		strings.ToUpper("Get"),
		"/api/forum/{slug}/threads",
		fdelivery.ForumGetThreads,
	},

	Route{
		"ForumGetUsers",
		strings.ToUpper("Get"),
		"/api/forum/{slug}/users",
		swagger.ForumGetUsers,
	},

	Route{
		"PostGetOne",
		strings.ToUpper("Get"),
		"/api/post/{id}/details",
		pdelivery.GetSinglePost,
	},

	Route{
		"PostUpdate",
		strings.ToUpper("Post"),
		"/api/post/{id}/details",
		pdelivery.UpdatePost,
	},

	Route{
		"PostsCreate",
		strings.ToUpper("Post"),
		"/api/thread/{slug_or_id}/create",
		pdelivery.PostsCreate,
	},

	Route{
		"Status",
		strings.ToUpper("Get"),
		"/api/service/status",
		delivery.Count,
	},

	Route{
		"ThreadCreate",
		strings.ToUpper("Post"),
		"/api/forum/{slug}/create",
		tdelivery.ThreadCreate,
	},

	Route{
		"ThreadGetOne",
		strings.ToUpper("Get"),
		"/api/thread/{slug_or_id}/details",
		tdelivery.GetThread,
	},

	Route{
		"ThreadGetPosts",
		strings.ToUpper("Get"),
		"/api/thread/{slug_or_id}/posts",
		pdelivery.GetPosts,
	},

	Route{
		"ThreadUpdate",
		strings.ToUpper("Post"),
		"/api/thread/{slug_or_id}/details",
		tdelivery.ThreadUpdate,
	},

	Route{
		"ThreadVote",
		strings.ToUpper("Post"),
		"/api/thread/{slug_or_id}/vote",
		tdelivery.ThreadVote,
	},

	Route{
		"UserCreate",
		strings.ToUpper("Post"),
		"/api/user/{nickname}/create",
		ud.UserCreate,
	},

	Route{
		"UserGetOne",
		strings.ToUpper("Get"),
		"/api/user/{nickname}/profile",
		ud.UserGetOne,
	},

	Route{
		"UserUpdate",
		strings.ToUpper("Post"),
		"/api/user/{nickname}/profile",
		ud.UserUpdate,
	},
}
