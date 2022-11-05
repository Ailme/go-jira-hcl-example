create "Task" {
  project          = "AG"             # required
  # required
  summary          = "${services.service_A.name} // Обновить библиотеку Library_A до актуальной версии"
  # optional
  description      = <<DESC
Нужно обновить библиотеку Library_A до актуальной версии.
После обновления проверить сервис на regress
DESC
  app_layer        = "Backend"          # optional
  components       = ["${services.service_A.name}"]      # optional
  sprint           = 100                # optional
  epic             = "AG-6815"          # optional
  labels           = ["need-regress"]   # optional
  story_point      = 2                  # optional
  qa_story_point   = 1                  # optional
  assignee         = env("DEVELOPER_ALEX")    # optional
  developer        = env("DEVELOPER_ALEX")    # optional
  team_lead        = team_lead          # optional
  tech_lead        = tech_lead          # optional
  release_engineer = release_engineer   # optional
  tester           = env("JIRA_TESTER")
}