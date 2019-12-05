# authorization

authorization is a handler plugin of mesher, it authorized for the requested user identity.

## Configurations
**In mesher.yaml**

**servicecomb.authorization.type**
>  *(required, string)* type is set by the way of authentication

### When choose oauth2 authorization, it's may need to configure
**servicecomb.authorization.endpoint**
>  *(required, string)* it is needed when you use a centralized auth system like oauth2, OpenID

**servicecomb.authorization.client**
>  *(required, string)* it is needed when you use a centralized auth system like oauth2, OpenID

**servicecomb.authorization.scopes**
>  *(optional, string)* Scope specifies optional requested permissions

**servicecomb.authorization.redirectURL**
>  *(optional, string)* 
> RedirectURL is the URL to redirect users going through the OAuth2 flow, after the resource owner's URLs.

## Example
```yaml
authorization:
 type: oauth2          #set by the way of authentication
 endpoint:
   authURL:   "your auth url"
   tokenURL:  "your token url"
 client:
   clientID:      "your client ID"
   clientSecret:  "your client secret"

```
## Step：

# Auth configure file init
**You must init authorization config file which will manage connection and report msg to authorization**
- For example:
- [1] Set your configurations in mesher.yaml. 
- [2] Implementation of the authorization interface, implemented by default in oauth2.go.
- [3] Add the auth package into cmd/mesher.go.




