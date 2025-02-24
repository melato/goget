# Simple go-get HTTP server

Implements an HTTP server that provides information about the location of
git repositories for go modules.

This allows you to publish Go modules using your own domain in the import path,
separate from the domain where the source code is hosted.

# Compiling
goget has no dependencies.
```
git clone https://github.com/melato/goget.git
cd goget/main
go build goget.go
```


# Configuration
The server is configured by a JSON file that has two sections:

## domains
The domains section maps a whole domain to a location for packages of this domain.

The configuration file below maps melato.org/M to https://github.com/melato/M.git
```
{
  "domains": {
     "melato.org": "https://github.com/melato/{{.path}}.git"
  }
}
```

## modules
For more detailed configuration, you can specify individual modules:

```
{
 "modules": [
  {
   "package": "melato.org/command",
   "repository": "https://github.com/melato/command.git"
  }
 ]
}
```

## Running

Run the server:
```
goget -port 8080 -c {config.json} server
```

This provides http service on port 8080.
To use https (TLS) on port 443, put it behind a reverse http proxy.

## Notes
- This web server ignores the go-get=1 parameter.  It produces the same output with or without it.
- Package names are matched exactly.  Nothing special is done for subpackages or versions.
- The response body is for humans to read.  The machine-readable information is in the "head" section.


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
	
# Pros 
By using your own domain and/or go-get server.

- You use your own domain/brand in your code, instead of using someone else's brand.
- You can move your code to a different code hosting service and it will work just the same.
- Your modules are automatically grouped in the import section, by your domain.


# Cons
On the other hand, hosting your own go-get server has the following disadvantages:

- Potential users of your module might not trust your domain.
- If your server goes down or your domain expires, your module will not be easily accessible.
- If your domain expires and someone else gets it, they can lie to the world about where the source code is.

These are risks with much of open source software.

# Mitigation
I don't know all the ways by which you can protect yourself from using Go from unknown locations.  The following mitigations below are problematic.

## Use code locally
You can bypass the go-get server and get the module directly from the code hosting service. 
You can then use it locally, by specifying its location in the "replace" section of their go.mod file.
Then Go will not try to get information from the internet.

This the simplest and safest mitigation, that can be used even with Go standard library modules. 

The major disadvantage is that you have to modify each go.mod file.

## Use your own go-get server
You could potentially configure your own go-get server to point to code, instead of relying on someone else's domain.  Unfortunately, this is not so easy because:
- You need a TLS certificate for the module's domain.
- You need to host this code somewhere accessible via https.
- Go may use cached information about Go modules, from Google's servers.

# Hosting Private modules
I don't know if you can use your own go-get server to make it easier to use your private modules.

I would have liked to be able to return a repository location like: git:mymodule, but this doesn't work.

A goproxy server might be the solution.  This is beyond the scope of goget.