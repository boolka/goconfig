# v1.3.0

- move vault source build to specific build tag (goconfig_vault)
- completely remove embedded vault client authorization
- refactor modules to simplify config package
- explicitly ignore dot prefixed files in config directory
- move out integration tests
- completely rewrite vault module
- remove vault cli functional - now only api usage can trigger vault server

# v1.2.1

- fix directory sources recognition
- fix sources post modification

# v1.2.0

- check vault server availability at startup
- add option to completely disable vault source
- test fixes
- update README.md

# v1.1.0

- add ability to find from certain file sources
- fix common test case config files
- update README.md

# v1.0.1

- add vault auth username/password & roleid/secretid params for cli usage
- flag support
- fix README.md Logger option description
- some minor fixes

# v1.0.0

- add vault support
- add context support
- add MustGet method
- embed config feature
- .editorconfig
- add logger
- concurrent test
- add top package for external imports
- fix concurrent issues
- fix documentation
- add 1.23 support

# v0.1.2

- number normalizations

# v0.1.1

- add github actions
- fix documentation

# v0.1.0

- Initial release
