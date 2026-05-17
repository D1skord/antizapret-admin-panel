---
name: release
description: Выкатить новый релиз проекта antizapret-admin-panel. Используй этот скилл всегда, когда пользователь говорит "сделай релиз", "выкатить релиз", "новый релиз", "make a release", "create release", "bump version", "выпусти версию" или что-то похожее. Скилл создаёт git-тег, который автоматически запускает GitHub Actions для сборки и публикации релиза.
---

# Release

Выпуск новой версии antizapret-admin-panel.

## Как устроен процесс

Сборка и публикация **полностью автоматические** — их выполняет GitHub Actions при появлении нового тега формата `v*`. Файл workflow: `.github/workflows/release.yml`.

Что делает Actions автоматически:
- Собирает frontend (`npm run build`)
- Компилирует Go-бинарник под linux/amd64
- Упаковывает архив `antizapret-admin-linux-amd64.tar.gz`
- Публикует релиз на GitHub с архивом и `deploy/install.sh`

Твоя задача: понять что изменилось → выбрать версию → создать и запушить тег.

## Шаги

### 1. Посмотри что изменилось с прошлого релиза

```bash
git fetch --tags
LAST_TAG=$(git describe --tags --abbrev=0)
echo "Последний тег: $LAST_TAG"
git log --oneline $LAST_TAG..HEAD
```

Покажи пользователю список изменений кратко.

### 2. Предложи версию (semver)

Ориентируйся на коммиты:
- **patch** (`1.1.0 → 1.1.1`) — только баг-фиксы
- **minor** (`1.1.0 → 1.2.0`) — новые фичи, обратно совместимые
- **major** (`1.1.0 → 2.0.0`) — breaking changes

Предложи версию сам, но спроси подтверждение у пользователя.

### 3. Убедись что main актуален

```bash
git status
git log origin/main..HEAD
```

Если есть незапушенные коммиты — запуши сначала.

### 4. Создай тег и запуши

```bash
git tag v<version>
git push origin v<version>
```

### 5. Сообщи результат

- Тег создан, GitHub Actions запущен автоматически
- Следить за сборкой: https://github.com/D1skord/antizapret-admin-panel/actions
- Готовый релиз появится здесь: https://github.com/D1skord/antizapret-admin-panel/releases

## Важно

- Не собирай бинарник вручную — это делает GitHub Actions
- Не добавляй строки `Co-Authored-By` в коммиты
- Заметки к релизу создай отдельно через `gh release edit` если нужно
