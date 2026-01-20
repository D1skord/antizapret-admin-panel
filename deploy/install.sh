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

SERVICE_OVERRIDE_DIR="$SYSTEMD_DIR/$SERVICE_NAME.d"
OVERRIDE_FILE="$SERVICE_OVERRIDE_DIR/env.conf"

# 2. Настройка пароля
if [ -n "$ADMIN_PASSWORD" ]; then
    ADMIN_PASSWORD_INPUT="$ADMIN_PASSWORD"
    echo_info "Используется пароль из переменной окружения ADMIN_PASSWORD."
else
    echo_info "Настройка пароля администратора."
    echo_info "При обновлении, вы можете оставить поле пустым, чтобы использовать существующий пароль."
    
    # Пытаемся считать пароль с tty
    # Используем echo для промпта, чтобы перенаправить stderr read в /dev/null (скрыть ошибки ввода-вывода)
    echo -n "Введите новый пароль для панели администрирования: "
    
    set +e
    if { read -s ADMIN_PASSWORD_INPUT < /dev/tty; } 2>/dev/null; then
        echo "" # Перенос строки после ввода
    else
        echo "" # Перенос строки
        # Чтение не удалось (нет TTY или ошибка).
        if [ -f "$OVERRIDE_FILE" ]; then
            echo_info "Интерактивный ввод недоступен. Используется существующий пароль."
            ADMIN_PASSWORD_INPUT=""
        else
            echo_error "Ошибка: Не удалось запросить пароль (интерактивный режим недоступен) и переменная ADMIN_PASSWORD не установлена."
            echo_error "Для первой установки используйте: curl ... | sudo ADMIN_PASSWORD='ваш_пароль' bash"
            exit 1
        fi
    fi
    set -e
fi

if [ -n "$ADMIN_PASSWORD_INPUT" ]; then
    mkdir -p "$SERVICE_OVERRIDE_DIR"
    # Если введен новый пароль, сохраняем его
    cat > "$OVERRIDE_FILE" << EOF
[Service]
Environment="ADMIN_PASSWORD=$ADMIN_PASSWORD_INPUT"
EOF
    echo_info "Новый пароль сохранен в конфигурации systemd."
elif [ ! -f "$OVERRIDE_FILE" ]; then
    # Если пароль не введен и файла конфигурации нет (первая установка)
    echo_error "Пароль не может быть пустым при первой установке."
else
    echo_info "Пароль не менялся. Используется существующий."
fi


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
echo_info "Пароль администратора: [сохранен в конфигурации сервиса]"
echo_info "Статус службы: systemctl status $SERVICE_NAME"
