1. Написать программу, которая использует мьютекс для безопасного доступа к данным
   из нескольких потоков. Выполните трассировку программы
> go run mutex_workers.go 2>trace_mutex_workers.out

> go tool trace trace_mutex_workers.out

2. Написать многопоточную программу, в которой будет использоваться явный вызов
   планировщика. Выполните трассировку программы:
> go run pitstop.go 2>pitstop.out
```
#5: 1634506
#8: 876814
#7: 1663825
#3: 1745733
#1: 1751412
#0: 873842
#2: 847500
#4: 882912
#6: 878238
#9: 1584829
```

>  go tool trace pitstop.out

3. Смоделировать ситуацию “гонки”, и проверить программу на наличии “гонки”
> go run -race race_example.go
```==================
WARNING: DATA RACE
Read at 0x00c000018130 by goroutine 7:
main.main.func1()
...race_example.go:24 +0xf9

Previous write at 0x00c000018130 by goroutine 1926:
sync/atomic.AddInt32()
/usr/local/go/src/runtime/race_amd64.s:292 +0xb
main.main.func2()
...race_example.go:29 +0x91

Goroutine 7 (running) created at:
main.main()
...race_example.go:19 +0x176

Goroutine 1926 (finished) created at:
main.main()
...race_example.go:26 +0x1ac
==================
final scores =  1
Found 1 data race(s)
exit status 66
```
