# Lab 在室把握システム

### これは何？
研究室に誰がいるかリアルタイムで把握できるシステムです．
オフィスとかでも十分に使用できると思います．

### 仕組み
ネットワークに流れるパケットを解析し，あらかじめ登録された mac アドレスと紐づくユーザーの名前をアプリ内で管理，そしてhtmlにちょいちょいと書き出す仕組みです．

### 使い方

1. 設定ファイルを curl で落とします．

```
curl "https://raw.githubusercontent.com/Issei0804-ie/who-is-in-a-lab/main/sample-address.json" > address.json
```

2. キャプチャしたい network interface を調べましょう!(mac だと ifconfig , manjaro だと ip address でいけます)．


3. install と実行を行います．実行時, network interface をキャプチャするために sudo で実行しています．

```
go install github.com/Issei0804-ie/who-is-in-a-lab@latest
sudo who-is-in-a-lab [network interface]
ex) sudo who-is-in-a-lab wlan0
```

4. HTMLに表示される名前 とmac アドレスを登録しましょう． curl からいけます．

```
curl --header "Content-Type:application/json" \
--request POST \
--data '{"name":"issei", "addresses":["ee:ee:ee:ee:ee:ee"]}' \
http://localhost/register
```

5. web(http://localhost) から簡単に在室状況が確認できます．


### 注意

どこからでも mac アドレスを登録できるようにしているので，悪意のあるユーザーから攻撃される恐れがあります．
