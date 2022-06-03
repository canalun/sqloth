CREATE TABLE `customer` (
  `id` int(10) AUTO_INCREMENT,
  `created_at` timestamp,
  `name` varchar(255) DEFAULT NULL,
  `material` JSON,
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`),
  KEY `id_name` (`id`, `name`),
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `product` (
  `id` int(14) AUTO_INCREMENT,
  `name` varchar(255) DEFAULT NULL,
  `owner` varchar(255) DEFAULT NULL,
  `description` TEXT DEFAULT NULL,
  `stock` tinyint(1),
  `sale_day` Datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`),
  KEY `id_name_stock` (`id`, `name`, `stock`)
  CONSTRAINT `owner` FOREIGN KEY (`owner`) REFERENCES `customer` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;