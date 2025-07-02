## 详细步骤：构建、运行及 Docker Compose 配置
### 1. 构建 Docker 镜像
```dockerfile
在项目根目录（含 `Dockerfile` 的目录）执行：
docker build -t myapp:latest .
```

+ **<font style="color:rgb(64, 64, 64);">-t myapp:latest</font>**<font style="color:rgb(64, 64, 64);">：指定镜像名称和标签</font>
+ **<font style="color:rgb(64, 64, 64);background-color:rgb(236, 236, 236);">.</font>**<font style="color:rgb(64, 64, 64);">：使用当前目录作为构建上下文</font>

**<font style="color:rgb(64, 64, 64);">优化建议</font>**<font style="color:rgb(64, 64, 64);">：若需清理构建缓存，添加</font><font style="color:rgb(64, 64, 64);"> </font>**<font style="color:rgb(64, 64, 64);background-color:rgb(236, 236, 236);">--no-cache</font>**<font style="color:rgb(64, 64, 64);"> </font><font style="color:rgb(64, 64, 64);">参数</font>

---

### <font style="color:rgb(64, 64, 64);">2. 运行容器</font>
```dockerfile
docker run -d \
  --name myapp_container \
  -p 8080:8080 \
  -v ./config:/app/config \  
  myapp:latest
```

+ **<font style="color:rgb(64, 64, 64);background-color:rgb(236, 236, 236);">-d</font>**<font style="color:rgb(64, 64, 64);">：后台运行</font>
+ **<font style="color:rgb(64, 64, 64);background-color:rgb(236, 236, 236);">-p 8080:8080</font>**<font style="color:rgb(64, 64, 64);">：映射容器端口到宿主机</font>
+ **<font style="color:rgb(64, 64, 64);background-color:rgb(236, 236, 236);">-v ./config:/app/config</font>**<font style="color:rgb(64, 64, 64);">：挂载配置文件（避免重建镜像）</font>

---

### <font style="color:rgb(64, 64, 64);">3.</font><font style="color:rgb(64, 64, 64);"> </font>**<font style="color:rgb(64, 64, 64);background-color:rgb(236, 236, 236);">.dockerignore</font>**<font style="color:rgb(64, 64, 64);"> </font><font style="color:rgb(64, 64, 64);">文件</font>
<font style="color:rgb(64, 64, 64);">在项目根目录创建 </font>**<font style="color:rgb(64, 64, 64);background-color:rgb(236, 236, 236);">.dockerignore</font>**<font style="color:rgb(64, 64, 64);">，忽略无关文件：</font>

```plain
# 忽略文件示例
.git
.gitignore
.dockerignore
Dockerfile
docker-compose.yml
**/*.log
**/tmp
**/.vscode
**/README.md
bin/
dist/
vendor/  # Go 项目通常不需要
```

**<font style="color:rgb(64, 64, 64);">作用</font>**<font style="color:rgb(64, 64, 64);">：</font>

+ <font style="color:rgb(64, 64, 64);">加速构建过程</font>
+ <font style="color:rgb(64, 64, 64);">减小镜像体积</font>
+ <font style="color:rgb(64, 64, 64);">避免泄露敏感文件</font>

---

### <font style="color:rgb(64, 64, 64);">4.</font><font style="color:rgb(64, 64, 64);"> </font>**<font style="color:rgb(64, 64, 64);background-color:rgb(236, 236, 236);">docker-compose.yml</font>**<font style="color:rgb(64, 64, 64);"> </font><font style="color:rgb(64, 64, 64);">配置</font>
<font style="color:rgb(64, 64, 64);">创建 </font>**<font style="color:rgb(64, 64, 64);background-color:rgb(236, 236, 236);">docker-compose.yml</font>**<font style="color:rgb(64, 64, 64);">：</font>

```yaml
version: '3.8'

services:
  myapp:
    image: myapp:latest
    build: 
      context: .  # 构建上下文路径
      dockerfile: Dockerfile
    container_name: myapp_prod
    restart: unless-stopped
    ports:
      - "8080:8080"
    volumes:
      - ./config:/app/config  # 持久化配置文件
    environment:
      - TZ=Asia/Shanghai  # 覆盖时区变量
    networks:
      - myapp_net

networks:
  myapp_net:
    driver: bridge
```

