-- truncate all table
SET FOREIGN_KEY_CHECKS = 0;
TRUNCATE TABLE `promotion`;
TRUNCATE TABLE `product_quantity`;
TRUNCATE TABLE `product`;

-- seed sample product
INSERT INTO `product` (`serial`, `name`, `price`) VALUES
('120P90', 'Google Home', 49.99),
('43N23P', 'MacBook Pro', 5399.99),
('A304SD', 'Alexa Speaker', 109.50),
('234234', 'Raspberry Pi B', 30.00);

-- seed sample product_quantity
INSERT INTO `product_quantity` (`product_id`, `quantity`) VALUES
(1, 10),
(2, 5),
(3, 10),
(4, 2);

-- seed promotion
INSERT INTO `promotion` (`type`, `product_id`, `match_quantity`, `promo_value`, `promo_product_id`) VALUES
(1, 2, 1, 1, 4),
(2, 1, 3, 2, 0),
(3, 3, 3, 10, 0);

SET FOREIGN_KEY_CHECKS = 1;
