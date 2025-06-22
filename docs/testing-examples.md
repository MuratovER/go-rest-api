# Примеры тестирования

## 🧪 Обзор тестирования

Проект использует стандартный пакет `testing` Go для модульных тестов и `testify` для дополнительных утилит.

## 📁 Структура тестов

```
├── internal/
│   ├── inventory_asset/
│   │   ├── delivery/http/
│   │   │   └── handlers_test.go
│   │   ├── repository/
│   │   │   └── pg_repository_test.go
│   │   └── usecase/
│   │       └── usecase_test.go
│   └── models/
│       └── inventory_asset_test.go
└── pkg/
    └── utils/
        └── pagination_test.go
```

## 🔧 Настройка тестов

### Установка зависимостей для тестирования

```bash
go get github.com/stretchr/testify/assert
go get github.com/stretchr/testify/mock
go get github.com/stretchr/testify/suite
```

### Конфигурация для тестов

```go
// test_config.go
package test

import (
    "inventory_assets/config"
    "os"
)

func GetTestConfig() *config.Config {
    return &config.Config{
        Server: config.ServerConfig{
            Port:       "8080",
            Mode:       "test",
            AppVersion: "1.0.0",
        },
        Database: config.DatabaseConfig{
            Host:     "localhost",
            Port:     "5432",
            User:     "test_user",
            Password: "test_password",
            Name:     "test_db",
            SSLMode:  "disable",
        },
        Logger: config.LoggerConfig{
            Level:    "debug",
            DevMode:  true,
        },
    }
}
```

## 📝 Примеры тестов

### 1. Тестирование моделей

```go
// internal/models/inventory_asset_test.go
package models

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestInventoryAsset_TableName(t *testing.T) {
    asset := InventoryAsset{}
    assert.Equal(t, "inventory_assets", asset.TableName())
}

func TestInventoryAsset_Validation(t *testing.T) {
    tests := []struct {
        name    string
        asset   InventoryAsset
        isValid bool
    }{
        {
            name: "Valid asset",
            asset: InventoryAsset{
                UserName:   "john_doe",
                Name:       "Laptop Dell XPS 13",
                SerialCode: "DLXPS13001",
                IsActive:   true,
                Price:      120000,
                URL:        "https://example.com/laptop1",
            },
            isValid: true,
        },
        {
            name: "Invalid asset - empty name",
            asset: InventoryAsset{
                UserName:   "john_doe",
                Name:       "",
                SerialCode: "DLXPS13001",
                IsActive:   true,
                Price:      120000,
                URL:        "https://example.com/laptop1",
            },
            isValid: false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if tt.isValid {
                assert.NotEmpty(t, tt.asset.Name)
                assert.NotEmpty(t, tt.asset.UserName)
                assert.NotEmpty(t, tt.asset.SerialCode)
                assert.NotEmpty(t, tt.asset.URL)
                assert.Greater(t, tt.asset.Price, 0)
            } else {
                assert.Empty(t, tt.asset.Name)
            }
        })
    }
}
```

### 2. Тестирование репозитория

