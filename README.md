# cloudflare-ddns

cloudflare-ddns is a simple tool which already integrated cron inside, that helps me auto-point my RaspberryPI IP to the domain which controls by Cloudflare. Beside of working as schedule worker, it can be use as server to register client IP to domain.

## Installation and Deploy
1. Clone repository

2. Install
- Install with command `go build cmd/ddns-<service_kind>/main.go -o cloudflare-ddns`. Eg:
    - To build service binary: `go build cmd/ddns-service/main.go -o cloudflare-ddns`
    - To build server binary: `go build cmd/ddns-server/main.go -o cloudflare-ddns`
- If you build from another machine not same os and(or) arch, you can build it with additional env params `GOOS=linux GOARCH=arm64 go build cmd/ddns-<service_kind>/main.go -o cloudflare-ddns`
    - Remember to copy built binary to destination machine

3. Test binary built
- Access to destination machine and from path where built binary place at run `./cloudflare-ddns --cf-api-key=<your_cf_api_key> --cf-api-email=<your_cf_email> --domains=<your_domains>` to test
- To test as service binary
    - Wait 5 - 10 second then check cloudflare record
- To test as server binary
    - Run `curl http://localhost:8008/register-ip`
    - Wait 5 - 10 second then check cloudflare record

4. Run binary service with machine startup
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
- Or you can run binary via nohup / tmux with the following command `cloudflare-ddns --cf-api-key=<your_cf_api_key> --cf-api-email=<your_cf_email> --domains=<your_domains>`

## Docker and docker-compose
- Another option to simplify all above step is running service via docker with simple command `docker-compose up`
    - Note that you should change configuration before execute command
    - And if image is wrong kind you need, just change build args and re-build image with `docker-compose build`

## Config flag
### Common config
- `--cf-api-key=<your_cf_api_key>`: required, cloudflare api key
- `--cf-api-email=<your_cf_email>`: required, cloudflare email
- `--domains=<your_domains>`: required, list domain need register, separate by comma
- `--ttl=<ttl_value>`: ttl value when submit to cloudflare, default value `1`
- `--log-dir=<path>`: path to put log file, default value `/var/log/cf-ddns`
- `--development=<false|true>`: run with development mode, default value `false`

### Service config
- `--mode=<cron|single>`: run service one time or cron time

### Server config
- `--port=<http_port>`: set service serve port, default value `8008`
- `--require-auth=<true|false>`: require auth to protect your record when publish service to internal, default value `true`
- `--auth-user=<username>`: username for basic authenticate
- `--auth-password=<password>`: password for basic authenticate

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License
[MIT](https://choosealicense.com/licenses/mit/)