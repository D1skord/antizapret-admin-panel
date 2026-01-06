#!/bin/bash

# Установщик Antizapret Admin Panel
#
# Делает следующее:
# 1. Скачивает последний релиз (tar.gz).
# 2. Запрашивает пароль и сохраняет его в systemd.
# 3. Останавливает службу (если уже установлена).
# 4. Распаковывает и обновляет бинарник и service-файл.
# 5. Перезапускает службу.

set -e

REPO="d1skord/antizapret-admin-panel"
SERVICE_NAME="antizapret-admin.service"
BINARY_NAME="antizapret-admin-panel"
INSTALL_DIR="/usr/local/bin"
SYSTEMD_DIR="/etc/systemd/system"
WORK_DIR="/usr/local/share/antizapret-admin"

echo_info() { echo -e "\033[34m[INFO]\033[0m $1"; }
echo_error() { echo -e "\033[31m[ERROR]\033[0m $1"; exit 1; }

# Проверка root прав
if [ "$EUID" -ne 0 ]; then
  echo_error "Пожалуйста, запустите скрипт с правами root (sudo)."
fi

# 1. Поиск URL последнего релиза
echo_info "Поиск последней версии..."
RELEASE_URL=$(curl -s "https://api.github.com/repos/$REPO/releases/latest" | grep "browser_download_url.*linux-amd64.tar.gz" | cut -d '"' -f 4)

if [ -z "$RELEASE_URL" ]; then
    echo_error "Не удалось найти ссылку на релиз (linux-amd64.tar.gz)."
fi

# 2. Настройка учетных данных
SERVICE_OVERRIDE_DIR="$SYSTEMD_DIR/$SERVICE_NAME.d"
OVERRIDE_FILE="$SERVICE_OVERRIDE_DIR/env.conf"

# --- Загрузка существующих значений ---
EXISTING_USERNAME=""
EXISTING_PASSWORD=""
if [ -f "$OVERRIDE_FILE" ]; then
    # Используем grep, чтобы найти строку, и cut, чтобы получить значение
    # Удаляем кавычки, которые могут быть вокруг значения
    EXISTING_USERNAME=$(grep 'ADMIN_USERNAME' "$OVERRIDE_FILE" | sed 's/.*ADMIN_USERNAME=//' | tr -d '"')
    EXISTING_PASSWORD=$(grep 'ADMIN_PASSWORD' "$OVERRIDE_FILE" | sed 's/.*ADMIN_PASSWORD=//' | tr -d '"')
    echo_info "Обнаружена существующая конфигурация."
fi

# --- Интерактивный ввод для Логина ---
if [ -n "$ADMIN_USERNAME" ]; then
    ADMIN_USERNAME_INPUT="$ADMIN_USERNAME"
    echo_info "Используется логин из переменной окружения ADMIN_USERNAME."
else
    echo_info "Настройка логина администратора (по умолчанию: admin)."
    if [ -n "$EXISTING_USERNAME" ]; then
        echo -n "Введите новый логин (или оставьте пустым, чтобы использовать '$EXISTING_USERNAME'): "
    else
        echo -n "Введите логин (или оставьте пустым, чтобы использовать 'admin'): "
    fi
    
    # Пытаемся читать с tty
    set +e
    if { read -r ADMIN_USERNAME_INPUT < /dev/tty; } 2>/dev/null; then
        echo ""
    else
        echo "" # Перенос строки после промпта
        echo_info "Интерактивный ввод недоступен для логина. Используется существующий или значение по умолчанию."
        ADMIN_USERNAME_INPUT=""
    fi
    set -e
fi

# --- Интерактивный ввод для Пароля ---
if [ -n "$ADMIN_PASSWORD" ]; then
    ADMIN_PASSWORD_INPUT="$ADMIN_PASSWORD"
    echo_info "Используется пароль из переменной окружения ADMIN_PASSWORD."
else
    echo_info "Настройка пароля администратора."
    if [ -n "$EXISTING_PASSWORD" ]; then
        echo -n "Введите новый пароль (или оставьте пустым, чтобы использовать существующий): "
    else
        echo -n "Введите новый пароль: "
    fi

    set +e
    if { read -s ADMIN_PASSWORD_INPUT < /dev/tty; } 2>/dev/null; then
        echo ""
    else
        echo "" # Перенос строки
        if [ ! -f "$OVERRIDE_FILE" ]; then
            echo_error "Ошибка: Не удалось запросить пароль (интерактивный режим недоступен) и ADMIN_PASSWORD не установлена."
        fi
        echo_info "Интерактивный ввод недоступен для пароля. Используется существующий пароль."
        ADMIN_PASSWORD_INPUT=""
    fi
    set -e
fi

