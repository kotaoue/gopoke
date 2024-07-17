# gopoke
My personal Pokédex

## Initialize
```bash
go run main.go --init
```

## Usage
```bash
go run main.go --height=170                       
ウツボット 170.0cm 15.5kg
ハンテール 170.0cm 27.0kg
エアームド 170.0cm 50.5kg
ジュカイン 170.0cm 52.2kg
ギルガルド 170.0cm 53.0kg
フリーザー 170.0cm 55.4kg
デオキシス 170.0cm 60.8kg
デスカーン 170.0cm 76.5kg
ゴルダック 170.0cm 76.6kg
バクフーン 170.0cm 79.5kg
```

``` bash
go run main.go --weight=80             
スターミー 110.0cm 80.0kg
マルノーム 170.0cm 80.0kg
ガルーラ 220.0cm 80.0kg
コータス 50.0cm 80.4kg
マギアナ 100.0cm 80.5kg
ゼブライカ 160.0cm 79.5kg
バクフーン 170.0cm 79.5kg
ヤドキング 200.0cm 79.5kg
シビルドン 210.0cm 80.5kg
ギギギアル 60.0cm 81.0kg
```

```bash
go run main.go --height=170 --weight=80
マルノーム 170.0cm 80.0kg
バクフーン 170.0cm 79.5kg
ゴルダック 170.0cm 76.6kg
デスカーン 170.0cm 76.5kg
エンペルト 170.0cm 84.5kg
リザードン 170.0cm 90.5kg
ゼブライカ 160.0cm 79.5kg
ゴーゴート 170.0cm 91.0kg
ゾロアーク 160.0cm 81.1kg
ヤドラン 160.0cm 78.5kg
```

```bash
go run main.go --name "%ドン" --limit=30 
ウオチルドン 200.0cm 175.0kg
ウツドン 100.0cm 6.4kg
カバルドン 200.0cm 300.0kg
グラードン 350.0cm 950.0kg
コライドン 250.0cm 303.0kg
サイドン 190.0cm 120.0kg
シビルドン 210.0cm 80.5kg
ジュラルドン 180.0cm 40.0kg
タンドン 30.0cm 12.0kg
トリトドン 90.0cm 29.9kg
ドサイドン 240.0cm 282.8kg
パッチルドン 230.0cm 150.0kg
ミライドン 350.0cm 240.0kg
ヤドン 120.0cm 36.0kg
リザードン 170.0cm 90.5kg
```

```bash
go run main.go --id=25
No.0025
ピカチュウ
分類: ねずみポケモン
高さ: 40.0cm    重さ: 6.0kg

尻尾を　立てて　まわりの　様子を
探っていると　ときどき
雷が　尻尾に　落ちてくる。
```

```bash
go run main.go --random
No.0302
ヤミラミ
分類: くらやみポケモン
高さ: 50.0cm    重さ: 11.0kg

洞窟の　暗闇に　潜む。
宝石を　食べているうちに
目が　宝石に　なってしまった。
```

## Thanks
- [PokéAPI](https://pokeapi.co/)
