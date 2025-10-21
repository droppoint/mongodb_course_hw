# Домашнее задание №6

## Генерирование keyfile
```bash
# Запуск кластера будет осуществляться из папки hw6
cd hw6
mkdir security
openssl rand -base64 756 > security/keyfile
chmod 400 security/keyfile
```

## Создание каталогов для volume'ов
```bash
mkdir -p db/{cfg-repl0,cfg-repl1,cfg-repl2,shard0-repl0,shard0-repl1,shard0-repl2,shard1-repl0,shard1-repl1,shard1-repl2}
mkdir -p backups/pbm-backups
# chown необходим из-за особенностей образа percona-server-mongodb
sudo chown -R 1001:1 db/* security/*
```

## Установка логина и пароля для root пользователя
```bash
export MONGO_ROOT_USERNAME="<username>"
export MONGO_ROOT_PASSWORD="<password>"
export MINIO_ROOT_USER="<username>"
export MINIO_ROOT_PASSWORD="<password>"
export PBM_USERNAME="<username>"
export PBM_PASSWORD="<password>"
```
Дополнительно нужно заполнить файл pbm-config.yaml логином и паролем из переменных MINIO_ROOT_USER и MINIO_ROOT_PASSWORD

## Развертывание кластера

Развертывание осуществляется при помощи docker compose. [Конфигурацию](https://github.com/droppoint/mongodb_course_hw/blob/main/hw6/docker-compose.yml) можно найти в папке [hw6](https://github.com/droppoint/mongodb_course_hw/blob/main/hw6/) и [скрипта](https://github.com/droppoint/mongodb_course_hw/blob/main/hw6/init-cluster.sh) для настройки развернутого кластера.

```bash
docker compose up -d
./init-cluster.sh
```

## Развертывание дампа stocks
```bash
wget https://dl.dropboxusercontent.com/s/p75zp1karqg6nnn/stocks.zip
unzip -qo stocks.zip
mongorestore "mongodb://$MONGO_ROOT_USERNAME:$MONGO_ROOT_PASSWORD@127.0.0.1:27017/test?authSource=admin" dump/stocks/values.bson 
```

## Шардирование коллекции
```bash
mongosh "mongodb://$MONGO_ROOT_USERNAME:$MONGO_ROOT_PASSWORD@127.0.0.1:27017/test?authSource=admin"
```

```javascript
db.values.createIndex({stock_symbol: 1});
sh.enableSharding("test");
sh.shardCollection("test.values",{ stock_symbol: 1 });
```

## Создание backup'а
```bash
docker compose exec pbm-agent-shard0-replica0 bash
[mongodb@f1612dbdc08c /]$ pbm backup
Starting backup '2025-10-22T09:20:09Z'....Backup '2025-10-22T09:20:09Z' to remote store 'http://minio:9000/pbm-backups'

... <через какое-то время> ...

[mongodb@f1612dbdc08c /]$ pbm list
Backup snapshots:
  2025-10-22T09:20:09Z <logical> [restore_to_time: 2025-10-22T09:20:15]

PITR <on>:
  2025-10-22T09:20:15 - 2025-10-22T09:38:39
```

## Удаление контейнеров и очистка
```bash
docker compose down --remove-orphans
sudo rm -r db dump backups security stocks.zip
```