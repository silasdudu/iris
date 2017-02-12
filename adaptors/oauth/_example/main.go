package main

import (
	"sort"
	"strings"

	"github.com/kataras/iris"
	"github.com/kataras/iris/adaptors/gorillamux"
	"github.com/kataras/iris/adaptors/oauth"
)

// IMPORTANT: Some website providers aren't allow localhost or 127.0.0.1 .
// Please emake sure that you try this with a
// domain registered to your provider's application.
// You may need to change the .Listen(":8080") to: app.Listen("mydomain.com")

// all configuration is ready, you have to set only:
//  RequestPath,CallbackRelativePath and the ___Key and __Secret
var configs = oauth.Config{
	RequestPath:          "/oauth/{provider}",
	RequestPathParam:     "provider",
	CallbackRelativePath: "/callback",
	RouteName:            "oauth", // for {{ url "oauth"}} in tempaltes

	GithubKey:    "YOUR_GITHUB_KEY",
	GithubSecret: "YOUR_GITHUB_SECRET",
	GithubName:   "github", // defaults to github

	FacebookKey:    "YOUR_FACEBOOK_KEY",
	FacebookSecret: "YOUR_FACEBOOK_KEY",
	FacebookName:   "facebook", // defaults to facebook,

	GplusKey:    "YOUR_GPLUS_KEY",
	GplusSecret: "YOUR_GPLUS_SECRET",
	GplusName:   "gplus",
}

// ProviderIndex ...
type ProviderIndex struct {
	Providers    []string
	ProvidersMap map[string]string
}

func main() {
	app := iris.New()
	app.Adapt(iris.DevLogger())
	app.Adapt(gorillamux.New()) // adapt a router, order doesn't matters but before Listen.
	// create the adaptor with our configs
	authentication := oauth.New(configs)
	// register the oauth/oauth2 adaptor
	app.Adapt(authentication)

	// set a  login success handler( you can use more than one handler)
	// if user succeed to logged in
	// client comes here from: localhost:8080/config.RouteName/lowercase_provider_name/callback 's first handler, but the  previous url is the localhost:8080/config.RouteName/lowercase_provider_name
	authentication.Success(func(ctx *iris.Context) {
		// if user couldn't validate then server sends StatusUnauthorized, which you can handle by:  authentication.Fail OR iris.OnError(iris.StatusUnauthorized, func(ctx *iris.Context){})
		user := authentication.User(ctx)

		// you can get the url by the named-route 'oauth' which you can change by Config's field: RouteName
		println("came from " + authentication.URL(strings.ToLower(user.Provider)))
		ctx.MustRender("user.html", user)
	})

	//
	// customize the oauth error page using:
	// authentication.Fail(func(ctx *iris.Context){....})
	//

	//  Note: on gorilla mux the {{ url }} and {{ path}} should give the key and the value, not only the values by order.
	//  {{ url "nameOfTheRoute" "parameterName" "parameterValue"}}.
	//
	// so: {{ url "providerLink" "facebook"}} should become
	// {{ url "providerLink" "provider" "facebook"}}
	//  for a path: "/auth/{provider}" with name 'providerLink'
	//
	// for the httprouter you do it like {{ url "nameOfTheRoute" "parameterValue" }}
	//
	// so here we're making a helper func (because we're using gorilla mux at this example)
	// which will fill the path parameter name and path parameter value
	app.Adapt(iris.TemplateFuncsPolicy{
		"providerURL": func(providerName string) string {
			// for route reversion

			// so here we prepend the configs.RequestPathParam before the providerName.
			return app.URL(configs.RouteName, configs.RequestPathParam, providerName)
		},
	})

	app.Get("/", func(ctx *iris.Context) {
		// show some providers to the template...
		m := make(map[string]string)
		m[configs.GithubName] = "Github"
		m[configs.FacebookName] = "Facebook"
		m[configs.GplusName] = "Gplus"

		var keys []string
		for k := range m {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		providerIndex := &ProviderIndex{Providers: keys, ProvidersMap: m}

		ctx.MustRender("index.html", providerIndex)
	})

	app.Listen("localhost:8080")
}
