# Minisnap

Minisnap is a small utility to create regular snapshot on BTRFS and ZFS volumes.

## Usage

For each volume managed by minisnap, a schedule has to be defined in the configuration file passed in to the program.

See `minisnap.conf` and `zfs.conf` for examples.

Invoking msnap is then as simple as running:

```
msnap -config /etc/minisnap.conf / /home
```

Passing in the `-dry_run` flag to the command will cause `msnap` to not perform any changes, but instead print out what would be done.

## Filesystem support notes

### Btrfs

BTRFS volumes are automatically detected and `msnap` will create and manage all snapshots in the `.snapshots` folder.

### ZFS

ZFS volumes are automatically detected and each created snapshot will be prefixed with `msnap_`.

Note that you *must* pass in the mountpoint of a filesystem, not its name.

The ZFS backend also supports taking recursive snapshots, if configured to do so (see `zfs.conf`).