```go
// internal/inventory_asset/repository/pg_repository_test.go
package repository

import (
    "context"
    "testing"
    "inventory_assets/internal/models"
    "inventory_assets/pkg/utils"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
)

type MockLogger struct {
    mock.Mock
}

func (m *MockLogger) Info(args ...interface{}) {}
func (m *MockLogger) Infof(format string, args ...interface{}) {}
func (m *MockLogger) Error(args ...interface{}) {}
func (m *MockLogger) Errorf(format string, args ...interface{}) {}
func (m *MockLogger) Fatal(args ...interface{}) {}
func (m *MockLogger) Fatalf(format string, args ...interface{}) {}
func (m *MockLogger) InitLogger() {}

func setupTestDB(t *testing.T) *gorm.DB {
    db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
    assert.NoError(t, err)
    
    // Миграция
    err = db.AutoMigrate(&models.InventoryAsset{})
    assert.NoError(t, err)
    
    return db
}

func TestInventoryAssetRepo_GetInventoryAssets(t *testing.T) {
    db := setupTestDB(t)
    logger := &MockLogger{}
    repo := NewInventoryAssetRepository(db, logger)
    
    // Создание тестовых данных
    testAssets := []models.InventoryAsset{
        {
            UserName:   "john_doe",
            Name:       "Laptop Dell XPS 13",
            SerialCode: "DLXPS13001",
            IsActive:   true,
            Price:      120000,
            URL:        "https://example.com/laptop1",
        },
        {
            UserName:   "jane_smith",
            Name:       "Monitor Samsung 27\"",
            SerialCode: "SMS27001",
            IsActive:   true,
            Price:      45000,
            URL:        "https://example.com/monitor1",
        },
    }
    
    for _, asset := range testAssets {
        db.Create(&asset)
    }
    
    // Тест
    pq := &utils.PaginationQuery{
        Page: 1,
        Size: 10,
    }
    
    result, err := repo.GetInventoryAssets(context.Background(), pq)
    
    assert.NoError(t, err)
    assert.NotNil(t, result)
    assert.Equal(t, int64(2), result.TotalCount)
    assert.Len(t, *result.Items, 2)
    assert.Equal(t, 1, result.Page)
    assert.Equal(t, 10, result.Size)
}

func TestInventoryAssetRepo_GetInventoryAssetsById(t *testing.T) {
    db := setupTestDB(t)
    logger := &MockLogger{}
    repo := NewInventoryAssetRepository(db, logger)
    
    // Создание тестового актива
    testAsset := models.InventoryAsset{
        UserName:   "john_doe",
        Name:       "Laptop Dell XPS 13",
        SerialCode: "DLXPS13001",
        IsActive:   true,
        Price:      120000,
        URL:        "https://example.com/laptop1",
    }
    
    db.Create(&testAsset)
    
    // Тест успешного получения
    pq := &utils.PaginationQuery{
        Page: 1,
        Size: 10,
    }
    
    result, err := repo.GetInventoryAssetsById(context.Background(), pq, testAsset.ID)
    
    assert.NoError(t, err)
    assert.NotNil(t, result)
    assert.Equal(t, int64(1), result.TotalCount)
    assert.Len(t, *result.Items, 1)
    assert.Equal(t, testAsset.Name, (*result.Items)[0].Name)
    
    // Тест несуществующего ID
    result, err = repo.GetInventoryAssetsById(context.Background(), pq, 999)
    
    assert.Error(t, err)
    assert.Nil(t, result)
    assert.Contains(t, err.Error(), "inventory asset with this id not found")
}
```

### 3. Тестирование Use Cases

```go
// internal/inventory_asset/usecase/usecase_test.go
package usecase

import (
    "context"
    "testing"
    "inventory_assets/internal/models"
    "inventory_assets/pkg/utils"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

type MockRepository struct {
    mock.Mock
}

func (m *MockRepository) GetInventoryAssetsById(ctx context.Context, pq *utils.PaginationQuery, id int) (*models.InventoryAssetList, error) {
    args := m.Called(ctx, pq, id)
    return args.Get(0).(*models.InventoryAssetList), args.Error(1)
}

func (m *MockRepository) GetInventoryAssets(ctx context.Context, pq *utils.PaginationQuery) (*models.InventoryAssetList, error) {
    args := m.Called(ctx, pq)
    return args.Get(0).(*models.InventoryAssetList), args.Error(1)
}

type MockLogger struct {
    mock.Mock
}

func (m *MockLogger) Info(args ...interface{}) {}
func (m *MockLogger) Infof(format string, args ...interface{}) {}
func (m *MockLogger) Error(args ...interface{}) {}
func (m *MockLogger) Errorf(format string, args ...interface{}) {}
func (m *MockLogger) Fatal(args ...interface{}) {}
func (m *MockLogger) Fatalf(format string, args ...interface{}) {}
func (m *MockLogger) InitLogger() {}

func TestInventoryAssetUC_GetInventoryAssets(t *testing.T) {
    mockRepo := new(MockRepository)
    mockLogger := new(MockLogger)
    cfg := &config.Config{}
    
    uc := NewInventoryAssetUseCase(cfg, mockRepo, mockLogger)
    
    ctx := context.Background()
    pq := &utils.PaginationQuery{Page: 1, Size: 10}
    
    expectedResult := &models.InventoryAssetList{
        Items:      &[]models.InventoryAsset{},
        TotalCount: 0,
        Page:       1,
        Size:       10,
        TotalPages: 0,
    }
    
    mockRepo.On("GetInventoryAssets", ctx, pq).Return(expectedResult, nil)
    
    result, err := uc.GetInventoryAssets(ctx, pq)
    
    assert.NoError(t, err)
    assert.Equal(t, expectedResult, result)
    mockRepo.AssertExpectations(t)
}

func TestInventoryAssetUC_GetInventoryAssetsById(t *testing.T) {
    mockRepo := new(MockRepository)
    mockLogger := new(MockLogger)
    cfg := &config.Config{}
    
    uc := NewInventoryAssetUseCase(cfg, mockRepo, mockLogger)
    
    ctx := context.Background()
    pq := &utils.PaginationQuery{Page: 1, Size: 10}
    assetID := 1
    
    expectedAsset := models.InventoryAsset{
        ID:         1,
        UserName:   "john_doe",
        Name:       "Laptop Dell XPS 13",
        SerialCode: "DLXPS13001",
        IsActive:   true,
        Price:      120000,
        URL:        "https://example.com/laptop1",
    }
    
    expectedResult := &models.InventoryAssetList{
        Items:      &[]models.InventoryAsset{expectedAsset},
        TotalCount: 1,
        Page:       1,
        Size:       10,
        TotalPages: 1,
    }
    
    mockRepo.On("GetInventoryAssetsById", ctx, pq, assetID).Return(expectedResult, nil)
    
    result, err := uc.GetInventoryAssetsById(ctx, pq, assetID)
    
    assert.NoError(t, err)
    assert.Equal(t, expectedResult, result)
    mockRepo.AssertExpectations(t)
}
```

