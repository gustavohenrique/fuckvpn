## FuckVPN

FuckVPN is an open-source proxy tool designed to bypass restrictive network environments. 
It allows users to run the tool on a corporate machine and transparently proxy connections 
from other machines on the same network. With FuckVPN, you can securely route traffic 
to external networks, ensuring privacy and access without the need for complicated VPN setups.

### Features

- Transparent Proxy: Seamlessly proxy traffic from machines within the same network.
- No Installation Hassle: Run FuckVPN directly on corporate machines without admin privileges.
- Secure Routing: Protects data by encrypting outbound traffic using a trusted proxy.
- Lightweight: Minimal resource usage, perfect for corporate environments where performance is key.
- Bypass Restrictions: Access external networks even in restricted corporate settings.

### Installation

```sh
go mod tidy
go build -o fuckvpn main.go
```

### Installing CA Certificate in Firefox

To enable HTTPS proxying using FuckVPN, the client machine needs to trust the proxy's SSL certificate. 
Here's how to install the CA certificate in Firefox:

1. Generate a CA certificate or use the `certs/ca.pem`
2. Install the CA Certificate in Firefox:
- Open Firefox.
- Click the menu icon (three horizontal lines) in the top-right corner and go to Settings.
- Scroll down to the Privacy & Security section.
- Under the Certificates section, click View Certificates.
- In the Authorities tab, click Import.
- Select the ca_cert.pem file you generated or provided by FuckVPN.
- After selecting the file, ensure the option Trust this CA to identify websites is checked.
- Click OK to complete the process.
- Now Firefox will trust the proxy’s CA certificate, allowing you to proxy HTTPS traffic without security warnings.

### Usage

```sh
./fuckvpn -cert certs/ca.pem -key certs/ca.key -addr :8080 -v
```

### Contributing

We welcome contributions! If you'd like to help improve FuckVPN, please follow these steps:

1. First, ask me if I'll accept your feature or bugfix.
2. Fork the repository.
3. Create a new branch with your feature or bugfix.
4. Commit your changes.
5. Submit a pull request.

### License

FuckVPN is released under the [MIT License](https://www.mit.edu/~amini/LICENSE.md)

### Disclaimer

FuckVPN is for educational purposes only. Use of this tool in any network, especially corporate environments,
should comply with all applicable laws and policies. We do not endorse or encourage the misuse of this tool in any unauthorized way.

