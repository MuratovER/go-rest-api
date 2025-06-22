# Руководство по развертыванию

## 🚀 Обзор развертывания

Данное руководство описывает различные способы развертывания Inventory Assets API.

## 📋 Предварительные требования

### Системные требования

- **Go**: 1.24+
- **PostgreSQL**: 12+
- **Docker**: 20.10+ (опционально)
- **Docker Compose**: 2.0+ (опционально)

### Минимальные ресурсы

- **CPU**: 1 ядро
- **RAM**: 512 MB
- **Диск**: 1 GB свободного места

### Рекомендуемые ресурсы

- **CPU**: 2 ядра
- **RAM**: 2 GB
- **Диск**: 5 GB свободного места

## 🔧 Локальное развертывание

### 1. Установка зависимостей

```bash
# Клонирование репозитория
git clone <repository-url>
cd go-rest-api

# Установка зависимостей Go
go mod download
go mod tidy
```

### 2. Настройка базы данных

```bash
# Установка PostgreSQL (Ubuntu/Debian)
sudo apt update
sudo apt install postgresql postgresql-contrib

# Создание пользователя и базы данных
sudo -u postgres psql
CREATE USER inventory_user WITH PASSWORD 'secure_password';
CREATE DATABASE inventory_assets OWNER inventory_user;
GRANT ALL PRIVILEGES ON DATABASE inventory_assets TO inventory_user;
\q
```

### 3. Настройка конфигурации

```bash
# Создание файла конфигурации
cp .env.example .env

# Редактирование конфигурации
nano .env
```

Пример `.env` файла:

```env
# Сервер
SERVER_PORT=8080
SERVER_MODE=production
SERVER_APP_VERSION=1.0.0

# База данных
DB_HOST=localhost
DB_PORT=5432
DB_USER=inventory_user
DB_PASSWORD=secure_password
DB_NAME=inventory_assets
DB_SSL_MODE=disable

# Логирование
LOGGER_LEVEL=info
LOGGER_DEV_MODE=false
```

### 4. Сборка и запуск

```bash
# Сборка приложения
make build

# Запуск приложения
make run

# Или напрямую
go run cmd/api/main.go
```

### 5. Проверка работоспособности

```bash
# Проверка здоровья сервера
curl http://localhost:8080/health

# Проверка API
curl http://localhost:8080/api/v1/inventory-assets
```

## 🐳 Развертывание с Docker

### 1. Сборка образа

```bash
# Сборка Docker образа
make docker-build

# Или напрямую
docker build -t inventory-assets-api:latest .
```

### 2. Запуск с Docker Compose

```bash
# Запуск всех сервисов
make compose-up

# Или напрямую
docker-compose up --build -d
```

### 3. Проверка контейнеров

```bash
# Просмотр запущенных контейнеров
docker ps

# Просмотр логов
docker-compose logs -f api

# Остановка сервисов
make compose-down
```

### 4. Production конфигурация

Создайте `docker-compose.prod.yml`:

```yaml
version: '3.8'

services:
  api:
    build: .
    ports:
      - "8080:8080"
    environment:
      - SERVER_MODE=production
      - DB_HOST=postgres
      - DB_PORT=5432
      - DB_USER=inventory_user
      - DB_PASSWORD=${DB_PASSWORD}
      - DB_NAME=inventory_assets
      - DB_SSL_MODE=require
      - LOGGER_LEVEL=warn
    depends_on:
      - postgres
    restart: unless-stopped
    networks:
      - app-network

  postgres:
    image: postgres:13
    environment:
      - POSTGRES_DB=inventory_assets
      - POSTGRES_USER=inventory_user
      - POSTGRES_PASSWORD=${DB_PASSWORD}
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./init.sql:/docker-entrypoint-initdb.d/init.sql
    restart: unless-stopped
    networks:
      - app-network

  nginx:
    image: nginx:alpine
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf
      - ./ssl:/etc/nginx/ssl
    depends_on:
      - api
    restart: unless-stopped
    networks:
      - app-network

volumes:
  postgres_data:

networks:
  app-network:
    driver: bridge
```

## ☁️ Облачное развертывание

### AWS (Amazon Web Services)

#### 1. EC2 развертывание

```bash
# Подключение к EC2
ssh -i your-key.pem ubuntu@your-ec2-ip

# Установка зависимостей
sudo apt update
sudo apt install golang-go postgresql postgresql-contrib nginx

# Клонирование и настройка
git clone <repository-url>
cd go-rest-api
make build

# Настройка systemd сервиса
sudo nano /etc/systemd/system/inventory-api.service
```

Содержимое `inventory-api.service`:

