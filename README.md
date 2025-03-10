# Финальный проект 1 семестра (простой уровень сложности)

REST API сервис для загрузки и выгрузки данных о ценах

## Требования к системе

- Аппаратные требования: 2ГБ ОЗУ и 8 ГБ дискового пространства или выше
- Операционная система: Linux (Ubuntu 22.04 или выше)
- СУБД: PostgreSQL 15+ (предполагается, что база данных уже создана и доступна на localhost:5432)
- go 1.23.3

## Установка и запуск

1. Клонируйте репозиторий проекта:
```bash
git clone https://github.com/avsinbox/itmo-devops-sem1-project-template.git
cd itmo-devops-sem1-project-template
```

2. Запустите скрипт подготовки:

```bash
chmod +x scripts/prepare.sh
./scripts/prepare.sh
```

3. Запустите приложение с помощью скрипта:

```bash
chmod +x scripts/run.sh
./scripts/run.sh
```

## Тестирование

Запустите скрипт тестирования:

```bash
chmod +x scripts/tests.sh
./scripts/tests.sh 1
```

## Контакт

[Пишите в Telegram](https://t.me/whereisal)
