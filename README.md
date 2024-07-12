# ETL

This is sort of a type driven ETL with database storage.

It currently has two data models, a concept of `commit`, and
`commit_output`. A taskfile is used along with several etl tool commands:

- `etl commit create`
- `etl commit list`
- `etl commit get <id>`
- `etl output save`
- `etl output list`
- `etl output restore`

The idea is to wire the relationships from one job to the next. The
ingest job just provides data for `etl commit create`.

~~~
# echo bar > output/foo.txt 
# etl commit get 1 | etl output save
Outputs created successfully: 2 files added, 0 skipped
~~~

As outputs would be created, the folder would be passed to etl output
save, and it will scan and store the files into the database based on the
commit id. A job that would produce outputs would invoke `etl output
save` as shown above, passing a valid commit record as input.

Each job would invoke `etl output restore` at the beginning, if they need
any previous inputs to do further processing. For example, an ETL job may
consume coverage reports to provide additional data.

```
# etl commit get 1 | etl output list
Filename: foo.txt
Contents:
bar

Filename: test.txt
Contents:
```

Just for convenience, an `etl output list` is provided.

Enstead of echo you can get the commit data by commit ID:

```
# etl commit get 1
{
  "id": 1,
  "commit_id": "db90d09fbbe01c9609a894e153131be8186055ca",
  "repository": "git@github.com:TykTechnologies/tyk.git",
  "created_at": "2024-07-13T00:27:56.579341028+02:00"
}
```