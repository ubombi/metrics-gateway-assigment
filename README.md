# Unfinished (Abadoned)

Clickhouse library API is strange:  *(If you reading the source code)*
*   `.Prepare` - execute query
*   `.Exec` - marshal data and send it within batch  *(clickhouse library default batch size 1Billion rows)*
*   `.Commit` - Finish query (like real transaction)



Try it with: `docker-compose build && docker-compose up`
