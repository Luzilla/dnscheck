# dnscheck

A little command line tool to check if your DNS servers are in sync.

```sh
$ dnscheck check --host github.com
2021/06/11 11:08:31 Discovered local resolver: 192.168.178.1:53
2021/06/11 11:08:31 Probing for NS records (nameservers) for github.com
2021/06/11 11:08:31 => host: github.com
2021/06/11 11:08:31 Found 8 nameservers for github.com.
+-------------------------+-----+------+--------------+
|           NS            | TTL | TYPE |     DATA     |
+-------------------------+-----+------+--------------+
| dns1.p08.nsone.net      |  60 | A    | 140.82.121.4 |
| dns2.p08.nsone.net      |  60 | A    | 140.82.121.4 |
| dns3.p08.nsone.net      |  60 | A    | 140.82.121.4 |
| dns4.p08.nsone.net      |  60 | A    | 140.82.121.4 |
| ns-421.awsdns-52.com    |  60 | A    | 140.82.121.4 |
| ns-520.awsdns-01.net    |  60 | A    | 140.82.121.3 |
| ns-1283.awsdns-32.org   |  60 | A    | 140.82.121.4 |
| ns-1707.awsdns-21.co.uk |  60 | A    | 140.82.121.4 |
+-------------------------+-----+------+--------------+
```

## License

Simplified BSD License