### 4. Тестирование HTTP Handlers

```go
// internal/inventory_asset/delivery/http/handlers_test.go
package http

import (
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "strconv"
    "testing"
    "inventory_assets/config"
    "inventory_assets/internal/models"
    "inventory_assets/pkg/utils"
    "github.com/labstack/echo/v4"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

type MockUseCase struct {
    mock.Mock
}

func (m *MockUseCase) GetInventoryAssetsById(ctx context.Context, pq *utils.PaginationQuery, id int) (*models.InventoryAssetList, error) {
    args := m.Called(ctx, pq, id)
    return args.Get(0).(*models.InventoryAssetList), args.Error(1)
}

func (m *MockUseCase) GetInventoryAssets(ctx context.Context, pq *utils.PaginationQuery) (*models.InventoryAssetList, error) {
    args := m.Called(ctx, pq)
    return args.Get(0).(*models.InventoryAssetList), args.Error(1)
}

type MockLogger struct {
    mock.Mock
}

func (m *MockLogger) Info(args ...interface{}) {}
func (m *MockLogger) Infof(format string, args ...interface{}) {}
func (m *MockLogger) Error(args ...interface{}) {}
func (m *MockLogger) Errorf(format string, args ...interface{}) {}
func (m *MockLogger) Fatal(args ...interface{}) {}
func (m *MockLogger) Fatalf(format string, args ...interface{}) {}
func (m *MockLogger) InitLogger() {}

func setupEcho() *echo.Echo {
    e := echo.New()
    return e
}

func TestInventoryAssetHandlers_GetInventoryAssets(t *testing.T) {
    e := setupEcho()
    mockUC := new(MockUseCase)
    mockLogger := new(MockLogger)
    cfg := &config.Config{}
    
    handler := NewInventoryAssetHandlers(cfg, mockUC, mockLogger)
    
    // Создание тестового запроса
    req := httptest.NewRequest(http.MethodGet, "/api/v1/inventory-assets?page=1&size=10", nil)
    rec := httptest.NewRecorder()
    c := e.NewContext(req, rec)
    
    expectedAsset := models.InventoryAsset{
        ID:         1,
        UserName:   "john_doe",
        Name:       "Laptop Dell XPS 13",
        SerialCode: "DLXPS13001",
        IsActive:   true,
        Price:      120000,
        URL:        "https://example.com/laptop1",
    }
    
    expectedResult := &models.InventoryAssetList{
        Items:      &[]models.InventoryAsset{expectedAsset},
        TotalCount: 1,
        Page:       1,
        Size:       10,
        TotalPages: 1,
    }
    
    mockUC.On("GetInventoryAssets", mock.Anything, mock.Anything).Return(expectedResult, nil)
    
    // Выполнение запроса
    err := handler.GetInventoryAssets()(c)
    
    assert.NoError(t, err)
    assert.Equal(t, http.StatusOK, rec.Code)
    
    // Проверка ответа
    var response models.InventoryAssetList
    err = json.Unmarshal(rec.Body.Bytes(), &response)
    assert.NoError(t, err)
    assert.Equal(t, expectedResult.TotalCount, response.TotalCount)
    assert.Len(t, *response.Items, 1)
    
    mockUC.AssertExpectations(t)
}

func TestInventoryAssetHandlers_GetInventoryAssetsById(t *testing.T) {
    e := setupEcho()
    mockUC := new(MockUseCase)
    mockLogger := new(MockLogger)
    cfg := &config.Config{}
    
    handler := NewInventoryAssetHandlers(cfg, mockUC, mockLogger)
    
    assetID := 1
    
    // Создание тестового запроса
    req := httptest.NewRequest(http.MethodGet, "/api/v1/inventory-assets/"+strconv.Itoa(assetID), nil)
    rec := httptest.NewRecorder()
    c := e.NewContext(req, rec)
    c.SetParamNames("id")
    c.SetParamValues(strconv.Itoa(assetID))
    
    expectedAsset := models.InventoryAsset{
        ID:         assetID,
        UserName:   "john_doe",
        Name:       "Laptop Dell XPS 13",
        SerialCode: "DLXPS13001",
        IsActive:   true,
        Price:      120000,
        URL:        "https://example.com/laptop1",
    }
    
    expectedResult := &models.InventoryAssetList{
        Items:      &[]models.InventoryAsset{expectedAsset},
        TotalCount: 1,
        Page:       1,
        Size:       10,
        TotalPages: 1,
    }
    
    mockUC.On("GetInventoryAssetsById", mock.Anything, mock.Anything, assetID).Return(expectedResult, nil)
    
    // Выполнение запроса
    err := handler.GetInventoryAssetsById()(c)
    
    assert.NoError(t, err)
    assert.Equal(t, http.StatusOK, rec.Code)
    
    // Проверка ответа
    var response models.InventoryAssetList
    err = json.Unmarshal(rec.Body.Bytes(), &response)
    assert.NoError(t, err)
    assert.Equal(t, expectedResult.TotalCount, response.TotalCount)
    assert.Len(t, *response.Items, 1)
    assert.Equal(t, expectedAsset.Name, (*response.Items)[0].Name)
    
    mockUC.AssertExpectations(t)
}
```

