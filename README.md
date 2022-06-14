# cloudflare-ddns

cloudflare-ddns is a simple tool that helps me auto-point my RaspberryPI IP to the domain which controls by Cloudflare. And it already integrates cron inside.

## Installation
- Clone repository
- Install with command `go build main.go -o cloudflare-ddns`
    - If you build from MacOS or amd64 arch you can build it with `GOOS=linux GOARCH=arm64 go build main.go -o cloudflare-ddns`
- Copy to your RaspberryPI and run `./cloudflare-ddns --cf-api-key=<your_cf_api_key> --cf-api-email=<your_cf_email> --domains=<your_domain_separate_by_comma>` to test

## How to run
- You can setup it via `systemd` for auto start whenever your PI start / re-start. Eg: mine file is `/etc/systemd/system/cf-ddns.service` with content
```
[Unit]
Description=DDNS using CF to point to domains
Requires=network.target
After=network.target

[Service]
User=root
Type=simple
ExecStart=/root/cloudflare-ddns --cf-api-key="<your_cf_api_key>" --cf-api-email="<your_cf_email>" --domains="<your_domain_separate_by_comma>"
Restart=always

[Install]
WantedBy=multi-user.target
```
- Or you can run binary via nohup / tmux with the following command `cloudflare-ddns --cf-api-key=<your_cf_api_key> --cf-api-email=<your_cf_email> --domains=<your_domain_separate_by_comma>`

## Service log
- Log from stdout and stderr is collect and save to dir `/var/log/cf-ddns`

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License
[MIT](https://choosealicense.com/licenses/mit/)