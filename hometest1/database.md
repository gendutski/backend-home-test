# Database

## Table Design

### Product
Table `product` is for storing list of all product

| Field      | Type          | Description                      |
| ---        | ---           | -----------                      |
| id         | bigint        | AUTO_INCREMENT, Primary Key      |
| serial     | varchar (20)  | Unique                           |
| name       | varchar (255) |                                  |
| price      | double (10,2) |                                  |
| updated_at | timestamp     | Default CURRENT_TIMESTAMP        |

### Product Quantity
Table `product_quantity` is for storing quantity of each product. It has one to one relation with table product.
The purpose this being split is:
- Flexibility: You can add more details regarding stock changes, such as date and time of change, reason for change (sale, return, etc.).
- Performance: Updates to the stock table will not lock the product table, thereby reducing contention in database operations.

| Field      | Type          | Description                         |
| ---        | ---           | -----------                         |
| id         | bigint        | AUTO_INCREMENT, Primary Key         |
| product_id | bigint        | Foreign key reference to product id |
| quantity   | int           | Default 0                           |
| updated_at | timestamp     | Default CURRENT_TIMESTAMP           |


### Promotion
Table `promotion` is for storing of promotion of each products<br />
Field `type` is enum for:
1. Free Item, will provide product items for free.
Products set as free items cannot be promoted.
There must be validation when the admin inputs promotional data.
2. Buy Items to Reduce Price, will provide a reduction in the price of the product
when a user purchases a certain number of items.<br />
Example: get the price value of 2 items if you buy 3 items.
3. Percent Discount, user will get a discount if user buy a number of items.



| Field            | Type          | Description                                    |
| ---              | ---           | -----------                                    |
| id               | bigint        | AUTO_INCREMENT, Primary Key                    |
| type             | int           | Is enum type that hard coded in source         |
| product_id       | bigint        | Foreign key reference to product               |
| match_quantity   | int           | Product quantity for get promotion             |
| promo_value      | float         | Promotion value, eg: discount value            |
| promo_product_id | bigint        | reference to product id, default: 0. indexed   |
| updated_at       | timestamp     | Default CURRENT_TIMESTAMP                      |

## Migrations
You can migrate table using sql files in `migration` folder.
You also can seed table data using `05-seed-data.sql`.
But beware, it will truncate all data

If you using linux, you can use srcipt `run-migration.sh` to run all migration sql.