**<font style="color:rgb(64, 64, 64);">关键配置说明</font>**<font style="color:rgb(64, 64, 64);">：</font>

+ **<font style="color:rgb(64, 64, 64);background-color:rgb(236, 236, 236);">build</font>**<font style="color:rgb(64, 64, 64);">：支持直接构建镜像</font>
+ **<font style="color:rgb(64, 64, 64);background-color:rgb(236, 236, 236);">volumes</font>**<font style="color:rgb(64, 64, 64);">：动态挂载配置文件目录</font>
+ **<font style="color:rgb(64, 64, 64);background-color:rgb(236, 236, 236);">restart: unless-stopped</font>**<font style="color:rgb(64, 64, 64);">：异常退出时自动重启</font>
+ <font style="color:rgb(64, 64, 64);">独立网络：隔离容器网络环境</font>

---

### <font style="color:rgb(64, 64, 64);">5. 使用 Docker Compose 操作</font>
```bash
# 启动服务（后台运行）
docker-compose up -d

# 查看日志
docker-compose logs -f

# 停止服务
docker-compose down

# 重新构建并启动（修改代码后）
docker-compose up -d --build
```

---

### <font style="color:rgb(64, 64, 64);">6. 配置文件注意事项</font>
1. <font style="color:rgb(64, 64, 64);">确保项目目录存在</font><font style="color:rgb(64, 64, 64);"> </font>**<font style="color:rgb(64, 64, 64);background-color:rgb(236, 236, 236);">config/config.ini</font>**<font style="color:rgb(64, 64, 64);"> </font><font style="color:rgb(64, 64, 64);">文件</font>

<font style="color:rgb(64, 64, 64);">目录结构示例：</font>

```latex
project-root/
├── Dockerfile
├── docker-compose.yml
├── .dockerignore
├── go.mod
├── main.go
└── config/
    └── config.ini  # 应用配置文件
```

<font style="color:rgb(64, 64, 64);">若配置文件需热更新，在容器内使用 </font>**<font style="color:rgb(64, 64, 64);background-color:rgb(236, 236, 236);">SIGHUP</font>**<font style="color:rgb(64, 64, 64);"> 信号通知应用重载：</font>

```bash
docker kill -s HUP myapp_prod
```

---

### <font style="color:rgb(64, 64, 64);">常见问题排查</font>
1. **<font style="color:rgb(64, 64, 64);">端口冲突</font>**<font style="color:rgb(64, 64, 64);">检查宿主机 8080 端口是否被占用：</font>**<font style="color:rgb(64, 64, 64);background-color:rgb(236, 236, 236);">netstat -tuln | grep 8080</font>**
2. **<font style="color:rgb(64, 64, 64);">配置文件权限</font>**<font style="color:rgb(64, 64, 64);">若容器启动失败，检查挂载目录权限：</font>**<font style="color:rgb(64, 64, 64);background-color:rgb(236, 236, 236);">chmod -R a+r ./config</font>**
3. **<font style="color:rgb(64, 64, 64);">时区问题</font>**<font style="color:rgb(64, 64, 64);">确保宿主机时区正确：</font>**<font style="color:rgb(64, 64, 64);background-color:rgb(236, 236, 236);">timedatectl set-timezone Asia/Shanghai</font>**
4. **<font style="color:rgb(64, 64, 64);">Go 依赖下载失败</font>**<font style="color:rgb(64, 64, 64);">在 </font>**<font style="color:rgb(64, 64, 64);background-color:rgb(236, 236, 236);">Dockerfile</font>**<font style="color:rgb(64, 64, 64);"> 中调整 </font>**<font style="color:rgb(64, 64, 64);background-color:rgb(236, 236, 236);">GOPROXY</font>**<font style="color:rgb(64, 64, 64);">：</font>

```bash
ENV GOPROXY=https://goproxy.cn,direct \
    GOPRIVATE=*.corp.example.com  # 私有模块配置
```

---

## <font style="color:rgb(64, 64, 64);">详细步骤：Node.js 应用镜像构建与部署</font>
---

