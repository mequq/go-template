
# Wire App together

Google wire is a code-generation tool for dependency injection.

[https://github.com/google/wire](https://github.com/google/wire)

I used this tool to wire different layers of application together.

on each folder of internal you can find a Provider tath defines witch function should use to generate wire

for example in the internal/data folder:

```go

var  DataProviderSet = wire.NewSet(

NewDataSource,

NewTenantRepo,

)
```

this line of code make wire to read function **NewDataSource** , **NewTenantRepo** part of wire structure, code will be generated if any other part of wire each function have some initiator ie:

```go
func  NewDataSource(logger *zap.Logger, cfg *config.ViperConfig) (*DataSource, error)
```

witch need logger , cfg and return *DataSource in return,

the wire will generate the necessary code to inject whom ever function is needed \*Datasource for functionality ie in the second function **NewTenantRepo** function initiator needed \*DataSource for initiation.

```go
func  NewTenantRepo(cfg *config.ViperConfig, logger *zap.Logger, data *DataSource)
```

the wire after generation will automatically put *DataSource to this function
code available at cmd/wire_gen.go this code is built automatically via cmd:

```bash
make generate
```

or manualy via:

```bash
go generate ./...
```

for details see make file, also you should make wire.go for generating code which makes wire input and output.

```go
func  wireApp(cfg *config.ViperConfig, logger *zap.Logger) (*app.App, error) {

panic(wire.Build(

app.NewApp,

server.ServerProviderSet,

service.ServiceProviderSet,

biz.BizProviderSet,

data.DataProviderSet,

))

}
```

This will make life easier during changes and the life cycle of the application.
