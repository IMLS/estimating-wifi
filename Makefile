docs: update docs-dir
	mkdocs build -f imls-docs/wifisess/mkdocs.yml -d $$PWD/docs

update: docs-dir
	git submodule update --recursive --remote


docs-dir:
	if [ ! -d docs ]; then mkdir docs; fi