### <font style="color:rgb(64, 64, 64);">1. 项目结构准备</font>
<font style="color:rgb(64, 64, 64);">确保项目目录结构如下：</font>

```latex
my-node-app/
├── src/              # 源代码目录
├── public/           # 静态资源
├── package.json
├── package-lock.json
├── Dockerfile
├── docker-compose.yml
├── .dockerignore
└── nginx.conf        # 可选自定义配置
```

---

### <font style="color:rgb(64, 64, 64);">2. 创建</font><font style="color:rgb(64, 64, 64);"> </font>**<font style="color:rgb(64, 64, 64);background-color:rgb(236, 236, 236);">.dockerignore</font>**<font style="color:rgb(64, 64, 64);"> </font><font style="color:rgb(64, 64, 64);">文件</font>
<font style="color:rgb(64, 64, 64);">在项目根目录创建 </font>**<font style="color:rgb(64, 64, 64);background-color:rgb(236, 236, 236);">.dockerignore</font>**<font style="color:rgb(64, 64, 64);">：</font>

```plain
# 基础忽略
.git
.gitignore
.dockerignore
node_modules
npm-debug.log
.DS_Store

# 构建过程忽略
dist
build

# 开发环境文件
.vscode
.idea
*.md
```

**<font style="color:rgb(64, 64, 64);">作用</font>**<font style="color:rgb(64, 64, 64);">：</font>

+ <font style="color:rgb(64, 64, 64);">减少构建上下文大小（加速构建）</font>
+ <font style="color:rgb(64, 64, 64);">避免泄露敏感文件</font>
+ <font style="color:rgb(64, 64, 64);">防止覆盖生产依赖</font>

---

### <font style="color:rgb(64, 64, 64);">3. 构建 Docker 镜像</font>
```plain
docker build -t my-node-app:1.0.0 .
```

**<font style="color:rgb(64, 64, 64);">构建过程解析</font>**<font style="color:rgb(64, 64, 64);">：</font>

1. <font style="color:rgb(64, 64, 64);">下载 Node.js 基础镜像</font>
2. <font style="color:rgb(64, 64, 64);">配置 npm 国内镜像源（加速依赖安装）</font>
3. <font style="color:rgb(64, 64, 64);">安装依赖（利用 Docker 层缓存优化）</font>
4. <font style="color:rgb(64, 64, 64);">构建应用（生成 dist 目录）</font>
5. <font style="color:rgb(64, 64, 64);">清理构建缓存（减小镜像体积）</font>
6. <font style="color:rgb(64, 64, 64);">复制构建产物到 Nginx 镜像</font>

---

### <font style="color:rgb(64, 64, 64);">4. 运行容器</font>
```bash
docker run -d \
  --name my-web-app \
  -p 80:80 \               
  -v $(pwd)/nginx.conf:/etc/nginx/conf.d/default.conf \ 
  my-node-app:1.0.0
```

