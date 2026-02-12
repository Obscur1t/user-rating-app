package main

//Table:
// id BIGSERIAL PRIMARY KEY
// name TEXT NOT NULL
// nickname TEXT NOT NULL UNIQUE
// likes INT NOT NULL DEFAULT 0
// viewers INT NOT NULL DEFAULT 0
// rating NUMERIC GENERATED ALWAYS AS (
//     CASE WHEN viewers > 0
//          THEN likes::NUMERIC / viewers
//          ELSE 0
//     END
// ) STORED

func main() {

}
