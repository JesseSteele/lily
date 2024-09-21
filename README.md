# lily
## Linux In Light Yogurt
*Currently in conceptual phase, changes are on developer branches*

A light-weight web server written in Go with the following features:
- SSL & [Diffie–Hellman](https://en.wikipedia.org/wiki/Diffie%E2%80%93Hellman_key_exchange) support
- Proxy pass
- Apps written in:
  - Go
  - Node.js
  - Python
  - BASH
    - Yes, with database connection method variables (`$_GET`, `$_POST`, `$_SESSION`, `$_SERVER`, etc) set from a unique environment and `.bashrc` functions
    - Intended to more easily manage "state" for AJAX calls
- Non-BASH scripts will automatically be run, so there is no need to create a `.service` config for Systemd; just drop in the main folder
- Cooperate with [standalone](https://eff-certbot.readthedocs.io/en/stable/using.html#standalone) support for Certbot/Letsencrypt
- Support Apache settings (including `RewriteEngine`) common in an `.htaccess` file or inside the server config (both supported)
- Simplifies settings by bringing most common defaults
