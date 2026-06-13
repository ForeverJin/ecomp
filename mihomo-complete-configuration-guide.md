# mihomo (MetaCubeX) 完整配置指南

> 来源：[https://wiki.metacubex.one/config/](https://wiki.metacubex.one/config/)
>
> 整理日期：2026-06-10
>
> 涵盖所有配置模块：全局配置、DNS、域名嗅探、入站、出站代理、代理集合、代理组、路由规则、规则集合、子规则、流量隧道、NTP、实验性配置

---

## 目录

- [1. 全局配置](#1-全局配置)
- [2. DNS 配置](#2-dns-配置)
- [3. 域名嗅探](#3-域名嗅探)
- [4. 入站配置](#4-入站配置)
  - [4.1 代理端口](#41-代理端口)
  - [4.2 Tun 入站](#42-tun-入站)
  - [4.3 Listeners 监听器](#43-listeners-监听器)
- [5. 出站代理](#5-出站代理)
  - [5.1 通用字段](#51-通用字段)
  - [5.2 TLS 配置](#52-tls-配置)
  - [5.3 传输层配置](#53-传输层配置)
  - [5.4 各协议配置](#54-各协议配置)
- [6. 代理集合](#6-代理集合)
- [7. 代理组](#7-代理组)
- [8. 路由规则](#8-路由规则)
- [9. 规则集合](#9-规则集合)
- [10. 子规则](#10-子规则)
- [11. 流量隧道](#11-流量隧道)
- [12. NTP](#12-ntp)
- [13. 实验性配置](#13-实验性配置)

---

# 1. 全局配置

## 允许局域网

允许其他设备经过 Clash 的代理端口访问互联网，可选值 `true/false`。

绑定地址，仅允许其他设备通过这个地址访问：

- `"*"` 绑定所有 IP 地址
- `192.168.31.31` 绑定单个 IPV4 地址
- `[aaaa::a8aa:ff:fe09:57d8]` 绑定单个 IPV6 地址

允许连接的 IP 地址段，仅作用于 `allow-lan` 为 `true`，默认值为 `0.0.0.0/0` 和 `::/0`：

```yaml
lan-allowed-ips:
  - 0.0.0.0/0
```

禁止连接的 IP 地址段，黑名单优先级高于白名单，默认值为空：

```yaml
lan-disallowed-ips:
  - 192.168.0.3/32
```

### 用户验证

`http(s)`/`socks`/`mixed` 代理的用户验证：

```yaml
authentication:
  - "user1:pass1"
  - "user2:pass2"
```

设置允许跳过验证的 IP 段：

```yaml
skip-auth-prefixes:
  - 127.0.0.1/8
  - ::1/128
```

## 运行模式

- `rule` 规则匹配
- `global` 全局代理（需要在 GLOBAL 策略组选择代理/策略）
- `direct` 全局直连

默认为规则模式。

## 日志级别

Clash 内核输出日志的等级：


| 级别        | 说明                              |
| --------- | ------------------------------- |
| `silent`  | 静默，不输出                          |
| `error`   | 仅输出错误至无法使用的日志                   |
| `warning` | 输出错误但不影响运行的日志，以及 error 级别内容     |
| `info`    | 输出一般运行的内容，以及 error 和 warning 级别 |
| `debug`   | 尽可能输出运行中所有信息                    |


## IPv6

是否允许内核接受 IPv6 流量，默认为 `true`。

## TCP Keep Alive 设置

修改此项以减少移动设备耗电问题：

- `keep-alive-interval`: TCP Keep Alive 包的间隔（秒）
- `keep-alive-max-idle`: TCP Keep Alive 的最大空闲时间
- `disable-keep-alive`: 禁用 TCP Keep Alive，在 Android 强制为 true

```yaml
disable-keep-alive: false
```

## 进程匹配模式

控制是否让 Clash 去匹配进程：


| 模式       | 说明                |
| -------- | ----------------- |
| `always` | 开启，强制匹配所有进程       |
| `strict` | 默认，由 Clash 判断是否开启 |
| `off`    | 不匹配进程，推荐在路由器上使用   |


```yaml
find-process-mode: strict
```

## 外部控制 (API)

外部控制器，可以使用 RESTful API 来控制你的 Clash 内核。

```yaml
external-controller: 127.0.0.1:9090
```

API CORS 标头配置：

```yaml
external-controller-cors:
  allow-origins:
    - '*'
  allow-private-network: true
```

API 的访问密钥（Secret）：

```yaml
secret: "your-secret"
```

Unix socket / Windows namedpipe API 监听（不验证 secret）：

```yaml
external-controller-unix: mihomo.sock
external-controller-pipe: \\.\pipe\mihomo
```

HTTPS-API 监听（需配置 TLS 证书）：

```yaml
external-controller-tls: 127.0.0.1:9443
```

在 RESTful API 端口上开启 DOH 服务器：

```yaml
external-doh-server: /dns-query
```

## 外部用户界面

可以将静态网页资源运行在 Clash API，路径为 API 地址/ui：

```yaml
external-ui: /path/to/ui/folder
external-ui-name: xd
external-ui-url: "https://github.com/MetaCubeX/metacubexd/archive/refs/heads/gh-pages.zip"
```

> **注意**：如果路径不在 Clash 工作目录，请手动设置 `SAFE_PATHS` 环境变量。

## 缓存

```yaml
profile:
  store-selected: true    # 储存 API 对策略组的选择
  store-fake-ip: true     # 储存 fakeip 映射表
```

## 统一延迟

开启统一延迟时，会计算 RTT，以消除连接握手等带来的不同类型节点的延迟差异。

## TCP 并发

启用 TCP 并发连接，将会使用 DNS 解析出的所有 IP 地址进行连接，使用第一个成功的连接。

## 出站接口

mihomo 的流量出站接口。

## 路由标记

为 Linux 下的出站连接提供默认流量标记。

## TLS

目前仅用于 API 的 https：

```yaml
tls:
  certificate: string    # 证书 PEM 格式，或者证书的路径
  private-key: string    # 证书对应的私钥 PEM 格式，或者私钥路径
  ech-key: |-
    -----BEGIN ECH KEYS-----
    ...
    -----END ECH KEYS-----
```

> 自 v1.19.18 版本开始，当 `certificate`、`private-key` 或 `ech-key` 为本地文件时支持自动重载。

## 全局客户端指纹

全局 TLS 客户端指纹伪装配置。

## GEOIP 数据模式

更改 geoip 使用文件，mmdb 或者 dat，可选 `true`/`false`，`true` 为 dat，默认 `false`。

## GEO 文件加载模式


| 模式                | 说明                 |
| ----------------- | ------------------ |
| `standard`        | 标准加载器              |
| `memconservative` | 专为内存受限设备优化的加载器（默认） |


```yaml
geodata-loader: memconservative
```

## 自动更新 GEO

更新间隔，单位为小时。

## 自定 GEO 下载地址

```yaml
geox-url:
  geoip: "https://testingcf.jsdelivr.net/gh/MetaCubeX/meta-rules-dat@release/geoip.dat"
  geosite: "https://testingcf.jsdelivr.net/gh/MetaCubeX/meta-rules-dat@release/geosite.dat"
  mmdb: "https://testingcf.jsdelivr.net/gh/MetaCubeX/meta-rules-dat@release/country.mmdb"
  asn: "https://github.com/xishang0128/geoip/releases/download/latest/GeoLite2-ASN.mmdb"
```

## 自定全局 UA

自定义外部资源下载时使用的 UA，默认为 `clash.meta`。

## ETag 支持

外部资源下载的 ETag 支持，默认为 `true`。

---

# 2. DNS 配置

> 来源：[https://wiki.metacubex.one/config/dns/](https://wiki.metacubex.one/config/dns/)
>
> 整理日期：2026-06-10

---

### 2.1 完整配置模板

```yaml
dns:
  enable: true
  cache-algorithm: arc
  prefer-h3: false
  use-hosts: true
  use-system-hosts: true
  respect-rules: false
  listen: 0.0.0.0:1053
  ipv6: false
  default-nameserver:
    - 223.5.5.5
  enhanced-mode: fake-ip
  fake-ip-range: 198.18.0.1/16
  # fake-ip-range6: fdfe:dcba:9876::1/64
  fake-ip-filter-mode: blacklist
  fake-ip-filter:
    - '*.lan'
  # fake-ip-ttl: 1
  nameserver-policy:
    '+.arpa': '10.0.0.1'
    'rule-set:cn':
      - https://doh.pub/dns-query
      - https://dns.alidns.com/dns-query
  nameserver:
    - https://doh.pub/dns-query
    - https://dns.alidns.com/dns-query
  fallback:
    - tls://8.8.4.4
    - tls://1.1.1.1
  proxy-server-nameserver:
    - https://doh.pub/dns-query
  proxy-server-nameserver-policy:
    'www.yournode.com': '114.114.114.114'
  direct-nameserver:
    - system
  direct-nameserver-follow-policy: false
  fallback-filter:
    geoip: true
    geoip-code: CN
    geosite:
      - gfw
    ipcidr:
      - 240.0.0.0/4
    domain:
      - '+.google.com'
      - '+.facebook.com'
      - '+.youtube.com'
```

---

### 2.2 DNS 支持的类型

mihomo 支持以下 DNS 协议类型：

### UDP

```yaml
- 223.5.5.5
# 或
- udp://223.5.5.5
```

### TCP

```yaml
- tcp://8.8.8.8
```

### DNS over TLS (DoT)

```yaml
- tls://1.1.1.1
```

### DNS over HTTPS (DoH)

```yaml
- https://doh.pub/dns-query
```

### DNS over QUIC (DoQ)

```yaml
- quic://dns.adguard.com:784
```

### system（系统 DNS）

```yaml
- system://
# 或
- system
```

### dhcp

```yaml
# 指定网卡接口
- dhcp://en0

# 仅限 cmfa，使用系统 dns
- dhcp://system
```

### rcode（返回码）

```yaml
- rcode://success            # No error
- rcode://format_error       # Format error
- rcode://server_failure     # Server failure
- rcode://name_error         # Non-existent domain
- rcode://not_implemented    # Not implemented
- rcode://refused            # Query refused
```

---

### 2.3 核心参数详解

### 3.1 enable

是否启用 DNS 模块。如为 `false`，则使用系统 DNS 解析。

### 3.2 cache-algorithm

支持的算法：


| 算法    | 全称                         | 说明             |
| ----- | -------------------------- | -------------- |
| `lru` | Least Recently Used        | 最近最少使用，**默认值** |
| `arc` | Adaptive Replacement Cache | 自适应替换缓存        |


### 3.3 prefer-h3

DOH 是否优先使用 HTTP/3。

> ⚠️ 强烈不建议和 `respect-rules` 一起使用。

### 3.4 listen

DNS 服务监听地址，支持 UDP 和 TCP。例如 `0.0.0.0:1053`。

### 3.5 ipv6

是否解析 IPv6。如为 `false`，则回应 AAAA 的空解析。

### 3.6 enhanced-mode

mihomo 的 DNS 处理模式。


| 模式           | 说明                               |
| ------------ | -------------------------------- |
| `fake-ip`    | 返回虚假 IP，连接时再通过映射获取真实 IP，**性能最佳** |
| `redir-host` | 传统的真实 IP 模式，兼容性更好，**默认值**        |


**推荐使用 `fake-ip` 模式**，响应更快，延迟更低。

### 3.7 fake-ip-range

fake-ip 下的 IPv4 地址段设置。**TUN 的默认 IPv4 地址也使用此值作为参考。**

例如：`198.18.0.1/16`

### 3.8 fake-ip-range6

fake-ip 下的 IPv6 地址段设置。

例如：`fdfe:dcba:9876::1/64`

### 3.9 fake-ip-filter

fake-ip 过滤列表。匹配到的域名不会下发 fake-ip 映射，而是返回真实 IP 用于连接。

值支持**域名通配**以及**引入域名集合**。

### 3.10 fake-ip-filter-mode

可选值 `blacklist` / `whitelist` / `rule`，默认 `blacklist`。


| 模式          | 说明                        |
| ----------- | ------------------------- |
| `blacklist` | 列表中的域名**不返回** fake-ip（默认） |
| `whitelist` | **只有**列表中的域名才返回 fake-ip   |
| `rule`      | 规则模式，支持更灵活的匹配语法           |


#### rule 模式示例

当 `fake-ip-filter-mode` 设置为 `rule` 时，`fake-ip-filter` 的写法发生改变：

```yaml
dns:
  fake-ip-filter-mode: rule
  fake-ip-filter:
    # 与路由 rules 匹配逻辑一致（自上而下），支持 GEOSITE、RuleSet、DOMAIN*、MATCH
    - RULE-SET,reject-domain,fake-ip    # RuleSet behavior 必须为 domain/classical
    - RULE-SET,proxy-domain,fake-ip
    - GEOSITE,gfw,fake-ip
    - DOMAIN,www.baidu.com,real-ip
    - DOMAIN-SUFFIX,qq.com,real-ip
    - DOMAIN-SUFFIX,jd.com,fake-ip
    - MATCH,fake-ip                      # 最终兜底：fake-ip 或 real-ip
```

> 当 RuleSet behavior 为 `classical` 时，仅会生效域名类规则。

### 3.11 fake-ip-ttl

配置 fake-ip 查询返回的 TTL 值。

> ⚠️ 非必要情况下请勿修改。

### 3.12 use-hosts

是否回应配置中的 hosts，默认 `true`。

### 3.13 use-system-hosts

是否查询系统 hosts，默认 `true`。

### 3.14 respect-rules

DNS 连接遵守路由规则，需配置 `proxy-server-nameserver`。

> ⚠️ 强烈不建议和 `prefer-h3` 一起使用。

---

### 2.4 DNS 服务器配置层级

DNS 解析的优先级顺序（从高到低）：

```
nameserver-policy > proxy-server-nameserver（节点域名）
> nameserver + fallback（并发查询 + fallback-filter 判定）
> direct-nameserver（直连出口重新解析）
```

### 4.1 default-nameserver

默认 DNS，用于**解析 DNS 服务器自身的域名**（例如 DOH 服务器的域名）。

- **必须为 IP 地址**，但可为加密 DNS
- 建议填写国内公共 DNS：`223.5.5.5`、`119.29.29.29`

```yaml
default-nameserver:
  - 223.5.5.5
```

### 4.2 nameserver-policy

指定域名查询的解析服务器，**优先级最高**，优先于 `nameserver` / `fallback` 查询。

- 键支持 geosite、域名通配（`+.`）
- 值支持字符串或数组

```yaml
nameserver-policy:
  '+.arpa': '10.0.0.1'
  'rule-set:cn':
    - https://doh.pub/dns-query
    - https://dns.alidns.com/dns-query
```

### 4.3 nameserver

默认的域名解析服务器。

```yaml
nameserver:
  - https://doh.pub/dns-query
  - https://dns.alidns.com/dns-query
```

### 4.4 fallback

后备域名解析服务器，一般情况下使用**境外 DNS**，保证结果可信。

配置 `fallback` 后默认启用 `fallback-filter`，`geoip-code` 默认为 `CN`。

```yaml
fallback:
  - tls://8.8.4.4
  - tls://1.1.1.1
```

### 4.5 proxy-server-nameserver

代理节点域名解析服务器，**仅用于解析代理节点的域名**。

如果不填则遵循 `nameserver-policy` → `nameserver` → `fallback` 的配置。

> 💡 防止「鸡蛋问题」：节点域名需要解析才能连接代理，但代理连接后才能解析。

```yaml
proxy-server-nameserver:
  - https://doh.pub/dns-query
```

### 4.6 proxy-server-nameserver-policy

格式同 `nameserver-policy`，仅用于节点域名解析。**当且仅当** `proxy-server-nameserver` 不为空时生效。

```yaml
proxy-server-nameserver-policy:
  'www.yournode.com': '114.114.114.114'
```

### 4.7 direct-nameserver

用于 `direct` 出口（直连）域名解析的 DNS 服务器。如果不填则遵循 `nameserver-policy` → `nameserver` → `fallback` 的配置。

> 📌 direct-nameserver 重新解析在 v1.19.10 后同样应用于 UDP 连接（对于 TUN 入站仅限 Fakeip 模式下）。

```yaml
direct-nameserver:
  - system
```

### 4.8 direct-nameserver-follow-policy

是否遵循 `nameserver-policy`，默认为不遵守（`false`）。仅当 `direct-nameserver` 不为空时生效。

```yaml
direct-nameserver-follow-policy: false
```

---

### 2.5 fallback-filter（防污染过滤器）（防污染过滤器）

配置 `fallback` 后默认启用，用于判断 `nameserver` 返回的结果是否被 DNS 污染。

满足条件的将使用 `fallback` 结果或只使用 `fallback` 解析。

```yaml
fallback-filter:
  geoip: true
  geoip-code: CN
  geosite:
    - gfw
  ipcidr:
    - 240.0.0.0/4
  domain:
    - '+.google.com'
    - '+.facebook.com'
    - '+.youtube.com'
```

### 5.1 geoip

是否启用 GeoIP 判断。

### 5.2 geoip-code

可选值为国家缩写，默认值为 `CN`。

- 除了 `geoip-code` 配置的国家 IP，其他的 IP 结果会被视为**污染**
- `geoip-code` 配置的国家的结果会**直接采用**，否则将采用 `fallback` 结果

### 5.3 geosite

可选值为对应的 geosite 内包含的集合。

- geosite 列表的内容被视为**已污染**
- 匹配到 geosite 的域名，将**只使用 `fallback` 解析**，不去使用 `nameserver`

### 5.4 ipcidr

书写内容为 `IP/掩码`。

这些网段的结果会被视为污染，`nameserver` 解析出这些结果时将会采用 `fallback` 的解析结果。

### 5.5 domain

这些域名被视为已污染。匹配到这些域名，会直接使用 `fallback` 解析，不去使用 `nameserver`。

---

### 2.6 附加参数

此部分可用于发向公网地址的 DNS 服务器，使用 `#` 附加，使用 `&` 连接不同的参数。

除了指定代理/接口和 ecs，其余项的值均为 bool（true/false）。

### 示例

```yaml
dns:
  nameserver:
    - 'https://8.8.8.8/dns-query#proxy&ecs=1.1.1.1/24&ecs-override=true'
proxies:
  - name: proxy
    type: ss
```

### 6.1 DNS 指定代理/接口进行连接

使用 `#proxy` 或 `#接口名` 指定通过代理/网络接口连接。

- 优先使用已有代理，如果不存在该名称的代理则指定接口连接
- `#RULES` 为遵守路由规则进行连接，等同于 `respect-rules`

> 💡 如需经过代理查询，应配置 `proxy-server-nameserver`，以防出现鸡蛋问题。

### 6.2 h3

强制 HTTP/3。

此选项与 `prefer-h3` 不冲突，填写后强制启用 HTTP/3 建立 DOH 连接。使用前需确保 DOH 服务器支持 HTTP/3。

### 6.3 skip-cert-verify

跳过 TLS 证书验证。

### 6.4 ecs

指定 DNS 查询的 subnet 地址。

### 6.5 ecs-override

强制覆盖 DNS 查询的 subnet 地址。

### 6.6 disable-ipv4

丢弃 A 记录回应。

### 6.7 disable-ipv6

丢弃 AAAA 记录回应。

### 6.8 disable-qtype-int

丢弃特定类型的回应。例如 `disable-qtype-65` 可以屏蔽 HTTPS（TYPE65）类型的 DNS 解析。

### 附加参数汇总表


| 参数                     | 说明                         | 值类型  |
| ---------------------- | -------------------------- | ---- |
| `#proxy` / `#接口名`      | 指定通过代理/网络接口连接              | 字符串  |
| `#RULES`               | 遵守路由规则，等同于 `respect-rules` | —    |
| `#h3`                  | 强制 HTTP/3 连接 DOH           | bool |
| `#skip-cert-verify`    | 跳过 TLS 证书验证                | bool |
| `#ecs=<IP/掩码>`         | 指定 DNS 查询的 subnet 地址       | 字符串  |
| `#ecs-override=true`   | 强制覆盖 subnet 地址             | bool |
| `#disable-ipv4`        | 丢弃 A 记录回应                  | bool |
| `#disable-ipv6`        | 丢弃 AAAA 记录回应               | bool |
| `#disable-qtype-<int>` | 丢弃特定类型回应                   | —    |


---

### 2.7 hosts 配置

```yaml
hosts:
  '*.clash.dev': 127.0.0.1
  'alpha.clash.dev': '::1'
  test.com: [1.1.1.1, 2.2.2.2]
  baidu.com: google.com
```

- 键支持**域名通配**
- 值支持字符串/数组，**域名重定向不支持数组**

> 📌 完整的域名优先级高于使用通配符的域名，例如：
> `foo.example.com` > `*.example.com` > `.example.com`

---

### 2.8 DNS 解析流程

### 示例配置

```yaml
dns:
  nameserver:
    - https://doh.pub/dns-query
  fallback:
    - https://8.8.8.8/dns-query
  direct-nameserver:
    - system
  nameserver-policy:
    "geosite:cn,private":
      - https://doh.pub/dns-query
      - https://dns.alidns.com/dns-query
  fallback-filter:
    geoip: true
    geoip-code: CN
    geosite:
      - gfw
    ipcidr:
      - 240.0.0.0/4
    domain:
      - '+.google.com'
      - '+.facebook.com'
      - '+.youtube.com'

rules:
  - DOMAIN-SUFFIX,google.com,PROXY
  - GEOIP,CN,DIRECT
  - MATCH,PROXY
```

### 解析流程图

```mermaid
flowchart TD
    Rule[匹配规则]
    Domain[匹配到域名规则]
    IP[匹配到目标 IP 规则]
    Resolve[解析域名]

    NameServer[使用 nameserver 查询]
    Policy[匹配 nameserver-policy]
    Concurrent[使用 nameserver 和 fallback 并发查询]
    Filter[匹配 fallback-filter]
    DirectNS[使用 direct-nameserver 重新解析]

    GetIP[查询得到IP]

    Proxy[发送域名给代理]
    Direct[使用 IP 直接连接]

    Rule -->  Domain
    Rule --> IP

    Domain -- 域名匹配到直连 --> Resolve
    Domain -- 域名匹配到直连并配置了 direct-nameserver --> DirectNS

    IP --> Resolve

    Resolve -- 配置了 nameserver-policy --> Policy
    Policy -- 未匹配到 --> NameServer
    Policy --> GetIP
    Resolve -- 未配置 nameserver-policy --> NameServer

    NameServer -- 配置了 fallback --> Concurrent
    Concurrent --> Filter
    Filter --> GetIP
    NameServer -- 未配置 fallback --> GetIP

    GetIP -- 匹配到直连并配置了 direct-nameserver --> DirectNS
    DirectNS --> Direct

    GetIP -- IP 匹配到代理 --> Proxy
    Domain -- 域名匹配到代理 --> Proxy

    GetIP -- IP 匹配到直连 --> Direct
```

### 流程说明

1. **匹配规则**：流量首先经过路由规则匹配，区分域名规则和 IP 规则
2. **域名匹配到直连**：需要解析域名，如果配置了 `direct-nameserver` 则使用它重新解析
3. **解析域名**：
  - 如果配置了 `nameserver-policy`，先匹配策略
  - 未匹配到策略或未配置策略，则使用 `nameserver` 查询
  - 如果配置了 `fallback`，则 `nameserver` 和 `fallback` **并发查询**
  - 查询结果经 `fallback-filter` 判断是否被污染
4. **发送到代理/直连**：
  - 匹配到代理 → 发送域名给代理服务器
  - 匹配到直连 → 使用解析得到的 IP 直接连接

> 📌 direct-nameserver 重新解析在 v1.19.10 后同样应用于 UDP 连接（对于 TUN 入站仅限 Fakeip 模式下）。

---

### 2.9 推荐配置方案（国内用户）（国内用户）

兼顾**解析速度**（国内 DOH）、**防污染**（fallback + filter）和**节点可达性**（独立的 proxy-server-nameserver）。

```yaml
dns:
  enable: true
  cache-algorithm: arc
  listen: 0.0.0.0:1053
  ipv6: false
  default-nameserver:
    - 223.5.5.5
    - 119.29.29.29
  enhanced-mode: fake-ip
  fake-ip-range: 198.18.0.1/16
  fake-ip-filter-mode: blacklist
  fake-ip-filter:
    - '*.lan'
    - 'localhost.ptlogin2.qq.com'
  nameserver:
    - https://doh.pub/dns-query
    - https://dns.alidns.com/dns-query
  fallback:
    - tls://8.8.4.4
    - tls://1.1.1.1
  proxy-server-nameserver:
    - https://doh.pub/dns-query
  fallback-filter:
    geoip: true
    geoip-code: CN
    geosite:
      - gfw
```

---

### 2.10 注意事项

1. **鸡蛋问题**：如需通过代理查询 DNS，务必配置 `proxy-server-nameserver`，避免节点域名无法解析
2. **respect-rules 与 prefer-h3**：强烈不建议同时使用
3. **fake-ip-ttl**：非必要情况请勿修改
4. **fallback-filter**：配置 `fallback` 后默认生效，`geoip-code` 默认为 `CN`
5. **default-nameserver**：必须为 IP 地址，用于解析其他 DNS 服务器的域名
6. **hosts 优先级**：完整域名 > 通配符域名（`*.example.com` > `.example.com`）
7. **direct-nameserver**：v1.19.10 后同样应用于 UDP 连接（TUN 入站仅限 Fakeip 模式下）

---

> 📖 原始文档：[https://wiki.metacubex.one/config/dns/](https://wiki.metacubex.one/config/dns/)

&nbsp;

# 3. 域名嗅探

```yaml
sniffer:
  enable: false
  force-dns-mapping: true
  parse-pure-ip: true
  override-destination: false
  sniff:
    HTTP:
      ports: [80, 8080-8880]
      override-destination: true
    TLS:
      ports: [443, 8443]
    QUIC:
      ports: [443, 8443]
  force-domain:
    - +.v2ex.com
  skip-domain:
    - Mijia Cloud
  skip-src-address:
    - 192.168.0.3/32
  skip-dst-address:
    - 192.168.0.3/32
```


| 参数                     | 说明                                |
| ---------------------- | --------------------------------- |
| `enable`               | 是否启用 sniffer                      |
| `force-dns-mapping`    | 对 `redir-host` 类型识别的流量进行强制嗅探      |
| `parse-pure-ip`        | 对所有未获取到域名的流量进行强制嗅探                |
| `override-destination` | 是否使用嗅探结果作为实际访问，默认为 true           |
| `sniff`                | 需要嗅探的协议设置，仅支持 `HTTP`/`TLS`/`QUIC` |
| `force-domain`         | 强制进行嗅探的域名列表，使用域名通配                |
| `skip-domain`          | 跳过嗅探的域名列表，使用域名通配                  |
| `skip-src-address`     | 跳过嗅探的来源 IP 段列表                    |
| `skip-dst-address`     | 跳过嗅探的目标 IP 段列表                    |


sniff 下的协议设置：

- `ports`: 端口范围
- `override-destination`: 覆盖全局 `override-destination` 设置

---

# 4. 入站配置

## 4.0 listeners 配置总览

```yaml
listeners:
  - name: socks5-in-1
    type: socks
    port: 10808
    #listen: 0.0.0.0
    # rule: sub-rule-name1
    # proxy: proxy
    # udp: false

  - name: http-in-1
    type: http
    port: 10809
    listen: 0.0.0.0

  - name: mixed-in-1
    type: mixed
    port: 10810
    listen: 0.0.0.0
    # udp: false

  - name: redir-in-1
    type: redir
    port: 10811
    listen: 0.0.0.0

  - name: tproxy-in-1
    type: tproxy
    port: 10812
    listen: 0.0.0.0
    # udp: false

  - name: tunnel-in-1
    type: tunnel
    port: 10816
    listen: 0.0.0.0
    network: [tcp, udp]
    target: target.com

  - name: tun-in-1
    type: tun
    stack: system
    dns-hijack:
      - 0.0.0.0:53
    inet4-address:
      - 198.19.0.1/30
    inet6-address:
      - "fdfe:dcba:9877::1/126"
```

通用字段说明：

- `name`: 入站名称
- `type`: 入站类型
- `port`: 监听端口
- `listen`: 监听地址，默认 `0.0.0.0`
- `rule`: 指定子规则，默认使用 rules
- `proxy`: 指定代理，非空则直接交由该 proxy 处理
- `udp`: 是否启用 UDP，默认 true

# 入站配置子页面

## 代理端口

> 来源: [https://wiki.metacubex.one/config/inbound/port/](https://wiki.metacubex.one/config/inbound/port/)

# 代理端口

http/socks/mixed端口验证与外部访问

## http(s) 代理端口

```
port: 7890
```

## socks4/4a/5 代理端口

```
socks-port: 7891
```

## 混合端口

> **Node

混合端口是一个特殊的端口，它同时支持 HTTP(S) 和 SOCKS5 协议。您可以使用任何支持 HTTP 或 SOCKS 代理的程序连接到这个端口

```
mixed-port: 7892
```

## 透明代理端口

> **Note

redirect 端口仅限 Linux(Android) 以及 macOS 适用，tproxy 端口仅限 linux(Android) 适用

redirect 透明代理端口，仅能代理 TCP 流量

```
redir-port: 7893
```

tproxy 透明代理端口，可代理 TCP 与 UDP 流量

```
tproxy-port: 7894
```

```
2024年4月1日
```

---

## Tun 入站

> 来源: [https://wiki.metacubex.one/config/inbound/tun/](https://wiki.metacubex.one/config/inbound/tun/)

# TUN

```
tun:
  enable: true
  stack: system
  auto-route: true
  auto-redirect: true
  auto-detect-interface: true
  dns-hijack:
    - any:53
    - tcp://any:53
  device: utun0
  mtu: 9000
  strict-route: true
  gso: true
  gso-max-size: 65536
  inet6-address: fdfe:dcba:9876::1/126
  udp-timeout: 300
  iproute2-table-index: 2022
  iproute2-rule-index: 9000
  endpoint-independent-nat: false
  route-address-set:
    - ruleset-1
  route-exclude-address-set:
    - ruleset-2
  route-address:
    - 0.0.0.0/1
    - 128.0.0.0/1
    - "::/1"
    - "8000::/1"
  route-exclude-address:
  - 192.168.0.0/16
  - fc00::/7
  include-interface:
  - eth0
  exclude-interface:
  - eth1
  include-uid:
  - 0
  include-uid-range:
  - 1000:9999
  exclude-uid:
  - 1000
  exclude-uid-range:
  - 1000:9999
  include-android-user:
  - 0
  - 10
  include-package:
  - com.android.chrome
  exclude-package:
  - com.android.captiveportallogin

## 旧写法
  inet4-route-address:
  - 0.0.0.0/1
  - 128.0.0.0/1
  inet6-route-address:
  - "::/1"
  - "8000::/1"
  inet4-route-exclude-address:
  - 192.168.0.0/16
  inet6-route-exclude-address:
  - fc00::/7
```

## enable

启用 tun

## stack

tun 模式堆栈，如无使用问题，建议使用 mixed栈，默认 gvisor

可用值： system/gvisor/mixed

> **协议栈之间的区别

- system 使用系统协议栈，可以提供更稳定/全面的 tun 体验，且占用相对其他堆栈更低
- gvisor 通过在用户空间中实现网络协议栈，可以提供更高的安全性和隔离性，同时可以避免操作系统内核和用户空间之间的切换，从而在特定情况下具有更好的网络处理性能
- mixed 混合堆栈，tcp 使用 system栈，udp 使用 gvisor栈，使用体验可能相对更好
- 简单性能测试
- 如果打开了防火墙，则无法使用 system 和 mixed 协议栈，通过以下方式放行内核：
- Windows: 设置 -> Windows 安全中心 -> 允许应用通过防火墙 -> 选中内核
- MacOS: 一般无需配置，防火墙默认放行签名软件，如果遇到开启防火墙无法使用的情况，可以尝试放行：系统设置 -> 网络 -> 防火墙 -> 选项 -> 添加 mihomo app
- Linux: 一般无需配置，防火墙默认不拦截应用，如果遇到开启防火墙无法使用的情况，可以尝试放行 TUN 网卡出站流量（假设 TUN 网卡为 Mihomo）: sudo iptables -A OUTPUT -o Mihomo -j ACCEPT

## device

指定 tun 网卡名称，MacOS 设备只能使用 utun 开头的网卡名

## auto-route

自动设置全局路由，可以自动将全局流量路由进入 tun 网卡。

## auto-redirect

仅支持 Linux，自动配置 iptables/nftables 以重定向 TCP 连接，需要auto-route已启用

*在 Android 中*：

仅转发本地 IPv4 连接。要通过热点或中继共享您的 VPN 连接，请使用 VPNHotspot。

*在 Linux 中*：

带有 auto-route 的 auto-redirect 现在可以在路由器上按预期工作，无需干预。

## auto-detect-interface

自动选择流量出口接口，多出口网卡同时连接的设备建议手动指定出口网卡

## dns-hijack

dns 劫持，将匹配到的连接导入内部 dns 模块，不书写协议则为 udp://

- 在 MacOS/Windows 无法自动劫持发往局域网的 dns 请求
- 在 Android 如开启 私人dns 则无法自动劫持 dns 请求

## strict-route

启用 auto-route 时执行严格的路由规则

*在 Linux 中*:

- 让不支持的网络无法到达
- 将所有连接路由到 tun

它可以防止地址泄漏，并使 DNS 劫持在 Android 上工作。

*在 Windows 中*:

- 添加防火墙规则以阻止 Windows
的 普通多宿主 DNS 解析行为
造成的 DNS 泄露

它可能会使某些应用程序（如 VirtualBox）在某些情况下无法正常工作。

## mtu

最大传输单元，会影响极限状态下的速率，一般用户默认即可。

## gso

启用通用分段卸载，仅支持 Linux

## gso-max-size

数据块的最大长度

## inet6-address

指定tun的ipv6地址

> **Note

- 目前会在程序启动时检查系统其他网卡是否有IPV6，如果不存在会禁用该功能，如果需要强制开启tun的v6地址，请手动设置SKIP_SYSTEM_IPV6_CHECK=1系统环境变量
- 需要同时将顶层配置中的ipv6设置为true才会生效

## udp-timeout

UDP NAT 过期时间，以秒为单位，默认为 300(5 分钟)

## iproute2-table-index

auto-route 生成的 iproute2 路由表索引，默认使用 2022

## iproute2-rule-index

auto-route 生成的 iproute2 规则起始索引，默认使用 9000

## endpoint-independent-nat

启用独立于端点的 NAT，性能可能会略有下降，所以不建议在不需要的时候开启。

## route-address-set

将指定规则集中的目标 IP CIDR 规则添加到防火墙，不匹配的流量将绕过路由
仅支持 Linux，且需要 nftables 以及auto-route 和 auto-redirect 已启用。

与任意配置中的 routing-mark 冲突

## route-exclude-address-set

将指定规则集中的目标 IP CIDR 规则添加到防火墙，匹配的流量将绕过路由
仅支持 Linux，且需要 nftables 以及auto-route 和 auto-redirect 已启用。

与任意配置中的 routing-mark 冲突

## route-address

启用 auto-route时路由自定义路由网段而不是默认路由，一般无需配置。

## route-exclude-address

启用 auto-route 时排除自定义网段

## include-interface

限制被路由的接口，默认不限制，与 exclude-interface 冲突，不可一起配置

## exclude-interface

排除路由的接口，与 include-interface 冲突，不可一起配置

## include-uid

包含的用户，使其被 Tun 路由流量，未被配置的用户不会被 Tun 路由流量，默认不限制

UID 规则仅在 Linux 下被支持,并且需要 auto-route

## include-uid-range

包含的用户范围，使其被 Tun 路由流量，未被配置的用户不会被 Tun 路由流量

## exclude-uid

排除用户，使其避免被 Tun 路由流量

## exclude-uid-range

排除用户范围，使其避免被 Tun 路由流量

## include-android-user

包含的 Android 用户，使其被 Tun 路由流量，未被配置的用户不会被 Tun 路由流量

Android 用户和应用规则仅在 Android 下被支持,并且需要 auto-route

| 常用用户 | ID |
| 机主 | 0 |
| 手机分身 | 10 |
| 应用多开 | 999 |

## include-package

包含的 Android 应用包名，使其被 Tun 路由流量，未配置的应用包不会被 Tun 路由流量

## exclude-package

排除 Android 应用包名，使其避免被 Tun 路由流量

## 旧写法，即将废弃

### inet4-route-address

启用 auto-route时路由自定义网段而不是默认路由，一般无需配置。

### inet6-route-address

启用 auto-route时路由自定义网段而不是默认路由，一般无需配置。

### inet4-route-exclude-address

启用 auto-route 时排除自定义网段

### inet6-route-exclude-address

启用 auto-route 时排除自定义网段

## Tun 的协议栈网络回环测试

从上到下分别为 system/gvisor/lwip,仅供参考，平台为 linux,Windows 和 MacOS 可能会有差异

```
2026年5月22日
```

---

# Listeners 监听器

## HTTP 监听器

> 来源: [https://wiki.metacubex.one/config/inbound/listeners/http/](https://wiki.metacubex.one/config/inbound/listeners/http/)

# HTTP

```
listeners:
- name: http-in
  type: http
  port: 7890
  listen: 0.0.0.0
  users:
    - username: username1
      password: password1
  # 下面两项如果填写则开启 tls（需要同时填写）
  # certificate: ./server.crt # 证书 PEM 格式，或者 证书的路径
  # private-key: ./server.key # 证书对应的私钥 PEM 格式，或者私钥路径
  # 下面两项为mTLS配置项，如果client-auth-type设置为 "verify-if-given" 或 "require-and-verify" 则client-auth-cert必须不为空
  # client-auth-type: "" # 可选值：""、"request"、"require-any"、"verify-if-given"、"require-and-verify"
  # client-auth-cert: string # 证书 PEM 格式，或者 证书的路径
  # 如果填写则开启ech（可由 mihomo generate ech-keypair <明文域名> 生成）
  # ech-key: |
  #   -----BEGIN ECH KEYS-----
  #   ACATwY30o/RKgD6hgeQxwrSiApLaCgU+HKh7B6SUrAHaDwBD/g0APwAAIAAgHjzK
  #   madSJjYQIf9o1N5GXjkW4DEEeb17qMxHdwMdNnwADAABAAEAAQACAAEAAwAIdGVz
  #   dC5jb20AAA==
  #   -----END ECH KEYS-----
```

## 通用字段

## 协议配置

### 用户验证

如果不填写 users 项，则遵从全局 用户验证 设置，如果填写会忽略全局设置，如想跳过该入站的验证可填写 users: []

```
2025年9月24日
```

---

## SOCKS 监听器

> 来源: [https://wiki.metacubex.one/config/inbound/listeners/socks/](https://wiki.metacubex.one/config/inbound/listeners/socks/)

# SOCKS

```
listeners:
- name: socks-in
  type: socks
  port: 7891
  listen: 0.0.0.0
  udp: true
  users:
    - username: username1
      password: password1
  # 下面两项如果填写则开启 tls（需要同时填写）
  # certificate: ./server.crt # 证书 PEM 格式，或者 证书的路径
  # private-key: ./server.key # 证书对应的私钥 PEM 格式，或者私钥路径
  # 下面两项为mTLS配置项，如果client-auth-type设置为 "verify-if-given" 或 "require-and-verify" 则client-auth-cert必须不为空
  # client-auth-type: "" # 可选值：""、"request"、"require-any"、"verify-if-given"、"require-and-verify"
  # client-auth-cert: string # 证书 PEM 格式，或者 证书的路径
  # 如果填写则开启ech（可由 mihomo generate ech-keypair <明文域名> 生成）
  # ech-key: |
  #   -----BEGIN ECH KEYS-----
  #   ACATwY30o/RKgD6hgeQxwrSiApLaCgU+HKh7B6SUrAHaDwBD/g0APwAAIAAgHjzK
  #   madSJjYQIf9o1N5GXjkW4DEEeb17qMxHdwMdNnwADAABAAEAAQACAAEAAwAIdGVz
  #   dC5jb20AAA==
  #   -----END ECH KEYS-----
```

## 通用字段

## 协议配置

### udp

是否监听 UDP

### 用户验证

如果不填写 users 项，则遵从全局 用户验证 设置，如果填写会忽略全局设置，如想跳过该入站的验证可填写 users: []

```
2025年9月24日
```

---

## Mixed 混合监听器

> 来源: [https://wiki.metacubex.one/config/inbound/listeners/mixed/](https://wiki.metacubex.one/config/inbound/listeners/mixed/)

# MIXED

```
listeners:
- name: mixed-in
  type: mixed
  port: 7892
  listen: 0.0.0.0
  udp: true
  users:
    - username: username1
      password: password1
  # 下面两项如果填写则开启 tls（需要同时填写）
  # certificate: ./server.crt # 证书 PEM 格式，或者 证书的路径
  # private-key: ./server.key # 证书对应的私钥 PEM 格式，或者私钥路径
  # 下面两项为mTLS配置项，如果client-auth-type设置为 "verify-if-given" 或 "require-and-verify" 则client-auth-cert必须不为空
  # client-auth-type: "" # 可选值：""、"request"、"require-any"、"verify-if-given"、"require-and-verify"
  # client-auth-cert: string # 证书 PEM 格式，或者 证书的路径
  # 如果填写则开启ech（可由 mihomo generate ech-keypair <明文域名> 生成）
  # ech-key: |
  #   -----BEGIN ECH KEYS-----
  #   ACATwY30o/RKgD6hgeQxwrSiApLaCgU+HKh7B6SUrAHaDwBD/g0APwAAIAAgHjzK
  #   madSJjYQIf9o1N5GXjkW4DEEeb17qMxHdwMdNnwADAABAAEAAQACAAEAAwAIdGVz
  #   dC5jb20AAA==
  #   -----END ECH KEYS-----
```

## 通用字段

## 协议配置

### udp

是否监听 UDP

### 用户验证

如果不填写 users 项，则遵从全局 用户验证 设置，如果填写会忽略全局设置，如想跳过该入站的验证可填写 users: []

```
2025年9月24日
```

---

## Redirect 重定向监听器

> 来源: [https://wiki.metacubex.one/config/inbound/listeners/redirect/](https://wiki.metacubex.one/config/inbound/listeners/redirect/)

# REDIRECT

```
listeners:
- name: redir-in
  type: redir
  port: 7893
  listen: 0.0.0.0
```

## 通用字段

```
2024年4月1日
```

---

## TProxy 监听器

> 来源: [https://wiki.metacubex.one/config/inbound/listeners/tproxy/](https://wiki.metacubex.one/config/inbound/listeners/tproxy/)

# TPROXY

```
listeners:
- name: tproxy-in
  type: tproxy
  port: 7894
  listen: 0.0.0.0
  udp: true
```

## 通用字段

## 协议配置

### udp

是否监听 UDP

```
2024年4月1日
```

---

## TUN 监听器

> 来源: [https://wiki.metacubex.one/config/inbound/listeners/tun/](https://wiki.metacubex.one/config/inbound/listeners/tun/)

# TUN

注意，listeners中的tun仅提供给高级用户使用，普通用户应使用顶层配置中的tun

```
listeners:
- name: tun-in
  type: tun
  stack: system
  dns-hijack:
  - 0.0.0.0:53
  # auto-detect-interface: false
  # auto-route: false
  # mtu: 9000
  inet4-address:
  - 198.19.0.1/30
  inet6-address:
  - "fdfe:dcba:9877::1/126"
  # strict-route: true
  # inet4-route-address:由
  # - 0.0.0.0/1
  # - 128.0.0.0/1
  # inet6-route-address:
  # - "::/1"
  # - "8000::/1"
  # endpoint-independent-nat: false
  # include-uid:
  # - 0
  # include-uid-range:
  # - 1000-99999
  # exclude-uid:
  # - 1000
  # exclude-uid-range:
  # - 1000-99999
  # include-android-user:
  # - 0
  # - 10
  # include-package:
  # - com.android.chrome
  # exclude-package:
  # - com.android.captiveportallogin
```

```
2025年8月19日
```

---

## VMess 监听器

> 来源: [https://wiki.metacubex.one/config/inbound/listeners/vmess/](https://wiki.metacubex.one/config/inbound/listeners/vmess/)

# VMESS

```
listeners:
- name: vmess-in-1
  type: vmess
  port: 10814 # 支持使用ports格式，例如200,302 or 200,204,401-429,501-503
  listen: 0.0.0.0
  # rule: sub-rule-name1 # 默认使用 rules，如果未找到 sub-rule 则直接使用 rules
  # proxy: proxy # 如果不为空则直接将该入站流量交由指定 proxy 处理 (当 proxy 不为空时，这里的 proxy 名称必须合法，否则会出错)
  users:
    - username: 1
      uuid: 9d0cb9d0-964f-4ef6-897d-6c6b3ccf9e68
      alterId: 1
  # ws-path: "/" # 如果不为空则开启 websocket 传输层
  # grpc-service-name: "GunService" # 如果不为空则开启 grpc 传输层
  # 下面两项如果填写则开启 tls（需要同时填写）
  # certificate: ./server.crt # 证书 PEM 格式，或者 证书的路径
  # private-key: ./server.key # 证书对应的私钥 PEM 格式，或者私钥路径
  # 下面两项为mTLS配置项，如果client-auth-type设置为 "verify-if-given" 或 "require-and-verify" 则client-auth-cert必须不为空
  # client-auth-type: "" # 可选值：""、"request"、"require-any"、"verify-if-given"、"require-and-verify"
  # client-auth-cert: string # 证书 PEM 格式，或者 证书的路径
  # 如果填写则开启ech（可由 mihomo generate ech-keypair <明文域名> 生成）
  # ech-key: |
  #   -----BEGIN ECH KEYS-----
  #   ACATwY30o/RKgD6hgeQxwrSiApLaCgU+HKh7B6SUrAHaDwBD/g0APwAAIAAgHjzK
  #   madSJjYQIf9o1N5GXjkW4DEEeb17qMxHdwMdNnwADAABAAEAAQACAAEAAwAIdGVz
  #   dC5jb20AAA==
  #   -----END ECH KEYS-----
  # 如果填写reality-config则开启reality（注意不可与certificate和private-key同时填写）
  # reality-config:
  #   dest: test.com:443
  #   private-key: jNXHt1yRo0vDuchQlIP6Z0ZvjT3KtzVI-T4E7RoLJS0 # 可由 mihomo generate reality-keypair 命令生成
  #   short-id:
  #     - 0123456789abcdef
  #   server-names:
  #     - test.com
  #   #下列两个 limit 为选填，可对未通过验证的回落连接限速，bytesPerSec 默认为 0 即不启用
  #   #回落限速是一种特征，不建议启用，如果您是面板/一键脚本开发者，务必让这些参数随机化
  #   limit-fallback-upload:
  #     after-bytes: 0 # 传输指定字节后开始限速
  #     bytes-per-sec: 0 # 基准速率（字节/秒）
  #     burst-bytes-per-sec: 0 # 突发速率（字节/秒），大于 bytesPerSec 时生效
  #   limit-fallback-download:
  #     after-bytes: 0 # 传输指定字节后开始限速
  #     bytes-per-sec: 0 # 基准速率（字节/秒）
  #     burst-bytes-per-sec: 0 # 突发速率（字节/秒），大于 bytesPerSec 时生效
```

```
2025年9月24日
```

---

## VLESS 监听器

> 来源: [https://wiki.metacubex.one/config/inbound/listeners/vless/](https://wiki.metacubex.one/config/inbound/listeners/vless/)

# VLESS

```
listeners:
- name: vless-in-1
  type: vless
  port: 10817 # 支持使用ports格式，例如200,302 or 200,204,401-429,501-503
  listen: 0.0.0.0
  # rule: sub-rule-name1 # 默认使用 rules，如果未找到 sub-rule 则直接使用 rules
  # proxy: proxy # 如果不为空则直接将该入站流量交由指定 proxy 处理 (当 proxy 不为空时，这里的 proxy 名称必须合法，否则会出错)
  users:
    - username: 1
      uuid: 9d0cb9d0-964f-4ef6-897d-6c6b3ccf9e68
      flow: xtls-rprx-vision
  # ws-path: "/" # 如果不为空则开启 websocket 传输层
  # grpc-service-name: "GunService" # 如果不为空则开启 grpc 传输层
  # xhttp-config: # 如果不为空则开启 xhttp 传输层
  #   path: "/"
  #   host: ""
  #   mode: auto # Available: "stream-one", "stream-up" or "packet-up"
  #   no-sse-header: false
  #   x-padding-bytes: "100-1000"
  #   x-padding-obfs-mode: false
  #   x-padding-key: x_padding
  #   x-padding-header: Referer
  #   x-padding-placement: queryInHeader # Available: queryInHeader, cookie, header, query
  #   x-padding-method: repeat-x # Available: repeat-x, tokenish
  #   uplink-http-method: POST # Available: POST, PUT, PATCH, DELETE
  #   session-placement: path # Available: path, query, cookie, header
  #   session-key: ""
  #   seq-placement: path # Available: path, query, cookie, header
  #   seq-key: ""
  #   uplink-data-placement: body # Available: body, cookie, header
  #   uplink-data-key: ""
  #   uplink-chunk-size: 0 # only applicable when uplink-data-placement is not body
  #   sc-max-buffered-posts: 30
  #   sc-stream-up-server-secs: "20-80"
  #   sc-max-each-post-bytes: 1000000
  # -------------------------
  # vless encryption服务端配置：
  # （原生外观 / 只 XOR 公钥 / 全随机数。1-RTT 每次下发随机 300 到 600 秒的 ticket 以便 0-RTT 复用 / 只允许 1-RTT）
  # 填写 "600s" 会每次随机取 50% 到 100%，即相当于填写 "300-600s"
  # / 是只能选一个，后面 base64 至少一个，无限串联，使用 mihomo generate vless-x25519 和 mihomo generate vless-mlkem768 生成，替换值时需去掉括号
  #
  # Padding 是可选的参数，仅作用于 1-RTT 以消除握手的长度特征，双端默认值均为 "100-111-1111.75-0-111.50-0-3333"：
  # 在 1-RTT client/server hello 后以 100% 的概率粘上随机 111 到 1111 字节的 padding
  # 以 75% 的概率等待随机 0 到 111 毫秒（"probability-from-to"）
  # 再次以 50% 的概率发送随机 0 到 3333 字节的 padding（若为 0 则不 Write()）
  # 服务端、客户端可以设置不同的 padding 参数，按 len、gap 的顺序无限串联，第一个 padding 需概率 100%、至少 35 字节
  # -------------------------
  # decryption: "mlkem768x25519plus.native/xorpub/random.600s(300-600s)/0s.(padding len).(padding gap).(X25519 PrivateKey).(ML-KEM-768 Seed)..."
  # 下面两项如果填写则开启 tls（需要同时填写）
  # certificate: ./server.crt # 证书 PEM 格式，或者 证书的路径
  # private-key: ./server.key # 证书对应的私钥 PEM 格式，或者私钥路径
  # 下面两项为mTLS配置项，如果client-auth-type设置为 "verify-if-given" 或 "require-and-verify" 则client-auth-cert必须不为空
  # client-auth-type: "" # 可选值：""、"request"、"require-any"、"verify-if-given"、"require-and-verify"
  # client-auth-cert: string # 证书 PEM 格式，或者 证书的路径
  # 如果填写则开启ech（可由 mihomo generate ech-keypair <明文域名> 生成）
  # ech-key: |
  #   -----BEGIN ECH KEYS-----
  #   ACATwY30o/RKgD6hgeQxwrSiApLaCgU+HKh7B6SUrAHaDwBD/g0APwAAIAAgHjzK
  #   madSJjYQIf9o1N5GXjkW4DEEeb17qMxHdwMdNnwADAABAAEAAQACAAEAAwAIdGVz
  #   dC5jb20AAA==
  #   -----END ECH KEYS-----
  # 如果填写reality-config则开启reality（注意不可与certificate和private-key同时填写）
  reality-config:
    dest: test.com:443
    private-key: jNXHt1yRo0vDuchQlIP6Z0ZvjT3KtzVI-T4E7RoLJS0 # 可由 mihomo generate reality-keypair 命令生成
    short-id:
      - 0123456789abcdef
    server-names:
      - test.com
    #下列两个 limit 为选填，可对未通过验证的回落连接限速，bytesPerSec 默认为 0 即不启用
    #回落限速是一种特征，不建议启用，如果您是面板/一键脚本开发者，务必让这些参数随机化
    limit-fallback-upload:
      after-bytes: 0 # 传输指定字节后开始限速
      bytes-per-sec: 0 # 基准速率（字节/秒）
      burst-bytes-per-sec: 0 # 突发速率（字节/秒），大于 bytesPerSec 时生效
    limit-fallback-download:
      after-bytes: 0 # 传输指定字节后开始限速
      bytes-per-sec: 0 # 基准速率（字节/秒）
      burst-bytes-per-sec: 0 # 突发速率（字节/秒），大于 bytesPerSec 时生效
  ### 注意，对于vless listener, 如果 "allow-insecure" 不为 true, 至少需要填写 “certificate和private-key” 或 “reality-config” 或 “decryption” 的其中一项 ###
  # allow-insecure: false # 是否允许不开启tls加密（注意：仅用于有 nginx, caddy 前置的情况）
```

```
2026年6月6日
```

---

## Trojan 监听器

> 来源: [https://wiki.metacubex.one/config/inbound/listeners/trojan/](https://wiki.metacubex.one/config/inbound/listeners/trojan/)

# Trojan

```
listeners:
- name: trojan-in-1
  type: trojan
  port: 10819 # 支持使用ports格式，例如200,302 or 200,204,401-429,501-503
  listen: 0.0.0.0
  # rule: sub-rule-name1 # 默认使用 rules，如果未找到 sub-rule 则直接使用 rules
  # proxy: proxy # 如果不为空则直接将该入站流量交由指定 proxy 处理 (当 proxy 不为空时，这里的 proxy 名称必须合法，否则会出错)
  users:
    - username: 1
      password: 9d0cb9d0-964f-4ef6-897d-6c6b3ccf9e68
  # ws-path: "/" # 如果不为空则开启 websocket 传输层
  # grpc-service-name: "GunService" # 如果不为空则开启 grpc 传输层
  # 下面两项如果填写则开启 tls（需要同时填写）
  certificate: ./server.crt # 证书 PEM 格式，或者 证书的路径
  private-key: ./server.key # 证书对应的私钥 PEM 格式，或者私钥路径
  # 下面两项为mTLS配置项，如果client-auth-type设置为 "verify-if-given" 或 "require-and-verify" 则client-auth-cert必须不为空
  # client-auth-type: "" # 可选值：""、"request"、"require-any"、"verify-if-given"、"require-and-verify"
  # client-auth-cert: string # 证书 PEM 格式，或者 证书的路径
  # 如果填写则开启ech（可由 mihomo generate ech-keypair <明文域名> 生成）
  # ech-key: |
  #   -----BEGIN ECH KEYS-----
  #   ACATwY30o/RKgD6hgeQxwrSiApLaCgU+HKh7B6SUrAHaDwBD/g0APwAAIAAgHjzK
  #   madSJjYQIf9o1N5GXjkW4DEEeb17qMxHdwMdNnwADAABAAEAAQACAAEAAwAIdGVz
  #   dC5jb20AAA==
  #   -----END ECH KEYS-----
  # 如果填写reality-config则开启reality（注意不可与certificate和private-key同时填写）
  # reality-config:
  #   dest: test.com:443
  #   private-key: jNXHt1yRo0vDuchQlIP6Z0ZvjT3KtzVI-T4E7RoLJS0 # 可由 mihomo generate reality-keypair 命令生成
  #   short-id:
  #     - 0123456789abcdef
  #   server-names:
  #     - test.com
  #   #下列两个 limit 为选填，可对未通过验证的回落连接限速，bytesPerSec 默认为 0 即不启用
  #   #回落限速是一种特征，不建议启用，如果您是面板/一键脚本开发者，务必让这些参数随机化
  #   limit-fallback-upload:
  #     after-bytes: 0 # 传输指定字节后开始限速
  #     bytes-per-sec: 0 # 基准速率（字节/秒）
  #     burst-bytes-per-sec: 0 # 突发速率（字节/秒），大于 bytesPerSec 时生效
  #   limit-fallback-download:
  #     after-bytes: 0 # 传输指定字节后开始限速
  #     bytes-per-sec: 0 # 基准速率（字节/秒）
  #     burst-bytes-per-sec: 0 # 突发速率（字节/秒），大于 bytesPerSec 时生效
  # ss-option: # like trojan-go's `shadowsocks` config
  #   enabled: false
  #   method: aes-128-gcm # aes-128-gcm/aes-256-gcm/chacha20-ietf-poly1305
  #   password: "example"
  ### 注意，对于trojan listener, 如果 "allow-insecure" 不为 true, 至少需要填写 “certificate和private-key” 或 “reality-config” 或 “ss-option” 的其中一项 ###
  # allow-insecure: false # 是否允许不开启tls加密（注意：仅用于有 nginx, caddy 前置的情况）
```

```
2026年6月6日
```

---

## AnyTLS 监听器

> 来源: [https://wiki.metacubex.one/config/inbound/listeners/anytls/](https://wiki.metacubex.one/config/inbound/listeners/anytls/)

# AnyTLS

```
listeners:
- name: anytls-in-1
  type: anytls
  port: 10818
  listen: 0.0.0.0
  users:
    username1: password1
    username2: password2
  certificate: ./server.crt # 证书 PEM 格式，或者 证书的路径
  private-key: ./server.key # 证书对应的私钥 PEM 格式，或者私钥路径
  # 下面两项为mTLS配置项，如果client-auth-type设置为 "verify-if-given" 或 "require-and-verify" 则client-auth-cert必须不为空
  # client-auth-type: "" # 可选值：""、"request"、"require-any"、"verify-if-given"、"require-and-verify"
  # client-auth-cert: string # 证书 PEM 格式，或者 证书的路径
  # 如果填写则开启ech（可由 mihomo generate ech-keypair <明文域名> 生成）
  # ech-key: |
  #   -----BEGIN ECH KEYS-----
  #   ACATwY30o/RKgD6hgeQxwrSiApLaCgU+HKh7B6SUrAHaDwBD/g0APwAAIAAgHjzK
  #   madSJjYQIf9o1N5GXjkW4DEEeb17qMxHdwMdNnwADAABAAEAAQACAAEAAwAIdGVz
  #   dC5jb20AAA==
  #   -----END ECH KEYS-----
  ### 注意，anytls listener, 如果 "allow-insecure" 不为 true, 必须填写 “certificate和private-key” ###
  # allow-insecure: false # 是否允许不开启tls加密（注意：仅用于有 nginx, caddy 前置的情况）
  padding-scheme: "" # https://github.com/anytls/anytls-go/blob/main/docs/protocol.md#cmdupdatepaddingscheme
```

通用字段

certificate 和 private-key 是必要的，除非 allow-insecure: true（注意：仅用于有 nginx, caddy 前置的情况）

## padding-scheme

参阅 [https://github.com/anytls/anytls-go/blob/main/docs/protocol.md#cmdupdatepaddingscheme](https://github.com/anytls/anytls-go/blob/main/docs/protocol.md#cmdupdatepaddingscheme)

```
2026年6月6日
```

---

## Mieru 监听器

> 来源: [https://wiki.metacubex.one/config/inbound/listeners/mieru/](https://wiki.metacubex.one/config/inbound/listeners/mieru/)

# Mieru

```
listeners:
- name: mieru-in-1
  type: mieru
  port: 10818
  listen: 0.0.0.0
  transport: TCP # 支持 TCP 或者 UDP
  users:
    username1: password1
    username2: password2
  # 一个 base64 字符串用于微调网络行为
  # traffic-pattern: ""
  # 如果开启，且客户端不发送用户提示，代理服务器将拒绝连接
  # user-hint-is-mandatory: false
```

通用字段

```
2026年4月20日
```

---

## Sudoku 监听器

> 来源: [https://wiki.metacubex.one/config/inbound/listeners/sudoku/](https://wiki.metacubex.one/config/inbound/listeners/sudoku/)

# Sudoku

```
listeners:
- name: sudoku-in-1
  type: sudoku
  port: 8443 # 仅支持单端口
  listen: 0.0.0.0
  key: "<server_key>" # 如果你使用sudoku生成的ED25519密钥对，此处是密钥对中的公钥，当然，你也可以仅仅使用任意uuid充当key
  aead-method: chacha20-poly1305 # 支持chacha20-poly1305或者aes-128-gcm以及none，sudoku的混淆层可以确保none情况下数据安全
  padding-min: 1 # 填充最小长度
  padding-max: 15 # 填充最大长度，均不建议过大
  table-type: prefer_ascii # 可选值：prefer_ascii、prefer_entropy 前者全ascii映射，后者保证熵值（汉明1）低于3
  # custom-table: xpxvvpvv # 可选，自定义字节布局，必须包含2个x、2个p、4个v，可随意组合。启用此处则需配置`table-type`为`prefer_entropy`
  # custom-tables: ["xpxvvpvv", "vxpvxvvp"] # 可选，自定义字节布局列表（x/v/p），用于 xvp 模式轮换；非空时覆盖 custom-table
  handshake-timeout: 5   # optional
  enable-pure-downlink: false # 是否启用混淆下行，false的情况下能在保证数据安全的前提下极大提升下行速度，与客户端保持相同(如果此处为false，则要求aead不可为none)
  httpmask:
    disable: false # true 禁用所有 HTTP 伪装/隧道
    mode: legacy # 可选：legacy（默认）、stream（split-stream）、poll、auto（先 stream 再 poll）、ws（WebSocket 隧道）
    # path_root: "" # 可选：HTTP 隧道端点一级路径前缀（双方需一致），例如 "aabbcc" 或 "/aabbcc/" => /aabbcc/session、/aabbcc/stream、/aabbcc/api/v1/upload、/aabbcc/ws
  #
  # 可选：当启用 HTTPMask 且识别到“像 HTTP 但不符合 tunnel/auth”的请求时，将原始字节透传给 fallback（常用于与其他服务共端口）：
  # fallback: "127.0.0.1:80"
  disable-http-mask: false # 可选：禁用 http 掩码/隧道（默认为 false）
```

通用字段

```
2026年3月9日
```

---

## TUIC v4 监听器

> 来源: [https://wiki.metacubex.one/config/inbound/listeners/tuic-v4/](https://wiki.metacubex.one/config/inbound/listeners/tuic-v4/)

# TUIC V4

```
listeners:
- name: tuicv4-in
  type: tuic
  port: 10003
  listen: 0.0.0.0
  token:
    - TOKEN
  certificate: ./server.crt # 证书 PEM 格式，或者 证书的路径
  private-key: ./server.key # 证书对应的私钥 PEM 格式，或者私钥路径
  # 下面两项为mTLS配置项，如果client-auth-type设置为 "verify-if-given" 或 "require-and-verify" 则client-auth-cert必须不为空
  # client-auth-type: "" # 可选值：""、"request"、"require-any"、"verify-if-given"、"require-and-verify"
  # client-auth-cert: string # 证书 PEM 格式，或者 证书的路径
  # 如果填写则开启ech（可由 mihomo generate ech-keypair <明文域名> 生成）
  # ech-key: |
  #   -----BEGIN ECH KEYS-----
  #   ACATwY30o/RKgD6hgeQxwrSiApLaCgU+HKh7B6SUrAHaDwBD/g0APwAAIAAgHjzK
  #   madSJjYQIf9o1N5GXjkW4DEEeb17qMxHdwMdNnwADAABAAEAAQACAAEAAwAIdGVz
  #   dC5jb20AAA==
  #   -----END ECH KEYS-----
  congestion-controller: bbr
  #  bbr-profile: "" # Available: "standard", "conservative", "aggressive". Default: "standard"
  max-idle-time: 15000
  authentication-timeout: 1000
  alpn:
    - h3
  max-udp-relay-packet-size: 1500
```

```
2026年4月20日
```

---

## TUIC v5 监听器

> 来源: [https://wiki.metacubex.one/config/inbound/listeners/tuic-v5/](https://wiki.metacubex.one/config/inbound/listeners/tuic-v5/)

# TUIC V5

```
listeners:
- name: tuicv5-in
  type: tuic
  port: 10004
  listen: 0.0.0.0
  users:
    UUID1: PASSWORD1
    UUID2: PASSWORD2
  certificate: ./server.crt # 证书 PEM 格式，或者 证书的路径
  private-key: ./server.key # 证书对应的私钥 PEM 格式，或者私钥路径
  # 下面两项为mTLS配置项，如果client-auth-type设置为 "verify-if-given" 或 "require-and-verify" 则client-auth-cert必须不为空
  # client-auth-type: "" # 可选值：""、"request"、"require-any"、"verify-if-given"、"require-and-verify"
  # client-auth-cert: string # 证书 PEM 格式，或者 证书的路径
  # 如果填写则开启ech（可由 mihomo generate ech-keypair <明文域名> 生成）
  # ech-key: |
  #   -----BEGIN ECH KEYS-----
  #   ACATwY30o/RKgD6hgeQxwrSiApLaCgU+HKh7B6SUrAHaDwBD/g0APwAAIAAgHjzK
  #   madSJjYQIf9o1N5GXjkW4DEEeb17qMxHdwMdNnwADAABAAEAAQACAAEAAwAIdGVz
  #   dC5jb20AAA==
  #   -----END ECH KEYS-----
  congestion-controller: bbr
  #  bbr-profile: "" # Available: "standard", "conservative", "aggressive". Default: "standard"
  max-idle-time: 15000
  authentication-timeout: 1000
  alpn:
    - h3
  max-udp-relay-packet-size: 1500
```

```
2026年4月20日
```

---

## Hysteria2 监听器

> 来源: [https://wiki.metacubex.one/config/inbound/listeners/hysteria2/](https://wiki.metacubex.one/config/inbound/listeners/hysteria2/)

# Hysteria2

```
listeners:
- name: hy-in
  type: hysteria2
  port: 8443
  listen: 0.0.0.0
  users:
    user1: password1
    user2: password2
  up: 1000
  down: 1000
  ignore-client-bandwidth: false
  obfs: salamander
  obfs-password: password
  masquerade: ""
  #  bbr-profile: "" # Available: "standard", "conservative", "aggressive". Default: "standard"
  #  realm-opts:
  #    enable: true # 必须手动开启
  #    server-url: https://realm.hy2.io
  #    token: public
  #    realm-id: my-cabin-1f3a8c2e9b
  #    stun-servers:
  #      - stun.nextcloud.com:3478
  #      - stun.sip.us:3478
  #      - global.stun.twilio.com:3478
  #    # proxy: DIRECT # 设置server-url通过哪个代理进行连接
  #    # 下面支持填写针对server-url的TLS配置(sni, skip-cert-verify, fingerprint, certificate, private-key, alpn)
  #    # skip-cert-verify： false
  #    # ......
  alpn:
  - h3
  certificate: ./server.crt # 证书 PEM 格式，或者 证书的路径
  private-key: ./server.key # 证书对应的私钥 PEM 格式，或者私钥路径
  # 下面两项为mTLS配置项，如果client-auth-type设置为 "verify-if-given" 或 "require-and-verify" 则client-auth-cert必须不为空
  # client-auth-type: "" # 可选值：""、"request"、"require-any"、"verify-if-given"、"require-and-verify"
  # client-auth-cert: string # 证书 PEM 格式，或者 证书的路径
  # 如果填写则开启ech（可由 mihomo generate ech-keypair <明文域名> 生成）
  # ech-key: |
  #   -----BEGIN ECH KEYS-----
  #   ACATwY30o/RKgD6hgeQxwrSiApLaCgU+HKh7B6SUrAHaDwBD/g0APwAAIAAgHjzK
  #   madSJjYQIf9o1N5GXjkW4DEEeb17qMxHdwMdNnwADAABAAEAAQACAAEAAwAIdGVz
  #   dC5jb20AAA==
  #   -----END ECH KEYS-----
```

## 通用字段

## 协议配置

### users

Hysteria 用户以及认证密码，格式为用户名: 密码

用户名不参与认证,仅用于入站规则匹配

### up/down

Hysteria 速率设置，默认 Mbps，具体可参考 Hysteria 文档

### ignore-client-bandwidth

可用值：true/false

启用后，服务器将忽略客户端设置的任何带宽，永远使用传统的拥塞控制算法 (BBR)

### obfs

QUIC 流量混淆器，仅可设为 salamander,如果为空则禁用

启用混淆将使服务器与标准的 QUIC 连接不兼容，失去 HTTP/3 伪装的能力

### obfs-password

QUIC 流量混淆器密码

### masquerade

伪装 HTTP/3 流量，仅支持file和http/https,如果为空，则始终返回 404 Not Found

具体可参考 Hysteria 文档

|  | 示例 | 描述 |
| --- | file | file:///var/www | 作为文件服务器 |
| http/https | [http://127.0.0.1:8080](http://127.0.0.1:8080) | 作为反向代理 |

### realm-opts

### alpn

TLS 应用层协议协商列表，客户端需写入服务端 alpn 列表中的一个，否则将认证失败，默认为h3

### certificate/private-key

TLS 证书文件路径

```
2026年5月15日
```

---

## TrustTunnel 监听器

> 来源: [https://wiki.metacubex.one/config/inbound/listeners/trust-tunnel/](https://wiki.metacubex.one/config/inbound/listeners/trust-tunnel/)

# 404 - Not found

---

## Tunnel 监听器

> 来源: [https://wiki.metacubex.one/config/inbound/listeners/tunnel/](https://wiki.metacubex.one/config/inbound/listeners/tunnel/)

# TUNNEL

```
listeners:
- name: tunnel-in
  type: tunnel
  port: 10816
  listen: 0.0.0.0
  network: [tcp, udp]
  target: target.com
```

```
2024年4月1日
```

---

## Snell 监听器

> 来源: [https://wiki.metacubex.one/config/inbound/listeners/snell/](https://wiki.metacubex.one/config/inbound/listeners/snell/)

# Snell

```
listeners:
- name: snell-in-1
  type: snell
  port: 10815 # 支持使用ports格式，例如200,302 or 200,204,401-429,501-503
  listen: 0.0.0.0
  psk: your-password
  version: 4 # 仅支持 4/5
  udp: true # UDP over TCP tunnel，默认 true
  # obfs-opts:
  #   mode: http # 可选：http / tls
  #   host: bing.com
  # rule: sub-rule-name1 # 默认使用 rules，如果未找到 sub-rule 则直接使用 rules
  # proxy: proxy # 如果不为空则直接将该入站流量交由指定 proxy 处理 (当 proxy 不为空时，这里的 proxy 名称必须合法，否则会出错)
```

```
2026年5月31日
```

---

---

# 5. 出站代理

# 出站代理子页面

## 出站代理通用字段

> 来源: [https://wiki.metacubex.one/config/proxies/](https://wiki.metacubex.one/config/proxies/)

# 通用字段

```
proxies:
- name: "ss"
  type: ss
  server: server
  port: 443
  ip-version: ipv4
  udp: true
  interface-name: eth0
  routing-mark: 1234
  tfo: false
  mptcp: false

  dialer-proxy: ss1

  smux:
    enabled: true
    protocol: smux
    max-connections: 4
    min-streams: 4
    max-streams: 0
    statistic: false
    only-tcp: false
    padding: true
    brutal-opts:
      enabled: true
      up: 50
      down: 100
```

代理节点，内容为数组

## name

必须，代理名称，不可重复

## type

必须，代理节点类型

## server

必须，代理节点服务器 (域名/ip)

## port

必须项，代理节点端口

## ip-version

代理软件出站使用的 ip 版本，如果不是 direct，则会影响 server 为域名时使用的 ip 地址

可选：dual/ipv4/ipv6/ipv4-prefer/ipv6-prefer ,默认使用 dual

- ipv4: 仅使用 IPv4
- ipv6: 仅使用 IPv6
- ipv4-prefer: 优先使用 IPv4，对于 TCP 会进行双栈解析，并发链接但是优先使用 IPv4 链接，UDP 则为双栈解析，获取结果中的第一个 IPv4
- ipv6-prefer:优先使用 IPv6，对于 TCP 会进行双栈解析，并发链接但是优先使用 IPv6 链接，UDP 则为双栈解析，获取结果中的第一个 IPv6

## udp

是否允许 UDP 通过代理，默认为 false

> **Note

此选项在 TUIC 等基于 UDP 的协议以及 direct 和 dns 类型中默认开启

## interface-name

指定节点绑定的接口，从此接口发起连接

## routing-mark

节点发起连接时附加的路由标记

## tfo

启用 TCP Fast Open, 仅生效于 TCP 协议

## mptcp

启用 TCP Multi Path, 仅生效于 TCP 协议

## dialer-proxy

指定当前 proxies 通过 dialer-proxy 建立网络连接，值可以为策略组/出站代理的 name，用法可参考dialer-proxy 文档

## smux

sing-mux，仅限使用 tcp 传输的协议

### smux.enabled

是否启用多路复用

### smux.protocol

多路复用协议，支持如下协议，默认使用 h2mux

| 协议 | 描述 |
| smux | [https://github.com/xtaci/smux](https://github.com/xtaci/smux) |
| yamux | [https://github.com/hashicorp/yamux](https://github.com/hashicorp/yamux) |
| h2mux | [https://golang.org/x/net/http2](https://golang.org/x/net/http2) |

### smux.max-connections

最大连接数量

与 max-streams 冲突

### smux.min-streams

在打开新连接之前，连接中的最小多路复用流数量

与 max-streams 冲突

### smux.max-streams

在打开新连接之前，连接中的最大多路复用流数量

与 max-connections 和 min-streams 冲突

### smux.statistic

控制是否将底层连接显示在面板中，方便打断底层连接

### smux.only-tcp

仅允许 tcp，如果设置为 true,smux 的设置将不会对 udp 生效，udp 连接会直接走节点默认 udp 协议传输

### smux.padding

启用填充

### smux.brutal-opts

TCP Brutal 设置

#### smux.brutal-opts.enabled

启用 TCP Brutal 拥塞控制算法

#### smux.brutal-opts.up/down

上传和下载带宽，以默认以 Mbps 为单位

```
2026年5月14日
```

---

## TLS 配置

> 来源: [https://wiki.metacubex.one/config/proxies/tls/](https://wiki.metacubex.one/config/proxies/tls/)

# TLS 配置

```
proxies:
- name: "tls-example"
  tls: true
  sni: example.com
  servername: example.com
  fingerprint: xxx
  alpn:
  - h2
  - http/1.1
  skip-cert-verify: true
  # certificate: xxxx
  # private-key: xxx
  client-fingerprint: chrome
  reality-opts:
    public-key: xxxx
    short-id: xxxx
    support-x25519mlkem768: true
  ech-opts:
    enable: true
    config: base64_encoded_config
    # query-server-name: xxx.com
```

## tls

启用 tls，仅适用于使用 tls 的协议，trojan 协议强制启用

## sni/servername

服务器名称指示，在 VMess/VLESS 中为 servername，如果为空，则为 server 中的地址

## fingerprint

证书指纹，仅适用于使用 tls 的协议，可使用 

```
openssl x509 -noout -fingerprint -sha256 -inform pem -in yourcert.pem
```

获取，也可以通过 Chrome 浏览器的“证书查看器”中“SHA256 指纹”的“证书”项获取

> **Warning

- 

当填入的是叶子证书（即包含 sni 名称的证书）时，只验证服务端发来的证书是否符合该指纹，不会做额外校验

- 

当填入的是其他类型的证书（如中级证书或根证书）的指纹，会验证服务端发来的证书链是否由该证书签发，自 v1.19.20 起还需符合 sni/servername 的要求

- 

此项中的指纹是完整证书的指纹，不是 HPKP 中定义的“证书公钥指纹”，请勿混淆

## alpn

支持的应用层协议协商列表，按优先顺序排列。

如果两个对等点都支持 ALPN，则选择的协议将是此列表中的一个，如果没有相互支持的协议则连接将失败。

参阅 Application-Layer Protocol Negotiation

## skip-cert-verify

跳过证书验证，仅适用于使用 tls 的协议

## certificate

如果填写则开启 mTLS（需要和 private-key 同时填写），内容为证书 PEM 格式，或者 证书的路径

## private-key

如果填写则开启 mTLS（需要和 certificate 同时填写），内容为证书对应的私钥 PEM 格式，或者私钥路径

## client-fingerprint

客户端 utls 指纹，仅适用于 VMess/VLESS/Trojan/AnyTLS 协议

> **Note

可选：chrome, firefox, safari, ios, android, edge, 360, qq, random, 若选择 random, 则按 Cloudflare Radar 数据按概率生成一个现代浏览器指纹。

## reality-opts

reality 配置，如果不为空，则启用 reality

### reality-opts.public-key

reality 服务端私钥对应的公钥

### reality-opts.short-id

服务端 short id 之一

### reality-opts.support-x25519mlkem768

支持 X25519-MLKEM768 密钥交换

## ech-opts

### ech-opts.enable

启用 ECH（Encrypted Client Hello），如果为 true，则启用 ECH

### ech-opts.config

ECH 配置，如果为空则通过 dns 解析，不为空则通过该值指定，格式为经过 base64 编码的 ech 参数（例如AEn+DQBFKwAgACABWIHUGj4u+PIggYXcR5JF0gYk3dCRioBW8uJq9H4mKAAIAAEAAQABAANAEnB1YmxpYy50bHMtZWNoLmRldgAA）

> **Info

您可以通过mihomo generate ech-keypair test.com命令为服务器端和客户端生成符合要求的自签名 ech 配置对，请将test.com自行替换为您想要对外展现的 SNI 域名，输出中Config:后的内容可填在此处，Key:后的内容应填在服务端的 ECH 配置（mihomo 的 listeners 中为ech-key）中

### ech-opts.query-server-name

可选项，不为空时用于指定通过 dns 解析时的域名

```
2026年5月14日
```

---

## 传输层配置

> 来源: [https://wiki.metacubex.one/config/proxies/transport/](https://wiki.metacubex.one/config/proxies/transport/)

# 传输层配置

```
proxies:
- name: "http-opts-example"
  type: xxxxx
  network: http
  http-opts:
    method: "GET"
    path:
    - '/'
    - '/video'
    headers:
      Connection:
      - keep-alive

proxies:
- name: "h2-opts-example"
  type: xxxxx
  network: h2
  h2-opts:
    host:
    - example.com
    path: /

proxies:
- name: "grpc-opts-example"
  type: xxxxx
  network: grpc
  grpc-opts:
    grpc-service-name: example
    # grpc-user-agent: 
    # ping-interval: 0
    # max-connections: 1
    # min-streams: 0
    # max-streams: 0

proxies:
- name: "ws-opts-example"
  type: xxxxx
  network: ws
  ws-opts:
    path: /path
    headers:
      Host: example.com
    max-early-data:
    early-data-header-name:
    v2ray-http-upgrade: false
    v2ray-http-upgrade-fast-open: false

proxies:
- name: "xhttp-opts-example"
  type: vless
  server: server
  port: 443
  uuid: uuid
  udp: true
  tls: true
  network: xhttp
  alpn: [h2] # By default, only h2 mode is supported. To enable h3 mode, you need to set alpn: [h3]; to enable HTTP/1.1 mode, you need to set alpn: [http/1.1].
  # ech-opts: ...
  # reality-opts: ...
  # skip-cert-verify: false
  # fingerprint: ...
  # certificate: ...
  # private-key: ...
  servername: xxx.com
  client-fingerprint: chrome
  encryption: ""
  xhttp-opts:
    path: "/"
    host: xxx.com
    # mode: "stream-one" # Available: "stream-one", "stream-up" or "packet-up"
    # headers:
    #   X-Forwarded-For: ""
    # no-grpc-header: false
    # x-padding-bytes: "100-1000"
    # x-padding-obfs-mode: false
    # x-padding-key: x_padding
    # x-padding-header: Referer
    # x-padding-placement: queryInHeader # Available: queryInHeader, cookie, header, query
    # x-padding-method: repeat-x # Available: repeat-x, tokenish
    # uplink-http-method: POST # Available: POST, PUT, PATCH, DELETE
    # session-placement: path # Available: path, query, cookie, header
    # session-key: ""
    # seq-placement: path # Available: path, query, cookie, header
    # seq-key: ""
    # uplink-data-placement: body # Available: body, cookie, header
    # uplink-data-key: ""
    # uplink-chunk-size: 0 # only applicable when uplink-data-placement is not body
    # sc-max-each-post-bytes: 1000000
    # sc-min-posts-interval-ms: 30
    # reuse-settings: # aka XMUX
    #   max-concurrency: "16-32"
    #   max-connections: "0"
    #   c-max-reuse-times: "0"
    #   h-max-request-times: "600-900"
    #   h-max-reusable-secs: "1800-3000"
    #   h-keep-alive-period: 0
    # download-settings:
    #   ## xhttp part
    #   path: "/"
    #   host: xxx.com
    #   headers:
    #     X-Forwarded-For: ""
    #   reuse-settings: # aka XMUX
    #     max-concurrency: "16-32"
    #     max-connections: "0"
    #     c-max-reuse-times: "0"
    #     h-max-request-times: "600-900"
    #     h-max-reusable-secs: "1800-3000"
    #     h-keep-alive-period: 0
    #   ## proxy part
    #   server: server
    #   port: 443
    #   tls: true
    #   alpn: ...
    #   ech-opts: ...
    #   reality-opts: ...
    #   skip-cert-verify: false
    #   fingerprint: ...
    #   certificate: ...
    #   private-key: ...
    #   servername: xxx.com
    #   client-fingerprint: chrome
```

## http-opts

http 传输层设置，仅传输层为 http 时生效

### http-opts.method

http 请求方法

### http-opts.path

http 请求路径

### http-opts.headers

http 请求头

> **Note

Mihomo 的 H2 传输层未实现多路复用功能，如果需要使用多路复用，在 Mihomo 中更建议使用 gRPC 协议，或者 sing-mux

## h2-opts

h2 传输层设置，仅传输层为 h2 时生效

### h2-opts.host

主机域名列表，如果设置，客户端将随机选择，服务器将验证

### h2-opts.path

http 请求路径

## grpc-opts

grpc 传输层设置，仅传输层为 grpc 时生效

### grpc-opts.grpc-service-name

gRPC 服务名称

### grpc-opts.grpc-user-agent

gRPC UserAgent

### grpc-opts.ping-interval

gRPC 心跳包间隔，默认关闭，单位为秒

### grpc-opts.max-connections

最大连接数量。默认值为 1 即只使用一条底层链接

与 max-streams 冲突

### grpc-opts.min-streams

在打开新连接之前，连接中的最小多路复用流数量

与 max-streams 冲突

### grpc-opts.max-streams

在打开新连接之前，连接中的最大多路复用流数量

与 max-connections 和 min-streams 冲突

## ws-opts

ws 传输层设置，仅传输层为 ws 时生效

### ws-opts.path

请求路径

### ws-opts.headers

请求头

### ws-opts.max-early-data

Early Data 首包长度阈值

### ws-opts.early-data-header-name

### ws-opts.v2ray-http-upgrade

使用 http upgrade

### ws-opts.v2ray-http-upgrade-fast-open

启用 http upgrade 的 fast open

## xhttp-opts

xhttp 传输层设置，仅传输层为 xhttp 时生效

默认仅支持 h2，如果开启 h3 模式需要设置alpn: [h3]，如果开启 http1.1 模式需要设置alpn: [http/1.1]

> **Note

仅 VLESS 支持 xhttp 传输层，请勿在其他协议上使用

### xhttp-opts.path

请求路径

### xhttp-opts.host

主机名

### xhttp-opts.mode

模式，可选值：auto, stream-one, stream-up or packet-up

### xhttp-opts.headers

请求头

### xhttp-opts.no-grpc-header

设置 stream-up/one 上行是否发送Content-Type: application/grpc头伪装成 gRPC

### xhttp-opts.x-padding-bytes

客户端发的 request header 均默认包含 Referer: ...?x_padding=XXX... ，默认长度为 100-1000，每次请求随机，服务端默认会检查它是否在服务端允许的范围内

### xhttp-opts.x-padding-obfs-mode

启用填充混淆。为向后兼容，默认为 false

### xhttp-opts.x-padding-key

用于存储填充值的键名。其含义取决于 x-padding-placement：

- URL 中查询参数的名称（如果 placement 为 queryInHeader）
- cookie 的名称
- HTTP 标头的名称
- URL 查询参数的名称

### xhttp-opts.x-padding-header

HTTP 标头的名称。仅当 x-padding-placement 为 header 或 queryInHeader 时才相关

### xhttp-opts.x-padding-placement

定义填充值的放置位置。可选值：queryInHeader、cookie、header、query。仅当 x-padding-obfs-mode 为 true 时才有效

### xhttp-opts.x-padding-method

定义填充值的生成方式

- repeat-x：默认方法（一个包含 X 个字符的长序列）
- tokenish：生成一个随机的 Base62 标记

### xhttp-opts.uplink-http-method

更改用于上行数据传输的 HTTP 方法。使用任何支持请求体（例如 PUT、PATCH）且 CDN 允许的方法

### xhttp-opts.session-placement

会话 ID 的放置位置。选项：path, query, cookie, header

### xhttp-opts.session-key

会话 ID 的键名（如果放置位置为path，则不适用）

### xhttp-opts.seq-placement

序列号的放置位置。选项：path、query、cookie、header。如果session-placement设置为path，则seq-placement也必须设置为path

### xhttp-opts.seq-key

序列号的键名（如果放置位置为path，则不适用）

### xhttp-opts.uplink-data-placement

将拆分后的上行链路数据片段放置位置（仅适用于packet-up模式）。选项：cookie 或 header

### xhttp-opts.uplink-data-key

用于传递数据片段的键的基本名称。客户端会自动将数据分割成多个数据块，服务器会将它们重新组装起来

### xhttp-opts.uplink-chunk-size

每个上行链路数据块的最大大小（以字节为单位）（仅当 uplink-data-placement 不是 body 时适用）

最小值：64 字节

推荐值：

- 对于cookie放置：3072 字节（默认 3KB）
- 对于header放置：4096 字节（默认 4KB）

如果未指定，则会根据 uplink-data-placement 自动选择合适的默认值

### xhttp-opts.sc-max-each-post-bytes

客户端每个 POST 最多携带多少数据，默认值 1000000 即 1MB，该值应小于 CDN 等 HTTP 中间盒所允许的最大值，服务端也会拒绝大于该值的 POST

### xhttp-opts.sc-min-posts-interval-ms

基于单个代理请求，客户端发起 POST 请求的最小间隔，默认值 30 毫秒

### xhttp-opts.reuse-settings

链接复用设置（即 XMUX）

注意：和原版实现不同，此项没有默认值，如果不填写则不开启链接复用，即每次打开一个新的底层链接

### xhttp-opts.reuse-settings.max-concurrency

每条底层连接中最多同时存在多少代理请求，连接中的代理请求数量达到该值后会建立新的连接，以容纳更多的代理请求

### xhttp-opts.reuse-settings.max-connections

最多同时存在多少条连接，连接数达到该值前每个新的代理请求都会开启一条新的连接，此后会开始复用已有的连接

该值与 max-concurrency 冲突，只能二选一

### xhttp-opts.reuse-settings.c-max-reuse-times

一条连接最多被复用几次，复用该次数后将不会被分配新的代理请求，将在内部最后一条代理请求关闭后断开

### xhttp-opts.reuse-settings.h-max-request-times

一条连接最多累计承载多少次，该项计数不严谨，且 Golang 的 GET 请求有自动重试，故不建议填写

### xhttp-opts.reuse-settings.h-max-reusable-secs

TCP/QUIC 连接持续该时间后将不会被分配新的 HTTP 请求，将在内部最后一个 HTTP 请求关闭后断开

### xhttp-opts.reuse-settings.h-keep-alive-period

H2 / H3 连接空闲时客户端每隔多少秒发一次保活包，默认 0，即 Chrome H2 的 45 秒，或 quic-go H3 的 10 秒。它是 XMUX 里唯一不允许填范围（该项取随机才是特征）且允许填负数（比如填 -1 关掉空闲保活包）的项，建议留 0。

### xhttp-opts.download-settings

上下行分离设置

注意：此项用于覆盖原始配置，每一项如果不填写则会沿用上行参数（仅支持覆盖实例中列出的选项，未列出的不会被覆写）

```
2026年5月14日
```

---

## Dialer Proxy

> 来源: [https://wiki.metacubex.one/config/proxies/dialer-proxy/](https://wiki.metacubex.one/config/proxies/dialer-proxy/)

# dialer-proxy

```
proxies:
- name: "ss1"
  dialer-proxy: dialer
  ...

- name: "ss2"
  ...

proxy-groups:
- name: dialer
  type: select
  proxies:
  - ss2

rules:
  - MATCH,ss1
```

指定当前 proxies 通过 dialer-proxy 建立网络连接，值可以为策略组/出站代理的 name

上述示例中，ss1 通过 ss2 建立连接

当通过 ss1 代理时，就组成了一条 內核 ---ss1--> ss2包裝器 ===ss2==> ss2服務端 ---ss1--> ss1服務端 --> 目标 的代理链

最终表现：

- 浏览器访问 IP 查询网站时只会显示 ss1 的 IP（即目标网站不知道 ss2 的存在）
- 内核对外看起来只是在访问 ss2（即你的宽带运营商不知道 ss1 的存在）
- 只有 ss2 的服务端知道在访问 ss1（ss2 的服务端也只知道在访问 ss1，并不知道在访问什么目标网站）

## 常见实例

### 通过订阅节点中转自己的 VPS 落地

```
proxies:
- name: "ss1"
  dialer-proxy: dialer
  ...

proxy-providers:
  provider1:
    type: http
    url: "http://test.com"

proxy-groups:
- name: dialer
  type: select
  use:
  - provider1

rules:
  - MATCH,ss1
```

这里将订阅地址填入 provider1 中，将你自己 VPS 中搭建的节点填入 ss1 中即可，此时通过浏览器访问时显示的是 ss1 的 IP

> **Note

没有特殊需求的情况下，在自己被中转的 VPS 落地中搭建的节点请勿选择任何 udp 类协议如 hy2/tuic/wg，以及带有 tls 伪装类协议如 reality/shadowtls，您的订阅节点可能不能正常通过这些协议，这里建议选择最简单的 ss aead 或者 vmess 协议

### 通过特定 socks 连接订阅节点

```
proxies:
- { name: "socks1", type: "socks" }

proxy-providers:
  provider1:
    type: http
    url: "http://test.com"
    override:
      dialer-proxy: socks1

proxy-groups:
- name: select1
  type: select
  use:
  - provider1

rules:
  - MATCH,select1
```

该实例适用于需要通过一个特定 socks 才能访问外网的情况（如内外网隔离环境），将环境提供的 socks 配置填入 socks1，将订阅地址填入 provider1 中，此时通过浏览器访问时显示的是订阅中节点的 IP

> **Note

这里只是示范了 dialer-proxy 的相关配置，您可能还需要更多的配置才能保证订阅下载，dns 解析等流程同样通过该 socks

## relay 迁移

relay 类型的 proxy-group 已经被废弃，而 proxy-group 并不直接支持 dialer-proxy，因此针对部分使用场景，给出参考方案

### relay 中包含多个 select

有如下配置

```
proxies:
  - { name: "proxy1", type: "socks" }
  - { name: "proxy2", type: "socks" }
  - { name: "proxy3", type: "socks" }
  - { name: "proxy4", type: "socks" }
proxy-groups:
  - { name: "relay-proxy", type: relay, proxies: ["select1", "select2"] }
  - { name: "select1", type: select, proxies: ["proxy1", "proxy2"] }
  - { name: "select2", type: select, proxies: ["proxy3", "proxy4"] }
rules:
  - MATCH,relay-proxy
```

迁移到 dialer-proxy 方案时，需要将 proxy3 和 proxy4 定义到 proxy-provider 中，并通过 override 设置该 provider 中所有 proxy 的 dialer-proxy，如下所示：

```
proxies:
  - { name: "proxy1", type: "socks" }
  - { name: "proxy2", type: "socks" }
proxy-groups:
  - { name: "select1", type: select, proxies: ["proxy1", "proxy2"] }
  - { name: "select2", type: select, use: ["provider1"] }
proxy-providers:
  provider1:
    type: inline
    override:
      dialer-proxy: select1
    payload:
      - { name: "proxy3", type: "socks" }
      - { name: "proxy4", type: "socks" }
rules:
  - MATCH,select2
```

```
2026年5月14日
```

---

## 内置代理策略

> 来源: [https://wiki.metacubex.one/config/proxies/built-in/](https://wiki.metacubex.one/config/proxies/built-in/)

# 预置出站

## DIRECT

直连，数据直接出站

## REJECT

拒绝，拦截数据出站

## REJECT-DROP

拒绝，与REJECT不同的是，该策略将静默抛弃请求

## PASS

绕过。匹配到 PASS 时，会将其视为已命中出站并跳过当前命中的规则分支，继续匹配后续规则。

若在 SUB-RULE 中命中 PASS，会跳出当前子规则并回到主规则继续匹配，不会继续匹配同一 sub-rules 内的后续规则。

## PASS-RULE

绕过。与PASS类似，区别是：在 SUB-RULE 中命中 PASS-RULE 时不会跳出子规则，仅跳过当前命中规则分支，继续在 sub-rules 中向后匹配规则。

## COMPATIBLE

兼容，在策略组筛选不出节点时出现，等效 DIRECT

```
2026年6月6日
```

---

## DIRECT

> 来源: [https://wiki.metacubex.one/config/proxies/direct/](https://wiki.metacubex.one/config/proxies/direct/)

# DIRECT

```
proxies:
- name: "direct"
  type: direct
  udp: true
  ip-version: ipv4
  interface-name: eth0
  routing-mark: 1234
```

通用字段

```
2026年5月14日
```

---

## DNS 代理

> 来源: [https://wiki.metacubex.one/config/proxies/dns/](https://wiki.metacubex.one/config/proxies/dns/)

# DNS

```
proxies:
- name: "dns-out"
  type: dns
```

dns出站会将请求劫持到内部dns模块，所有请求均在内部处理

```
2026年5月14日
```

---

## HTTP 代理

> 来源: [https://wiki.metacubex.one/config/proxies/http/](https://wiki.metacubex.one/config/proxies/http/)

# HTTP

```
proxies:
- name: "http"
  type: http
  server: server
  port: 443
  # username: username
  # password: password
  # tls: true # https
  # skip-cert-verify: true
  # sni: custom.com
  # fingerprint: xxxx # 同 experimental.fingerprints 使用 sha256 指纹,配置协议独立的指纹,将忽略 experimental.fingerprints
  # ip-version: dual
  headers:
```

通用字段

TLS 字段

```
2026年5月14日
```

---

## SOCKS 代理

> 来源: [https://wiki.metacubex.one/config/proxies/socks/](https://wiki.metacubex.one/config/proxies/socks/)

# SOCKS

```
proxies:
- name: "socks"
  type: socks5
  server: server
  port: 443
  # username: username
  # password: password
  # tls: true
  # fingerprint: xxxx
  # skip-cert-verify: true
  # udp: true
  # ip-version: ipv6
```

通用字段

TLS 字段

```
2026年5月14日
```

---

## Shadowsocks

> 来源: [https://wiki.metacubex.one/config/proxies/ss/](https://wiki.metacubex.one/config/proxies/ss/)

# Shadowsocks

```
proxies:
- name: "ss1"
  type: ss
  server: server
  port: 443
  cipher: aes-128-gcm
  password: "password"
  udp: true
  udp-over-tcp: false
  udp-over-tcp-version: 2
  ip-version: ipv4
  plugin: obfs
  plugin-opts:
    mode: tls
  smux:
    enabled: false
```

通用字段

## Cipher

| 方法 | --- | aes-128-ctr | aes-192-ctr | aes-256-ctr |
| aes-128-cfb | aes-192-cfb | aes-256-cfb |
| aes-128-gcm | aes-192-gcm | aes-256-gcm |
| aes-128-ccm | aes-192-ccm | aes-256-ccm |
| aes-128-gcm-siv |  | aes-256-gcm-siv |

| 方法 |  |
| chacha20-ietf |  |
| chacha20 | xchacha20 |
| chacha20-ietf-poly1305 | xchacha20-ietf-poly1305 |
| chacha8-ietf-poly1305 | xchacha8-ietf-poly1305 |


| 方法                            |
| ----------------------------- |
| 2022-blake3-aes-128-gcm       |
| 2022-blake3-aes-256-gcm       |
| 2022-blake3-chacha20-poly1305 |



| 方法          |
| ----------- |
| lea-128-gcm |
| lea-192-gcm |
| lea-256-gcm |



| 方法                 |
| ------------------ |
| rabbit128-poly1305 |
| aegis-128l         |
| aegis-256          |
| aez-384            |
| deoxys-ii-256-128  |
| rc4-md5            |
| none               |


## password

Shadowsocks 密码

## udp-over-tcp

启用 UDP over TCP，默认 false

## udp-over-tcp-version

UDP over TCP 的协议版本，默认 1。可选值 1/2。

## 插件

### plugin

插件，支持 obfs/v2ray-plugin/gost-plugin/shadow-tls/restls/kcptun

### plugin-opts

插件设置

```
plugin: obfs
plugin-opts:
  mode: tls
  host: bing.com

plugin: v2ray-plugin
  plugin-opts:
    mode: websocket # no QUIC now
      # tls: true # wss
      # 可使用 openssl x509 -noout -fingerprint -sha256 -inform pem -in yourcert.pem 获取
      # 配置指纹将实现 SSL Pining 效果
      # fingerprint: xxxx
      # skip-cert-verify: true
      # host: bing.com
      # path: "/"
      # mux: true
      # headers:
      #   custom: value
      # v2ray-http-upgrade: false

plugin: gost-plugin
  plugin-opts:
    mode: websocket
      # tls: true # wss
      # 可使用 openssl x509 -noout -fingerprint -sha256 -inform pem -in yourcert.pem 获取
      # 配置指纹将实现 SSL Pining 效果
      # fingerprint: xxxx
      # skip-cert-verify: true
      # host: bing.com
      # path: "/"
      # mux: true
      # headers:
      #   custom: value

plugin: shadow-tls
  client-fingerprint: chrome
  plugin-opts:
    host: "cloud.tencent.com"
    password: "shadow_tls_password"
    version: 2 # support 1/2/3

plugin: restls
  client-fingerprint: chrome  # 可以是chrome, ios, firefox, safari中的一个
  plugin-opts:
    host: "www.microsoft.com" # 应当是一个TLS 1.3 服务器
    password: [YOUR_RESTLS_PASSWORD]
    version-hint: "tls13"
    # Control your post-handshake traffic through restls-script
    # Hide proxy behaviors like "tls in tls".
    # see https://github.com/3andne/restls/blob/main/Restls-Script:%20Hide%20Your%20Proxy%20Traffic%20Behavior.md
    # 用restls剧本来控制握手后的行为，隐藏"tls in tls"等特征
    # 详情：https://github.com/3andne/restls/blob/main/Restls-Script:%20%E9%9A%90%E8%97%8F%E4%BD%A0%E7%9A%84%E4%BB%A3%E7%90%86%E8%A1%8C%E4%B8%BA.md
    restls-script: "300?100<1,400~100,350~100,600~100,300~200,300~100"

plugin: kcptun
  plugin-opts:
    key: it's a secrect # pre-shared secret between client and server
    crypt: aes # aes, aes-128, aes-128-gcm, aes-192, salsa20, blowfish, twofish, cast5, 3des, tea, xtea, xor, none, null
    mode: fast # profiles: fast3, fast2, fast, normal, manual
    conn: 1 # set num of UDP connections to server
    autoexpire: 0 # set auto expiration time(in seconds) for a single UDP connection, 0 to disable
    scavengettl: 600 # set how long an expired connection can live (in seconds)
    mtu: 1350 # set maximum transmission unit for UDP packets
    ratelimit: 0 # set maximum outgoing speed (in bytes per second) for a single KCP connection, 0 to disable. Also known as packet pacing
    sndwnd: 128 # set send window size(num of packets)
    rcvwnd: 512 # set receive window size(num of packets)
    datashard: 10 # set reed-solomon erasure coding - datashard
    parityshard: 3 # set reed-solomon erasure coding - parityshard
    dscp: 0 # set DSCP(6bit)
    nocomp: false # disable compression
    acknodelay: false # flush ack immediately when a packet is received
    nodelay: 0
    interval: 50
    resend: 0
    sockbuf: 4194304 # per-socket buffer in bytes
    smuxver: 1 # specify smux version, available 1,2
    smuxbuf: 4194304 # the overall de-mux buffer in bytes
    framesize: 8192 # smux max frame size
    streambuf: 2097152 # per stream receive buffer in bytes, smux v2+
    keepalive: 10 # seconds between heartbeats
```

```
2026年1月17日
```

---

## ShadowsocksR

> 来源: [https://wiki.metacubex.one/config/proxies/ssr/](https://wiki.metacubex.one/config/proxies/ssr/)

# ShadowsocksR

```
proxies:
  - name: "ssr"
    type: ssr
    server: server
    port: 443
    cipher: chacha20-ietf
    password: "password"
    obfs: tls1.2_ticket_auth
    protocol: auth_sha1_v4
    # obfs-param: domain.tld
    # protocol-param: "#"
    # udp: true
```

通用字段

```
2026年5月14日
```

---

## Snell

> 来源: [https://wiki.metacubex.one/config/proxies/snell/](https://wiki.metacubex.one/config/proxies/snell/)

# Snell

```
proxies
- name: "snell"
  type: snell
  server: server
  port: 44046
  psk: yourpsk
  # version: 4
  # udp: true
  # reuse: false
  # obfs-opts:
  #   mode: http
  #   host: bing.com
```

通用字段

## psk

必须，Snell 预共享密钥

## version

snell 版本，支持 v1/2/3/4/5，仅 v3/4/5 支持 udp

## reuse

可选，支持 v4/5，默认为 false

## obfs-opts

Snell 混淆设置

### obfs-opts.mode

Snell 混淆模式，支持 http/tls

### obfs-opts.host

Snell 混淆域名

```
2026年5月31日
```

---

## VMess

> 来源: [https://wiki.metacubex.one/config/proxies/vmess/](https://wiki.metacubex.one/config/proxies/vmess/)

# VMess

```
proxies:
- name: "vmess"
  type: vmess
  server: server
  port: 443
  udp: true
  uuid: uuid
  alterId: 0
  cipher: auto
  packet-encoding: packetaddr
  global-padding: false
  authenticated-length: false

  tls: true
  servername: example.com
  alpn:
  - h2
  - http/1.1
  fingerprint: xxxx
  client-fingerprint: chrome
  skip-cert-verify: true
  reality-opts:
    public-key: xxxx
    short-id: xxxx

  network: tcp

  smux:
    enabled: false
```

通用字段

TLS 字段

## uuid

必须，VMess 用户 ID

## alterId

必须，如果不为 0，则启用旧协议

## cipher

必须，加密方法，支持 auto/none/zero/aes-128-gcm/chacha20-poly1305

## packet-encoding

UDP 包编码，为空则使用原始编码，可选 packetaddr (由 v2ray 5+ 支持)/xudp (由 xray 支持)

## global-padding

协议参数。如果启用会随机浪费流量（在 v2ray 中默认启用并且无法禁用）。

## authenticated-length

协议参数。启用长度块加密

## network

传输层，支持 ws/http/h2/grpc，不配置或配置其他值则为 tcp

参阅 传输层配置

```
2026年5月14日
```

---

## VLESS

> 来源: [https://wiki.metacubex.one/config/proxies/vless/](https://wiki.metacubex.one/config/proxies/vless/)

# VLESS

```
proxies:
- name: "vless"
  type: vless
  server: server
  port: 443
  udp: true
  uuid: uuid
  flow: xtls-rprx-vision
  packet-encoding: xudp

  tls: true
  servername: example.com
  alpn:
  - h2
  - http/1.1
  fingerprint: xxxx
  client-fingerprint: chrome
  skip-cert-verify: true
  reality-opts:
    public-key: xxxx
    short-id: xxxx
  encryption: ""

  network: tcp

  smux:
    enabled: false
```

> **Note

Meta 的 xtls-* 流控实际上与 Xray-core 中的 xtls-*-udp443 等效，如需拦截 443 端口的 UDP 流量，请使用逻辑规则：AND,((NETWORK,UDP),(DST-PORT,443)),REJECT

通用字段

TLS 字段

## uuid

必须，VLESS 用户 ID

## flow

VLESS 子协议，可用值为 xtls-rprx-vision

## packet-encoding

UDP 包编码，为空则使用原始编码，可选 packetaddr (由 v2ray 5+ 支持)/ xudp (由 xray 支持)

## encryption

vless encryption 客户端配置：

encryption: "mlkem768x25519plus.native/xorpub/random.1rtt/0rtt.(padding len).(padding gap).(X25519 Password).(ML-KEM-768 Client)..."

（native/xorpub 的 XTLS Vision 可以 Splice。只使用 1-RTT 模式 / 若服务端发的 ticket 中秒数不为零则 0-RTT 复用）

/ 是只能选一个，后面 base64 至少一个，无限串联，使用  mihomo generate vless-x25519 和 mihomo generate vless-mlkem768 生成，替换值时需去掉括号

Padding 是可选的参数，仅作用于 1-RTT 以消除握手的长度特征，双端默认值均为 "100-111-1111.75-0-111.50-0-3333"：

- 

在 1-RTT client/server hello 后以 100% 的概率粘上随机 111 到 1111 字节的 padding

- 

以 75% 的概率等待随机 0 到 111 毫秒（"probability-from-to"）

- 

再次以 50% 的概率发送随机 0 到 3333 字节的 padding（若为 0 则不 Write()）

- 

服务端、客户端可以设置不同的 padding 参数，按 len、gap 的顺序无限串联，第一个 padding 需概率 100%、至少 35 字节

## network

传输层，支持 ws/http/h2/grpc/xhttp，不配置或配置其他值则为 tcp

参阅 传输层配置

```
2026年5月14日
```

---

## Trojan

> 来源: [https://wiki.metacubex.one/config/proxies/trojan/](https://wiki.metacubex.one/config/proxies/trojan/)

# Trojan

```
proxies:
- name: "trojan"
  type: trojan
  server: server
  port: 443
  password: yourpsk
  udp: true

  sni: example.com
  alpn:
  - h2
  - http/1.1
  client-fingerprint: random
  fingerprint: xxxx
  skip-cert-verify: true
  ss-opts:
    enabled: false
    method: aes-128-gcm
    password: "example"
  reality-opts:
    public-key: xxxx
    short-id: xxxx

  network: tcp

  smux:
    enabled: false
```

通用字段

TLS 字段

## password

必须，trojan 服务器密码

## ss-opts

### ss-opfs.enabled

启用 trojan-go 的 shadowsocks AEAD 加密

### ss-opfs.method

加密方法，支持 aes-128-gcm/aes-256-gcm/chacha20-ietf-poly1305

### ss-opfs.password

trojan-go 的 shadowsocks AEAD 加密密码

## network

传输层，支持 ws/grpc，不配置或配置其他值则为 tcp

参阅 传输层配置

```
2026年5月14日
```

---

## AnyTLS

> 来源: [https://wiki.metacubex.one/config/proxies/anytls/](https://wiki.metacubex.one/config/proxies/anytls/)

# AnyTLS

```
proxies:
- name: anytls
  type: anytls
  server: 1.2.3.4
  port: 443
  password: "<your password>"
  client-fingerprint: chrome
  udp: true
  idle-session-check-interval: 30
  idle-session-timeout: 30
  min-idle-session: 0
  sni: "example.com"
  alpn:
    - h2
    - http/1.1
  skip-cert-verify: true
```

通用字段

TLS 字段

> **Tip

Mihomo 不支持 AnyTLS+Reality 的组合（未来也不会支持），如果您想隐藏 SNI 请配合ECH使用，如果您非要使用 Reality 请选择Vmess、VLESS、Trojan协议。

## idle-session-check-interval

检查空闲会话的时间间隔。默认值：30 秒。

## idle-session-timeout

在检查中，关闭闲置时间超过此值的会话。默认值：30 秒。

## min-idle-session

在检查中，至少前 n 个空闲会话保持打开状态。默认值：n=0

```
2026年5月14日
```

---

## Mieru

> 来源: [https://wiki.metacubex.one/config/proxies/mieru/](https://wiki.metacubex.one/config/proxies/mieru/)

# Mieru

```
proxies:
  - name: mieru
    type: mieru
    server: server
    port: 2999
    port-range: 2090-2099
    transport: TCP
    username: user
    password: password
    multiplexing: MULTIPLEXING_LOW
    traffic-pattern: ""
```

通用字段

## port-range

端口范围，不可与 port 同时书写

## transport

协议，目前支持 TCP 和 UDP

## multiplexing

多路复用，可以使用的值包括 MULTIPLEXING_OFF, MULTIPLEXING_LOW, MULTIPLEXING_MIDDLE, MULTIPLEXING_HIGH。其中 MULTIPLEXING_OFF 会关闭多路复用功能。默认值为 MULTIPLEXING_LOW

## traffic-pattern

一个 base64 字符串用于微调网络行为，格式请参考mieru 官方文档

```
2026年5月14日
```

---

## Sudoku

> 来源: [https://wiki.metacubex.one/config/proxies/sudoku/](https://wiki.metacubex.one/config/proxies/sudoku/)

# Sudoku

```
proxies:
  - name: sudoku
    type: sudoku
    server: 1.2.3.4
    port: 443 
    key: "<client_key>"
    aead-method: chacha20-poly1305
    padding-min: 2
    padding-max: 7
    table-type: prefer_ascii
    # custom-table: xpxvvpvv
    # custom-tables: ["xpxvvpvv", "vxpvxvvp"]
    httpmask:
      disable: false
      mode: legacy
      tls: true
      mask-host: ""
      path-root: ""
      multiplex: off
    enable-pure-downlink: false
```

通用字段

## key

如果你使用 sudoku 生成的 ED25519 密钥对，请填写密钥对中的私钥，否则填入和服务端相同的 uuid

## aead-method

可选值：chacha20-poly1305、aes-128-gcm、none 我们保证在 none 的情况下 sudoku 混淆层仍然确保安全

## padding-min

最小填充字节数

## padding-max

最大填充字节数

## table-type

可选值：prefer_ascii、prefer_entropy、up_ascii_down_entropy、up_entropy_down_ascii

## custom-table

可选，自定义字节布局，必须包含 2 个 x、2 个 p、4 个 v，可随意组合；只对 entropy 方向生效

## custom-tables

可选，自定义字节布局列表（x/v/p），非空时覆盖 custom-table

## httpmask.disable

是否禁用所有 HTTP 伪装/隧道

## httpmask.mode

可选：legacy（默认）、stream、poll、auto、ws；stream/poll/auto/ws 支持走 CDN/反代

## httpmask.tls

可选：仅在 mode 为 stream/poll/auto/ws 时生效；true 强制 https；false 强制 http（不会根据端口自动推断）

## httpmask.host

可选：覆盖 Host/SNI（支持 example.com 或 example.com:443）；仅在 mode 为 stream/poll/auto/ws 时生效

## httpmask.path-root

可选：HTTP 隧道端点一级路径前缀（双方需一致），例如 "aabbcc" => /aabbcc/session、/aabbcc/stream、/aabbcc/api/v1/upload

## httpmask.multiplex

可选：off（默认）、auto（复用 h1.1 keep-alive / h2 连接，减少每次建链 RTT）、on（单条隧道内多路复用多个目标连接；仅在 mode=stream/poll/auto 生效；ws 强制 off）

## enable-pure-downlink

是否启用混淆下行，false 的情况下能在保证数据安全的前提下极大提升下行速度，与服务端端保持相同 (如果此处为 false，则要求 aead 不可为 none)

```
2026年5月14日
```

---

## Hysteria

> 来源: [https://wiki.metacubex.one/config/proxies/hysteria/](https://wiki.metacubex.one/config/proxies/hysteria/)

# Hysteria

```
proxies:
- name: "hysteria"
  type: hysteria
  server: server.com
  port: 443
  # ports: 1000,2000-3000,4000 # port 不可省略
  auth-str: yourpassword
  # obfs: obfs_str
  # alpn:
  #   - h3
  protocol: udp # 支持 udp/wechat-video/faketcp
  up: "30 Mbps" # 若不写单位,默认为 Mbps
  down: "200 Mbps" # 若不写单位,默认为 Mbps
  # sni: server.com
  # skip-cert-verify: false
  # recv-window-conn: 12582912
  # recv-window: 52428800
  # disable_mtu_discovery: false
  # fingerprint: xxxx # 配置指纹将实现 SSL Pining 效果, 可使用 openssl x509 -noout -fingerprint -sha256 -inform pem -in yourcert.pem 获取
  # fast-open: true # 启用 Fast Open (降低连接建立延迟),默认为 false
```

通用字段

TLS 字段

```
2025年9月19日
```

---

## Hysteria2

> 来源: [https://wiki.metacubex.one/config/proxies/hysteria2/](https://wiki.metacubex.one/config/proxies/hysteria2/)

# Hysteria2

配置参考

```
proxies:
- name: "hysteria2"
  type: hysteria2
  server: server.com
  port: 443
  ports: 443-8443
  hop-interval: 30
  password: yourpassword
  up: "30 Mbps"
  down: "200 Mbps"
  # bbr-profile: "" # Available: "standard", "conservative", "aggressive". Default: "standard"
  obfs: salamander # 默认为空，如果填写则开启obfs，目前仅支持salamander
  obfs-password: yourpassword

  sni: server.com
  skip-cert-verify: false
  fingerprint: xxxx # 配置指纹将实现 SSL Pining 效果, 可使用 openssl x509 -noout -fingerprint -sha256 -inform pem -in yourcert.pem 获取
  alpn:
    - h3
  # realm-opts:
  #   enable: true # 必须手动开启
  #   server-url: https://realm.hy2.io
  #   token: public
  #   realm-id: my-cabin-1f3a8c2e9b
  #   stun-servers:
  #     - stun.nextcloud.com:3478
  #     - stun.sip.us:3478
  #     - global.stun.twilio.com:3478
  #   # 下面支持填写针对server-url的TLS配置(sni, skip-cert-verify, fingerprint, certificate, private-key, alpn)
  #   # skip-cert-verify： false
  #   # ......
  ###quic-go特殊配置项，不要随意修改除非你知道你在干什么###
  # initial-stream-receive-window： 8388608
  # max-stream-receive-window： 8388608
  # initial-connection-receive-window： 20971520
  # max-connection-receive-window： 20971520
```

通用字段

TLS 字段

## ports

配置则启用端口跳跃，忽略port，格式参考端口范围

## hop-interval

端口跳跃的间隔，单位为秒，默认为 30

支持填写"15-30"会每次随机选取其中一个值作为切换间隔，仅支持写一个范围（即不允许出现逗号）

## password

认证密码

## up/down

brutal 速率控制，若不写单位，默认为 Mbps

## obfs

QUIC 流量混淆器类型，仅可设为 salamander，如果为空则禁用

## obfs-password

## realm-opts

QUIC 流量混淆器密码

```
2026年5月15日
```

---

## TUIC

> 来源: [https://wiki.metacubex.one/config/proxies/tuic/](https://wiki.metacubex.one/config/proxies/tuic/)

# TUIC

```
proxies:
- name: tuic
  server: www.example.com
  port: 10443
  type: tuic
  token: TOKEN
  uuid: 00000000-0000-0000-0000-000000000001
  password: PASSWORD_1
  # ip: 127.0.0.1
  # heartbeat-interval: 10000
  # alpn: [h3]
  disable-sni: true
  reduce-rtt: true
  request-timeout: 8000
  udp-relay-mode: native
  # congestion-controller: bbr
  # bbr-profile: "" # Available: "standard", "conservative", "aggressive". Default: "standard"
  # max-udp-relay-packet-size: 1500
  # fast-open: true
  # skip-cert-verify: true
  # max-open-streams: 20
  # sni: example.com
```

通用字段

TLS 字段

## token

必须，用于 TUIC V4 的用户标识，使用 TUIC V5 时不可书写

## uuid

必须，用于 TUICV5 的用户唯一识别码，使用 TUIC V4 时不可书写

## password

必须，用于 TUICV5 的用户密码，使用 TUIC V4 时不可书写

## ip

用于覆盖“server”选项中设置的服务器地址的 DNS 查找结果

## heartbeat-interval

发送保持连接活动的心跳包的间隔时间，单位为毫秒

## disable-sni

设置是否在 TLS 握手中禁用 SNI（服务器名称指示）SNI 用于在同一 IP 地址上承载多个 HTTPS 站点

## reduce-rtt

设置是否在客户端启用 QUIC 的 0-RTT 握手这可以减少连接建立时间，但可能增加重放攻击的风险

## request-timeout

设置建立到 TUIC 代理服务器的连接的超时时间，单位为毫秒

## udp-relay-mode

设置 UDP 数据包中继模式，可以是 native/quic

## congestion-controller

设置拥塞控制算法，可选项为 cubic/new_reno/bbr

## max-udp-relay-packet-size

设置最大的 UDP 数据包中继大小，单位为字节

## fast-open

设置是否启用 Fast Open，这可以减少连接建立时间

## max-open-streams

设置最大打开流的数量过多的打开流可能会影响性能

```
2026年4月20日
```

---

## WireGuard

> 来源: [https://wiki.metacubex.one/config/proxies/wg/](https://wiki.metacubex.one/config/proxies/wg/)

# WireGuard

## 简化写法

如果只有一个 peer，可以使用简化写法。

```
proxies:
- name: "wg"
  type: wireguard
  private-key: eCtXsJZ27+4PbhDkHnB923tkUn2Gj59wZw5wFA75MnU=
  server: 162.159.192.1
  port: 2480
  ip: 172.16.0.2
  ipv6: fd01:5ca1:ab1e:80fa:ab85:6eea:213f:f4a5
  public-key: Cr8hWlKvtDt7nrvf+f0brNQQzabAqrjfBvas9pmowjo=
  allowed-ips: ['0.0.0.0/0']
  # pre-shared-key: 31aIhAPwktDGpH4JDhA8GNvjFXEf/a6+UaQRyOAiyfM=
  # reserved: [209,98,59]  # 字符串格式也是合法的，如"U4An"
  # persistent-keepalive: 0
  udp: true
  # mtu: 1408
  # dialer-proxy: "ss1"  # 一个出站代理的标识。当值不为空时，将使用指定的 proxy/proxy-group 发出连接
  # remote-dns-resolve: true # 强制dns远程解析，默认值为false
  # dns: [ 1.1.1.1, 8.8.8.8 ] # 仅在remote-dns-resolve为true时生效
  # 如果存在则开启AmneziaWG功能
  # amnezia-wg-option:
  #   jc: 5
  #   jmin: 500
  #   jmax: 501
  #   s1: 30
  #   s2: 40
  #   s3: 50                                            # AmneziaWG v1.5 and v2
  #   s4: 5                                             # AmneziaWG v1.5 and v2
  #   h1: 123456                                        # AmneziaWG v1.0 and v1.5
  #   h2: 67543                                         # AmneziaWG v1.0 and v1.5
  #   h3: 123123                                        # AmneziaWG v1.0 and v1.5
  #   h4: 32345                                         # AmneziaWG v1.0 and v1.5
  #   h1: 123456-123500                                 # AmneziaWG v2.0 only
  #   h2: 67543-67550                                   # AmneziaWG v2.0 only
  #   h3: 123123-123200                                 # AmneziaWG v2.0 only
  #   h4: 32345-32350                                   # AmneziaWG v2.0 only
  #   i1: <b 0xf6ab3267fa><c><b 0xf6ab><t><r 10><wt 10> # AmneziaWG v1.5 and v2
  #   i2: <b 0xf6ab3267fa><r 100>                       # AmneziaWG v1.5 and v2
  #   i3: ""                                            # AmneziaWG v1.5 and v2
  #   i4: ""                                            # AmneziaWG v1.5 and v2
  #   i5: ""                                            # AmneziaWG v1.5 and v2
  #   j1: <b 0xffffffff><c><b 0xf6ab><t><r 10>          # AmneziaWG v1.5 only (removed in v2)
  #   j2: <c><b 0xf6ab><t><wt 1000>                     # AmneziaWG v1.5 only (removed in v2)
  #   j3: <t><b 0xf6ab><c><r 10>                        # AmneziaWG v1.5 only (removed in v2)
  #   itime: 60                                         # AmneziaWG v1.5 only (removed in v2)
```

## 完整写法

完整写法可以指定多个 peer。

如果使用多个 peer，每一个 peer 的allowed-ips需要做区分；此时顶层段落的server, port, public-key, pre-shared-key, reserved等字段均会被忽略，不过private-key仍然在顶层指定。

```
proxies:
- name: "wg"
  type: wireguard
  ip: 172.16.0.2
  ipv6: fd01:5ca1:ab1e:80fa:ab85:6eea:213f:f4a5
  private-key: eCtXsJZ27+4PbhDkHnB923tkUn2Gj59wZw5wFA75MnU=
  peers:
    - server: 162.159.192.1
      port: 2480
      public-key: Cr8hWlKvtDt7nrvf+f0brNQQzabAqrjfBvas9pmowjo=
      allowed-ips: ['0.0.0.0/0']
      # pre-shared-key: 31aIhAPwktDGpH4JDhA8GNvjFXEf/a6+UaQRyOAiyfM=
      # reserved: [209,98,59]  # 字符串格式也是合法的，如"U4An"
  udp: true
  # mtu: 1408
  # dialer-proxy: "ss1"  # 一个出站代理的标识。当值不为空时，将使用指定的 proxy/proxy-group 发出连接
  # remote-dns-resolve: true # 强制dns远程解析，默认值为false
  # dns: [ 1.1.1.1, 8.8.8.8 ] # 仅在remote-dns-resolve为true时生效
```

通用字段

### ip

本机在 Wireguard 网络中使用的 IPv4 地址

### ipv6

可选字段，本机在 Wireguard 网络中使用的 IPv6 地址

### private-key

base64 编码的 Wireguard 客户端私钥

可以使用wg genkey | tee privatekey | wg pubkey > publickey命令生成一对可用的公私钥文件

### public-key

base64 编码的 Wireguard 服务端公钥

### allowed-ips

可选字段，限制客户端的哪些 IP 段的流量由服务端进行转发。一般情况下可填['0.0.0.0/0']

### pre-shared-key

可选字段，预共享密钥

### reserved

可选字段，Wireguard 协议保留字段的值，部分 WARP 节点需要使用

### persistent-keepalive

可选字段，定期发送数据包来维持连接的持久性

### mtu

可选字段，设置 MTU 值

### remote-dns-resolve

可选字段，是否强制 dns 远程解析，默认值为 false

### dns

可选字段，当remote-dns-resolve为 true 时生效，指定远程解析使用的 dns 服务器

## 从 Wireguard 标准配置文件翻译

假设有如下 Wireguard 标准配置文件：

```
[Interface]
Address = <本机组网IP>
ListenPort = <本地监听端口>
PrivateKey = <本机私钥>
DNS = <使用的DNS>
MTU = <预设MTU>

[Peer]
AllowedIPs = <转发IP段>
Endpoint = <远端地址>:<远端端口>
PublicKey = <远端公钥>
```

对应的 clash 节点配置为：

```
- name: "wg"
  type: wireguard
  ip: <本机组网IP，IPv4往这里填>
  ipv6: <本机组网IP，IPv6往这里填>  # 没有v6地址直接删除
  private-key: <本机私钥>
  peers:
    - server: <远端地址>
      port: <远端端口>
      public-key: <远端公钥>
      allowed-ips: ['0.0.0.0/0']     # 分流由clash处理
      # reserved: [209,98,59]        # 如果需要自己填
  udp: true
  mtu: <预设MTU>               # 按需设置，不用直接删除
  remote-dns-resolve: true    # 按需设置，不用直接删除
  dns: <使用的DNS>             # 按需设置，不用直接删除
```

```
2026年3月20日
```

---

## Tailscale

> 来源: [https://wiki.metacubex.one/config/proxies/tailscale/](https://wiki.metacubex.one/config/proxies/tailscale/)

# Tailscale

```
proxies:
  - name: "tailscale"
    type: tailscale
    hostname: mihomo
    auth-key: tskey-auth-xxxx
    control-url: https://controlplane.tailscale.com
    state-dir: ./tailscale
    ephemeral: false
    udp: true
    accept-routes: true
    exit-node: 100.64.0.1
    exit-node-allow-lan-access: true
    dialer-proxy: "ss1"
    interface-name: "WLAN"
    routing-mark: 6666
    ip-version: ipv4-prefer
```

> **Note

和其他代理协议一样，tailscale 只是一个可选代理节点。只有当你通过规则、策略组或其他路由方式把流量指向这个节点时，对应流量才会走 Tailscale。

此外，所有类型的 proxy 都只会在其匹配上的第一个链接触发的时候才会开始链接（即并不会随 Mihomo 同步启动），所以第一次访问 Tailscale 节点超时是正常情况，请多次尝试直到 Tailscale 正常加载完成。

## name

必须，代理名称，不可重复。

## type

必须，固定为 tailscale。

## hostname

可选，Tailscale 设备名。不填写时由 tsnet 处理。

## auth-key

可选，Tailscale 或 Headscale 的登录密钥。不填写时，首次启动会在日志中输出交互式登录 URL。

## control-url

可选，自定义 Tailscale control server 或 Headscale 地址。

## state-dir

可选，tsnet 状态目录，默认值为 tailscale。

## ephemeral

可选，是否作为 ephemeral node 登录，默认值为 false。

## udp

可选，是否启用 UDP，默认值为 false。

## accept-routes

可选，是否接受 Tailnet 中发布的 subnet routes。

## exit-node

可选，使用指定 exit node。可以填写节点 IP，也可以填写 auto:any。

当目标不在 Tailscale 路由内时，连接会直接报错，不会回退到直连。访问公网需要配置可用的 exit-node，或接受覆盖目标网段的 subnet routes。

## exit-node-allow-lan-access

可选，使用 exit node 时是否允许访问本地 LAN。

## dialer-proxy

可选，一个出站代理的标识。当值不为空时，将使用指定的 proxy 发出 Tailscale 控制面和 DERP/STUN 等连接。

## interface-name

可选，指定出站网卡。

## routing-mark

可选，Linux 下配置 fwmark。

## ip-version

可选，指定出站使用的 IP 版本。

可选值：dual/ipv4/ipv6/ipv4-prefer/ipv6-prefer。

```
2026年5月22日
```

---

## SSH

> 来源: [https://wiki.metacubex.one/config/proxies/ssh/](https://wiki.metacubex.one/config/proxies/ssh/)

# SSH

```
proxies:
- name: "ssh-out"
  type: ssh

  server: 127.0.0.1
  port: 22
  username: root
  password: password
  private-key: key
  private-key-passphrase: key_password
  host-key:
  - "ssh-rsa AAAAB3NzaC1yc2EAA..."
  host-key-algorithms:
  - rsa
```

通用字段

## username

SSH 用户

## password

SSH 密码

## private-key

密钥内容/路径

## private_key_passphrase

密钥密码

## host-key

主机密钥，留空接受所有

## host-key-algorithms

主机密钥算法

```
2026年5月14日
```

---

## MASQUE

> 来源: [https://wiki.metacubex.one/config/proxies/masque/](https://wiki.metacubex.one/config/proxies/masque/)

# MASQUE

```
proxies:
# masque
- name: "masque"
  type: masque
  server: server.com
  port: 443
  private-key: BASE64_ENCODED_PRIVATE_KEY
  public-key: BASE64_ENCODED_PUBLIC_KEY
  ip: 172.16.0.2/32
  ipv6: fd00::2/128
  mtu: 1280
  udp: true
  # 一个出站代理的标识。当值不为空时，将使用指定的 proxy 发出连接
  # dialer-proxy: "ss1"
  # remote-dns-resolve: true # 强制 dns 远程解析，默认值为 false
  # dns: [ 1.1.1.1, 8.8.8.8 ] # 仅在 remote-dns-resolve 为 true 时生效
  # congestion-controller: bbr # 默认不开启
  # bbr-profile: "" # Available: "standard", "conservative", "aggressive". Default: "standard"

# masque-h2
- name: "masque-h2"
  type: masque
  server: server.com
  port: 443
  private-key: BASE64_ENCODED_PRIVATE_KEY
  public-key: BASE64_ENCODED_PUBLIC_KEY
  ip: 172.16.0.2/32
  ipv6: fd00::2/128
  mtu: 1280
  udp: true
  network: h2
  # 一个出站代理的标识。当值不为空时，将使用指定的 proxy 发出连接
  # dialer-proxy: "ss1"
  # remote-dns-resolve: true # 强制 dns 远程解析，默认值为 false
  # dns: [ 1.1.1.1, 8.8.8.8 ] # 仅在 remote-dns-resolve 为 true 时生效
```

## 获取 masque 配置

通过 usque 工具生成 masque 配置。

通用字段

## private-key

必须，base64 编码的 ECDSA 私钥

## public-key

必须，base64 编码的 ECDSA 公钥（服务器端公钥）

> **Note

注意去掉 PEM 格式中的头尾标识符，如 -----BEGIN PUBLIC KEY----- ， -----END PUBLIC KEY----- 和 PEM 格式中的换行符 \n。

## ip

本地 IPv4 地址，CIDR 格式（如 172.16.0.2/32）

## ipv6

本地 IPv6 地址，CIDR 格式（如 fd00::2/128）

## mtu

TUN 设备的 MTU 大小，默认为 1280

## udp

是否启用 UDP 支持，默认为 false

## remote-dns-resolve

是否启用远程 DNS 解析，通过 MASQUE 隧道解析 DNS

## dns

远程 DNS 服务器列表，在启用 remote-dns-resolve 后生效

## congestion-controller

修改默认的拥塞控制算法，默认为不启用，可选值包括 bbr

## network

optional, 默认为 quic，masque-h2 需要设置为h2

```
2026年5月14日
```

---

## TrustTunnel

> 来源: [https://wiki.metacubex.one/config/proxies/trust-tunnel/](https://wiki.metacubex.one/config/proxies/trust-tunnel/)

# 404 - Not found

---

## OpenVPN

> 来源: [https://wiki.metacubex.one/config/proxies/openvpn/](https://wiki.metacubex.one/config/proxies/openvpn/)

# OpenVPN

```
proxies:
  - name: "openvpn"
    type: openvpn
    server: vpn.example.com
    port: 1194
    proto: udp
    # auth-user-pass 模式
    username: "user"
    password: "pass"
    # 证书模式 (与上方二选一)
    # cert: |
    #   -----BEGIN CERTIFICATE-----
    #   ...
    #   -----END CERTIFICATE-----
    # key: |
    #   -----BEGIN PRIVATE KEY-----
    #   ...
    #   -----END PRIVATE KEY-----
    ca: |
      -----BEGIN CERTIFICATE-----
      MIIB...example
      -----END CERTIFICATE-----
    # tls-crypt: |
    #  -----BEGIN OpenVPN Static key V1-----
    #  ...
    #  -----END OpenVPN Static key V1-----
    # ping: 10
    # ping-restart: 60
    # dev: tun
    # cipher: AES-128-GCM
    # auth: SHA256
    # comp-lzo: "no"
    udp: true
    # mtu: 1500
    # dialer-proxy: "ss1"
    # remote-dns-resolve: true
    # dns: [ 1.1.1.1, 8.8.8.8 ]
```

通用字段

## server

**必须**，OpenVPN 服务器地址。

## port

**必须**，OpenVPN 服务器端口，默认 1194。

## proto

可选，协议类型，支持 udp 或 tcp，默认 udp。

## username / password

**可选**，auth-user-pass 认证模式所需的用户名与密码。



**注意**：需与 cert / key 二选一，两者不能同时为空。

## ca

**必须**，CA 证书内容。请从 .ovpn 文件的 &nbsp; 标签中复制。

## cert

**可选**，客户端证书内容。从 .ovpn 文件的 &nbsp; 标签中复制。使用用户名/密码认证时可省略。

## key

**可选**，客户端私钥内容。从 .ovpn 文件的 &nbsp; 标签中复制。使用用户名/密码认证时可省略。

## tls-crypt

**可选**，TLS 加密密钥。从 .ovpn 文件的 &nbsp; 标签中复制，不需要保留标签。

## ping

可选，默认值为0。

## ping-restart

可选，默认值为0。

## dev

可选，虚拟网卡类型，当前仅支持 tun，默认 tun。

## cipher

可选，加密方式，支持 AES-128-GCM / AES-256-GCM/ AES-128-CBC / AES-256-CBC /CHACHA20-POLY1305默认 AES-128-GCM， AES-CBC 会按 AES-128-CBC 处理。

## auth

可选，数据验证算法，支持 MD5 /SHA1 / SHA256/ SHA384 / SHA512，默认 SHA256；AEAD cipher 会忽略 auth。

## comp-lzo

可选，数据压缩方式。可选值：yes、no、adaptive。

## udp

可选，是否启用 UDP。true 表示使用 UDP，false 表示使用 TCP。

## mtu

可选，最大传输单元，默认 1500。

## dialer-proxy

可选，指定出站代理的标识。若设置，OpenVPN 连接将通过该代理发出。

## remote-dns-resolve

可选，是否强制 DNS 远程解析，默认 false。

## dns

可选，仅在 remote-dns-resolve 为 true 时生效，指定远程解析的 DNS 地址列表。

```
2026年6月6日
```

---

---

# 6. 代理集合

# 代理集合子页面

## 代理集合

> 来源: [https://wiki.metacubex.one/config/proxy-providers/](https://wiki.metacubex.one/config/proxy-providers/)

# 代理集合

```
proxy-providers:
  provider1:
    type: http
    url: "http://test.com"
    path: ./proxy_providers/provider1.yaml
    interval: 3600
    proxy: DIRECT
    size-limit: 0
    age-secret-key: AGE-SECRET-KEY-1ZTQLLN0A4U3ZTT3DCZKYN0CGZEZQLWX2DFTXUWMT4ZHR0N2UG6LSW9NT0N
    header:
      User-Agent:
      - "mihomo/1.18.3"
      Authorization:
      - 'token 1231231'
      # X-Age-Public-Key:
      # - 'age1xh86kh9v23vattr58yedspm3f57sxvnswu9krr6ns438amekx5gsd09uma'
    health-check:
      enable: true
      url: https://www.gstatic.com/generate_204
      interval: 300
      timeout: 5000
      lazy: true
      expected-status: 204
    override:
      tfo: false
      mptcp: false
      udp: true
      udp-over-tcp: false
      down: "50 Mbps"
      up: "10 Mbps"
      skip-cert-verify: true
      dialer-proxy: proxy
      interface-name: tailscale0
      routing-mark: 233
      ip-version: ipv4-prefer
      additional-prefix: "provider1 prefix |"
      additional-suffix: "| provider1 suffix"
      proxy-name:
      - pattern: "IPLC-(.*?)倍"
        target: "iplc x $1"
    filter: "(?i)港|hk|hongkong|hong kong"
    exclude-filter: "xxx"
    exclude-type: "ss|http"
    payload:
      - name: "ss1"
        type: ss
        server: server
        port: 443
        cipher: chacha20-ietf-poly1305
        password: "password"
```

## name

必须，如provider1，不能重复，建议不要和策略组名称重复

## type

必须，provider类型，可选http / file / inline

## url

类型为http是则需要配置

## path

可选，文件路径，不可重复，不填写时会使用 url 的 MD5 作为此文件的文件名

由于安全问题，此路径将限制只允许在 HomeDir(有启动参数 -d 配置) 中，如果想存储到其他位置，请通过设置 SAFE_PATHS 环境变量指定额外的安全路径。该环境变量的语法同本操作系统的 PATH 环境变量解析规则（即 Windows 下以分号分割，其他系统下以冒号分割）

## interval

更新provider的时间，单位为秒

## proxy

经过指定代理进行下载/更新

## size-limit

限制下载文件的最大大小，默认为 0 即不限制文件大小，单位为字节 (b)

## age-secret-key

如果设置会age-secret-key尝试通过此secret解密age armor格式加密的配置文件

注意：

- 对于加密内容，目前仅支持 age-encryption.org/v1 的 official ASCII "armor" format
- 对于key格式，目前仅支持 age-encryption.org/v1 的 x25519 recipient type 和 The mlkem768-x25519 hybrid post-quantum recipient type
- 目前核心不会主动将公钥发送给服务器，用户需要自行设置header里X-Age-Public-Key或者通过其他方式将公钥上传
- 核心同样支持命令行参数-age-secret-key或CLASH_AGE_SECRET_KEY环境变量来加载加密配置文件

实用工具：

- 您可以通过 mihomo age keygen 生成符合要求的 x25519 key
- 您可以通过 mihomo age keygen-pq 生成符合要求的 mlkem768-x25519 key
- 您可以通过 mihomo age convert &nbsp; 从 age-secret-key 导出 age-public-key
- 您可以通过 mihomo age decrypt &nbsp; &nbsp; &nbsp; 将已加密文件解密， 为 - 时会从标准输入读取， 为 - 时会往标准输出写入
- 您可以通过 mihomo age encrypt &nbsp; &nbsp; &nbsp; 将未加密文件加密， 为 - 时会从标准输入读取， 为 - 时会往标准输出写入

参考实现：

- Golang: FiloSottile/age
- Rust: str4d/rage
- Typescript: FiloSottile/typage

## header

自定义 http 请求头

## health-check

健康检查 (延迟测试)

### health-check.enable

是否启用，可选 true/false

### health-check.url

健康检查地址，推荐使用以下地址之一

```
https://cp.cloudflare.com
https://www.gstatic.com/generate_204
```

### health-check.interval

健康检查间隔时间，单位为秒

### health-check.timeout

健康检查超时时间，单位为毫秒

### health-check.lazy

懒惰状态，默认为true,不使用该集合节点时，不进行测试

### health-check.expected-status

参阅 期望状态

## override

覆写节点内容，以下为支持的字段

### override.additional-prefix

为节点名称添加固定前缀

### override.additional-suffix

为节点名称添加固定后缀

### override.proxy-name

对节点名称内容进行替换，支持正则表达式，pattern 为替换内容，target 为替换目标

### override.其余配置项

参阅通用字段  tfo

参阅通用字段  mptcp

参阅通用字段  udp

参阅Shadowsocks  udp-over-tcp

参阅Hysteria/Hysteria2  up

参阅Hysteria/Hysteria2  down

参阅通用字段  skip-cert-verify

参阅通用字段  dialer-proxy

参阅通用字段  interface-name

参阅通用字段  routing-mark

参阅通用字段  ip-version

## filter

筛选满足关键词或正则表达式的节点，可以使用 ` 区分多个正则表达式

## exclude-filter

排除满足关键词或正则表达式的节点，可以使用 ` 区分多个正则表达式

## exclude-type

不支持正则表达式，通过 | 分割，根据节点类型排除

provider 的 exclude-type 使用配置文件中的 type 类型进行排除

## payload

内容，仅 type 为 inline 时生效

当 http 或 file 解析失败时，也可以使用 payload 作为备用代理

```
2026年6月6日
```

---

## 代理集合内容

> 来源: [https://wiki.metacubex.one/config/proxy-providers/content/](https://wiki.metacubex.one/config/proxy-providers/content/)

# proxy-providers 文件内容

```
proxies:
- name: "ss1"
  type: ss
  server: server
  port: 443
  cipher: chacha20-ietf-poly1305
  password: "password"
- name: "ss2"
  type: ss
  server: server
  port: 443
  cipher: chacha20-ietf-poly1305
  password: "password"

ss://YWVzLTI1Ni1nY206bWV0YUAxMjcuMC4wLjE6NDQz#home
vmess://eyJhZGQiOiIxMjcuMC4wLjEiLCJhaWQiOiIwIiwiYWxwbiI6IiIsImZwIjoiIiwiaG9zdCI6IiIsImlkIjoiMTIyMzQ1Njc4OSIsIm5ldCI6InRjcCIsInBhdGgiOiIiLCJwb3J0IjoiNDQzIiwicHMiOiJ2bWVzcyIsInNjeSI6ImF1dG8iLCJzbmkiOiIiLCJ0bHMiOiIiLCJ0eXBlIjoibm9uZSIsInYiOiIyIn0=

c3M6Ly9ZV1Z6TFRJMU5pMW5ZMjA2YldWMFlVQXhNamN1TUM0d0xqRTZORFF6I2hvbWUKdm1lc3M6Ly9leUpoWkdRaU9pSXhNamN1TUM0d0xqRWlMQ0poYVdRaU9pSXdJaXdpWVd4d2JpSTZJaUlzSW1ad0lqb2lJaXdpYUc5emRDSTZJaUlzSW1sa0lqb2lNVEl5TXpRMU5qYzRPU0lzSW01bGRDSTZJblJqY0NJc0luQmhkR2dpT2lJaUxDSndiM0owSWpvaU5EUXpJaXdpY0hNaU9pSjJiV1Z6Y3lJc0luTmplU0k2SW1GMWRHOGlMQ0p6Ym1raU9pSWlMQ0owYkhNaU9pSWlMQ0owZVhCbElqb2libTl1WlNJc0luWWlPaUl5SW4wPQ==
```

> **Note

base64 的 uri 通常为提供商提供给 v2ray/xray 的订阅链接内容

YAML/uri/base64不可写在同一文件,uri/base64不需要proxies:字段,直接书写即可

```
2024年4月1日
```

---

---

# 7. 代理组

## 7.1 通用字段

```yaml
proxy-groups:
  - name: "proxy"
    type: select
    proxies:
      - DIRECT
      - ss
    use:
      - provider1
    url: 'https://www.gstatic.com/generate_204'
    interval: 300
    lazy: true
    empty-fallback: COMPATIBLE
    timeout: 5000
    max-failed-times: 5
    disable-udp: true
    interface-name: en0
    routing-mark: 11451
    include-all: false
    include-all-proxies: false
    include-all-providers: false
    filter: "(?i)港|hk|hongkong|hong kong"
    exclude-filter: "美|日"
    exclude-type: "Shadowsocks|Http"
    expected-status: 204
    hidden: true
    icon: xxx
```


| 参数                      | 说明                           |
| ----------------------- | ---------------------------- |
| `name`                  | 必须，策略组的名字（如有特殊符号，应使用引号包裹）    |
| `type`                  | 必须，策略组的类型                    |
| `proxies`               | 引入出站代理或其他策略组                 |
| `use`                   | 引入代理集合                       |
| `url`                   | 健康检查测试地址                     |
| `interval`              | 健康检查间隔（秒），不为 0 则启用定时测试       |
| `lazy`                  | 懒惰状态，默认 true，未选择到当前策略组时不进行测试 |
| `empty-fallback`        | 当组为空时的回退 proxy，默认 COMPATIBLE |
| `timeout`               | 健康检查超时时间                     |
| `max-failed-times`      | 最大失败次数                       |
| `disable-udp`           | 禁用 UDP                       |
| `interface-name`        | 出站接口                         |
| `routing-mark`          | 路由标记                         |
| `include-all`           | 引入所有代理和集合                    |
| `include-all-proxies`   | 引入所有代理                       |
| `include-all-providers` | 引入所有集合                       |
| `filter`                | 正则过滤代理名                      |
| `exclude-filter`        | 正则排除代理名                      |
| `exclude-type`          | 排除代理类型                       |
| `expected-status`       | 期望的健康检查状态码                   |
| `hidden`                | 在 UI 中隐藏                     |
| `icon`                  | 图标                           |


# 代理组子页面

## 内置代理组

> 来源: [https://wiki.metacubex.one/config/proxy-groups/built-in/](https://wiki.metacubex.one/config/proxy-groups/built-in/)

# 预置代理组

## GLOBAL

默认填充所有代理组和代理节点

该策略组也可以在配置文件中自定义

但 web 面板以及部分客户端使用 GLOBAL 策略组内的 策略组顺序 进行排序

建议在自定义 GLOBAL 策略组时，书写完整当前配置所有的策略组 (除了 GLOBAL 本身)

```
2025年1月14日
```

---

## 手动选择 (Select)

> 来源: [https://wiki.metacubex.one/config/proxy-groups/select/](https://wiki.metacubex.one/config/proxy-groups/select/)

# Select

```
proxy-groups:
- name: Proxy
  type: select
  proxies:
  - ss
  - vmess
  - auto
  #disable-udp: true
```

## 通用字段

参阅 通用字段

```
2024年3月13日
```

---

## 自动选择 (URLTest)

> 来源: [https://wiki.metacubex.one/config/proxy-groups/auto/](https://wiki.metacubex.one/config/proxy-groups/auto/)

# 404 - Not found

---

## 自动回退 (Fallback)

> 来源: [https://wiki.metacubex.one/config/proxy-groups/fallback/](https://wiki.metacubex.one/config/proxy-groups/fallback/)

# Fallback

```
proxy-groups:
- name: "fallback"
  type: fallback
  proxies:
  - ss
  - vmess
  url: 'https://www.gstatic.com/generate_204'
  interval: 300
  #lazy: true
```

当前节点超时时，则会按代理顺序选择第一个可用节点

## 通用字段

参阅 通用字段

```
2025年3月19日
```

---

## 负载均衡 (LoadBalance)

> 来源: [https://wiki.metacubex.one/config/proxy-groups/load-balance/](https://wiki.metacubex.one/config/proxy-groups/load-balance/)

# Load-balance

```
proxy-groups:
- name: "load-balance"
  type: load-balance
  proxies:
  - ss1
  - ss2
  - vmess1
  url: 'https://www.gstatic.com/generate_204'
  interval: 300
  #lazy: true
  #strategy: consistent-hashing
```

## 通用字段

参阅 通用字段

## strategy

负载均衡策略

- 

round-robin 将会把所有的请求分配给策略组内不同的代理节点

- 

consistent-hashing 将相同的 目标地址 的请求分配给策略组内的同一个代理节点

- 

sticky-sessions: 将相同的 来源地址 和 目标地址 的请求分配给策略组内的同一个代理节点，缓存 10 分钟过期

> **Note

目标地址 为域名时，使用顶级域名匹配

```
2024年10月12日
```

---

## 链式代理 (Relay)

> 来源: [https://wiki.metacubex.one/config/proxy-groups/relay/](https://wiki.metacubex.one/config/proxy-groups/relay/)

# Relay

> **Warning

relay 策略已经被弃用，请使用dialer-proxy

```
2026年5月14日
```

---

---

# 8. 路由规则

```yaml
rules:
  - DOMAIN,ad.com,REJECT
  - DOMAIN-SUFFIX,google.com,auto
  - DOMAIN-KEYWORD,google,auto
  - DOMAIN-WILDCARD,*.google.com,auto
  - DOMAIN-REGEX,^abc.*com,PROXY
  - GEOSITE,youtube,PROXY
  - IP-CIDR,127.0.0.0/8,DIRECT,no-resolve
  - IP-CIDR6,2620:0:2d0:200::7/32,auto
  - IP-SUFFIX,8.8.8.8/24,PROXY
  - IP-ASN,13335,DIRECT
  - GEOIP,CN,DIRECT
  - SRC-GEOIP,cn,DIRECT
  - SRC-IP-ASN,9808,DIRECT
  - SRC-IP-CIDR,192.168.1.201/32,DIRECT
  - SRC-IP-SUFFIX,192.168.1.201/8,DIRECT
  - DST-PORT,80,DIRECT
  - SRC-PORT,7777,DIRECT
  - IN-PORT,7890,PROXY
  - IN-TYPE,SOCKS/HTTP,PROXY
  - IN-USER,mihomo,PROXY
  - IN-NAME,ss,PROXY
  - PROCESS-PATH,/usr/bin/wget,PROXY
  - PROCESS-PATH,C:\Program Files\Google\Chrome\Application\chrome.exe,PROXY
  - PROCESS-PATH-WILDCARD,/usr/*/wget,PROXY
  - PROCESS-PATH-REGEX,.*bin/wget,PROXY
  - PROCESS-NAME,curl,PROXY
  - PROCESS-NAME,com.termux,PROXY
  - PROCESS-NAME-WILDCARD,*telegram*,PROXY
  - PROCESS-NAME-REGEX,curl$,PROXY
  - UID,1001,DIRECT
  - NETWORK,udp,DIRECT
  - DSCP,4,DIRECT
  - RULE-SET,providername,proxy
  - AND,((DOMAIN,baidu.com),(NETWORK,UDP)),DIRECT
  - OR,((NETWORK,UDP),(DOMAIN,baidu.com)),REJECT
  - NOT,((DOMAIN,baidu.com)),PROXY
  - SUB-RULE,(NETWORK,tcp),sub-rule
  - MATCH,auto
```

规则将按照从上到下的顺序匹配，列表顶部的规则优先级高于其底下的规则。
如请求为 udp，而代理节点没有 udp 支持，则会继续向下匹配。

### 域名规则


| 规则类型              | 说明                                         |
| ----------------- | ------------------------------------------ |
| `DOMAIN`          | 匹配完整域名                                     |
| `DOMAIN-SUFFIX`   | 匹配域名后缀（例：`google.com` 匹配 `www.google.com`） |
| `DOMAIN-KEYWORD`  | 域名关键字匹配                                    |
| `DOMAIN-WILDCARD` | 通配符匹配，仅支持 `*` 和 `?`                        |
| `DOMAIN-REGEX`    | 域名正则表达式匹配                                  |
| `GEOSITE`         | 匹配 Geosite 内的域名                            |


### IP 规则


| 规则类型                   | 说明             |
| ---------------------- | -------------- |
| `IP-CIDR` / `IP-CIDR6` | 匹配 IP 地址范围     |
| `IP-SUFFIX`            | 匹配 IP 后缀范围     |
| `IP-ASN`               | 匹配 IP 所属 ASN   |
| `GEOIP`                | 匹配 IP 所属国家代码   |
| `SRC-GEOIP`            | 匹配来源 IP 所属国家代码 |
| `SRC-IP-ASN`           | 匹配来源 IP 所属 ASN |
| `SRC-IP-CIDR`          | 匹配来源 IP 地址范围   |
| `SRC-IP-SUFFIX`        | 匹配来源 IP 后缀范围   |


### 端口规则


| 规则类型       | 说明         |
| ---------- | ---------- |
| `DST-PORT` | 匹配请求目标端口范围 |
| `SRC-PORT` | 匹配请求来源端口范围 |
| `IN-PORT`  | 匹配入站端口     |


### 入站规则


| 规则类型      | 说明                  |
| --------- | ------------------- |
| `IN-TYPE` | 匹配入站类型              |
| `IN-USER` | 匹配入站用户名，支持 `/` 分隔多个 |
| `IN-NAME` | 匹配入站名称              |


### 进程规则


| 规则类型                    | 说明                     |
| ----------------------- | ---------------------- |
| `PROCESS-PATH`          | 使用完整进程路径匹配             |
| `PROCESS-PATH-WILDCARD` | 使用进程路径通配符匹配            |
| `PROCESS-PATH-REGEX`    | 使用进程路径正则匹配             |
| `PROCESS-NAME`          | 使用进程名匹配（Android 可匹配包名） |
| `PROCESS-NAME-WILDCARD` | 使用进程名通配符匹配             |
| `PROCESS-NAME-REGEX`    | 使用进程名正则匹配              |


### 其他规则


| 规则类型                 | 说明                                              |
| -------------------- | ----------------------------------------------- |
| `UID`                | 匹配 Linux USER ID                                |
| `NETWORK`            | 匹配 `tcp` 或 `udp`                                |
| `DSCP`               | 匹配 DSCP 标记（仅限 tproxy udp 入站）                    |
| `RULE-SET`           | 引用规则集合，需配置 rule-providers                       |
| `AND` & `OR` & `NOT` | 逻辑规则：`LOGIC_TYPE,((payload1),(payload2)),Proxy` |
| `SUB-RULE`           | 匹配至子规则                                          |
| `MATCH`              | 匹配所有请求，无需条件                                     |


### 附加参数


| 参数           | 说明                       |
| ------------ | ------------------------ |
| `no-resolve` | 跳过 DNS 解析，仅支持关于目标 IP 的规则 |
| `src`        | 将目标 IP 匹配转为来源 IP 匹配      |


---

# 9. 规则集合

# 规则集合子页面

## 规则集合

> 来源: [https://wiki.metacubex.one/config/rule-providers/](https://wiki.metacubex.one/config/rule-providers/)

# rule-provider

```
rule-providers:
  google:
    type: http
    path: ./rule1.yaml
    url: "https://raw.githubusercontent.com/../Google.yaml"
    interval: 600
    proxy: DIRECT
    behavior: classical
    format: yaml
    # path-in-bundle: "geo/geosite/cn.mrs"
    size-limit: 0
    header:
      User-Agent:
      - "mihomo/1.18.3"
      Authorization:
      - 'token 1231231'
    payload:
      - 'DOMAIN-SUFFIX,google.com'
```

## name

必须，如google,不能重复

## type

必须，provider类型，可选http / file / inline

## url

类型为http则必须配置

## path

可选，文件路径，不可重复，不填写时会使用 url 的 MD5 作为此文件的文件名

由于安全问题，此路径将限制只允许在 HomeDir(有启动参数 -d 配置) 中，如果想存储到其他位置，请通过设置 SAFE_PATHS 环境变量指定额外的安全路径。该环境变量的语法同本操作系统的PATH环境变量解析规则（即Windows下以分号分割，其他系统下以冒号分割）

## interval

更新provider的时间，单位为秒

## proxy

经过指定代理进行下载/更新

## behavior

行为，可选domain/ipcidr/classical，对应不同格式的 rule-provider 文件格式，请按实际格式填写

## format

格式，可选 yaml/text/mrs，默认 yaml

mrs目前 behavior 仅支持 domain/ipcidr，可以通过mihomo convert-ruleset domain/ipcidr yaml/text XXX.yaml XXX.mrs转换得到

## path-in-bundle

当设置该选项时，如果本地文件不存在，会优先从 Home Dir 的 BundleMRS.7z 中解压该文件，此项的值用于指定在 BundleMRS.7z 中的路径

BundleMRS.7z 可在 [https://github.com/MetaCubeX/meta-rules-dat](https://github.com/MetaCubeX/meta-rules-dat) 处下载

## size-limit

限制下载文件的最大大小，默认为 0 即不限制文件大小，单位为字节 (b)

## header

自定义 http 请求头

## payload

内容，仅 type 为 inline 时生效

```
2026年6月6日
```

---

## 规则集合内容

> 来源: [https://wiki.metacubex.one/config/rule-providers/content/](https://wiki.metacubex.one/config/rule-providers/content/)

# rule-providers 文件内容

## classical

classical 支持路由规则的全部类型 (rule-set/sub-rule 除外)

```
payload:
- DOMAIN-SUFFIX,google.com
- DOMAIN-KEYWORD,google
- DOMAIN,ad.com
- SRC-IP-CIDR,192.168.1.201/32
- IP-CIDR,127.0.0.0/8
- GEOIP,CN
- DST-PORT,80
- SRC-PORT,7777
DOMAIN-SUFFIX,google.com
DOMAIN-KEYWORD,google
DOMAIN,ad.com
SRC-IP-CIDR,192.168.1.201/32
IP-CIDR,127.0.0.0/8
GEOIP,CN
DST-PORT,80
SRC-PORT,7777
```

## domain

domain类规则集合内容通配应遵守clash 通配符

```
payload:
- '.blogger.com'
- '*.*.microsoft.com'
- 'books.itunes.apple.com'
.blogger.com
*.*.microsoft.com
books.itunes.apple.com
```

## ipcidr

```
payload:
- '192.168.1.0/24'
- '10.0.0.0.1/32'
192.168.1.0/24
10.0.0.0.1/32
```

```
2026年4月19日
```

---

---

# 10. 子规则

# 子规则

## 子规则

> 来源: [https://wiki.metacubex.one/config/sub-rules/](https://wiki.metacubex.one/config/sub-rules/)

# 404 - Not found

---

---

# 11. 流量隧道

# 流量隧道

## 流量隧道

> 来源: [https://wiki.metacubex.one/config/tunnels/](https://wiki.metacubex.one/config/tunnels/)

# TUNNEL

流量转发隧道，可以转发 tcp/udp 流量，也可经过代理转发

```
tunnels:
- tcp/udp,127.0.0.1:6553,114.114.114.114:53,proxy
- network: [tcp, udp]
  address: 127.0.0.1:7777
  target: target.com
  proxy: proxy
```

## 单行

```
tunnels:
- tcp/udp,127.0.0.1:6553,8.8.8.8:53,proxy
```

## 多行

```
tunnels:
- network: [tcp, udp]
  address: 127.0.0.1:6553
  target: 8.8.8.8:53
  proxy: proxy
```

单行顺序分别对应多行的 network/address/target/proxy

### network

需要监听的网络类型，可为 tcp/udp

### address

本地监听地址

### target

转发的目标地址

### proxy

可选项，经过某个 proxies/proxy-groups 发送流量

如上示例为   访问127.0.0.1:6553为经过proxy这个proxies/proxy-groups访问8.8.8.8:53

```
2024年8月5日
```

---

---

# 12. NTP

```yaml
ntp:
  enable: true
  write-to-system: true
  server: time.apple.com
  port: 123
  interval: 30
```


| 参数                | 说明                           |
| ----------------- | ---------------------------- |
| `enable`          | 是否启用 NTP 服务                  |
| `write-to-system` | 是否同步至系统时间，需要 root/管理员模式运行    |
| `server`          | NTP 服务地址，默认 `time.apple.com` |
| `port`            | NTP 服务端口，默认 `123`            |
| `interval`        | 同步时间间隔，单位（分），默认 30 分         |


---

# 13. 实验性配置

```yaml
experimental:
  quic-go-disable-gso: false
  quic-go-disable-ecn: false
  dialer-ip4p-convert: false
```


| 参数                    | 说明           |
| --------------------- | ------------ |
| `quic-go-disable-gso` | 禁用 GSO       |
| `quic-go-disable-ecn` | 禁用 ECN       |
| `dialer-ip4p-convert` | 启用 IP4P 地址转换 |


---

&nbsp;