# ow-exporter Implementation Plan

## ✅ Completed Analysis

### Platform Structure Discovered
- **PC (mouseKeyboard-view)**: клавиатура + мышь
- **Console (controller-view)**: геймпад
- Переключение через табы: `#mouseKeyboardFilter` / `#controllerFilter`

### Game Modes Found
- **Quick Play**: `.quickPlay-view.is-active`
- **Competitive**: `.competitive-view.is-active`

### Metrics Identified
15 основных метрик для всех героев:
- Time Played, Games Won, Win %, Weapon Accuracy
- Eliminations per Life, Kill Streak, Multikill
- Avg per 10min: Eliminations, Deaths, Final Blows, Solo Kills, Objective Kills, Objective Time, Hero Damage, Healing

### Prometheus Label Structure
```
ow_hero_time_played_seconds{
    username="player1",
    hero="mercy",
    platform="pc",        # "pc" | "console"
    gamemode="competitive" # "quickplay" | "competitive"
}
```

## 📁 Files Created in tmp/
- `profile_de5bb4aca17492e0.html` - первый профиль
- `profile_c14aad9eba729abe.html` - второй профиль
- `platform_analysis.md` - анализ структуры платформ
- `metrics_mapping.go` - хардкод мапинг метрик
- `parser_example.go` - пример парсера с поддержкой платформ

## 🎯 Next Implementation Steps

### 1. Project Structure
```
cmd/ow-exporter/
├── main.go                    # HTTP server + CLI
├── api/
│   ├── handlers.go           # REST API для пользователей
│   └── routes.go             # Маршруты API
├── scraper/
│   ├── client.go             # HTTP client + rate limiting
│   ├── parser.go             # HTML парсинг (из tmp/parser_example.go)
│   └── scheduler.go          # Фоновое обновление каждые 5 минут
├── storage/
│   ├── interface.go          # Storage abstraction
│   ├── sqlite.go             # SQLite для пользователей
│   └── memory.go             # In-memory cache для статистики
├── metrics/
│   └── prometheus.go         # Prometheus exporter
├── models/
│   ├── user.go               # User model
│   ├── stats.go              # Stats models (из tmp/metrics_mapping.go)
│   └── hero_metrics.go       # Hardcoded metrics
└── config/
    └── config.go             # Configuration
```

### 2. Core Features
- ✅ **HTTP парсинг** (без headless браузера)
- ✅ **Platform support** (PC/Console разделение)
- ✅ **Game mode support** (QuickPlay/Competitive)
- ✅ **Hardcoded metrics** (15 основных метрик)
- 🔄 **REST API** для управления пользователями
- 🔄 **SQLite storage** для персистентности
- 🔄 **Prometheus exporter** с полными лейблами
- 🔄 **Background scheduler** для автообновления

### 3. API Endpoints
```
GET    /api/users              # Список пользователей
POST   /api/users              # Добавить пользователя
PUT    /api/users/{username}   # Обновить пользователя
DELETE /api/users/{username}   # Удалить пользователя
GET    /metrics                # Prometheus метрики
GET    /health                 # Health check
```

### 4. Configuration
```yaml
database:
  type: "sqlite"
  path: "/data/ow-exporter.db"
scraper:
  interval: "5m"
  timeout: "30s"
  user_agent: "ow-exporter/1.0"
server:
  port: 9420
  enable_api: true
```

### 5. Example Usage
```bash
# Добавить пользователя
curl -X POST /api/users \
  -d '{"username": "player1", "profile_id": "de5bb4aca17492e0bba120a1d1%7Ca92a11ef8d304356fccfff8df12e1dc6"}'

# Получить метрики
curl /metrics
```

## 🚀 Ready for Implementation

Все необходимые компоненты спланированы и подготовлены:
- ✅ Структура платформ понята
- ✅ Метрики замаплены
- ✅ Парсер спроектирован
- ✅ API спланировано
- ✅ Архитектура готова

Можно начинать создание `cmd/ow-exporter/` с реального кода!