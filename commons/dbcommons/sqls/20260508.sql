CREATE TABLE IF NOT EXISTS `usersubrels` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `user_id` varchar(32) DEFAULT '',
  `subscriber_id` varchar(32) DEFAULT '',
  `subscriber_device_id` varchar(50) DEFAULT '',
  `created_time` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
  `app_key` varchar(20) DEFAULT '',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_subscriber` (`app_key`,`user_id`,`subscriber_id`,`subscriber_device_id`),
  KEY `idx_user` (`app_key`,`subscriber_id`,`subscriber_device_id`,`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;