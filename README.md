# Spread love

via QR codes.

## Setup

Create a systemd service on your VPS:

```ini
[Unit]
Description=Spreadlove main service
After=network.target

[Service]
Type=simple
WorkingDirectory=/home/<VPS_USER>/spreadlove
ExecStart=/home/<VPS_USER>/spreadlove/spreadlove
User=<VPS_USER>
Group=<VPS_USER>
Restart=always

[Install]
WantedBy=multi-user.target
```

Start:

```sh
sudo systemctl daemon-reload
sudo systemctl enable spreadlove.service --now
```

Also ensure that you can execute `systemctl restart spreadlove.service` without SUDO:

1. `sudo EDITOR=vim visudo`
2. vps_user ALL=(ALL) NOPASSWD: /usr/bin/systemctl restart spreadlove.service

This is needed for the final step of deploy, where we restart the systemd service.

You will also need to add necessary secrets to GH secrets. Better to generate a new public/private pairing for the rsync
transfer, since it cannot be password protected:
<https://github.com/easingthemes/ssh-deploy?tab=readme-ov-file#configuration>

## TODO

- [x] submit message
- [x] admin panel
- [x] ability to mark pending messages as qualified to move into actual messages
- [x] QR code generation
- [ ] printing
