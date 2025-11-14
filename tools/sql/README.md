
# SQL Tools

The SQL tools are useful to work with large SQL file dumps that won't open in other tools because
they're too large. With these tools you can split them into tables, so that you can import or
analyze each table's contents separately.

## Dependencies

none (Pure Go implementation)

## Usage

```bash
# Show list of available tables
sql-tables dump-1337.sql:

# Extract users_table into separate SQL file
sql-extract theleak-1337.sql users_table > users.sql;
```
