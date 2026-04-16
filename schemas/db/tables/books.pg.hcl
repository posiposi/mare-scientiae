table "books" {
  schema = schema.public

  column "id" {
    type    = uuid
    null    = false
    default = sql("gen_random_uuid()")
  }

  column "google_books_id" {
    type = varchar(50)
    null = false
  }

  column "title" {
    type = varchar(500)
    null = false
  }

  column "subtitle" {
    type = varchar(500)
    null = true
  }

  column "authors" {
    type = sql("text[]")
    null = false
  }

  column "created_at" {
    type    = timestamptz
    null    = false
    default = sql("now()")
  }

  column "updated_at" {
    type    = timestamptz
    null    = false
    default = sql("now()")
  }

  primary_key {
    columns = [column.id]
  }

  index "idx_books_google_books_id" {
    columns = [column.google_books_id]
    unique  = true
  }
}
