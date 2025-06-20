# 热加载功能测试指南

## 测试步骤

1. **启动热加载服务器**
   ```bash
   # Windows PowerShell
   .\scripts\dev.ps1
   
   # 或者直接运行
   air
   ```

2. **验证服务器启动**
   - 检查控制台输出，应该看到类似以下信息：
     ```
     Watching for changes in /path/to/project
     Building...
     Running...
     ```

3. **测试代码修改**
   - 修改任意 `.go` 文件（比如在 `cmd/server/main.go` 中添加一行日志）
   - 保存文件
   - 观察控制台输出，应该看到自动重新编译和重启

4. **测试配置修改**
   - 修改 `.yaml` 或 `.json` 配置文件
   - 保存文件
   - 观察是否触发重新加载

## 预期行为

- ✅ 文件修改后自动检测
- ✅ 自动重新编译
- ✅ 自动重启服务器
- ✅ 保持日志输出
- ✅ 错误时显示编译错误

## 故障排除

如果热加载不工作：

1. **检查air是否正确安装**
   ```bash
   air -v
   ```

2. **检查配置文件**
   - 确认 `.air.toml` 存在且配置正确

3. **检查文件权限**
   - 确保有写入 `tmp/` 目录的权限

4. **手动安装air**
   ```bash
   go install github.com/air-verse/air@latest
   ```

## 支持的开发环境

- ✅ Windows (PowerShell)
- ✅ Linux/macOS (Bash)
- ✅ 支持所有Go文件类型
- ✅ 支持配置文件修改 