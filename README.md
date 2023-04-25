# Simple go-get HTTP server

Implements an HTTP server that provides information about the location of
git repositories for go modules.

This allows you to publish Go modules using your own domain in the import path,
unrelated to the domain where the source code is hosted.

## Compiling
```
git clone https://github.com/melato/goget.git
cd goget/main
date > version
go build goget.go
```


## Configuration
The server is configured by a yaml file that has two sections:

### domains
The domains section maps a whole domain to a location for packages of this domain.
For example, the configuration file below maps melato.org/X to https://github.com/melato/X.git
```
domains:
	melato.org: "https://github.com/melato/{{.path}}.git"
```
The domains configuration is used if the module path is not found in the modules section, described below.

### modules
```
modules:
- path: melato.org/command
  repository: https://github.com/melato/command.git
- path: melato.org/goget
  repository: https://github.com/melato/goget.git
```

Then run the server:
```
goget -port 8080 -c {config.yaml} server
```

This provides http service on port 8080.
To use https (TLS) on port 443, put it behind a reverse http proxy.


# Example Generated html
```
<!DOCTYPE html>
<html>
<head>
<title>melato.org/command</title>
<meta name="go-import" content="melato.org/command git https://github.com/melato/command.git">
</head>
<body>
<div>module: melato.org/command</div>
<div>repository: https://github.com/melato/command.git</div>
</body>
</html>

```
You can specify a different HTML template, using the -template flag.


# How Go finds modules
	
Suppose you have the following Go program, main.go:
```
package main
import "gopkg.in/yaml.v2"

func main() {
	yaml.Marshal("")
}
```

You can compile and run this program like this:
```
	go mod init main
	go mod tidy
	go run main.go
```

`go mod tidy` uses the import path "gopkg.in/yaml.v2" to make an HTTP request to https://gopkg.in/yaml.v2?go-get=1

You can make this request yourself to see what it returns:
```
curl https://gopkg.in/yaml.v2?go-get=1
```

Or you can use a browser, but make sure to look at the source html of the resulting page.
The important information is in the html head section:

```
<meta name="go-import" content="gopkg.in/yaml.v2 git https://gopkg.in/yaml.v2">
```

This program allows you to do the same with packages that use your own domain.
	
## Notes
- This web server ignores the go-get=1 parameter.  It produces the same output with or without it.
- Package names are matched exactly.  Nothing special is done for subpackages or versions.
- The response body is for humans to read.  The machine-readable information is in the "head" section.

## Pros and Cons

There are pros and cons when you use your own domain and/or go-get server.
### Pros 
- Your modules are automatically grouped in the import section.
- You use your own domain/brand in your code, instead of using someone else's brand.
- You can move your code to a different code hosting service and it will work just the same.


### Cons
- If your server goes down or your domain disappears, your module will not be as easily accessible.
But anyone can configure their own go-get server to point to the code,
point your dead domain to their go-get server and use your code as before.
- Potential users of your module might not trust your domain.
- If your domain expires and someone else gets it, they can lie to the world about where the source code is.
That is a danger with much of open source software.

In all these cases, someone can bypass your go-get server and get your module directly from the code hosting service.
Someone can then use it from their disk by specifying its location in the "replace" section of their go.mod file.		
Then Go will not try to get information from your domain.


## Hosting Private modules
I don't know if you can use your own go-get server to make it easier to use your private modules.yaml
I would have liked to be able to return a repository location like: git:mymodule, but this doesn't work.
If the code is accessible via https, then it could work.  Be careful not to leak secret urls to the world.
A GOPROXY server may help.

## GOPROXY
The GOPROXY functionality used by go tools is related to go-get.

