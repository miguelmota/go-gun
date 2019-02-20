# go-gun

> [GunDB](https://github.com/amark/gun) implementation in Go (golang)

NOTE: is this is an early work-in-progress, highly unstable.

## Getting started

```go
package main

import (
	"log"

	"github.com/miguelmota/go-gun/server"
)

func main() {
	srv := server.NewServer(&server.Config{
		Port: 8080,
	})

	log.Fatal(srv.Start())
}
```

```bash
go run main.go
```

Frontend

```html
<script src="https://cdn.jsdelivr.net/npm/gun/gun.js"></script>
<script>
const gun = Gun('ws://localhost:8080')

gun.get('person').put({
  name: 'Alice',
  email: 'alice@example.com'
})

gun.get('person').on((data, key) => {
  console.log('update:', data)
})
</script>
```

## Test

```bash
make test
```

## Resources

- [Porting GUN documentation](https://gun.eco/docs/Porting-GUN)

## License

[MIT](LICENSE)
