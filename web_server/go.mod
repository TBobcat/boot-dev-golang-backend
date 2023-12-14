module web_server

go 1.21.4

require github.com/go-chi/chi/v5 v5.0.10 // indirect
require internal/dblogic v1.0.0
replace internal/dblogic => ./internal/dblogic