### 5. Интеграционные тесты

```go
// tests/integration/inventory_assets_test.go
package integration

import (
    "context"
    "testing"
    "inventory_assets/internal/inventory_asset"
    "inventory_assets/internal/inventory_asset/repository"
    "inventory_assets/internal/inventory_asset/usecase"
    "inventory_assets/internal/models"
    "inventory_assets/pkg/utils"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/suite"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
)

type InventoryAssetsIntegrationTestSuite struct {
    suite.Suite
    db     *gorm.DB
    repo   inventory_asset.Repository
    useCase inventory_asset.UseCase
}

func (suite *InventoryAssetsIntegrationTestSuite) SetupSuite() {
    db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
    assert.NoError(suite.T(), err)
    
    err = db.AutoMigrate(&models.InventoryAsset{})
    assert.NoError(suite.T(), err)
    
    suite.db = db
    suite.repo = repository.NewInventoryAssetRepository(db, nil)
    suite.useCase = usecase.NewInventoryAssetUseCase(nil, suite.repo, nil)
}

func (suite *InventoryAssetsIntegrationTestSuite) TearDownSuite() {
    sqlDB, err := suite.db.DB()
    assert.NoError(suite.T(), err)
    sqlDB.Close()
}

func (suite *InventoryAssetsIntegrationTestSuite) SetupTest() {
    // Очистка таблицы перед каждым тестом
    suite.db.Exec("DELETE FROM inventory_assets")
}

func (suite *InventoryAssetsIntegrationTestSuite) TestFullFlow() {
    ctx := context.Background()
    pq := &utils.PaginationQuery{Page: 1, Size: 10}
    
    // Создание тестового актива
    testAsset := models.InventoryAsset{
        UserName:   "john_doe",
        Name:       "Laptop Dell XPS 13",
        SerialCode: "DLXPS13001",
        IsActive:   true,
        Price:      120000,
        URL:        "https://example.com/laptop1",
    }
    
    suite.db.Create(&testAsset)
    
    // Тест получения всех активов
    result, err := suite.useCase.GetInventoryAssets(ctx, pq)
    suite.NoError(err)
    suite.NotNil(result)
    suite.Equal(int64(1), result.TotalCount)
    suite.Len(*result.Items, 1)
    
    // Тест получения актива по ID
    result, err = suite.useCase.GetInventoryAssetsById(ctx, pq, testAsset.ID)
    suite.NoError(err)
    suite.NotNil(result)
    suite.Equal(int64(1), result.TotalCount)
    suite.Len(*result.Items, 1)
    suite.Equal(testAsset.Name, (*result.Items)[0].Name)
}

func TestInventoryAssetsIntegrationTestSuite(t *testing.T) {
    suite.Run(t, new(InventoryAssetsIntegrationTestSuite))
}
```

