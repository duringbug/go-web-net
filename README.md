go.mod: Go模块管理文件，定义了项目的模块名及其依赖的版本。Go模块系统通过这个文件管理依赖包。

go.sum: 记录项目中使用的每一个依赖包的校验和及版本信息。Go编译器通过它确保依赖包的完整性和正确性。

```bash
go mod init github.com/duringbug/go-web-net
go get github.com/eclipse/paho.mqtt.golang
go mod tidy # 移除不需要的包
go clean -modcache
```

# cell 

```bash
./scripts/build_cell.sh 
./build/cell -conf ./configs/cells_config/cell_config01.json 
```

# organ system
```bash
./scripts/build_organsys.sh
./build/organsys -conf ./configs/organsys_config/organ_config01.json
```
