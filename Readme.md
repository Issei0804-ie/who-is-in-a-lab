# Lab 在室把握システム

### これは何？
研究室に誰がいるかリアルタイムで把握できるシステムです．
オフィスとかでも十分に使用できると思います．

### 仕組み
ネットワークに流れるパケットを解析し，あらかじめ登録された mac アドレスと紐づくユーザーの名前をアプリ内で管理，そしてhtmlにちょいちょいと書き出す仕組みです．

### 使い方

1. main.go でハードコーディングしている `device` に対して適切な network interface を書いてください(mac だと ifconfig , linux だと ip address でいけます)．

2. mac アドレスを登録しましょう． curl からいけます．

```
curl --header "Content-Type:application/json" \
--request POST \
--data '{"name":"issei", "addresses":["ee:ee:ee:ee:ee:ee"]}' \
http://localhost/register
```

3. build と実行を行います．network interface をキャプチャするために sudo で実行しています．

```
go build -o hoge
sudo ./hoge
```

4. web(http://localhost) から簡単に在室状況が確認できます．


### 注意

どこからでも mac アドレスを登録できるようにしているので，悪意のあるユーザーから攻撃される恐れがあります．
