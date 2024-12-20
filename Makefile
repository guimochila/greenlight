# Make file
#
ifneq (,$(wildcard .env))
    include .env
    export
endif

migration_dir = "./config/db/migrations"
goose_migrate_cmd = goose postgres "$(DB_DSN)" -dir $(migration_dir)

run-dev:
	air

migration-status:
	$(goose_migrate_cmd) status
migration-up:
	$(goose_migrate_cmd) up
migration-down:
	$(goose_migrate_cmd) down
