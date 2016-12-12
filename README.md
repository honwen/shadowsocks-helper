# shadowsocks-helper

```shell
NAME:
   ss(r)-helper - shadowsocks(R)-helper

USAGE:
   ssr-helper_linux_amd64 [global options] command [command options] [arguments...]

VERSION:
   20161212

COMMANDS:
     help, h  Shows a list of commands or help for one command
   CONVERTER:
     json2ssr  convert FROM[gui-config.json] to URI[ssr://host:port:protocol:method:obfs:pass]
     json2ss   convert FROM[gui-config.json] to URI[ss://method:pass@host:port]
     ss2ssr    convert URI[ss://method:pass@host:port] to URI[ssr://host:port:protocol:method:obfs:pass]
     ss2json   convert URI[ss://method:pass@host:port] to JSON
     ssr2json  convert URI[ssr://host:port:protocol:method:obfs:pass] to JSON
   HELPER:
     dnsmasq  generate DNSMASQ(server/ipset) from GFWLIST(online) with SSR-Proxy(optional)
     ssrrank  speed-test for shadowsocksr from list of FROM[ssr://host:port:protocol:method:obfs:pass]

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```
