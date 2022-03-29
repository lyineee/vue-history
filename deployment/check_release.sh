#!/bin/sh
etag=""
dist_file_name=dist.tar.gz;
use_proxy=false;
proxy_url=https://ghproxy.com/;
extract_path=/usr/share/nginx/html
while true; do
    status_code=$(curl -sS -I -H "If-None-Match: $etag" https://api.github.com/repos/lyineee/vue-history/releases/latest | head -n 1 | awk '{print $2}');
    if [ $status_code -eq 200 ];then
        echo "new release found";

        jd=$(curl -sS -H "Accept: application/vnd.github.v3+json" https://api.github.com/repos/lyineee/vue-history/releases/latest);
        rm $dist_file_name;
        rm -rf dist
        tag_name=$(echo $jd|grep -Eo '"tag_name": "[^"]+"'|awk '{print $2}'|sed -e 's/"//g');
        echo "tag name:" $tag_name;
        dist_url="https://github.com/lyineee/vue-history/releases/download/${tag_name}/${dist_file_name}";
        echo "assert url:" $dist_url;
        wget -q $dist_url;
        tar -xzf $dist_file_name;
        rm -rf $extract_path
        cp -r dist/* $extract_path
        echo "update complete"
        etag=$(curl -sS --head https://api.github.com/repos/lyineee/vue-history/releases/latest|grep etag|grep -Eo '".+"');
        echo "new etag:" $etag;
    fi
    sleep 1;
done