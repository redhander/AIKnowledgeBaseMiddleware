# 热加载功能设置完成

## 🎉 功能已成功添加

你的AI Knowledge Base Middleware项目现在已经支持热加载功能！每次修改代码后，服务器会自动重新编译并重启。

## 📁 新增文件

1. **`.air.toml`** - Air热加载工具配置文件
2. **`Makefile`** - 包含构建、运行、热加载等命令
3. **`scripts/dev.ps1`** - Windows PowerShell开发脚本
4. **`scripts/dev.sh`** - Linux/macOS Bash开发脚本
5. **`dev.bat`** - Windows批处理文件
6. **`test_hot_reload.md`** - 测试指南

## 🚀 使用方法

### 方法1：使用批处理文件（Windows推荐）
```cmd
dev.bat
```

### 方法2：使用PowerShell脚本
```powershell
.\scripts\dev.ps1
```

### 方法3：使用Bash脚本（Linux/macOS）
```bash
./scripts/dev.sh
```

### 方法4：直接使用air命令
```bash
air
```

## ⚙️ 配置说明

### Air配置 (`.air.toml`)
- **监听文件类型**: `.go`, `.yaml`, `.yml`, `.json`
- **排除目录**: `tmp`, `vendor`, `testdata`, `deployments/volumes`
- **构建延迟**: 1000ms（防止频繁重建）
- **自动清理**: 退出时清理临时文件

### 支持的开发环境
- ✅ Windows (PowerShell + CMD)
- ✅ Linux/macOS (Bash)
- ✅ 所有Go文件类型
- ✅ 配置文件修改

## 🔧 可用命令

### Make命令（如果安装了make）
```bash
make help      # 显示所有命令
make build     # 构建应用
make run       # 正常运行
make dev       # 热加载模式
make clean     # 清理构建文件
make test      # 运行测试
```

### 直接Go命令
```bash
go run cmd/server/main.go    # 正常运行
go build cmd/server/main.go  # 构建
```

## 🧪 测试热加载

1. 启动热加载服务器
2. 修改任意 `.go` 文件
3. 保存文件
4. 观察控制台输出，应该看到自动重新编译和重启

## 🐛 故障排除

### 如果热加载不工作：

1. **检查air安装**
   ```bash
   air -v
   ```

2. **手动安装air**
   ```bash
   go install github.com/air-verse/air@latest
   ```

3. **检查文件权限**
   - 确保有写入 `tmp/` 目录的权限

4. **检查配置文件**
   - 确认 `.air.toml` 存在且配置正确

## 📝 开发建议

1. **使用热加载进行开发** - 提高开发效率
2. **生产环境使用正常构建** - 确保性能
3. **定期清理tmp目录** - 释放磁盘空间
4. **查看air日志** - 调试构建问题

## 🎯 下一步

现在你可以：
1. 使用 `dev.bat` 启动热加载开发模式
2. 修改代码并观察自动重新加载
3. 享受高效的开发体验！

---

**注意**: 热加载功能仅用于开发环境，生产环境请使用正常的构建和部署流程。 