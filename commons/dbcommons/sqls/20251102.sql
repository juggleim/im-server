ALTER TABLE `p_hismsgs` 
ADD COLUMN `read_time` BIGINT NULL DEFAULT 0 AFTER `is_read`;

INSERT INTO `globalconfs` (`conf_key`,`conf_value`)VALUES('jimdb_version','20251102') ON DUPLICATE KEY UPDATE conf_value=VALUES(conf_value);