---
title: Configuration
---

## AppArmor

As there are a lot of rules, it is recommended to enable caching AppArmor profiles.
In `/etc/apparmor/parser.conf`, add `write-cache` and `Optimize=compress-fast`.

```sh
echo 'write-cache' | sudo tee -a /etc/apparmor/parser.conf
echo 'Optimize=compress-fast' | sudo tee -a /etc/apparmor/parser.conf
```

!!! info

    See [Speed up AppArmor Start] on the Arch Wiki for more information:
    [Speed up AppArmor Start]: https://wiki.archlinux.org/title/AppArmor#Speed-up_AppArmor_start_by_caching_profiles


## Personal directories

This project is designed in such a way that it is easy to personalize the
directories your programs have access by defining a few variables.

The profiles heavily use the (largely extended) XDG directory variables defined
in the **[Variables Reference](variables.md)** page.

??? note "XDG variables overview"

    See **[Variables Reference](variables.md)** page for more.

    | Description | Name | Value |
    |-------------|:----:|---------|
    | Desktop | `@{XDG_DESKTOP_DIR}` | `Desktop` |
    | Download | `@{XDG_DOWNLOAD_DIR}` | `Downloads` |
    | Templates | `@{XDG_TEMPLATES_DIR}` | `Templates` |
    | Public | `@{XDG_PUBLICSHARE_DIR}` | `Public` |
    | Documents | `@{XDG_DOCUMENTS_DIR}` | `Documents` |
    | Music | `@{XDG_MUSIC_DIR}` | `Music` |
    | Pictures | `@{XDG_PICTURES_DIR}` | `Pictures` |
    | Videos | `@{XDG_VIDEOS_DIR}` | `Videos` |
    | Books | `@{XDG_BOOKS_DIR}` | `Books` |
    | Projects | `@{XDG_PROJECTS_DIR}` | `Projects` |
    | Screenshots | `@{XDG_SCREENSHOTS_DIR}` | `@{XDG_PICTURES_DIR}/Screenshots` |
    | Sync | `@{XDG_SYNC_DIR}` | `Sync` |
    | Torrents | `@{XDG_TORRENTS_DIR}` | `Torrents` |
    | Vm | `@{XDG_VM_DIR}` | `.vm`
    | Wallpapers | `@{XDG_WALLPAPERS_DIR}` | `@{XDG_PICTURES_DIR}/Wallpapers` |

You can personalize these values by creating a file such as:
`/etc/apparmor.d/tunables/xdg-user-dirs.d/local` where you define your own
personal directories. Example:
```sh
@{XDG_VIDEOS_DIR}+="Films"
@{XDG_MUSIC_DIR}+="Musique"
@{XDG_PICTURES_DIR}+="Images"
@{XDG_BOOKS_DIR}+="BD" "Comics"
@{XDG_PROJECTS_DIR}+="Git" "Papers"
```

Then restart the apparmor service to reload the profiles in the kernel:
```sh
sudo systemctl restart apparmor.service
```

**Examples**

- For git support, you may want to add your `GO_PATH` in the `XDG_PROJECTS_DIR`:
    ```sh
    @{XDG_PROJECTS_DIR}+="go"
    ```
- If you use Keepass, personalize `XDG_PASSWORD_STORE_DIR` with your password directory. Eg:
    ```sh
    @{XDG_PASSWORD_STORE_DIR}+="@{HOME}/.keepass/"
    ```
- Add pacman integration with your AUR helper. Eg for `yay`:
    ```sh
    @{user_pkg_dirs}+=@{user_cache_dirs}/yay/
    ```

## Local profile extensions

You can extend any profile with your own rules by creating a file in the 
`/etc/apparmor.d/local/` directory with the name of your profile. For example,
to extend the `foo` profile, create a file `/etc/apparmor.d/local/foo` and add
your rules in it.

**Example**

- `child-open`, a profile that allows other program to open resources (URL, 
  picture, books...) with some predefined GUI application. To allow it to open
  URLs with Firefox, create the file `/etc/apparmor.d/local/child-open` with:
  ```sh
    @{bin}/firefox rPx,
  ```

!!! note

    This is an example, no need to add Firefox into `child-open`, it is already there.

!!! info

    `rPx` allows transition to the Firefox profile. Use `rPUx` to allow
    transition to an unconfined state if you do not have the profile for a
    given program.


Then, reload the apparmor rules with `sudo systemctl restart apparmor`.
