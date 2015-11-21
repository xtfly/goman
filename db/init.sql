INSERT INTO `system_setting` (`name`, `value`) VALUES
('site_name', 'GoMan社区'),
('description', 'GoMan社区 社交化知识社区'),
('keywords', 'GoMan,知识社区,社交社区,问答社区'),
('upload_url', ''),
('upload_dir', ''),
('icp_beian', ''),
('sensitive_words', ''),
('def_focus_uids', '1'),
('answer_edit_time', '30'),
('unread_flush_interval', '100'),
('newer_invitation_num', '5'),
('index_per_page', '20'),
('from_email', ''),
('img_url', ''),
('ui_style', 'default'),
('uninterested_fold', '5'),
('answer_unique', 'false'),
('notifications_per_page', '10'),
('contents_per_page', '10'),
('hot_question_period', '7'),
('recommend_users_number', '6'),
('ucenter_enabled', 'false'),
('register_valid_type', 'Y'),
('best_answer_day', '30'),
('answer_self_question', 'true'),
('censoruser', 'admin'),
('best_answer_min_count', '3'),
('reputation_function', '[最佳答案]*3+[赞同]*1-[反对]*1+[发起者赞同]*2-[发起者反对]*1'),
('statistic_code', ''),
('upload_enable', 'true'),
('answer_length_lower', '2'),
('quick_publish', 'true'),
('register_type', '0'),
('question_title_limit', '100'),
('register_seccode', 'true'),
('admin_login_seccode', 'true'),
('comment_limit', '0'),
('best_answer_reput', '20'),
('publisher_reputation_factor', '10'),
('request_route_custom', ''),
('upload_size_limit', '512'),
('upload_avatar_size_limit', '512'),
('topic_title_limit', '12'),
('url_rewrite_enable', 'N'),
('best_agree_min_count', '3'),
('site_close', 'N'),
('close_notice', '站点已关闭，管理员请登录。'),
('integral_system_enabled', 'N'),
('integral_system_config_register', '2000'),
('integral_system_config_profile', '100'),
('integral_system_config_invite', '200'),
('integral_system_config_best_answer', '200'),
('integral_system_config_answer_fold', '-50'),
('integral_system_config_new_question', '-20'),
('integral_system_config_new_answer', '-5'),
('integral_system_config_thanks', '-10'),
('integral_system_config_invite_answer', '-10'),
('username_rule', '1'),
('username_length_min', '2'),
('username_length_max', '14'),
('category_enable', 'Y'),
('integral_unit', '金币'),
('nav_menu_show_child', '1'),
('anonymous_enable', 'Y'),
('report_reason', '广告/SPAM\n恶意灌水\n违规内容\n文不对题\n重复发问'),
('allowed_upload_types', 'jpg,jpeg,png,gif,zip,doc,docx,rar,pdf,psd'),
('site_announce', ''),
('report_message_uid', '1'),
('today_topics', ''),
('welcome_recommend_users', ''),
('welcome_message_pm', '尊敬的{username}，您已经注册成为{sitename}的会员，请您在发表言论时，遵守当地法律法规。\n如果您有什么疑问可以联系管理员。\n\n{sitename}'),
('time_style', '0'),
('reputation_log_factor', '3'),
('advanced_editor_enable', 'Y'),
('auto_question_lock_day', '0'),
('default_timezone', 'Etc/GMT-8'),
('reader_questions_last_days', '30'),
('reader_questions_agree_count', '10'),
('new_user_email_setting', '{"follow_man":10, "new_answer": false}'),
('new_user_notification_setting', '{}'),
('user_action_history_fresh_upgrade', 'Y'),
('ucenter_charset', 'UTF-8'),
('question_topics_limit', '10'),
('auto_create_social_topics', 'N'),
('new_question_force_add_topic', 'N'),
('unfold_question_comments', 'N'),
('report_diagnostics', 'Y'),
('admin_notifications', ''),
('slave_mail_config', ''),
('receiving_email_global_config', '{"enabled:12,"publish_user":false}'),
('last_sent_valid_email_id', '1'),
('integral_system_config_answer_change_source', 'Y'),
('enable_help_center', 'N'),
('ucenter_path', '');