# --- Определение финальных значений и запись ---
FINAL_USERNAME=$ADMIN_USERNAME_INPUT
if [ -z "$FINAL_USERNAME" ]; then
    FINAL_USERNAME=$EXISTING_USERNAME
    if [ -z "$FINAL_USERNAME" ]; then
        FINAL_USERNAME="admin"
    fi
fi

FINAL_PASSWORD=$ADMIN_PASSWORD_INPUT
if [ -z "$FINAL_PASSWORD" ]; then
    FINAL_PASSWORD=$EXISTING_PASSWORD
fi

# Проверка, что пароль не пустой при первой установке
if [ -z "$FINAL_PASSWORD" ]; then
    echo_error "Пароль не может быть пустым. Установите его через интерактивный ввод или переменную ADMIN_PASSWORD."
fi

# Создаем директорию и записываем обе переменные
mkdir -p "$SERVICE_OVERRIDE_DIR"
cat > "$OVERRIDE_FILE" << EOF
[Service]
Environment="ADMIN_USERNAME=$FINAL_USERNAME"
Environment="ADMIN_PASSWORD=$FINAL_PASSWORD"
Environment="OPENVPN_CLIENTS_PATH=/root/antizapret/client/openvpn/vpn-udp/"
Environment="OPENVPN_ANTIZAPRET_PATH=/root/antizapret/client/openvpn/antizapret-udp/"
Environment="CLIENT_SCRIPT_PATH=/root/antizapret/client.sh"
EOF

echo_info "Учетные данные сохранены в конфигурации systemd."

# 3. Скачивание
TMP_DIR=$(mktemp -d)
trap 'rm -rf -- "$TMP_DIR"' EXIT
ARCHIVE_PATH="$TMP_DIR/release.tar.gz"

echo_info "Скачивание: $RELEASE_URL"
curl -L -o "$ARCHIVE_PATH" "$RELEASE_URL"

echo_info "Распаковка..."
tar -xzf "$ARCHIVE_PATH" -C "$TMP_DIR"

# Предполагаем структуру архива:
# /antizapret-admin-panel (бинарник)
# /antizapret-admin.service
# /uninstall.sh

# 4. Остановка текущей службы (обновление)
if systemctl is-active --quiet $SERVICE_NAME; then
    echo_info "Остановка текущей службы..."
    systemctl stop $SERVICE_NAME
fi

# 5. Установка файлов
echo_info "Установка бинарника в $INSTALL_DIR..."
# Ищем бинарник в распакованной папке (он может быть в корне архива или в подпапке)
# Мы ищем файл с именем antizapret-admin
FOUND_BINARY=$(find "$TMP_DIR" -type f -name "$BINARY_NAME" | head -n 1)
if [ -z "$FOUND_BINARY" ]; then
    echo_error "Бинарный файл не найден в архиве."
fi
cp "$FOUND_BINARY" "$INSTALL_DIR/$BINARY_NAME"
chmod +x "$INSTALL_DIR/$BINARY_NAME"

echo_info "Настройка systemd сервиса..."
FOUND_SERVICE=$(find "$TMP_DIR" -type f -name "$SERVICE_NAME" | head -n 1)
if [ -z "$FOUND_SERVICE" ]; then
    echo_error "Файл сервиса не найден в архиве."
fi
cp "$FOUND_SERVICE" "$SYSTEMD_DIR/$SERVICE_NAME"

# Установка скрипта удаления
FOUND_UNINSTALL=$(find "$TMP_DIR" -type f -name "uninstall.sh" | head -n 1)
if [ -n "$FOUND_UNINSTALL" ]; then
    cp "$FOUND_UNINSTALL" "$INSTALL_DIR/antizapret-admin-uninstall"
    chmod +x "$INSTALL_DIR/antizapret-admin-uninstall"
    echo_info "Скрипт удаления установлен: $INSTALL_DIR/antizapret-admin-uninstall"
fi

# Создание рабочей директории (если нет)
mkdir -p "$WORK_DIR"

# 6. Запуск
echo_info "Перезагрузка демона systemd..."
systemctl daemon-reload
systemctl enable $SERVICE_NAME
systemctl start $SERVICE_NAME

PUBLIC_IP=$(curl -s ifconfig.me || echo "IP-ВАШЕГО-СЕРВЕРА")

echo_info "---"
echo_info "Установка/Обновление завершено!"
echo_info "Панель доступна: http://$PUBLIC_IP:8080"
echo_info "Логин администратора: $FINAL_USERNAME"
echo_info "Пароль администратора: $FINAL_PASSWORD"
echo_info "Пожалуйста, сохраните этот пароль в безопасном месте. Он больше не будет показан."
echo_info "Статус службы: systemctl status $SERVICE_NAME"
