# ==============================================================================
# Makefile helper functions for release
#
#

.PHONY: release.run
release.run: release.verify release.ensure-tag
	@scripts/release.sh

.PHONY: release.verify
release.verify: tools.verify.git-chglog tools.verify.github-release

.PHONY: release.tag
release.tag: tools.verify.gsemver release.ensure-tag
	@git push origin `git describe --tags --abbrev=0`

.PHONY: release.ensure-tag
release.ensure-tag: tools.verify.gsemver
	@scripts/ensure_tag.sh
