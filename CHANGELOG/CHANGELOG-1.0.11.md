 
<a name="v1.0.11"></a>
## [v1.0.11] - 2024-04-01
### Other Changes
- fix to skip contaminated version numbers on `pkg.go.dev`.

### BREAKING CHANGE

disable contaminated product package versions `v1.0.9` and `v1.0.10` on `pkg.go.dev`.

Code migration requires changing the version of the 'bifrost' package from v1.0.9 or v1.0.10 to v1.0.11, as shown in the following example:

`go.mod`:

```
require (
	github.com/ClessLi/bifrost v1.0.11
)
```

[v1.0.11]: https://github.com/ClessLi/bifrost/compare/v1.0.9...v1.0.11
