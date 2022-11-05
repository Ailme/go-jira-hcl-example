variables {
  developers       = ["jira_user_2", "jira_user_3", "jira_user_4"]
  tester           = "jira_user_1"
  team_lead        = "jira_user_5"
  tech_lead        = "jira_user_5"
  release_engineer = "jira_user_6"
  services         = [
    { name = "service_A"},
    { name = "service_B"},
    { name = "service_C"},
  ]
}

create "Task" {
  project          = "AG"             # required
  # required
  summary          = "${services.0.name} // Обновить библиотеку Library_A до актуальной версии"
  # optional
  description      = <<DESC
Нужно обновить библиотеку Library_A до актуальной версии.
После обновления проверить сервис на regress
DESC
  app_layer        = "Backend"          # optional
  components       = ["${services.0.name}"]      # optional
  sprint           = 100                # optional
  epic             = "AG-6815"          # optional
  labels           = ["need-regress"]   # optional
  story_point      = 2                  # optional
  qa_story_point   = 1                  # optional
  assignee         = developers.0       # optional
  developer        = developers.0       # optional
  team_lead        = team_lead          # optional
  tech_lead        = tech_lead          # optional
  release_engineer = release_engineer   # optional
  tester           = tester             # optional
}