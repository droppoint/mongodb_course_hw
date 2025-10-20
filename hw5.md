# Домашнее задание №5

## Развертывание MongoDB

```bash
docker network create mongo-net
docker run --name mongo-courses-server --network mongo-net -d -p 27017:27017 -v /home/droppoint/mongo:/data/db -e MONGO_INITDB_ROOT_USERNAME="testroot" -e MONGO_INITDB_ROOT_PASSWORD="***" mongo:8.0.14
docker run -it --rm --name mongo-courses-client --network mongo-net mongo:8.0.14 mongosh "mongodb://testroot:***@mongo-courses-server:27017"
```

В качестве коллекции для тестов был выбран[справочник покемонов](https://raw.githubusercontent.com/Biuni/PokemonGO-Pokedex/master/pokedex.json), также известный как Pokedex.

## Создание коллекции со схемой валидации 

*Схема валидации*
```javascript
{
  $jsonSchema: {
    bsonType: "object",
    required: [
      "id",
      "num",
      "name",
      "img",
      "type",
      "height",
      "weight",
      "candy",
      "egg",
      "spawn_chance",
      "avg_spawns",
      "spawn_time",
      "weaknesses"
    ],
    properties: {
      id: { bsonType: "int" },
      num: {
        bsonType: "string",
        pattern: "^[0-9]{3}$"
      },
      name: {
        bsonType: "string",
        minLength: 1
      },
      img: {
        bsonType: "string",
        pattern: "^https?://"
      },
      type: {
        bsonType: "array",
        minItems: 1,
        items: { bsonType: "string" }
      },
      height: { bsonType: "string" },
      weight: { bsonType: "string" },
      candy: { bsonType: "string" },
      candy_count: { bsonType: "int" },
      egg: { bsonType: "string" },
      spawn_chance: {
        bsonType: ["double", "int"],
        minimum: 0
      },
      avg_spawns: {
        bsonType: ["double", "int"],
        minimum: 0
      },
      spawn_time: { bsonType: "string" },
      multipliers: {
        bsonType: ["array", "null"],
        items: { bsonType: "double" }
      },
      weaknesses: {
        bsonType: "array",
        minItems: 1,
        items: { bsonType: "string" }
      },
      next_evolution: {
        bsonType: "array",
        items: {
          bsonType: "object",
          required: ["num", "name"],
          properties: {
            num: { bsonType: "string" },
            name: { bsonType: "string" }
          }
        }
      },
      prev_evolution: {
        bsonType: "array",
        items: {
          bsonType: "object",
          required: ["num", "name"],
          properties: {
            num: { bsonType: "string" },
            name: { bsonType: "string" }
          }
        }
      }
    }
  }
}
```

*Создание коллекции*
```javascript
db.createCollection("pokedex", {
  validator: {
    $jsonSchema: {
      <схема валидации>
    }
  },
  validationLevel: "strict",
  validationAction: "error"
});
```

## Проверка работы схемы валидации

*Невалидный документ*
```javascript
test> db.pokedex.insertOne({ name: "InvalidMon" })

Uncaught:
MongoServerError: Document failed validation
Additional information: {
  failingDocumentId: ObjectId('68f6071f5d96f4184ece5f47'),
  details: {
    operatorName: '$jsonSchema',
    schemaRulesNotSatisfied: [
      {
        operatorName: 'required',
        specifiedAs: {
          required: [
            'id',         'num',
            'name',       'img',
            'type',       'height',
            'weight',     'candy',
            'egg',        'spawn_chance',
            'avg_spawns', 'spawn_time',
            'weaknesses'
          ]
        },
        missingProperties: [
          'avg_spawns', 'candy',
          'egg',        'height',
          'id',         'img',
          'num',        'spawn_chance',
          'spawn_time', 'type',
          'weaknesses', 'weight'
        ]
      }
    ]
  }
}
```

*Вставка коллекции из справочника pokedex*

Так как нам нужен на вход массив, нужно убрать поле "pokemon" из json документа
```bash
curl -O https://raw.githubusercontent.com/Biuni/PokemonGO-Pokedex/master/pokedex.json
cat pokedex.json | jq '.pokemon' > pokemon.json
```

Загрузка исправленного документа в коллекцию
```bash
docker run -it --rm --name mongo-courses-importer -v ./pokemon.json:/pokemon.json --network mongo-net mongo:8.0.14 mongoimport -vvvvv --uri "mongodb://testroot:***@mongo-courses-server:27017/test" --collection pokedex --file /pokemon.json --jsonArray --stopOnError
2025-10-20T10:00:00.000+0000	using write concern: &{majority <nil> 0s}
2025-10-20T10:00:00.000+0000	using 24 decoding workers
2025-10-20T10:00:00.000+0000	using 1 insert workers
2025-10-20T10:00:00.000+0000	will listen for SIGTERM, SIGINT, and SIGKILL
2025-10-20T10:00:00.000+0000	filesize: 85645 bytes
2025-10-20T10:00:00.000+0000	using fields: 
2025-10-20T10:00:00.000+0000	connected to: mongodb://[**REDACTED**]@mongo-courses-server:27017/test?authSource=admin
2025-10-20T10:00:00.000+0000	ns: test.pokedex
2025-10-20T10:00:00.000+0000	connected to node type: standalone
2025-10-20T10:00:00.000+0000	151 document(s) imported successfully. 0 document(s) failed to import.
```


## Удаление контейнера
  ```bash
$ docker stop mongo-courses-server
mongo-courses-server
$ docker rm mongo-courses-server mongo-courses-importer
mongo-courses-server
```


