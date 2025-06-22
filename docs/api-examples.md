# API Examples

Примеры использования Inventory Assets API.

## 🔧 Настройка

### Базовый URL
```
http://localhost:8080/api/v1
```

### Заголовки
```http
Content-Type: application/json
Accept: application/json
```

## 📋 Примеры запросов

### 1. Получить все инвентарные активы

#### Базовый запрос
```bash
curl -X GET "http://localhost:8080/api/v1/inventory-assets" \
  -H "Accept: application/json"
```

#### С пагинацией
```bash
curl -X GET "http://localhost:8080/api/v1/inventory-assets?page=1&size=5" \
  -H "Accept: application/json"
```

#### С сортировкой
```bash
curl -X GET "http://localhost:8080/api/v1/inventory-assets?page=1&size=10&orderBy=name" \
  -H "Accept: application/json"
```

**Ответ:**
```json
{
  "items": [
    {
      "id": 1,
      "username": "john_doe",
      "name": "Laptop Dell XPS 13",
      "serial_code": "DLXPS13001",
      "is_active": true,
      "price": 120000,
      "url": "https://example.com/laptop1"
    },
    {
      "id": 2,
      "username": "jane_smith",
      "name": "Monitor Samsung 27\"",
      "serial_code": "SMS27001",
      "is_active": true,
      "price": 45000,
      "url": "https://example.com/monitor1"
    }
  ],
  "total": 25,
  "page": 1,
  "size": 10,
  "pages": 3
}
```

### 2. Получить инвентарный актив по ID

```bash
curl -X GET "http://localhost:8080/api/v1/inventory-assets/1" \
  -H "Accept: application/json"
```

**Ответ:**
```json
{
  "items": [
    {
      "id": 1,
      "username": "john_doe",
      "name": "Laptop Dell XPS 13",
      "serial_code": "DLXPS13001",
      "is_active": true,
      "price": 120000,
      "url": "https://example.com/laptop1"
    }
  ],
  "total": 1,
  "page": 1,
  "size": 10,
  "pages": 1
}
```

### 3. Обработка ошибок

#### Неверный ID
```bash
curl -X GET "http://localhost:8080/api/v1/inventory-assets/999" \
  -H "Accept: application/json"
```

**Ответ:**
```json
{
  "error": "inventory asset with this id not found",
  "status": 404
}
```

#### Неверные параметры пагинации
```bash
curl -X GET "http://localhost:8080/api/v1/inventory-assets?page=invalid&size=abc" \
  -H "Accept: application/json"
```

**Ответ:**
```json
{
  "error": "invalid pagination parameters",
  "status": 400
}
```

## 🐍 Python примеры

### Использование requests

```python
import requests
import json

# Базовый URL
BASE_URL = "http://localhost:8080/api/v1"

# Получить все активы
def get_all_assets(page=1, size=10):
    url = f"{BASE_URL}/inventory-assets"
    params = {"page": page, "size": size}
    
    response = requests.get(url, params=params)
    response.raise_for_status()
    
    return response.json()

# Получить актив по ID
def get_asset_by_id(asset_id):
    url = f"{BASE_URL}/inventory-assets/{asset_id}"
    
    response = requests.get(url)
    response.raise_for_status()
    
    return response.json()

# Пример использования
try:
    # Получить первую страницу активов
    assets = get_all_assets(page=1, size=5)
    print(f"Найдено активов: {assets['total']}")
    
    for asset in assets['items']:
        print(f"- {asset['name']} (ID: {asset['id']})")
    
    # Получить конкретный актив
    asset = get_asset_by_id(1)
    print(f"\nДетали актива: {asset['items'][0]['name']}")
    
except requests.exceptions.RequestException as e:
    print(f"Ошибка запроса: {e}")
```

## 🔄 JavaScript примеры

### Использование fetch

