# git-credential-superstore

git-credential-superstore is a git credential-helper that thinly wraps the existing git-credential-store helper.

It allows the storing of credentials in separate store files, based on the host and/or path of the repository.

## Usage

1. Build

```bash
go build
```

2. Install:

```bash
sudo install git-credential-superstore /usr/lib/git-core/
```

3. Update `.gitconfig`

```ini
[credential]
    helper = superstore --file foo.com/repo1=~/.git-creds-foo-repo1 bar.com=~/.git-creds-bar specialrepo=~/.git-creds-special =~/.git-credentials
    useHttpPath = true
```

This will store and retrieve credentials for:
- `foo.com/repo1` in `~/.git-creds-foo-repo1`
- `specialrepo` in `.git-creds-special`
- `bar.com` in `.git-creds-bar`
- any other URL in `~/.git-credentials`

In case of multiple matches, the longest match is used.
