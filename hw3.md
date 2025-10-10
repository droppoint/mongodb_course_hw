# Домашнее задание №3

## Развертывание MongoDB

  ```bash
docker network create mongo-net
docker run --name mongo-courses-server --network mongo-net -d -p 27017:27017 -v /home/droppoint/mongo:/data/db -e MONGO_INITDB_ROOT_USERNAME="testroot" -e MONGO_INITDB_ROOT_PASSWORD="***" mongo:8.0.14
docker run -it --rm --name mongo-courses-client --network mongo-net mongo:8.0.14 mongosh "mongodb://testroot:***@mongo-courses-server:27017"
```

## Создание коллекций с товарами

```javascript
db.products.insertMany(
  [
    {sku: "10001", name: "Мыло", attributes: {weight: 100}},
    {sku: "10002", name: "Шампунь", attributes: {volume: 200}},
    {sku: "10003", name: "Перчатки", attributes:{materials: ["nitrile"]}},
    {sku: "10004", name: "Средство для мытья посуды", attributes:{volume: 200}},
    {sku: "10005", name: "Набор тарелок", attributes: {materials: ["porcelain"]}},
    {sku: "10006", name: "Кошелек", attributes: {materials: ["leather"]}},
    {sku: "10007", name: "Вантуз",attributes: {materials: ["rubber", "wood"]}},
    {sku: "10008", name: "Набор гвоздей", attributes: {materials: ["stainless steel"]}},
    {sku: "10009", name: "Молоток", attributes: {materials: ["stainless steel", "wood"]}},
    {sku: "10010", name: "Стремянка 2м", attributes: {materials: ["aluminum"]}},
    {sku: "10011", name: "Ведро пластиковое", attributes: {materials: ["plastic"], volume: 5000}},
    {sku: "10012", name: "Набор прищепок", attributes: {materials: ["plastic"]}},
    {sku: "10013", name: "Шланг садовый", attributes: {materials: ["plastic", "rubber"], length: 20000}},
    {sku: "10014", name: "Насос погружной", attributes: {power: 300}},
    {sku: "10015", name: "Розетка электрическая", attributes: {voltage: 230}},
    {sku: "10016", name: "Зарядное устройство", attributes: {ports: ["USB-A", "USB-C"]}},
    {sku: "10017", name: "Варенье вишневое", attributes: {volume: 200}},
    {sku: "10018", name: "Набор конфет", attributes: {weight: 150}},
    {sku: "10019", name: "Набор инструментов", attributes: {weight: 2000}},
    {sku: "10020", name: "Блокнот A4", attributes: {weight: 80}},
  ]
)
```

## Создание wildcard индекса
```javascript
db.products.createIndex({"$**": "text"})
```

## Анализ запроса
```javascript
db.products.find({$text:{$search: "plastic"}}).explain("executionStats")
```
В таком запросе wildcard индекс задействуется
```javascript
{
  explainVersion: '1',
  ...
    winningPlan: {
      isCached: false,
      stage: 'TEXT_MATCH',
      indexPrefix: {},
      indexName: '$**_text',
      parsedTextQuery: {
        terms: [ 'plastic' ],
        negatedTerms: [],
        phrases: [],
        negatedPhrases: []
      },
      textIndexVersion: 3,
      inputStage: {
        stage: 'FETCH',
        inputStage: {
          stage: 'IXSCAN',
          keyPattern: { _fts: 'text', _ftsx: 1 },
          indexName: '$**_text',
          isMultiKey: true,
          isUnique: false,
          isSparse: false,
          isPartial: false,
          indexVersion: 2,
          direction: 'backward',
          indexBounds: {}
        }
      }
    },
    ...
}
```
Но при изменении запроса на более узконаправленный запрос, wildcard индекс уже не используется:
```javascript
db.products.find({"attributes.materials": "plastic"}).explain("executionStats")

{
  explainVersion: '1',
  ...
    winningPlan: {
      isCached: false,
      stage: 'COLLSCAN',
      filter: { 'attributes.materials': { '$eq': 'plastic' } },
      direction: 'forward'
    },
  ...
}
```
для такого запроса нужен иной индекс
```javascript
db.products.createIndex({"attributes.materials": 1})
db.products.find({"attributes.materials": "plastic"}).explain("executionStats")
{
  explainVersion: '1',
  ...
    winningPlan: {
      isCached: false,
      stage: 'FETCH',
      inputStage: {
        stage: 'IXSCAN',
        keyPattern: { 'attributes.materials': 1 },
        indexName: 'attributes.materials_1',
        isMultiKey: true,
        multiKeyPaths: { 'attributes.materials': [ 'attributes.materials' ] },
        isUnique: false,
        isSparse: false,
        isPartial: false,
        indexVersion: 2,
        direction: 'forward',
        indexBounds: { 'attributes.materials': [ '["plastic", "plastic"]' ] }
      }
    },
  ...
}
```


## Удаление контейнера
  ```bash
$ docker stop mongo-courses-server
mongo-courses-server
$ docker rm mongo-courses-server
mongo-courses-server
```


