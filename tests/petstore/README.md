# PetStore test

This test suite sets up an Sqlite test suite, that:

- creates a `petstore.db` with data from `petstore.sql`,
- runs ovh/venom integration test with `task`,
- installs ovh/venom with `task setup`.

It requires `etl` be running on port 3000. To quickly put this together,
run `task up` in this projects root directory.
