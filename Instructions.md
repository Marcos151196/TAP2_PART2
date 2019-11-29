# TASK 3 (EC2 PSEUDODISTRIBUTED)

## Important commands

``` bash
/home/ubuntu/hadoop/bin/hdfs dfs -rm -r output.task3
rm -rf /tmp/hadoop*
hdfs namenode -format
start-dfs.sh
start-yarn.sh
```

``` bash
/home/ubuntu/hadoop/bin/mapred streaming -input s3a://datasets.elasticmapreduce/ngrams/books/20090715/spa-all/1gram/data -output output.task3 -mapper "/home/ubuntu/task1 -task 0 -phase map" -reducer "/home/ubuntu/task1 -task 0 -phase reduce" -io typedbytes -inputformat SequenceFileInputFormat
```

## Launch task 3
``` bash
cd /home/ubuntu/TAP2_PART2/launchtask3
./launchtask3
```
## Check output in DFS

``` bash
/home/ubuntu/hadoop/bin/hdfs dfs -cat output.task3/part-00000
```
