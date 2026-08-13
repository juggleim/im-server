ALTER TABLE `useractivities` 
ADD INDEX `idx_app_key_userid` (`app_key` ASC, `user_id` ASC, `time_mark` ASC);

ALTER TABLE `useractivities` 
ADD INDEX `idx_time_mark_id` (`time_mark` ASC, `id` ASC);