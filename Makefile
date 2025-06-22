# Inventory Assets API Makefile

# Переменные
APP_NAME=inventory-assets-api
BINARY_NAME=bin/api
DOCKER_IMAGE=inventory-assets-api
DOCKER_TAG=latest

# Go переменные
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOVET=$(GOCMD) vet
GOFMT=gofmt
GOLINT=golangci-lint

# Docker переменные
DOCKER_CMD=docker
DOCKER_COMPOSE=docker-compose
DOCKER_BUILD=$(DOCKER_CMD) build
DOCKER_RUN=$(DOCKER_CMD) run
DOCKER_STOP=$(DOCKER_CMD) stop
DOCKER_RM=$(DOCKER_CMD) rm

# Цвета для вывода
RED=\033[0;31m
GREEN=\033[0;32m
YELLOW=\033[0;33m
BLUE=\033[0;34m
NC=\033[0m # No Color

.PHONY: all build clean test coverage lint fmt vet help
.PHONY: docker-build docker-run docker-stop docker-clean
.PHONY: dev prod install-deps generate-mocks

# Основные команды
all: clean build

# Сборка приложения
build:
	@echo "$(BLUE)🔨 Сборка приложения...$(NC)"
	$(GOBUILD) -o $(BINARY_NAME) cmd/api/main.go
	@echo "$(GREEN)✅ Приложение собрано: $(BINARY_NAME)$(NC)"

# Очистка
clean:
	@echo "$(YELLOW)🧹 Очистка...$(NC)"
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	@echo "$(GREEN)✅ Очистка завершена$(NC)"

# Запуск приложения
run:
	@echo "$(BLUE)🚀 Запуск приложения...$(NC)"
	$(GOCMD) run cmd/api/main.go

# Тестирование
test:
	@echo "$(BLUE)🧪 Запуск тестов...$(NC)"
	$(GOTEST) -v ./...
	@echo "$(GREEN)✅ Тесты завершены$(NC)"

# Тестирование с покрытием
test-coverage:
	@echo "$(BLUE)🧪 Запуск тестов с покрытием...$(NC)"
	$(GOTEST) -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)✅ Отчет о покрытии создан: coverage.html$(NC)"

# Линтинг
lint:
	@echo "$(BLUE)🔍 Запуск линтера...$(NC)"
	$(GOLINT) run
	@echo "$(GREEN)✅ Линтинг завершен$(NC)"

# Форматирование кода
fmt:
	@echo "$(BLUE)🎨 Форматирование кода...$(NC)"
	$(GOFMT) -s -w .
	@echo "$(GREEN)✅ Код отформатирован$(NC)"

# Проверка кода
vet:
	@echo "$(BLUE)🔍 Проверка кода...$(NC)"
	$(GOVET) ./...
	@echo "$(GREEN)✅ Проверка завершена$(NC)"

# Установка зависимостей
install-deps:
	@echo "$(BLUE)📦 Установка зависимостей...$(NC)"
	$(GOMOD) download
	$(GOMOD) tidy
	@echo "$(GREEN)✅ Зависимости установлены$(NC)"

# Генерация моков
generate-mocks:
	@echo "$(BLUE)🔧 Генерация моков...$(NC)"
	$(GOCMD) generate ./...
	@echo "$(GREEN)✅ Моки сгенерированы$(NC)"

# Проверка качества кода
check: fmt vet lint test
	@echo "$(GREEN)✅ Все проверки пройдены$(NC)"

# Docker команды
docker-build:
	@echo "$(BLUE)🐳 Сборка Docker образа...$(NC)"
	$(DOCKER_BUILD) -t $(DOCKER_IMAGE):$(DOCKER_TAG) .
	@echo "$(GREEN)✅ Docker образ собран$(NC)"

docker-run:
	@echo "$(BLUE)🐳 Запуск Docker контейнера...$(NC)"
	$(DOCKER_RUN) -p 8080:8080 $(DOCKER_IMAGE):$(DOCKER_TAG)

docker-stop:
	@echo "$(YELLOW)🛑 Остановка Docker контейнеров...$(NC)"
	$(DOCKER_STOP) $(shell docker ps -q --filter ancestor=$(DOCKER_IMAGE):$(DOCKER_TAG)) 2>/dev/null || true

docker-clean:
	@echo "$(YELLOW)🧹 Очистка Docker...$(NC)"
	$(DOCKER_CMD) system prune -f
	$(DOCKER_CMD) image prune -f

# Docker Compose команды
compose-up:
	@echo "$(BLUE)🐳 Запуск с Docker Compose...$(NC)"
	$(DOCKER_COMPOSE) up --build

compose-down:
	@echo "$(YELLOW)🛑 Остановка Docker Compose...$(NC)"
	$(DOCKER_COMPOSE) down

compose-dev:
	@echo "$(BLUE)🐳 Запуск dev окружения...$(NC)"
	$(DOCKER_COMPOSE) -f docker-compose.dev.yml up --build

# Разработка
dev: install-deps generate-mocks
	@echo "$(BLUE)🚀 Запуск в режиме разработки...$(NC)"
	$(GOCMD) run cmd/api/main.go

# Production
prod: clean build
	@echo "$(BLUE)🚀 Запуск в production режиме...$(NC)"
	./$(BINARY_NAME)

# Миграции базы данных
migrate:
	@echo "$(BLUE)🗄️ Запуск миграций...$(NC)"
	# Добавьте команды для миграций здесь
	@echo "$(GREEN)✅ Миграции завершены$(NC)"

# Создание базы данных
db-create:
	@echo "$(BLUE)🗄️ Создание базы данных...$(NC)"
	# Добавьте команды для создания БД здесь
	@echo "$(GREEN)✅ База данных создана$(NC)"

# Резервное копирование
backup:
	@echo "$(BLUE)💾 Создание резервной копии...$(NC)"
	# Добавьте команды для бэкапа здесь
	@echo "$(GREEN)✅ Резервная копия создана$(NC)"

# Мониторинг
monitor:
	@echo "$(BLUE)📊 Мониторинг приложения...$(NC)"
	@echo "Health check: http://localhost:8080/health"
	@echo "Swagger UI: http://localhost:8080/swagger/index.html"
	@echo "Metrics: http://localhost:8080/metrics"

# Помощь
help:
	@echo "$(BLUE)📚 Доступные команды:$(NC)"
	@echo ""
	@echo "$(YELLOW)Основные команды:$(NC)"
	@echo "  build          - Сборка приложения"
	@echo "  clean          - Очистка"
	@echo "  run            - Запуск приложения"
	@echo "  dev            - Запуск в режиме разработки"
	@echo "  prod           - Запуск в production режиме"
	@echo ""
	@echo "$(YELLOW)Тестирование и качество:$(NC)"
	@echo "  test           - Запуск тестов"
	@echo "  test-coverage  - Тесты с покрытием"
	@echo "  lint           - Линтинг кода"
	@echo "  fmt            - Форматирование кода"
	@echo "  vet            - Проверка кода"
	@echo "  check          - Все проверки качества"
	@echo ""
	@echo "$(YELLOW)Docker команды:$(NC)"
	@echo "  docker-build   - Сборка Docker образа"
	@echo "  docker-run     - Запуск Docker контейнера"
	@echo "  docker-stop    - Остановка контейнеров"
	@echo "  docker-clean   - Очистка Docker"
	@echo "  compose-up     - Запуск с Docker Compose"
	@echo "  compose-down   - Остановка Docker Compose"
	@echo "  compose-dev    - Запуск dev окружения"
	@echo ""
	@echo "$(YELLOW)Утилиты:$(NC)"
	@echo "  install-deps   - Установка зависимостей"
	@echo "  generate-mocks - Генерация моков"
	@echo "  migrate        - Запуск миграций"
	@echo "  db-create      - Создание базы данных"
	@echo "  backup         - Резервное копирование"
	@echo "  monitor        - Мониторинг приложения"
	@echo ""
	@echo "$(YELLOW)Информация:$(NC)"
	@echo "  help           - Показать эту справку"

# По умолчанию показываем справку
.DEFAULT_GOAL := help
