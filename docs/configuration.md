---
title: Configuration
---

This project is designed in such a way that it is easy to personalize it to fit any system.
It is mostly done by setting personalized XDG like directories in AppArmor tunables. More advanced configuration can be done by adding your own rules in local profile addition.

!!! danger

    You need to ensure that all personal directories you are using are well-defined XDG directory. You may need to edit these variables to your own settings. 

    This part is vital to ensure that the profiles are correctly configured for your system. It will lead to breakage if not done correctly.


## Personalize Apparmor

### Tunables

The profiles heavily use the **largely extended** [XDG directory variables](#xdg-variables). All the variables are list you can append with your own values.

1. First create the directory `/etc/apparmor.d/tunables/xdg-user-dirs.d/apparmor.d.d`:
  ```sh
  sudo mkdir -p /etc/apparmor.d/tunables/xdg-user-dirs.d/apparmor.d.d
  ```
2. Then create a `local` addition file in it where you define your own personal directories. *Example:*
  ```sh
  @{XDG_VIDEOS_DIR}+="Films"
  @{XDG_MUSIC_DIR}+="Musique"
  @{XDG_PICTURES_DIR}+="Images"
  @{XDG_BOOKS_DIR}+="BD" "Comics"
  @{XDG_PROJECTS_DIR}+="Git" "Papers"
  ```
3. Then restart the AppArmor service to reload the profiles in the kernel:
  ```sh
  sudo systemctl reload apparmor.service
  ```

### Profile Additions

You can extend any profile with your own rules by creating a file in the `/etc/apparmor.d/local/` directory with the name of the profile you want to personalize.

**Example**

By default, `nautilus` (and any file browser) only allows access to user files. Thus, your cannot browse system files such as `/etc/`, `/srv/`, `/var/`. You can change this behavior by creating a local profile addition file for `nautilus`:

1. Create the file `/etc/apparmor.d/local/nautilus` and add the following rules in it:
  ```sh
    /** r,
  ```
  You call also restrict this to specific directories:
  ```sh
    /etc/** r,
    /srv/** r,
    /var/** r,
  ```
2. Then restart the AppArmor service to reload the profiles in the kernel:
  ```sh
  sudo systemctl reload apparmor.service
  ```

### XDG variables

Please ensure that all personal directories you are using are well-defined XDG directory defined below. If not, personalize the [variables](#tunables) to your own settings.

??? quote "**User directories**"

    <figure markdown>

      | Description | Name | Default Value(s) |
      |-------------|------|---------------|
      | Desktop | `@{XDG_DESKTOP_DIR}` | `Desktop` |
      | Documents | `@{XDG_DOCUMENTS_DIR}` | `Documents` |
      | Downloads | `@{XDG_DOWNLOAD_DIR}` | `Downloads` |
      | Music | `@{XDG_MUSIC_DIR}` | `Music` |
      | Pictures | `@{XDG_PICTURES_DIR}` | `Pictures` |
      | Videos | `@{XDG_VIDEOS_DIR}` | `Videos` |
      | Screenshots | `@{XDG_SCREENSHOTS_DIR}` | `@{XDG_PICTURES_DIR}/Screenshots` |
      | Wallpapers | `@{XDG_WALLPAPERS_DIR}` | `@{XDG_PICTURES_DIR}/Wallpapers` |
      | Books | `@{XDG_BOOKS_DIR}` | `Books` |
      | Games | `@{XDG_GAMES_DIR}` | `.games` |
      | Templates | `@{XDG_TEMPLATES_DIR}` | `Templates` |
      | Public | `@{XDG_PUBLICSHARE_DIR}` | `Public` |
      | Projects | `@{XDG_PROJECTS_DIR}` | `Projects` |
      | Private | `@{XDG_PRIVATE_DIR}` | `.{p,P}rivate {p,P}rivate` |
      | Work | `@{XDG_WORK_DIR}` | `Work` |
      | Mail | `@{XDG_MAIL_DIR}` | `Mail .{m,M}ail` |
      | Sync | `@{XDG_SYNC_DIR}` | `Sync` |
      | Torrents | `@{XDG_TORRENTS_DIR}` | `Torrents` |
      | Vm | `@{XDG_VM_DIR}` | `.vm` |
      | Vm Shares | `@{XDG_VM_SHARES_DIR}` | `VM_Shares` |
      | Disk images | `@{XDG_IMG_DIR}` | `images` |
      | Games Studio | `@{XDG_GAMESSTUDIO_DIR}` | `.unity3d` |

    </figure>

??? quote "**Dotfiles**"

    <figure markdown>

      | Description | Name | Default Value(s) |
      |-------------|------|------------------|
      | Cache | ` @{XDG_CACHE_DIR}` | `.cache` |
      | Config | `@{XDG_CONFIG_DIR}` | `.config` |
      | Data | `@{XDG_DATA_DIR}` | `.local/share` |
      | State | `@{XDG_STATE_DIR}` | `.local/state` |
      | Bin | `@{XDG_BIN_DIR}` | `.local/bin` |
      | Lib | `@{XDG_LIB_DIR}` | `.local/lib` |
      | GPG | `@{XDG_GPG_DIR}` | `.gnupg` |
      | SSH | `@{XDG_SSH_DIR}` | `.ssh` |
      | Private | `@{XDG_PRIVATE_DIR}` | `.{p,P}rivate {p,P}rivate` |
      | Passwords | `@{XDG_PASSWORD_STORE_DIR}` | `.password-store` |

    </figure>

??? quote "**Full configuration path**"

    <figure markdown>

      | Description | Name | Default Value(s) |
      |-------------|:----:|---------------|
      | Cache | `@{user_cache_dirs}` | `@{HOME}/@{XDG_CACHE_DIR}` |
      | Config | `@{user_config_dirs}` | `@{HOME}/@{XDG_CONFIG_DIR}` |
      | Bin | `@{user_bin_dirs}` | `@{HOME}/@{XDG_BIN_DIR}` |
      | Lib | `@{user_lib_dirs}` | `@{HOME}/@{XDG_LIB_DIR}` |
      | Share | `@{user_share_dirs}` | ` @{HOME}/@{XDG_DATA_DIR}` |
      | State | `@{user_state_dirs}` | ` @{HOME}/@{XDG_STATE_DIR}` |
      | Build | `@{user_build_dirs}` | `/tmp/build/` |
      | Packages | `@{user_pkg_dirs}` | `/tmp/pkg/` |

    </figure>

??? quote "**Full user path**"

    <figure markdown>

      | Description | Name | Default Value(s) |
      |-------------|:----:|---------------|
      | Documents | `@{user_documents_dirs}` | `@{HOME}/@{XDG_DOCUMENTS_DIR} @{MOUNTS}/@{XDG_DOCUMENTS_DIR}` |
      | Downloads | `@{user_download_dirs}` | `@{HOME}/@{XDG_DOWNLOAD_DIR} @{MOUNTS}/@{XDG_DOWNLOAD_DIR}` |
      | Music | `@{user_music_dirs}` | `@{HOME}/@{XDG_MUSIC_DIR} @{MOUNTS}/@{XDG_MUSIC_DIR}` |
      | Pictures | `@{user_pictures_dirs}` | `@{HOME}/@{XDG_PICTURES_DIR} @{MOUNTS}/@{XDG_PICTURES_DIR}` |
      | Videos | `@{user_videos_dirs}` | `@{HOME}/@{XDG_VIDEOS_DIR} @{MOUNTS}/@{XDG_VIDEOS_DIR}` |
      | Books | `@{user_books_dirs}` | `@{HOME}/@{XDG_BOOKS_DIR} @{MOUNTS}/@{XDG_BOOKS_DIR}` |
      | Games | `@{user_games_dirs}` | `@{HOME}/@{XDG_GAMES_DIR} @{MOUNTS}/@{XDG_GAMES_DIR}` |
      | Private | `@{user_private_dirs}` | `@{HOME}/@{XDG_PRIVATE_DIR} @{MOUNTS}/@{XDG_PRIVATE_DIR}` |
      | Passwords | `@{user_passwordstore_dirs}` | `@{HOME}/@{XDG_PASSWORD_STORE_DIR} @{MOUNTS}/@{XDG_PASSWORD_STORE_DIR}` |
      | Work | `@{user_work_dirs}` | `@{HOME}/@{XDG_WORK_DIR} @{MOUNTS}/@{XDG_WORK_DIR}` |
      | Mail | `@{user_mail_dirs}` | `@{HOME}/@{XDG_MAIL_DIR} @{MOUNTS}/@{XDG_MAIL_DIR}` |
      | Projects | `@{user_projects_dirs}` | `@{HOME}/@{XDG_PROJECTS_DIR} @{MOUNTS}/@{XDG_PROJECTS_DIR}` |
      | Public | `@{user_publicshare_dirs}` | `@{HOME}/@{XDG_PUBLICSHARE_DIR} @{MOUNTS}/@{XDG_PUBLICSHARE_DIR}` |
      | Templates | `@{user_templates_dirs}` | `@{HOME}/@{XDG_TEMPLATES_DIR} @{MOUNTS}/@{XDG_TEMPLATES_DIR}` |
      | Torrents | `@{user_torrents_dirs}` | `@{HOME}/@{XDG_TORRENTS_DIR} @{MOUNTS}/@{XDG_TORRENTS_DIR}` |
      | Sync | `@{user_sync_dirs}` | `@{HOME}/@{XDG_SYNC_DIR} @{MOUNTS}/*/@{XDG_SYNC_DIR}` |
      | Vm | `@{user_vm_dirs}` | `@{HOME}/@{XDG_VM_DIR} @{MOUNTS}/@{XDG_VM_DIR}` |
      | Vm Shares | `@{user_vmshare_dirs}` | `@{HOME}/@{XDG_VM_SHARES_DIR} @{MOUNTS}/@{XDG_VM_SHARES_DIR}` |
      | Disk images | `@{user_img_dirs}` | `@{HOME}/@{XDG_IMG_DIR} @{MOUNTS}/@{XDG_IMG_DIR}` |

    </figure>

System variables can also be personalized, they are defined in the **[Variables Reference](variables.md)** page.


## Program Personalization

### Examples

All profiles use the variables defined above. Therefore, you can personalize them by setting your own values in `/etc/apparmor.d/tunables/xdg-user-dirs.d/apparmor.d.d/local`.

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

### Mount points

Common mount points are defined in the `@{MOUNTS}` variable. If you mount a disk on a different location, you can add it to the `@{MOUNTS}` variable.

**Example**

If you mount a disk on `/ssd/`, add the following to `/etc/apparmor.d/tunables/xdg-user-dirs.d/apparmor.d.d/local`:
```sh
@{MOUNTS}+=/ssd/
```

<!-- ### User data

!!! warning "TODO" -->

### File browsers

All supported file browsers (`nautilus`, `dolphin`, `thunar`) are configured to only access user files. If you want to allow access to system files, you can create a local profile addition file for the file browser you are using.

### Games

In order to not allow access to user data, game profiles use the `@{XDG_GAMESSTUDIO_DIR}` variable. It may need to be expanded with other game studio directory. The default is `@{XDG_GAMESSTUDIO_DIR}="unity3d"`.

The `@{XDG_GAMES_DIR}` variable is used to define the game directory such as steam storage directory. If your steam storage is on another drive, you should personalize `@{user_games_dirs}` instead.
