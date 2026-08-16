# bridge-shm

Go CLI. Use the commands below. Do not overwrite the project's own README.md.

## 构建 / 运行 / 测试

```text
go build ./...
go run .
go test ./...
```

## 评测镜像

- `benzhi.Dockerfile`
- `build_benzhi_docker.sh`
- `BENZHI_README.md`（本文件）

```bash
chmod +x build_benzhi_docker.sh
./build_benzhi_docker.sh bridge-shm linux/arm64
./build_benzhi_docker.sh bridge-shm linux/amd64
docker run -it bridge-shm:latest
```
