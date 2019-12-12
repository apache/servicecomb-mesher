# OAuth2

Mesher provides a high-level general-purpose middleware abstraction layer. One of the abstractions is oauth2, which is free  user learning [complexity inside handler chain](https://docs.go-chassis.com/dev-guides/how-to-implement-handler.html), so that users only need to focus on the development of their own business.

## configuration

Writing business code

When you use authorization code model, you need to implement the follow parameters. Otherwise you need to implement the interface of config in /oauth2manage.

For example, implement the authorization code model **In config.go**

```go

	auth.Use(&auth.Config{
		UseConfig: func() *oauth2.Config {
			// set your config info
			// For example
			config := oauth2.Config{
				ClientID:     "",           // (required, string) your client_ID
				ClientSecret: "",           // (required, string) your client_Secret
				Scopes:       []string{""}, // (optional, string) scope specifies requested permissions
				RedirectURL:  "http://localhost:30101/",  // (required, string) URL to redirect users going through the OAuth2 flow
				Endpoint: oauth2.Endpoint{ // (required, string) your auth server endpoint
					AuthURL:  "",
					TokenURL: "",
				},
			}
			return &config
		},
	})
```

Registration grant type: The default is the authorization code model

**In oauth2_handler.go**
```go
gt, err := oauth2manage.NewType(grantType)
		if err != nil {
			openlogging.Error("grant_type error: " + err.Error())
			cb(oauth2manage.InvResponse(http.StatusUnauthorized, err))
			return
		}
```

Change the configuration file and add the oauth2 handler to the chain. Note that as authentication, generally speaking,it is a server function, it must be placed in the provider chain.

```yaml
handler:
    chain:
      Consumer:
        outgoing: 
      Provider:
        incoming: oauth2 #provider handlers
```

## How to use

**oauth2-handler Init**
- [1] Registration grant type manually to init oauth2 manager, default is the authorization code model.
- [2] Implement the config file of your grant type, the interface definition in config.go.
- [3] Adding oauth2's provider handler name oauth2 defined in /oauth2 to providerChain.
- [4] You must import proxy/handler/oauth2 to init oauth2 handler. All the handlers which are customized for mesher are defined in file cmd/mesher/mesher.go.
- more details about handler chains in [go-chassis](https://github.com/go-chassis/go-chassis#readme)



