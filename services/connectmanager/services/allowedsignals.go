package services

// These allowlists describe the standalone client's upstream protocol, not the
// actor registry. Server API and internal methods must not be added here.
func CheckClientTopic(topic string) bool {
	switch topic {
	case "p_msg", "g_msg":
		return true
	default:
		return false
	}
}

func CheckClientQueryTopic(topic string) bool {
	switch topic {
	case "recall_msg", "modify_msg", "qry_hismsgs", "qry_hismsg_by_ids",
		"qry_convers", "qry_top_convers", "sync_convers", "sync_msgs",
		"mark_read", "qry_read_detail", "del_convers", "qry_conver",
		"qry_total_unread_count", "clear_unread", "undisturb_convers", "top_convers",
		"set_top_msg", "del_top_msg", "get_top_msg", "add_favorite_msgs",
		"del_favorite_msgs", "qry_favorite_msgs", "qry_merged_msgs", "reg_push_token",
		"qry_mention_msgs", "clear_total_unread", "del_msg", "clean_hismsg",
		"add_conver", "file_cred", "batch_trans", "set_user_undisturb", "get_user_undisturb",
		"mark_unread", "push_switch", "upd_log_state", "qry_first_unread_msg",
		"msg_exset", "del_msg_exset", "qry_msg_exset", "tag_add_convers", "tag_del_convers",
		"create_user_conver_tags", "del_user_conver_tags", "qry_user_conver_tags",
		"qry_user_info", "sub_user_status", "unsub_user_status", "qry_group_info",
		"qry_friend_infos", "qry_user_status", "rtc_create", "rtc_invite", "rtc_hangup",
		"rtc_accept", "rtc_quit", "rtc_upd_state", "rtc_join", "qry_conver_conf",
		"rtc_member_rooms", "rtc_qry", "set_user_settings", "get_user_settings", "rtc_ping":
		return true
	default:
		return false
	}
}
