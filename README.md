# Emit

Minimal event emitter using channels.

- Type Safe: built on generics
- Batching: timeout and size

### Installation
```
go get github.com/jonashiltl/emit
```

### Usage
```go
e := new(emit.Emitter) // or &emit.Emmiter{}

emit.On(e, func(batch []MyEvent) {
    fmt.Println(batch)
}, WithBatchSize(10), WithTimeout(10*time.Second))

emit.Emit(e, MyEvent{})
```

Inspired by [mint](https://github.com/btvoidx/mint)