```ini
[Unit]
Description=Inventory Assets API
After=network.target postgresql.service

[Service]
Type=simple
User=ubuntu
WorkingDirectory=/home/ubuntu/go-rest-api
ExecStart=/home/ubuntu/go-rest-api/bin/api
Restart=always
RestartSec=5
Environment=DB_HOST=localhost
Environment=DB_PORT=5432
Environment=DB_USER=inventory_user
Environment=DB_PASSWORD=secure_password
Environment=DB_NAME=inventory_assets

[Install]
WantedBy=multi-user.target
```

```bash
# Запуск сервиса
sudo systemctl daemon-reload
sudo systemctl enable inventory-api
sudo systemctl start inventory-api
sudo systemctl status inventory-api
```

#### 2. ECS (Elastic Container Service)

```yaml
# task-definition.json
{
  "family": "inventory-assets-api",
  "networkMode": "awsvpc",
  "requiresCompatibilities": ["FARGATE"],
  "cpu": "256",
  "memory": "512",
  "executionRoleArn": "arn:aws:iam::account:role/ecsTaskExecutionRole",
  "containerDefinitions": [
    {
      "name": "api",
      "image": "your-account.dkr.ecr.region.amazonaws.com/inventory-assets-api:latest",
      "portMappings": [
        {
          "containerPort": 8080,
          "protocol": "tcp"
        }
      ],
      "environment": [
        {
          "name": "SERVER_MODE",
          "value": "production"
        },
        {
          "name": "DB_HOST",
          "value": "your-rds-endpoint"
        }
      ],
      "logConfiguration": {
        "logDriver": "awslogs",
        "options": {
          "awslogs-group": "/ecs/inventory-assets-api",
          "awslogs-region": "us-east-1",
          "awslogs-stream-prefix": "ecs"
        }
      }
    }
  ]
}
```

### Google Cloud Platform (GCP)

#### 1. Cloud Run

```bash
# Сборка и отправка образа
gcloud builds submit --tag gcr.io/PROJECT_ID/inventory-assets-api

# Развертывание
gcloud run deploy inventory-assets-api \
  --image gcr.io/PROJECT_ID/inventory-assets-api \
  --platform managed \
  --region us-central1 \
  --allow-unauthenticated \
  --set-env-vars SERVER_MODE=production,DB_HOST=your-cloud-sql-ip
```

#### 2. GKE (Google Kubernetes Engine)

```yaml
# deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: inventory-assets-api
spec:
  replicas: 3
  selector:
    matchLabels:
      app: inventory-assets-api
  template:
    metadata:
      labels:
        app: inventory-assets-api
    spec:
      containers:
      - name: api
        image: gcr.io/PROJECT_ID/inventory-assets-api:latest
        ports:
        - containerPort: 8080
        env:
        - name: SERVER_MODE
          value: "production"
        - name: DB_HOST
          value: "your-cloud-sql-ip"
        resources:
          requests:
            memory: "256Mi"
            cpu: "250m"
          limits:
            memory: "512Mi"
            cpu: "500m"
---
apiVersion: v1
kind: Service
metadata:
  name: inventory-assets-api-service
spec:
  selector:
    app: inventory-assets-api
  ports:
  - protocol: TCP
    port: 80
    targetPort: 8080
  type: LoadBalancer
```

### Azure

#### 1. Azure Container Instances

```bash
# Сборка и отправка образа
az acr build --registry your-registry --image inventory-assets-api .

# Развертывание
az container create \
  --resource-group your-rg \
  --name inventory-assets-api \
  --image your-registry.azurecr.io/inventory-assets-api:latest \
  --dns-name-label inventory-assets-api \
  --ports 8080 \
  --environment-variables SERVER_MODE=production DB_HOST=your-sql-server
```

## 🔒 Безопасность

### 1. SSL/TLS сертификаты

```bash
# Получение сертификата Let's Encrypt
sudo certbot --nginx -d your-domain.com

# Автоматическое обновление
sudo crontab -e
# Добавить строку:
0 12 * * * /usr/bin/certbot renew --quiet
```

### 2. Firewall настройки

```bash
# UFW (Ubuntu)
sudo ufw allow 22/tcp
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
sudo ufw enable

# iptables
sudo iptables -A INPUT -p tcp --dport 22 -j ACCEPT
sudo iptables -A INPUT -p tcp --dport 80 -j ACCEPT
sudo iptables -A INPUT -p tcp --dport 443 -j ACCEPT
sudo iptables -A INPUT -j DROP
```

### 3. База данных безопасность

```sql
-- Создание пользователя только для чтения
CREATE USER read_only_user WITH PASSWORD 'read_only_password';
GRANT CONNECT ON DATABASE inventory_assets TO read_only_user;
GRANT USAGE ON SCHEMA public TO read_only_user;
GRANT SELECT ON ALL TABLES IN SCHEMA public TO read_only_user;

-- Ограничение подключений
ALTER USER inventory_user CONNECTION LIMIT 10;
```

