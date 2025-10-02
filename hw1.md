# Домашнее задание №1

## Развертывание MongoDB

```bash
docker network create mongo-net
docker run --name mongo-courses-server --network mongo-net -d -p 27017:27017 -v /home/droppoint/mongo:/data/db -e MONGO_INITDB_ROOT_USERNAME="testroot" -e MONGO_INITDB_ROOT_PASSWORD="***" mongo:8.0.14
docker run -it --rm --name mongo-courses-client --network mongo-net mongo:8.0.14 mongosh "mongodb://testroot:***@mongo-courses-server:27017"
```

## Создание коллекции со случайным количеством элементов

```javascript
test> for ( i = 0; i < Math.random()*100; ++i ) {
...     db.docs.insertOne( 
...         {
...             documentNumber: i+1,
...             text: `This is document number: ${i+1}` 
...         } 
...     );
... };
{
  acknowledged: true,
  insertedId: ObjectId('68de72adcd5ca76825ce5f50')
}
```

## Подсчет количества документов в коллекции
```javascript
test> db.docs.countDocuments({});
10
```

## Удаление контейнера
```bash
test> exit

$ docker stop mongo-courses-server
mongo-courses-server
$ docker rm mongo-courses-server
mongo-courses-server
```