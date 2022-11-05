create "Task" {
  project          = "AG"             # required
  # required
  summary          = "service_A // Обновить библиотеку Library_A до актуальной версии"
  # optional
  description      = <<DESC
Нужно обновить библиотеку Library_A до актуальной версии.
После обновления проверить сервис на regress
DESC
  app_layer        = "Backend"        # optional
  components       = ["service_A"]    # optional
  sprint           = 100              # optional
  epic             = "AG-6815"        # optional
  # optional
  labels           = ["need-regress"]
  story_point      = 2                # optional
  qa_story_point   = 1                # optional
  assignee         = "user_A"         # optional
  developer        = "user_A"         # optional
  team_lead        = "user_B"         # optional
  tech_lead        = "user_B"         # optional
  release_engineer = "user_C"         # optional
  tester           = "user_D"         # optional
  #  parent           = "AA-1234"     # optional, используется для Sub-Task
}