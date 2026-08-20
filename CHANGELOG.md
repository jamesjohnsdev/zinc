# Changelog

## [0.1.1](https://github.com/jamesjohnsdev/zinc/compare/v0.1.0...v0.1.1) (2026-08-20)


### Features

* **gh:** add gh CLI wrapper for repos, PRs, issues, and runs ([36bc235](https://github.com/jamesjohnsdev/zinc/commit/36bc2351f7dad2fdee5f64e29267821d1b953b50))
* **git:** add commit log retrieval ([b0af0e2](https://github.com/jamesjohnsdev/zinc/commit/b0af0e2e0fd70e0f7603a89d254ff174e9e0056b))
* **git:** add diff retrieval for files ([b9c323c](https://github.com/jamesjohnsdev/zinc/commit/b9c323c143f9e067c05b60bde3c1ddfc1e1a8f11))
* **git:** add git command runner wrapper ([c931e24](https://github.com/jamesjohnsdev/zinc/commit/c931e24b28b5f1a5768b91c585172fc235c56249))
* **git:** checkout, create, and delete branches ([d058e34](https://github.com/jamesjohnsdev/zinc/commit/d058e34edd7188d1fd04657b18a8b0eebfa625a4))
* **git:** commit staged changes ([4d7ecfc](https://github.com/jamesjohnsdev/zinc/commit/4d7ecfcbad31e09ddd801bcaca4c82555be5d610))
* **git:** list local and remote branches ([0b63fd3](https://github.com/jamesjohnsdev/zinc/commit/0b63fd3fd11a02f9c76d0e0086013d38a34130bb))
* **git:** parse working tree status via porcelain v2 ([5b7abe2](https://github.com/jamesjohnsdev/zinc/commit/5b7abe2e9871e027670ad5d243940327a56059e6))
* **git:** stage and unstage files ([f0e2e91](https://github.com/jamesjohnsdev/zinc/commit/f0e2e919455a9751460bd8528954e2480fe2d21d))
* **git:** stash push, pop, list, and drop ([326b62a](https://github.com/jamesjohnsdev/zinc/commit/326b62ad897b798a5457d053afc6cdcb4e2b1a88))
* scaffold bubbletea application entrypoint ([2e02d00](https://github.com/jamesjohnsdev/zinc/commit/2e02d001e67a497f1d02955ec092f8ea62e25a98))
* **ui:** add base lipgloss theme and layout styles ([a73acb3](https://github.com/jamesjohnsdev/zinc/commit/a73acb3abf0b7b05180a09606a5ddaa669c9be88))
* **ui:** add branches panel ([f805cd4](https://github.com/jamesjohnsdev/zinc/commit/f805cd45c318411896581f1562d486cebd4f9d15))
* **ui:** add commit log panel ([c619980](https://github.com/jamesjohnsdev/zinc/commit/c6199800ac17ab68a9229a928ed91aa5abf30958))
* **ui:** add commit message input panel ([5f82f37](https://github.com/jamesjohnsdev/zinc/commit/5f82f379e36f047556760ce835fdb7c673a3a878))
* **ui:** add diff view panel with scrolling ([0e86396](https://github.com/jamesjohnsdev/zinc/commit/0e86396e8dd261399643c7ddecc58fd1f55dec8b))
* **ui:** add discard changes confirmation ([82af3c9](https://github.com/jamesjohnsdev/zinc/commit/82af3c9b8800e8a0bb5d693c5c050665f842396a))
* **ui:** add GitHub screen with PRs, issues, and workflow runs ([2983045](https://github.com/jamesjohnsdev/zinc/commit/298304532d54ca967cf8bcef53eda858b1d373d1))
* **ui:** add number key shortcuts to jump between panels ([4f69acc](https://github.com/jamesjohnsdev/zinc/commit/4f69acc7e57178aea3dc59a41e272316e25b9427))
* **ui:** add panel focus navigation ([0dfed8c](https://github.com/jamesjohnsdev/zinc/commit/0dfed8c895eb977257fe91c66ef2f0fc487abced))
* **ui:** add refresh and error toast handling ([7b54144](https://github.com/jamesjohnsdev/zinc/commit/7b541446254ca9cb440db8370dea20eb6431399c))
* **ui:** add stash panel ([9f74f7c](https://github.com/jamesjohnsdev/zinc/commit/9f74f7c13fc2c5fb4f1a3940fd2a540760bc2472))
* **ui:** add status bar with keybinding help ([fb003d0](https://github.com/jamesjohnsdev/zinc/commit/fb003d09ea6fbf59cdea3328abb66c2609666ba9))
* **ui:** polish visual theme ([2c84151](https://github.com/jamesjohnsdev/zinc/commit/2c84151257120e1ef36299f76508634966c45264))
* **ui:** render file status panel ([4fd070f](https://github.com/jamesjohnsdev/zinc/commit/4fd070fb30143b8ddb11b59175c9a34586671f80))
* **ui:** wire branch actions with confirmation dialogs ([8a6747f](https://github.com/jamesjohnsdev/zinc/commit/8a6747f522fbb804acff8fed39758942c1d198b4))
* **ui:** wire staging keybindings to status panel ([f27bf42](https://github.com/jamesjohnsdev/zinc/commit/f27bf42ac8f2828e41741ab498617e7b84a6be19))


### Bug Fixes

* **gh:** pass repo as a positional arg to 'gh repo view' ([555d55f](https://github.com/jamesjohnsdev/zinc/commit/555d55f06ed5821bc3a4d9586e41d2234262b9b5))
* **ui:** correct viewport scrolling in diff panel ([b63e9ed](https://github.com/jamesjohnsdev/zinc/commit/b63e9edf020cca4f275f3e00855fb51b72d83aaa))
* **ui:** stop the layout overflowing the terminal on startup ([8b78e2d](https://github.com/jamesjohnsdev/zinc/commit/8b78e2d76d4c944523999425aea600c151385ab9))
* **ui:** truncate list panel lines instead of letting them wrap ([18a3e14](https://github.com/jamesjohnsdev/zinc/commit/18a3e14ca4fe0b7f399a6de49396d5adbe29dd60))
* **ui:** use non-deprecated viewport paging and drop unused color ([0774eab](https://github.com/jamesjohnsdev/zinc/commit/0774eab16922427c145cdd7c4d828cd8cad1a752))
