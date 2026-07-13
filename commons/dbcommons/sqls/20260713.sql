CREATE TABLE IF NOT EXISTS `converconfs` (
  `id` int NOT NULL AUTO_INCREMENT,
  `conver_id` varchar(100) DEFAULT '',
  `conver_type` tinyint DEFAULT '0',
  `sub_channel` varchar(32) DEFAULT '',
  `item_key` varchar(100) DEFAULT '',
  `item_value` varchar(2000) DEFAULT '',
  `item_type` tinyint DEFAULT '0',
  `created_time` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
  `updated_time` datetime(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  `app_key` varchar(20) DEFAULT '',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_key` (`app_key`,`conver_id`,`conver_type`,`sub_channel`,`item_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
