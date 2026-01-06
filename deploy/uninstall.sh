#!/bin/bash

# Скрипт удаления Antizapret Admin Panel

SERVICE_NAME="antizapret-admin.service"
BINARY_PATH="/usr/local/bin/antizapret-admin-panel"
UNINSTALL_SCRIPT_PATH="/usr/local/bin/antizapret-admin-uninstall"
WORK_DIR="/usr/local/share/antizapret-admin"

echo "Остановка сервиса..."
systemctl stop $SERVICE_NAME || true
systemctl disable $SERVICE_NAME || true

echo "Удаление unit-файла..."
rm -f /etc/systemd/system/$SERVICE_NAME
systemctl daemon-reload

echo "Удаление бинарного файла..."
rm -f $BINARY_PATH

echo "Удаление скрипта удаления..."
rm -f $UNINSTALL_SCRIPT_PATH

# Опционально: удаление рабочей директории (раскомментировать, если нужно удалять данные)
# echo "Удаление рабочих данных..."
# rm -rf $WORK_DIR

echo "Antizapret Admin Panel успешно удалена."
