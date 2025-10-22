# Домашнее задание №7

## Генерирование keyfile
```bash
# Запуск кластера будет осуществляться из папки hw6
cd hw7
mkdir security
openssl rand -base64 756 > security/keyfile
chmod 400 security/keyfile
```

## Создание каталогов для volume'ов
```bash
mkdir -p db/{repl0,repl1,repl2}
# chown необходим из-за особенностей образа percona-server-mongodb
sudo chown -R 1001:1 db/* security/*
```

## Установка логина и пароля для root пользователя
```bash
export MONGO_INITDB_ROOT_USERNAME="<username>"
export MONGO_INITDB_ROOT_PASSWORD="<password>"
```

## Развертывание replicaset и PMM

Развертывание осуществляется при помощи docker compose. [Конфигурацию](https://github.com/droppoint/mongodb_course_hw/blob/main/hw7/docker-compose.yml) можно найти в папке [hw7](https://github.com/droppoint/mongodb_course_hw/blob/main/hw7/) и [скрипта](https://github.com/droppoint/mongodb_course_hw/blob/main/hw7/init.sh) для настройки развернутого replicaset'а.

```bash
docker compose up -d
./init.sh
```

## Развертывание дампа stocks
```bash
wget https://dl.dropboxusercontent.com/s/p75zp1karqg6nnn/stocks.zip
unzip -qo stocks.zip
mongorestore "mongodb://$MONGO_INITDB_ROOT_USERNAME:$MONGO_INITDB_ROOT_PASSWORD@127.0.0.1:27017/test?authSource=admin" dump/stocks/values.bson
```

## Нагрузка БД
Нагрузка на БД осуществляется через работу программы loadgenerator из контейнера mongo-load. Код можно найти [здесь](https://github.com/droppoint/mongodb_course_hw/blob/main/hw7/loadgenerator/main.go).

```bash
docker compose run --rm mongo-load
```
![Percona Monitoring and Management](hw7/pmm.png?raw=true "Результат нагрузки")


## Удаление контейнеров и очистка
```bash
docker compose down --remove-orphans
sudo rm -r db dump security stocks.zip
```