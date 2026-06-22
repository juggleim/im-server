CREATE TABLE IF NOT EXISTS `performance_metrics` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `node_name` varchar(128) NOT NULL DEFAULT '',
  `collect_time` bigint NOT NULL DEFAULT 0,
  `metric_type` varchar(128) NOT NULL DEFAULT '',
  `metric_value` double NOT NULL DEFAULT 0,
  `created_time` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  KEY `idx_node_collect_time` (`node_name`,`collect_time`),
  KEY `idx_metric_type_collect_time` (`metric_type`,`collect_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `dailyactivities` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `time_mark` BIGINT NULL,
  `count` INT NULL,
  `app_key` VARCHAR(20) NULL,
  `created_time` DATETIME(3) NULL DEFAULT CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  UNIQUE INDEX `uniq_mark` (`app_key` ASC, `time_mark` ASC)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `msgrealtimestats` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `stat_type` tinyint DEFAULT '0',
  `channel_type` tinyint DEFAULT '0',
  `time_mark` bigint DEFAULT '0',
  `count` int DEFAULT '0',
  `app_key` varchar(20) DEFAULT '',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_mark` (`app_key`,`stat_type`,`channel_type`,`time_mark`),
  KEY `idx_time_mark_id` (`time_mark`,`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

ALTER TABLE `useractivities` ADD INDEX `idx_time_mark_id` (`time_mark`,`id`);
ALTER TABLE `connectcounts` ADD INDEX `idx_time_mark_id` (`time_mark`,`id`);