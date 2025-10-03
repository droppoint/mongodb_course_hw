# Домашнее задание №2

## Развертывание MongoDB

  ```bash
docker network create mongo-net
docker run --name mongo-courses-server --network mongo-net -d -p 27017:27017 -v /home/droppoint/mongo:/data/db -e MONGO_INITDB_ROOT_USERNAME="testroot" -e MONGO_INITDB_ROOT_PASSWORD="***" mongo:8.0.14
docker run -it --rm --name mongo-courses-client --network mongo-net mongo:8.0.14 mongosh "mongodb://testroot:***@mongo-courses-server:27017"
```

## Создание коллекций с товарами и складами

```javascript
db.products.insertMany(
  [
    { sku: "10001", name: "Мыло" },
    { sku: "10002", name: "Шампунь" },
    { sku: "10003", name: "Перчатки" },
    { sku: "10004", name: "Средство для мытья посуды" },
    { sku: "10005", name: "Набор тарелок" },
    { sku: "10006", name: "Кошелек" },
    { sku: "10007", name: "Вантуз" },
    { sku: "10008", name: "Набор гвоздей" },
    { sku: "10009", name: "Молоток" },
    { sku: "10010", name: "Стремянка 2м" },
    { sku: "10011", name: "Ведро пластиковое" },
    { sku: "10012", name: "Набор прищепок" },
    { sku: "10013", name: "Шланг садовый" },
    { sku: "10014", name: "Насос погружной" },
    { sku: "10015", name: "Розетка электрическая" }
  ]
)

db.warehouses.insertMany(
  [
    {
      id: "101",
      address: "ул.Строителей, д.5",
      stocks: [
        { sku: "10001", qty: 1 },
        { sku: "10002", qty: 2 },
        { sku: "10003", qty: 3 },
      ]
    },
    {
      id: "102",
      address: "б-р Раскольникова, д.100",
      stocks: [
        { sku: "10003", qty: 1 },
        { sku: "10004", qty: 2 },
        { sku: "10005", qty: 3 },
      ]
    },
    {
      id: "103",
      address: "ул. Ленина, д.50",
      stocks: [
        { sku: "10005", qty: 1 },
        { sku: "10006", qty: 2 },
        { sku: "10007", qty: 3 },
      ]
    },
    {
      id: "104",
      address: "ул. Маршала Жукова, д.1",
      stocks: [
        { sku: "10007", qty: 1 },
        { sku: "10008", qty: 2 },
        { sku: "10009", qty: 3 },
      ]
    },
    {
      id: "105",
      address: "ул. товарища Кржижановского, д.12",
      stocks: [
        { sku: "10009", qty: 1 },
        { sku: "10010", qty: 2 },
        { sku: "10011", qty: 3 },
      ]
    }
  ]
)
```

## Подсчет количества товаров через aggregation framework
```javascript
db.warehouses.aggregate([
  { $unwind: "$stocks" },
  {
    $group: {
      _id: "$stocks.sku",
      qty: { $sum: "$stocks.qty" }
    }
  },
  {
    $lookup: {
      from: "products",
      localField: "_id",
      foreignField: "sku",
      as: "product"
    }
  },
  { $unwind: "$product" },
  {
    $project: {
      _id: 0,
      sku: "$_id",
      name: "$product.name",
      qty: 1
    }
  }
])
```

## Подсчет количества товаров через mapReduce
```javascript
var productsMap = {};
db.products.find().forEach(function (p) {
  productsMap[p.sku] = p.name;
});

var mapFunction = function () {
  this.stocks.forEach(function (item) {
    emit(item.sku, { qty: item.qty, name: productsMap[item.sku] });
  });
};

var reduceFunction = function (key, values) {
  var total = 0;
  var name = "";
  values.forEach(function (v) {
    total += v.qty;
    name = name || v.name;
  });
  return { qty: total, name: name };
};

db.warehouses.mapReduce(
  mapFunction,
  reduceFunction,
  {
    out: "stock_totals",
    scope: { productsMap: productsMap }
  }
);
```

## Удаление контейнера
  ```bash
$ docker stop mongo-courses-server
mongo-courses-server
$ docker rm mongo-courses-server
mongo-courses-server
```