## 🚀 Запуск тестов

### Базовые команды

```bash
# Запуск всех тестов
go test ./...

# Запуск тестов с подробным выводом
go test -v ./...

# Запуск тестов с покрытием
go test -cover ./...

# Запуск тестов с HTML отчетом о покрытии
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html

# Запуск конкретного теста
go test -v ./internal/models -run TestInventoryAsset

# Запуск тестов с таймаутом
go test -timeout 30s ./...
```

### Использование Makefile

```bash
# Запуск всех тестов
make test

# Тесты с покрытием
make test-coverage

# Все проверки качества
make check
```

## 📊 Отчеты о покрытии

### Генерация отчета

```bash
# Создание отчета о покрытии
go test -coverprofile=coverage.out ./...

# HTML отчет
go tool cover -html=coverage.out -o coverage.html

# Консольный отчет
go tool cover -func=coverage.out
```

### Анализ покрытия

```bash
# Показать покрытие по функциям
go tool cover -func=coverage.out | grep -E "(TOTAL|inventory_asset)"

# Показать покрытие по строкам
go tool cover -func=coverage.out | grep -v "TOTAL"
```

## 🔧 Настройка CI/CD

### GitHub Actions

```yaml
# .github/workflows/test.yml
name: Tests

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    
    services:
      postgres:
        image: postgres:13
        env:
          POSTGRES_PASSWORD: postgres
          POSTGRES_DB: test_db
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5
        ports:
          - 5432:5432
    
    steps:
    - uses: actions/checkout@v2
    
    - name: Set up Go
      uses: actions/setup-go@v2
      with:
        go-version: 1.24
    
    - name: Install dependencies
      run: go mod download
    
    - name: Run tests
      run: go test -v -coverprofile=coverage.out ./...
    
    - name: Upload coverage to Codecov
      uses: codecov/codecov-action@v1
      with:
        file: ./coverage.out
```

## 📝 Лучшие практики

### 1. Именование тестов

```go
func TestFunctionName_Scenario_ExpectedResult(t *testing.T) {
    // test implementation
}
```

### 2. Структура тестов

```go
func TestExample(t *testing.T) {
    // Arrange (подготовка)
    expected := "expected result"
    input := "test input"
    
    // Act (действие)
    result := functionToTest(input)
    
    // Assert (проверка)
    assert.Equal(t, expected, result)
}
```

### 3. Использование табличных тестов

```go
func TestFunction_TableDriven(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
    }{
        {"normal case", "input", "expected"},
        {"edge case", "", ""},
        {"special case", "special", "special_result"},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := functionToTest(tt.input)
            assert.Equal(t, tt.expected, result)
        })
    }
}
```

### 4. Моки и стабы

```go
// Использование интерфейсов для моков
type MockRepository struct {
    mock.Mock
}

func (m *MockRepository) GetByID(id int) (*Model, error) {
    args := m.Called(id)
    return args.Get(0).(*Model), args.Error(1)
}

// В тесте
mockRepo := new(MockRepository)
mockRepo.On("GetByID", 1).Return(&Model{ID: 1}, nil)
```

### 5. Тестовые данные

```go
// Создание тестовых данных
func createTestAsset() models.InventoryAsset {
    return models.InventoryAsset{
        UserName:   "test_user",
        Name:       "Test Asset",
        SerialCode: "TEST001",
        IsActive:   true,
        Price:      1000,
        URL:        "https://example.com/test",
    }
}
``` 