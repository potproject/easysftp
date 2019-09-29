# potproject/easysftp
SFTP connection easy. [pkg/sftp](https://github.com/pkg/sftp) wrapper Golang Library

**Work in Progress**

## Usage

```sh
import "github.com/potproject/easysftp"
```

## Example
```go
package main
import "github.com/potproject/easysftp"

func main(){
  // SFTP File Download
  easysftp.Download("USERNAME", "example.hostname.local", 22, "~/.ssh/id_rsa", "/tmp/localfile.txt" "/tmp/remotefile.txt")
  
  // SFTP File Upload
  easysftp.Upload("USERNAME", "example.hostname.local", 22, "~/.ssh/id_rsa", "/tmp/localfile.txt" "/tmp/remotefile.txt")
}
```

## LICENSE
MIT