INSERT INTO `jobs` (`id`, `name`) VALUES
(1, '销售'),
(2, '市场/市场拓展/公关'),
(3, '商务/采购/贸易'),
(4, '计算机软、硬件/互联网/IT'),
(5, '电子/半导体/仪表仪器'),
(6, '通信技术'),
(7, '客户服务/技术支持'),
(8, '行政/后勤'),
(9, '人力资源'),
(10, '高级管理'),
(11, '生产/加工/制造'),
(12, '质控/安检'),
(13, '工程机械'),
(14, '技工'),
(15, '财会/审计/统计'),
(16, '金融/银行/保险/证券/投资'),
(17, '建筑/房地产/装修/物业'),
(18, '交通/仓储/物流'),
(19, '普通劳动力/家政服务'),
(20, '零售业'),
(21, '教育/培训'),
(22, '咨询/顾问'),
(23, '学术/科研'),
(24, '法律'),
(25, '美术/设计/创意'),
(26, '编辑/文案/传媒/影视/新闻'),
(27, '酒店/餐饮/旅游/娱乐'),
(28, '化工'),
(29, '能源/矿产/地质勘查'),
(30, '医疗/护理/保健/美容'),
(31, '生物/制药/医疗器械'),
(32, '翻译（口译与笔译）'),
(33, '公务员'),
(34, '环境科学/环保'),
(35, '农/林/牧/渔业'),
(36, '兼职/临时/培训生/储备干部'),
(37, '在校学生'),
(38, '其他');

INSERT INTO `users_group` (`id`, `type`, `custom`, `name`, `repu_lower`, `repu_higer`, `repu_factor`, `permission`) VALUES
(1, 0, 0, '超级管理员', 0, 0, 5, '"is_admin":1,"is_moderator":1,"publish_question":1,"edit_question":1,"edit_topic":1,"manage_topic":1,"create_topic":1,"redirect_question":1,"upload_attach":1,"publish_url":1,"publish_article":1,"edit_article":1,"edit_question_topic":1,"publish_comment":1'),
(2, 0, 0, '前台管理员', 0, 0, 4, '"is_moderator":1,"publish_question":1,"edit_question":1,"edit_topic":1,"manage_topic":1,"create_topic":1,"redirect_question":1,"upload_attach":1,"publish_url":1,"publish_article":1,"edit_article":1,"edit_question_topic":1,"publish_comment":1'),
(3, 0, 0, '未验证会员', 0, 0, 0, '"publish_question":1,"human_valid":1,"question_valid_hour":2,"answer_valid_hour":2'),
(4, 0, 0, '普通会员', 0, 0, 0, '"publish_question":1,"question_valid_hour":10,"answer_valid_hour":10'),
(5, 1, 0, '注册会员', 0, 100, 1, '"publish_question":1,"upload_attach":1,"publish_url":1,"publish_article":1,"edit_question_topic":1,"question_valid_hour":5,"answer_valid_hour":5'),
(6, 1, 0, '初级会员', 100, 200, 1, '"publish_question":1,"upload_attach":1,"publish_url":1,"publish_article":1,"edit_question_topic":1,"edit_question_topic":1,"question_valid_hour":5,"answer_valid_hour":5'),
(7, 1, 0, '中级会员', 200, 500, 1, '"publish_question":1,"edit_topic":1,"create_topic":1,"redirect_question":1,"upload_attach":1,"publish_url":1,"publish_article":1,"publish_comment":1'),
(8, 1, 0, '高级会员', 500, 1000, 1, '"publish_question":1,"edit_question":1,"edit_topic":1,"create_topic":1,"redirect_question":1,"upload_attach":1,"publish_url":1,"publish_article":1,"edit_question_topic":1,"publish_comment":1'),
(9, 1, 0, '核心会员', 1000, 999999, 1, '"publish_question":1,"edit_question":1,"edit_topic":1,"create_topic":1,"redirect_question":1,"upload_attach":1,"publish_url":1,"publish_article":1,"edit_question_topic":1,"publish_comment":1'),
(10, 0, 0, '游客', 0, 0, 0, '"visit_site":1,"visit_explore":1,"search_avail":1,"visit_question":1,"visit_topic":1,"visit_feature":1,"visit_people":1,"visit_chapter":5,"answer_show":1');
