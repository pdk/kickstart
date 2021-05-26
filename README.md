# kickstart

a kick starter for go webapp projects

Check this out, rename all the "kickstart" to your appname, to make it your own.

## iterating

This application contains a built-in file watcher that (if activated) will
terminate the app on file change. Use the `run.sh` script to iterate quickly:

    $ ./scripts/run.sh

If there is compile error after a change, you'll need to re-invoke the script.

## db migration

Database migrations are managed by a simple migrate package.

    $ go run cmd/kickstart.go -migrate sql/*.sql
    2021/05/25 21:23:20 executed sql/2021-05-25-create-thing.sql
    2021/05/25 21:23:20 skipped 0 scripts
    2021/05/25 21:23:20 executed 1 scripts

