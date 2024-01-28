# Emit

Minimal event emitter using channels.

- Type Safe: built on generics
- Batching: timeout and length limit

### Installation
```
go get github.com/jonashiltl/emit
```

### Usage
```
e := new(emit.Emitter) // or &emit.Emmiter{}

emit.On(e, func(t []MyEvent) {
    fmt.Println(t)
}, WithBatchSize(10), WithTimeout(10*time.Second))

emit.Emit(e, MyEvent{})
```

Inspired by [mint](https://github.com/btvoidx/mint)