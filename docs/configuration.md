# Configuration

You can use environment variables to configure the application.
The supported parameters and their corresponding environment variables are listed below in the following format:

```go
// ENV_VARIABLE_NAME DEFAULT_VALUE
```

In order to access the configuration web interface, you need to authenticate. By default the `admin` user and `jetspotter` password are used. You can change the password by setting the `AUTH_PASSWORD` environment variable.

```go
{%
   include-markdown "snippets/config.snippet"
   comments=false
%}
```
