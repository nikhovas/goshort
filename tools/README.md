# Instructions

## Warning

Paths to files can be different depending on your system state. This instruction works fine with Vagrant image from
`https://github.com/nikhovas/goshort-dev-env`.

## Requirements

`sudo bash tools/install.sh` from project root.

## Supervisor

### Do:

```bash
sudo supervisorctl reread
sudo supervisorctl update
```

### Check:

```bash
sudo supervisorctl status
```

### Goshort service control:

```bash
sudo supervisorctl restart goshort
sudo supervisorctl start goshort
sudo supervisorctl stop goshort
```

## Logrotate

It is running after executing `install.sh`.

To check what will be done when chron will launch logrotate, do:

```bash
logrotate -d /etc/logrotate.d/goshort
```