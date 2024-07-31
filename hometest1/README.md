# Backend Home Test #1

Pernahkah Anda berbelanja online? Bayangkan Anda perlu membangun layanan pembayaran di bagian belakang yang akan mendukung berbagai promosi dengan inventaris yang diberikan.

Bangun sistem pembayaran dengan item-item berikut:

| SKU    | Name           | Price   | Inventory Qty |
| ---    | ---            | ---:    | ---:         |
| 120P90 | Google Home    | 49.99   | 10            |
| 43N23P | MacBook Pro    | 5399.99 | 5             |
| A304SD | Alexa Speaker  | 109.50  | 10            |
| 234234 | Raspberry Pi B | 20.00   | 2             |

### Sistem harus memiliki metode promosi berikut:
- Setiap penjualan MacBook Pro, disertai Raspberry Pi B gratis
- Beli 3 Google Homes dengan harga 2
- Pembelian lebih dari 3 Speaker Alexa, akan mendapatkan diskon 10% untuk semua speaker Alexa

### Contoh Skenario:
- Barang yang Dipindai: MacBook Pro, Raspberry Pi B<br/>
Total: $5.399,99

- Barang yang Dipindai: Google Home, Google Home, Google Home<br/>
Total: $99,98

- Barang yang Dipindai: Speaker Alexa, Speaker Alexa, Speaker Alexa<br />
Total: $295,6

## Dokumentasi

- [Dokumentasi Database](database.md)
- [API Coontract](api-contract.md)

## How to run
### 1. Migrate database
- Read [Database document](database.md)

### 2. Using go run
- Set `.env` file like `.env-example`
- Run command:
```
go run main.go -loadDotEnv=true
```

### 3. Build docker file
-  Build docker image
```
docker build -t "$DOCKER_NAME" .
```

- Delete existing container
```
docker container rm "$CONTAINER_NAME"
```

- Create container<br>Don't use localhost for mysql host
```
docker container create --name "$CONTAINER_NAME" -e HTTP_PORT=$HTTP_PORT -e MYSQL_HOST="$MYSQL_HOST" -e MYSQL_USERNAME="$MYSQL_USERNAME" -e MYSQL_DB_NAME="$MYSQL_DB_NAME" -e MYSQL_PASSWORD="$MYSQL_PASSWORD" -p $HTTP_PORT:$HTTP_PORT $DOCKER_NAME
```

- Start container
```
docker container start "$CONTAINER_NAME"
```
