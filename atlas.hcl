locals {
  databases = {
    local = {
      user     = getenv("POSTGRES_USER")
      password = getenv("POSTGRES_PASSWORD")
      host     = "db"
      port     = "5432"
      db_name  = getenv("POSTGRES_DB")
    }
  }

  postgres_url = {
    for env_name, config in local.databases :
    env_name => "postgres://${config.user}:${config.password}@${config.host}:${config.port}/${config.db_name}?sslmode=disable"
  }
}

env "local" {
  src = [
    "file://schemas/db",
    "file://schemas/db/tables",
  ]
  url = local.postgres_url.local
}
