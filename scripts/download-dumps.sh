#!/bin/bash
set -euo pipefail
IFS=$'\n\t'

url="https://ftp.ripe.net/ripe/dbase/split"
dir="tests/data"

mkdir -p "$dir"

wget -qO- "$url/ripe.db.as-block.gz" | gzip -dc > "$dir/ripe.db.as-block"
wget -qO- "$url/ripe.db.as-set.gz" | gzip -dc > "$dir/ripe.db.as-set"
# wget -qO- "$url/ripe.db.aut-num.gz" | gzip -dc > "$dir/ripe.db.aut-num"
# wget -qO- "$url/ripe.db.domain.gz" | gzip -dc > "$dir/ripe.db.domain"
wget -qO- "$url/ripe.db.filter-set.gz" | gzip -dc > "$dir/ripe.db.filter-set"
wget -qO- "$url/ripe.db.inet-rtr.gz" | gzip -dc > "$dir/ripe.db.inet-rtr"
# wget -qO- "$url/ripe.db.inet6num.gz" | gzip -dc > "$dir/ripe.db.inet6num"
# wget -qO- "$url/ripe.db.inetnum.gz" | gzip -dc > "$dir/ripe.db.inetnum"
wget -qO- "$url/ripe.db.irt.gz" | gzip -dc > "$dir/ripe.db.irt"
# wget -qO- "$url/ripe.db.key-cert.gz" | gzip -dc > "$dir/ripe.db.key-cert"
wget -qO- "$url/ripe.db.mntner.gz" | gzip -dc > "$dir/ripe.db.mntner"
# wget -qO- "$url/ripe.db.organisation.gz" | gzip -dc > "$dir/ripe.db.organisation"
wget -qO- "$url/ripe.db.peering-set.gz" | gzip -dc > "$dir/ripe.db.peering-set"
wget -qO- "$url/ripe.db.person.gz" | gzip -dc > "$dir/ripe.db.person"
wget -qO- "$url/ripe.db.poem.gz" | gzip -dc > "$dir/ripe.db.poem"
wget -qO- "$url/ripe.db.poetic-form.gz" | gzip -dc > "$dir/ripe.db.poetic-form"
# wget -qO- "$url/ripe.db.role.gz" | gzip -dc > "$dir/ripe.db.role"
wget -qO- "$url/ripe.db.route-set.gz" | gzip -dc > "$dir/ripe.db.route-set"
# wget -qO- "$url/ripe.db.route.gz" | gzip -dc > "$dir/ripe.db.route"
# wget -qO- "$url/ripe.db.route6.gz" | gzip -dc > "$dir/ripe.db.route6"
wget -qO- "$url/ripe.db.rtr-set.gz" | gzip -dc > "$dir/ripe.db.rtr-set"

wget -qO- "ftp://ftp.radb.net/radb/dbase/archive/2021/radb.db.210801.gz" | gzip -dc > "$dir/radb.db"
