ALTER TABLE `userconvertags` 
ADD COLUMN `tag_order` INT NULL DEFAULT 0 AFTER `tag_name`,
ADD INDEX `idx_tag` (`app_key`, `user_id`, `tag_order`);

CREATE TABLE IF NOT EXISTS `binddevices` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `app_key` varchar(20) DEFAULT '',
  `user_id` varchar(32) DEFAULT '',
  `platform` varchar(20) DEFAULT '',
  `device_id` varchar(100) DEFAULT '',
  `device_company` varchar(45) DEFAULT NULL,
  `device_model` varchar(45) DEFAULT NULL,
  `created_time` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_deviceid` (`app_key`,`user_id`,`device_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;