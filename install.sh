#!/bin/bash

# Скрипт для установки последней версии antizapret-admin.
#
# Использование:
# curl -sSL https://raw.githubusercontent.com/USERNAME/REPO/main/install.sh | bash
#
# Замените USERNAME/REPO на ваш репозиторий на GitHub.

set -e

# --- Переменные (замените на ваши) ---
GITHUB_REPO="USERNAME/REPO"
# ---

# --- Вспомогательные функции ---
echo_info() {
    echo "[INFO] $1"
}

echo_error() {
    echo "[ERROR] $1" >&2
    exit 1
}

command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# --- Проверка зависимостей ---
if ! command_exists curl; then
    echo_error "curl не найден. Пожалуйста, установите curl и попробуйте снова."
fi
if ! command_exists tar; then
    echo_error "tar не найден. Пожалуйста, установите tar и попробуйте снова."
fi

# --- Основная логика ---
echo_info "Поиск последней версии для репозитория $GITHUB_REPO..."

# Получаем URL для скачивания последнего релиза
LATEST_RELEASE_URL=$(curl -s "https://api.github.com/repos/$GITHUB_REPO/releases/latest" | grep "browser_download_url.*linux-amd64.tar.gz" | cut -d '"' -f 4)

if [ -z "$LATEST_RELEASE_URL" ]; then
    echo_error "Не удалось найти URL для скачивания последнего релиза. Убедитесь, что репозиторий существует и содержит релизы."
fi

echo_info "Найдена версия: $LATEST_RELEASE_URL"

# Создаем временную директорию для скачивания
TMP_DIR=$(mktemp -d)
trap 'rm -rf -- "$TMP_DIR"' EXIT # Очистка временной директории при выходе

DOWNLOAD_PATH="$TMP_DIR/antizapret-admin.tar.gz"

echo_info "Скачивание в $DOWNLOAD_PATH..."
curl -L -o "$DOWNLOAD_PATH" "$LATEST_RELEASE_URL"

echo_info "Распаковка архива..."
tar -xzf "$DOWNLOAD_PATH" -C "$TMP_DIR"

INSTALL_PATH="/usr/local/bin"
BINARY_NAME="antizapret-admin"

echo_info "Установка бинарного файла в $INSTALL_PATH..."
if [ -w "$INSTALL_PATH" ]; then
    # Есть права на запись
    mv "$TMP_DIR/$BINARY_NAME" "$INSTALL_PATH/$BINARY_NAME"
else
    # Нет прав, используем sudo
    echo_info "Требуются права суперпользователя для перемещения файла в $INSTALL_PATH."
    sudo mv "$TMP_DIR/$BINARY_NAME" "$INSTALL_PATH/$BINARY_NAME"
fi

echo_info "---"
echo_info "Установка успешно завершена!"
echo_info ""

# Запрос пароля
read -p "Пожалуйста, введите пароль для ADMIN_PASSWORD: " ADMIN_PASSWORD_INPUT
echo_info ""

# Запуск приложения
PUBLIC_IP=$(curl -s ifconfig.me || echo "<IP-адрес-вашего-сервера>")

echo_info "Запуск приложения antizapret-admin..."
echo_info "Панель будет доступна по адресу http://$PUBLIC_IP:8080"
echo_info "Нажмите Ctrl+C для остановки."
echo "---"

ADMIN_PASSWORD="$ADMIN_PASSWORD_INPUT" $INSTALL_PATH/$BINARY_NAME
