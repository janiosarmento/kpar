.PHONY: install

# Builds and deploys kpar in a single atomic step, bumping the patch version.
# This is the only supported way to build+install: it guarantees ~/.local/bin/kpar
# is never left pointing at a stale binary.
install:
	$(eval CURRENT := $(shell cat VERSION))
	$(eval NEW := $(shell echo $(CURRENT) | awk -F. '{print $$1"."$$2"."$$3+1}'))
	@echo $(NEW) > VERSION
	go install -ldflags "-X main.version=$(NEW)" ./cmd/kpar
	/bin/cp -f ~/go/bin/kpar ~/.local/bin/kpar
	@echo "Installed kpar $(NEW) -> ~/.local/bin/kpar"
