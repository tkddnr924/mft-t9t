# mft-t9t

github.com/t9t/gomft를 사용하여 개발. $MFT 덤프하는 기능

## TODO
1. $MFT NFTF 시간 속성 등 추가
2. Windows XP 기능 호환
   - 해당 기능은 [tkddnr924/gomft](https://github.com/tkddnr924/gomft) 에서 진행 예정

## Docker

### pull
```shell
docker pull golang:1.19.6
```

### build
```shell
docker build .
```

### file dump
```shell
docker cp <container_name>:/go/src/app/mft-t9t.exe <output_directory>
```


