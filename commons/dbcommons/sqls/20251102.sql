ALTER TABLE `p_hismsgs` 
ADD COLUMN `read_time` BIGINT NULL DEFAULT 0 AFTER `is_read`;
