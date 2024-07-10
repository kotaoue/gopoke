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
$ go run main.go --height=170 --weight=80
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

## Thanks
- [PokéAPI](https://pokeapi.co/)
