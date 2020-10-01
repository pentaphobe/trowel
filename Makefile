VERSION=$(shell cat VERSION)

update_version_tag:
	@echo Updating $(VERSION) tag
	git push origin :refs/tags/$(VERSION)
	git tag -f $(VERSION)
	git push origin --tags

.PHONY: move_version