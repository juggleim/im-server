ALTER TABLE `g_hismsgs` 
ADD INDEX `idx_sender_time` (`app_key`, `conver_id`, `sub_channel`, `sender_id`, `send_time`);

ALTER TABLE `p_hismsgs` 
ADD INDEX `idx_sender_time` (`app_key`, `conver_id`, `sub_channel`, `sender_id`, `send_time`);