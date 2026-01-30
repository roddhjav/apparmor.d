---
title: Report AppArmor logs
---

The **[aa-log](usage.md#apparmor-log)** tool reports all AppArmor `DENIED` and `ALLOWED`. It should be used to fix AppArmor related issues.

While testing, if something get wrong, you need to put the profile in complain mode, so that you can investigate, and it does not block your program.

When creating [an issue on Github][newissue], please post a link to the [paste] of the audit log generated with:
```sh
aa-log -R
```

!!! question "No logs with `aa-log`?"

    If the log file is empty, check that Auditd is running:

    ```sh
    sudo systemctl status auditd.service
    ```

    If Auditd is disabled aa-log will not have new results, you can enable Auditd with:

    ```sh
    sudo systemctl enable auditd.service --now
    ```

If this command produces nothing, use `-s` to provide all logs since boot time (provided that `journalctl` collected them):
```sh
aa-log -s -R
```

!!! question "No logs with `aa-log -s`?"

    On certain distributions/configurations, AppArmor logs in journal could be taken over by *auditd* when it is installed. To overcome this, `systemd-journald-audit.socket` could be enabled:

    ```sh
    sudo systemctl enable systemd-journald-audit.socket
    ```

You can get older logs with:

```sh
aa-log -R -f <nb>
```

Where `<nb>` is `1`, `2`, `3` and `4` (the rotated audit log file).

[newissue]: https://github.com/roddhjav/apparmor.d/issues/new
[paste]: https://pastebin.com/
