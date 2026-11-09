-- Skip only exact definitions. Do not hide an incompatible same-name index.
SET @index_ddl = (SELECT IF(COUNT(*) = 3 AND SUM(NON_UNIQUE = 1 AND INDEX_TYPE = 'BTREE' AND SUB_PART IS NULL AND COLLATION = 'A' AND IS_VISIBLE = 'YES' AND ((SEQ_IN_INDEX = 1 AND COLUMN_NAME = 'app_key') OR (SEQ_IN_INDEX = 2 AND COLUMN_NAME = 'user_id') OR (SEQ_IN_INDEX = 3 AND COLUMN_NAME = 'time_mark'))) = 3, 'SELECT 1', 'ALTER TABLE `useractivities` ADD INDEX `idx_app_key_userid` (`app_key` ASC, `user_id` ASC, `time_mark` ASC)') FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'useractivities' AND INDEX_NAME = 'idx_app_key_userid');
PREPARE index_stmt FROM @index_ddl;
EXECUTE index_stmt;
DEALLOCATE PREPARE index_stmt;

SET @index_ddl = (SELECT IF(COUNT(*) = 2 AND SUM(NON_UNIQUE = 1 AND INDEX_TYPE = 'BTREE' AND SUB_PART IS NULL AND COLLATION = 'A' AND IS_VISIBLE = 'YES' AND ((SEQ_IN_INDEX = 1 AND COLUMN_NAME = 'time_mark') OR (SEQ_IN_INDEX = 2 AND COLUMN_NAME = 'id'))) = 2, 'SELECT 1', 'ALTER TABLE `useractivities` ADD INDEX `idx_time_mark_id` (`time_mark` ASC, `id` ASC)') FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'useractivities' AND INDEX_NAME = 'idx_time_mark_id');
PREPARE index_stmt FROM @index_ddl;
EXECUTE index_stmt;
DEALLOCATE PREPARE index_stmt;
