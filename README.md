# mpu_transfer


Upload/download large object to S3 or compatible cloud storage by leveraging MPU.

## Install

```
go get https://github.com/zhanga5/mpu_transfer.git
```

## Run

```
export REGION=<REGION>
export ACCESS_SERVER=<SERVER_ACCESS_ADDRESS>
export ACCESS_KEY=<MY_ACCESS_KEY>
export ACCESS_SECRET=<MY_ACCESS_SECRET>

./mpu_transfer -bucket zhanga5 -key 100MB -file 100MB -upload -part-size 10
./mpu_transfer -bucket zhanga5 -key 100MB -file 100MB -concurrent 10
```
