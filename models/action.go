package models

func HomeActivity(uid int64, limit int) ([]*Topic, bool) {
	// // 我关注的话题
	// focus_topics_ids, ok := GetFoucsTopicIdsByUid(uid)
	// if !ok {
	//
	// }
	//
	// // 关注话题相关的问题列表
	// focus_topics_questions_ids, ok := GetItemIdsByTopicIds(focus_topics_ids, "question", 1000)
	// if !ok {
	//
	// }
	//
	// // 问题对应的话题列表
	// //map[int64]Topics
	// focus_topics, ok := GetTopicsByItemIds(focus_topics_questions_ids)
	// if !ok {
	//
	// }
	// for k, v:= range focus_topics {
	//   for _, q := range v {
	//     if focus_topics_ids.
	//   }
	// }
	//
	//
	// // 关注话题相关的文章列表
	// focus_topics_article_ids, ok := GetItemIdsByTopicIds(focus_topics_ids, "article", 1000)

	return nil, false
}

/*
public function home_activity($uid, $limit = 10)
{
  // 我关注的话题
  if ($user_focus_topics_ids = $this->model('topic')->get_focus_topic_ids_by_uid($uid))
  {
    if ($user_focus_topics_questions_ids = $this->model('topic')->get_item_ids_by_topics_ids($user_focus_topics_ids, 'question', 1000))
    {
      if ($user_focus_topics_info = $this->model('question')->get_topic_info_by_question_ids($user_focus_topics_questions_ids))
      {
        foreach ($user_focus_topics_info AS $key => $user_focus_topics_info_by_question)
        {
          foreach ($user_focus_topics_info_by_question AS $_key => $_val)
          {
            if (!in_array($_val['topic_id'], $user_focus_topics_ids))
            {
              unset($user_focus_topics_info[$key][$_key]);
            }
          }
        }
      }
    }

    $user_focus_topics_article_ids = $this->model('topic')->get_item_ids_by_topics_ids($user_focus_topics_ids, 'article', 1000);
  }
*/