## 📊 Мониторинг

### 1. Prometheus метрики

```go
// Добавьте в main.go
import (
    "github.com/prometheus/client_golang/prometheus/promhttp"
    "net/http"
)

func main() {
    // ... существующий код ...
    
    // Prometheus метрики
    go func() {
        http.Handle("/metrics", promhttp.Handler())
        http.ListenAndServe(":9090", nil)
    }()
}
```

### 2. Grafana дашборд

```json
{
  "dashboard": {
    "title": "Inventory Assets API",
    "panels": [
      {
        "title": "HTTP Requests",
        "type": "graph",
        "targets": [
          {
            "expr": "rate(http_requests_total[5m])",
            "legendFormat": "{{method}} {{endpoint}}"
          }
        ]
      },
      {
        "title": "Response Time",
        "type": "graph",
        "targets": [
          {
            "expr": "histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m]))",
            "legendFormat": "95th percentile"
          }
        ]
      }
    ]
  }
}
```

### 3. Логирование

```bash
# Настройка logrotate
sudo nano /etc/logrotate.d/inventory-api

# Содержимое:
/home/ubuntu/go-rest-api/logs/*.log {
    daily
    missingok
    rotate 30
    compress
    delaycompress
    notifempty
    create 644 ubuntu ubuntu
    postrotate
        systemctl reload inventory-api
    endscript
}
```

## 🔄 CI/CD Pipeline

### GitHub Actions

```yaml
# .github/workflows/deploy.yml
name: Deploy

on:
  push:
    branches: [main]

jobs:
  deploy:
    runs-on: ubuntu-latest
    
    steps:
    - uses: actions/checkout@v2
    
    - name: Set up Go
      uses: actions/setup-go@v2
      with:
        go-version: 1.24
    
    - name: Build and push Docker image
      run: |
        docker build -t your-registry/inventory-assets-api:${{ github.sha }} .
        docker push your-registry/inventory-assets-api:${{ github.sha }}
    
    - name: Deploy to production
      run: |
        # Команды развертывания
        ssh user@server "docker pull your-registry/inventory-assets-api:${{ github.sha }}"
        ssh user@server "docker-compose up -d"
```

### GitLab CI

```yaml
# .gitlab-ci.yml
stages:
  - test
  - build
  - deploy

test:
  stage: test
  script:
    - go test ./...

build:
  stage: build
  script:
    - docker build -t $CI_REGISTRY_IMAGE:$CI_COMMIT_SHA .
    - docker push $CI_REGISTRY_IMAGE:$CI_COMMIT_SHA

deploy:
  stage: deploy
  script:
    - ssh user@server "docker pull $CI_REGISTRY_IMAGE:$CI_COMMIT_SHA"
    - ssh user@server "docker-compose up -d"
```

## 🚨 Troubleshooting

### Частые проблемы

#### 1. Проблемы с подключением к БД

```bash
# Проверка подключения
psql -h localhost -U inventory_user -d inventory_assets

# Проверка логов PostgreSQL
sudo tail -f /var/log/postgresql/postgresql-*.log
```

#### 2. Проблемы с портами

```bash
# Проверка занятых портов
sudo netstat -tulpn | grep :8080

# Освобождение порта
sudo fuser -k 8080/tcp
```

#### 3. Проблемы с правами доступа

```bash
# Проверка прав на файлы
ls -la bin/api

# Установка правильных прав
chmod +x bin/api
chown ubuntu:ubuntu bin/api
```

### Логи и отладка

```bash
# Просмотр логов приложения
journalctl -u inventory-api -f

# Просмотр логов Docker
docker logs -f container_name

# Отладка с delve
dlv exec bin/api -- -config=config.yaml
```

## 📈 Масштабирование

### Горизонтальное масштабирование

```bash
# Docker Swarm
docker service scale inventory-api=3

# Kubernetes
kubectl scale deployment inventory-assets-api --replicas=5
```

### Вертикальное масштабирование

```bash
# Увеличение ресурсов контейнера
docker run -m 2g -c 2 inventory-assets-api

# Настройка лимитов в Kubernetes
resources:
  requests:
    memory: "1Gi"
    cpu: "500m"
  limits:
    memory: "2Gi"
    cpu: "1000m"
```

## 📞 Поддержка

### Контакты

- **Email**: support@example.com
- **Slack**: #inventory-api-support
- **Documentation**: https://docs.example.com

### Полезные команды

```bash
# Проверка статуса сервисов
make monitor

# Резервное копирование
make backup

# Обновление приложения
git pull
make build
sudo systemctl restart inventory-api
``` 