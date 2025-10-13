# Домашнее задание №4

## Генерирование keyfile
```bash
# Запуск кластера будет осуществляться из папки hw4
cd hw4
mkdir security
openssl rand -base64 1024 > security/keyfile
chmod 400 security/keyfile
```

## Создание каталогов для volume'ов
```bash
mkdir -p db/{cfg-repl0,cfg-repl1,cfg-repl2,shard0-repl0,shard0-repl1,shard0-repl2,shard1-repl0,shard1-repl1,shard1-repl2,shard2-repl0,shard2-repl1,shard2-repl2}
# chown необходим из-за особенностей образа percona-server-mongodb
sudo chown -R 1001:1 db/*
```

## Установка логина и пароля для root пользователя
```bash
export MONGO_ROOT_USERNAME="<username>"
export MONGO_ROOT_PASSWORD="<password>"
```

## Развертывание кластера

Развертывание осуществляется при помощи docker compose. [Конфигурацию](https://github.com/droppoint/mongodb_course_hw/blob/main/hw4/docker-compose.yml) можно найти в папке [hw4](https://github.com/droppoint/mongodb_course_hw/blob/main/hw4/) и [скрипта](https://github.com/droppoint/mongodb_course_hw/blob/main/hw4/init-cluster.sh) для настройки развернутого кластера.

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

## Map Reduce запрос для вычисления максимальной дневной разницы
```javascript
db.values.mapReduce(
  function() {
    emit(this.stock_symbol, this.high - this.low);
  },
  function(symbol, diffs) {
    return Math.max.apply(null, diffs);
  },
  {
    out: { inline: 1 }
  }
);
```