```javascript
// Базовый URL
const BASE_URL = 'http://localhost:8080/api/v1';

// Получить все активы
async function getAllAssets(page = 1, size = 10) {
    try {
        const url = `${BASE_URL}/inventory-assets?page=${page}&size=${size}`;
        const response = await fetch(url);
        
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        
        return await response.json();
    } catch (error) {
        console.error('Ошибка при получении активов:', error);
        throw error;
    }
}

// Получить актив по ID
async function getAssetById(assetId) {
    try {
        const url = `${BASE_URL}/inventory-assets/${assetId}`;
        const response = await fetch(url);
        
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        
        return await response.json();
    } catch (error) {
        console.error(`Ошибка при получении актива ${assetId}:`, error);
        throw error;
    }
}

// Пример использования
async function main() {
    try {
        // Получить первую страницу активов
        const assets = await getAllAssets(1, 5);
        console.log(`Найдено активов: ${assets.total}`);
        
        assets.items.forEach(asset => {
            console.log(`- ${asset.name} (ID: ${asset.id})`);
        });
        
        // Получить конкретный актив
        const asset = await getAssetById(1);
        console.log(`\nДетали актива: ${asset.items[0].name}`);
        
    } catch (error) {
        console.error('Ошибка:', error);
    }
}

// Запуск
main();
```

## 📊 Postman коллекция

### Импорт в Postman

1. Создайте новую коллекцию "Inventory Assets API"
2. Добавьте переменную окружения `base_url` со значением `http://localhost:8080/api/v1`

### Запросы

#### 1. Get All Assets
```
GET {{base_url}}/inventory-assets?page=1&size=10
```

#### 2. Get Asset by ID
```
GET {{base_url}}/inventory-assets/1
```

## 🔍 Тестирование с помощью curl

### Скрипт для тестирования

```bash
#!/bin/bash

BASE_URL="http://localhost:8080/api/v1"

echo "🧪 Тестирование Inventory Assets API"
echo "=================================="

# Тест 1: Получить все активы
echo "1. Получение всех активов..."
response=$(curl -s -w "\n%{http_code}" "${BASE_URL}/inventory-assets")
http_code=$(echo "$response" | tail -n1)
body=$(echo "$response" | head -n -1)

if [ "$http_code" -eq 200 ]; then
    echo "✅ Успешно получены активы"
    echo "📊 Ответ: $body" | jq '.'
else
    echo "❌ Ошибка: HTTP $http_code"
    echo "📄 Ответ: $body"
fi

echo ""

# Тест 2: Получить актив по ID
echo "2. Получение актива по ID 1..."
response=$(curl -s -w "\n%{http_code}" "${BASE_URL}/inventory-assets/1")
http_code=$(echo "$response" | tail -n1)
body=$(echo "$response" | head -n -1)

if [ "$http_code" -eq 200 ]; then
    echo "✅ Успешно получен актив"
    echo "📊 Ответ: $body" | jq '.'
else
    echo "❌ Ошибка: HTTP $http_code"
    echo "📄 Ответ: $body"
fi

echo ""

# Тест 3: Неверный ID
echo "3. Тест с неверным ID..."
response=$(curl -s -w "\n%{http_code}" "${BASE_URL}/inventory-assets/999")
http_code=$(echo "$response" | tail -n1)
body=$(echo "$response" | head -n -1)

if [ "$http_code" -eq 404 ]; then
    echo "✅ Правильно обработана ошибка 404"
    echo "📊 Ответ: $body" | jq '.'
else
    echo "❌ Неожиданный код ответа: HTTP $http_code"
    echo "📄 Ответ: $body"
fi

echo ""
echo "🏁 Тестирование завершено"
```

## 📈 Мониторинг и логи

### Проверка логов

```bash
# Просмотр логов приложения
docker logs inventory-assets-api

# Фильтрация по уровню
docker logs inventory-assets-api | grep "ERROR"

# Мониторинг в реальном времени
docker logs -f inventory-assets-api
```

### Метрики

```bash
# Проверка состояния сервера
curl -X GET "http://localhost:8080/health"

# Статистика базы данных (если доступна)
curl -X GET "http://localhost:8080/metrics"
```

## 🔧 Отладка

### Включение подробного логирования

```bash
# Установка переменной окружения для отладки
export LOGGER_LEVEL=debug

# Перезапуск приложения
docker-compose restart api
```

### Проверка подключения к базе данных

```bash
# Подключение к PostgreSQL
docker exec -it postgres psql -U postgres -d inventory_assets

# Проверка таблиц
\dt

# Проверка данных
SELECT * FROM inventory_assets LIMIT 5;
``` 