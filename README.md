# Paramleak

### Description:
`A tool for extract parameters from Websites (HTML/JSON)`


### Installation:
```
go install github.com/SharokhAtaie/paramleak@latest
```

### Flags:
```ruby
██████╗  █████╗ ██████╗  █████╗ ███╗   ███╗██╗     ███████╗ █████╗ ██╗  ██╗
██╔══██╗██╔══██╗██╔══██╗██╔══██╗████╗ ████║██║     ██╔════╝██╔══██╗██║ ██╔╝
██████╔╝███████║██████╔╝███████║██╔████╔██║██║     █████╗  ███████║█████╔╝ 
██╔═══╝ ██╔══██║██╔══██╗██╔══██║██║╚██╔╝██║██║     ██╔══╝  ██╔══██║██╔═██╗ 
██║     ██║  ██║██║  ██║██║  ██║██║ ╚═╝ ██║███████╗███████╗██║  ██║██║  ██╗
╚═╝     ╚═╝  ╚═╝╚═╝  ╚═╝╚═╝  ╚═╝╚═╝     ╚═╝╚══════╝╚══════╝╚═╝  ╚═╝╚═╝  ╚═╝
                  Created by Sharo_k_h :)

Flags:
        -url,           -u              Url for Get All parameters
        -list,          -l              List of Url for Get All parameters
        -method,        -X              Http Method for requests (GET/POST/PATCH/DELETE/PUT)
        -body,          -d              Body data for Post/Patch/Delete/Put Requests
        -header,        -H              Custom Header (You can set only 1 custom header)
        -delay,         -p              Time for delay example: 1000 Millisecond (1 second)
        -verbose,       -v              Verbose mode
        -silent,        -s              Silent mode

```

### Usage:
```ruby
➜ ✗ paramleak -u "https://example.com/params.html"

██████╗  █████╗ ██████╗  █████╗ ███╗   ███╗██╗     ███████╗ █████╗ ██╗  ██╗
██╔══██╗██╔══██╗██╔══██╗██╔══██╗████╗ ████║██║     ██╔════╝██╔══██╗██║ ██╔╝
██████╔╝███████║██████╔╝███████║██╔████╔██║██║     █████╗  ███████║█████╔╝ 
██╔═══╝ ██╔══██║██╔══██╗██╔══██║██║╚██╔╝██║██║     ██╔══╝  ██╔══██║██╔═██╗ 
██║     ██║  ██║██║  ██║██║  ██║██║ ╚═╝ ██║███████╗███████╗██║  ██║██║  ██╗
╚═╝     ╚═╝  ╚═╝╚═╝  ╚═╝╚═╝  ╚═╝╚═╝     ╚═╝╚══════╝╚══════╝╚═╝  ╚═╝╚═╝  ╚═╝
                Created by Sharo_k_h :)

something
test_var
user_id_i
name
rdt_to
obj_key1
val1
obj_key2
val2
test_obj
empty_var
param1
method
param2
```

#### Thanks to ProjectDiscovery for best libraries
