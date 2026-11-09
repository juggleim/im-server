ALTER TABLE `ioscertificates`
  DROP INDEX `uniq_package`,
  ADD UNIQUE KEY `uniq_package` (`app_key`,`package`);

ALTER TABLE `androidpushconfs`
  DROP INDEX `uniq_channel`,
  ADD UNIQUE KEY `uniq_channel` (`app_key`,`push_channel`,`package`);

ALTER TABLE `androidpushconfs`
CHANGE COLUMN `push_conf` `push_conf` VARCHAR(5000) NULL DEFAULT NULL ;

-- Skip only an exact existing index; incompatible same-name indexes must fail.
SET @index_ddl = (SELECT IF(COUNT(*) = 2 AND SUM(NON_UNIQUE = 1 AND INDEX_TYPE = 'BTREE' AND SUB_PART IS NULL AND COLLATION = 'A' AND IS_VISIBLE = 'YES' AND ((SEQ_IN_INDEX = 1 AND COLUMN_NAME = 'app_key') OR (SEQ_IN_INDEX = 2 AND COLUMN_NAME = 'push_token'))) = 2, 'SELECT 1', 'ALTER TABLE `pushtokens` ADD KEY `idx_pushtoken` (`app_key`,`push_token`)') FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'pushtokens' AND INDEX_NAME = 'idx_pushtoken');
PREPARE index_stmt FROM @index_ddl;
EXECUTE index_stmt;
DEALLOCATE PREPARE index_stmt;

-- The single-node baseline already has nullable varchar(500). Preserve it;
-- also accept the nullable varchar(200) definition used by older upgrades.
SET @column_ddl = (SELECT IF(COUNT(*) = 1 AND SUM(DATA_TYPE = 'varchar' AND CHARACTER_MAXIMUM_LENGTH >= 200 AND IS_NULLABLE = 'YES' AND COLUMN_DEFAULT IS NULL) = 1, 'SELECT 1', 'ALTER TABLE `appnavs` ADD COLUMN `admin_url` varchar(200) DEFAULT NULL AFTER `alias_no`') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'appnavs' AND COLUMN_NAME = 'admin_url');
PREPARE column_stmt FROM @column_ddl;
EXECUTE column_stmt;
DEALLOCATE PREPARE column_stmt;
