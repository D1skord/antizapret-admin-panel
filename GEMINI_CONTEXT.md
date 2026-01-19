# Контекст проекта Antizapret Admin Panel

Этот файл предназначен для AI-агентов, чтобы быстро погрузиться в суть проекта.

## Обзор
Проект представляет собой админ-панель для управления VPN-сервером (OpenVPN/WireGuard/AmneziaWG) в контексте решения AntiZapret.
Архитектура: Монолит. Go (Backend) + Vue 3 (Frontend).
Фронтенд собирается в статику и вшивается в Go-бинарник с помощью `//go:embed`.

## Структура директорий
- `/frontend`: Исходный код Vue 3 приложения. Использует Vite.
- `/internal`: Go код.
  - `/api`: HTTP хендлеры (Gin).
  - `/entity`: Структуры данных.
  - `/repository`: Работа с файловой системой (чтение/запись конфигов VPN).
  - `/service`: Бизнес-логика.
- `/mock_fs`: Эмуляция файловой системы сервера для локальной разработки (ключи, конфиги).
- `/deploy`: Скрипты и конфиги для Linux-деплоя (systemd service, uninstall).
- `.github/workflows`: CI/CD для сборки релизов.

## Процесс сборки и релиза
**Важно:** Проект распространяется как `.tar.gz` архив через GitHub Releases.

1. **GitHub Actions (`release.yml`)**:
   - При пуше тега `v*` запускается сборка.
   - Собирается фронтенд (`npm run build`).
   - Собирается бэкенд для Linux (`GOOS=linux GOARCH=amd64`).
   - Формируется архив со структурой:
     - `antizapret-admin` (бинарник)
     - `deploy/antizapret-admin.service` (systemd unit)
     - `deploy/uninstall.sh` (скрипт удаления)
     - `README.md`
   - Архив загружается в Release assets.

2. **Установка (`install.sh`)**:
   - Пользователь качает скрипт: `curl ... | sudo bash`.
   - Скрипт парсит последний релиз с GitHub API.
   - Качает архив, распаковывает.
   - Останавливает сервис `antizapret-admin.service` (если есть).
   - Обновляет бинарник в `/usr/local/bin`.
   - Обновляет Unit-файл в `/etc/systemd/system`.
   - Кладёт `uninstall.sh` в `/usr/local/bin/antizapret-admin-uninstall`.
   - Запускает сервис.

## Особенности разработки
- **Локальный запуск**: `make dev-backend` (использует `air` для хот-релоада) + `make dev-frontend`.
- **Переменные окружения**:
  - `OPENVPN_CLIENTS_PATH`: Путь к папке с клиентскими конфигами.
  - `CLIENT_SCRIPT_PATH`: Путь к баш-скрипту управления клиентами.
  - В локальной разработке они указывают на `mock_fs/`.

## Текущее состояние (2025-01-17)
- Реализована полноценная схема авто-обновления через install.sh.
- Добавлен systemd сервис.
- Добавлен скрипт удаления.
- Makefile поддерживает `build-release`.
