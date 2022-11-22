# px2csv-go

## Code

https://github.com/statisticssweden/PxWeb/discussions

https://github.com/statisticssweden/PxWin

https://github.com/statisticssweden/PxWeb

https://github.com/statisticssweden/PCAxis.Core

https://github.com/statisticssweden/PCAxis.Core/blob/master/PCAxis.Core/Parsers/PXFileParser.vb

https://github.com/statisticssweden/PCAxis.Core/blob/master/PCAxis.Core/Serializers/PXFileSerializer.vb

https://github.com/statisticssweden/PCAxis.Core/blob/master/PCAxis.Core/PaxiOM/Variable.vb

https://github.com/statisticssweden/PCAxis.Core/blob/master/PCAxis.Core/PaxiOM/Misc/PXKeywords.vb

### Documentation

https://www.scb.se/globalassets/vara-tjanster/px-programmen/px-file_format_specification_2013.pdf

https://www.scb.se/en/services/statistical-programs-for-px-files/px-file-format/

https://www.scb.se/en/services/statistical-programs-for-px-files/px-file-format/keywords-with-long-texts-256-per-line/

https://www.scb.se/px-en

https://groups.google.com/g/pcaxis

### Stat.fi materials

https://pxnet2.stat.fi/database/StatFin/StatFin_rap.csv

https://statfin.stat.fi/api1.html

https://statfin.stat.fi/database/StatFin/StatFin_rap.csv

### Stat.fi maintainer for PxEdit

https://groups.google.com/g/pcaxis/c/mBH_2jh5rN0/m/liFBobXWCwAJ


### Test data

https://github.com/search?q=AXIS-VERSION+KEYS+extension%3Apx&type=Code

https://github.com/ofurkusi/pxr-legacy/tree/master/pkg/tests

https://github.com/r-forge/pxr/tree/master/pkg/tests

https://github.com/cran/qmrparser/tree/master/inst/extdata

https://github.com/statisticssweden/PxWeb/tree/master/PXWeb/Resources/PX/Databases/Example

## Performance measuring

We want to use `bpftrace` to make sure that our reads and writes are page-sized.

```text
sudo bpftrace \
  -e 'tracepoint:syscalls:sys_exit_write /pid == cpid/ { @[comm] = hist(args->ret); }' \
  -c "./bin/pcaxis2parquet-linux-amd64 --px ./data/statfin_vtp_pxt_124l.px --csv /dev/null"
Attaching 1 probe...


@[pcaxis2parquet-]: 
[1K, 2K)               1 |                                                    |
[2K, 4K)               0 |                                                    |
[4K, 8K)           45127 |@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@|
[8K, 16K)              1 |                                                    |

sudo bpftrace \
  -e 'tracepoint:syscalls:sys_exit_write /pid == cpid/ { @[comm] = hist(args->ret); }' \
  -c "./bin/pcaxis2parquet-linux-amd64 --px ./data/statfin_vtp_pxt_124l.px --csv /dev/null"
Attaching 1 probe...


@[bpftrace]: 
[8, 16)                1 |@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@|

@[pcaxis2parquet-]: 
[0]                    1 |                                                    |
[1]                    0 |                                                    |
[2, 4)                 0 |                                                    |
[4, 8)                 0 |                                                    |
[8, 16)                1 |                                                    |
[16, 32)               0 |                                                    |
[32, 64)               0 |                                                    |
[64, 128)              0 |                                                    |
[128, 256)             0 |                                                    |
[256, 512)             0 |                                                    |
[512, 1K)              0 |                                                    |
[1K, 2K)               0 |                                                    |
[2K, 4K)               1 |                                                    |
[4K, 8K)           56545 |@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@|
```

```
hyperfine './bin/px2csv-linux-amd64 --px ./data/statfin_vtp_pxt_124l.px --csv /dev/null' '/home/mikael/.sdkman/candidates/java/22.3.r19-grl/bin/java -jar ../px2csv-java/build/libs/px2csv.jar ./data/statfin_vtp_pxt_124l.px /dev/null' '/home/mikael/.sdkman/candidates/java/19.0.1-amzn/bin/java -jar ../px2csv-java/build/libs/px2csv.jar ./data/statfin_vtp_pxt_124l.px /dev/null'

Benchmark 1: ./bin/px2csv-linux-amd64 --px ./data/statfin_vtp_pxt_124l.px --csv /dev/null
  Time (mean ± σ):      1.279 s ±  0.015 s    [User: 1.229 s, System: 0.054 s]
  Range (min … max):    1.262 s …  1.317 s    10 runs
 
Benchmark 2: /home/mikael/.sdkman/candidates/java/22.3.r19-grl/bin/java -jar ../px2csv-java/build/libs/px2csv.jar ./data/statfin_vtp_pxt_124l.px /dev/null
  Time (mean ± σ):      1.346 s ±  0.012 s    [User: 1.748 s, System: 0.181 s]
  Range (min … max):    1.328 s …  1.368 s    10 runs
 
Benchmark 3: /home/mikael/.sdkman/candidates/java/19.0.1-amzn/bin/java -jar ../px2csv-java/build/libs/px2csv.jar ./data/statfin_vtp_pxt_124l.px /dev/null
  Time (mean ± σ):      1.499 s ±  0.013 s    [User: 1.788 s, System: 0.124 s]
  Range (min … max):    1.484 s …  1.520 s    10 runs
 
Summary
  './bin/px2csv-linux-amd64 --px ./data/statfin_vtp_pxt_124l.px --csv /dev/null' ran
    1.05 ± 0.02 times faster than '/home/mikael/.sdkman/candidates/java/22.3.r19-grl/bin/java -jar ../px2csv-java/build/libs/px2csv.jar ./data/statfin_vtp_pxt_124l.px /dev/null'
    1.17 ± 0.02 times faster than '/home/mikael/.sdkman/candidates/java/19.0.1-amzn/bin/java -jar ../px2csv-java/build/libs/px2csv.jar ./data/statfin_vtp_pxt_124l.px /dev/null'

```