**<font style="color:rgb(64, 64, 64);">访问应用</font>**<font style="color:rgb(64, 64, 64);">：浏览器打开</font><font style="color:rgb(64, 64, 64);"> </font>[**<font style="color:rgb(64, 64, 64);background-color:rgb(236, 236, 236);">http://localhost:80</font>**](http://localhost:80)

---

### <font style="color:rgb(64, 64, 64);">5. 创建 </font>**<font style="color:rgb(64, 64, 64);background-color:rgb(236, 236, 236);">docker-compose.yml</font>**
```yaml
version: '3.8'

services:
  webapp:
    image: my-node-app:1.0.0
    build: 
      context: .   # 构建上下文路径
      dockerfile: Dockerfile
    container_name: node_app_prod
    restart: always
    ports:
      - "80:80"
    volumes:
      - ./nginx.conf:/etc/nginx/conf.d/default.conf  # 挂载自定义配置
    environment:
      - TZ=Asia/Shanghai
    networks:
      - app_network

networks:
  app_network:
    driver: bridge
```

---

### <font style="color:rgb(64, 64, 64);">6. 使用 Docker Compose 管理</font>
```bash
# 启动服务（后台运行）
docker-compose up -d

# 查看实时日志
docker-compose logs -f

# 停止并删除服务
docker-compose down

# 重新构建镜像并启动（代码更新后）
docker-compose up -d --build
```

---

### <font style="color:rgb(64, 64, 64);">7. 自定义 Nginx 配置（可选）</font>
<font style="color:rgb(64, 64, 64);">创建 </font>**<font style="color:rgb(64, 64, 64);background-color:rgb(236, 236, 236);">nginx.conf</font>**<font style="color:rgb(64, 64, 64);"> 文件优化性能：</font>

```bash
server {
  listen 80;
  server_name localhost;
  
  root /usr/share/nginx/html;
  
  # 静态资源缓存
  location ~* \.(js|css|png|jpg|jpeg|gif|ico|svg)$ {
    expires 30d;
    add_header Cache-Control "public, no-transform";
  }

  # 前端路由支持
  location / {
    try_files $uri $uri/ /index.html;
  }

  # 禁止访问 .env 文件
  location ~ /\.env {
    deny all;
    return 404;
  }

  # Gzip 压缩
  gzip on;
  gzip_types text/docekrfile text/css application/json application/javascript text/xml;
}
```

---

### <font style="color:rgb(64, 64, 64);">常见问题解决方案</font>
#### **<font style="color:rgb(64, 64, 64);">问题1：npm install 超时</font>**
**<font style="color:rgb(64, 64, 64);">解决方案</font>**<font style="color:rgb(64, 64, 64);">：在 Dockerfile 中增加超时设置</font>

```bash
RUN npm config set fetch-retry-mintimeout 20000 \
  && npm config set fetch-retry-maxtimeout 120000 \
  && npm install --force
```

#### **<font style="color:rgb(64, 64, 64);">问题2：容器启动后显示 403 Forbidden</font>**
**<font style="color:rgb(64, 64, 64);">原因</font>**<font style="color:rgb(64, 64, 64);">：Nginx 权限问题 </font>**<font style="color:rgb(64, 64, 64);">修复</font>**<font style="color:rgb(64, 64, 64);">：在 Dockerfile 添加权限设置</font>

```bash
RUN chown -R nginx:nginx /usr/share/nginx/html
```

#### **<font style="color:rgb(64, 64, 64);">问题3：构建产物过大</font>**
**<font style="color:rgb(64, 64, 64);">优化方案</font>**<font style="color:rgb(64, 64, 64);">：多阶段构建优化</font>

```bash
# 最终阶段删除缓存
RUN rm -rf /usr/share/nginx/html/node_modules \
  && rm -rf /usr/share/nginx/html/.npm
```

---

### <font style="color:rgb(64, 64, 64);">最佳实践建议</font>
1. **<font style="color:rgb(64, 64, 64);">镜像标签</font>**<font style="color:rgb(64, 64, 64);">使用语义化版本</font><font style="color:rgb(64, 64, 64);"> </font>**<font style="color:rgb(64, 64, 64);background-color:rgb(236, 236, 236);">my-app:1.2.3</font>**<font style="color:rgb(64, 64, 64);"> </font><font style="color:rgb(64, 64, 64);">代替</font><font style="color:rgb(64, 64, 64);"> </font>**<font style="color:rgb(64, 64, 64);background-color:rgb(236, 236, 236);">latest</font>**
2. **<font style="color:rgb(64, 64, 64);">安全扫描</font>**<font style="color:rgb(64, 64, 64);">构建后扫描镜像：</font>**<font style="color:rgb(64, 64, 64);background-color:rgb(236, 236, 236);">docker scan my-node-app:1.0.0</font>**

**<font style="color:rgb(64, 64, 64);">资源限制</font>**<font style="color:rgb(64, 64, 64);">在 docker-compose 中配置资源限制：</font>

```yaml
deploy:
  resources:
    limits:
      cpus: '0.5'
      memory: 512M
```

**<font style="color:rgb(64, 64, 64);">健康检查</font>**<font style="color:rgb(64, 64, 64);">添加容器健康监测：</font>

```yaml
healthcheck:
  test: ["CMD", "curl", "-f", "http://localhost:80"]
  interval: 30s
  timeout: 5s
  retries: 3
```

<font style="color:rgb(64, 64, 64);">通过以上配置，可实现高效开发部署流程，同时确保生产环境稳定运行。</font>

