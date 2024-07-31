CREATE TABLE `promotion` (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `type` int UNSIGNED NOT NULL,
  `product_id` bigint UNSIGNED NOT NULL,
  `match_quantity` int UNSIGNED NOT NULL DEFAULT 0,
  `promo_value` int UNSIGNED NOT NULL DEFAULT 0,
  `promo_product_id` bigint UNSIGNED NOT NULL DEFAULT 0,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `deleted_at` timestamp NULL DEFAULT NULL,

  PRIMARY KEY (`id`),
  FOREIGN KEY `promotion_FK1` (`product_id`) REFERENCES `product` (`id`),
  KEY `promotion_IDX1` (`promo_product_id`)
);