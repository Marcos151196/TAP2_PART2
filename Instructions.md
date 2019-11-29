# TASK 3 (pseudo-distributed on AWS)
hadoop 3.2.1
Java 8
Ubuntu 18.04 en instancia EC2 t2.medium


## 1 - Añadir soporte para compresión LZO
Para hacer el maven install hace falta tener gcc instalado.

```bash
sudo apt-get install build-essential liblzo2-dev maven
git clone https://github.com/twitter/hadoop-lzo.git
cd hadoop-lzo/
mvn clean package
sudo cp ~/hadoop-lzo/target/hadoop-lzo-0.4.21-SNAPSHOT.jar /home/ubuntu/hadoop/share/hadoop/mapreduce/lib 
```

## 2 - Añadir librerias de AWS para Hadoop

```bash
wget http://central.maven.org/maven2/org/apache/hadoop/hadoop-aws/3.2.1/hadoop-aws-3.2.1.jar
wget http://central.maven.org/maven2/com/amazonaws/aws-java-sdk-bundle/1.11.683/aws-java-sdk-bundle-1.11.683.jar

mv hadoop-aws-3.2.1.jar /home/ubuntu/hadoop/share/hadoop/common
mv aws-java-sdk-bundle-1.11.683.jar /home/ubuntu/hadoop/share/hadoop/common
```

## 3 - Configurar XMLs

### core-site.xml

```xml
<configuration>
	<property>
		<name>fs.defaultFS</name>
		<value>hdfs://localhost:9000</value>
	</property>
	<property>
		<name>fs.s3a.access.key</name>
		<value>AKIAQXBZZNNDWV6V43ED</value>
	</property>
	<property>
		<name>fs.s3a.secret.key</name>
		<value>NVCvQGLHoNoqymQis/PMf/FjwhPRV63oknvmMDvk</value>
	</property>
	<property>
		<name>fs.AbstractFileSystem.s3a.imp</name>
		<value>org.apache.hadoop.fs.s3a.S3A</value>
	</property>
	<property>
		<name>fs.s3a.endpoint</name>
		<value>s3.us-east-1.amazonaws.com</value>
	</property>
	<property>
		<name>io.compression.codecs</name>
		<value>org.apache.hadoop.io.compress.GzipCodec, org.apache.hadoop.io.compress.DefaultCodec, org.apache.hadoop.io.compress.BZip2Codec, com.hadoop.compression.lzo.LzoCodec, com.hadoop.compression.lzo.LzopCodec</value>
	</property>
	<property>
		<name>io.compression.codec.lzo.class</name>
	<value>com.hadoop.compression.lzo.LzoCodec</value>
	</property>
	<property>
		<name>fs.s3a.block.size</name>
		<value>3221225472</value>
	</property>
</configuration>
```

### hdfs-site.xml
```xml
<configuration>
	<property>
      <name>dfs.replication</name>
      <value>1</value>
  </property>
</configuration>
```

### yarn-site.xml
```xml
<configuration>
	<property>
		<name>yarn.nodemanager.aux-services</name>
		<value>mapreduce_shuffle</value>
  	</property>
  	<property>
		<name>yarn.nodemanager.env-whitelist</name>
		<value>JAVA_HOME,HADOOP_COMMON_HOME,HADOOP_HDFS_HOME,HADOOP_CONF_DIR,CLASSPATH_PREPEND_DISTCACHE,HADOOP_YARN_HOME,HADOOP_MAPRED_HOME</value>
  	</property>
  	<property>
		<name>yarn.nodemanager.vmem-check-enabled</name>
		<value>false</value>
  	</property>
</configuration>
```

### mapred-site.xml
```xml
<configuration>
	<property>
		<name>mapreduce.framework.name</name>
		<value>yarn</value>
	</property>
</configuration>
```

### ~/.bashrc
```bash
export PDSH_RCMD_TYPE=ssh

# JAVA
export JAVA_HOME=/usr/lib/jvm/java-8-openjdk-amd64
export PATH=$PATH:$JAVA_HOME/bin

# HADOOP
export HADOOP_HOME=/home/ubuntu/hadoop
export PATH=$PATH:$HADOOP_HOME/bin:$HADOOP_HOME/sbin


export HADOOP_CLASSPATH=$HADOOP_CLASSPATH:/home/ubuntu/hadoop/lib:/home/ubuntu/hadoop-lzo/target/hadoop-lzo-0.4.21-SNAPSHOT.jar
export JAVA_LIBRARY_PATH=/home/ubuntu/hadoop-lzo/target/native/Linux-amd64-64:/home/ubuntu/hadoop/lib/native

# GO
export GOROOT=/usr/local/go
export PATH=$GOPATH/bin:$GOROOT/bin:$PATH
```

### hadoop-env.sh
```bash
export JAVA_HOME=/usr/lib/jvm/java-8-openjdk-amd64/
export HADOOP_CLASSPATH=$HADOOP_CLASSPATH:/home/ubuntu/hadoop/lib:/home/ubuntu/hadoop-lzo/target/hadoop-lzo-0.4.21-SNAPSHOT.jar
export JAVA_LIBRARY_PATH=/home/ubuntu/hadoop-lzo/target/native/Linux-amd64-64:/home/ubuntu/hadoop/lib/native
export HADOOP_HOME=/home/ubuntu/hadoop
```
## 4 - Formatear hdfs

```bash
rm -rf /tmp/hadoop*
~/hadoop/bin/hdfs namenode -format
```

## 5 - Lanzar mapred y comprobar resultados

```bash
start-dfs.sh
start-yarn.sh

hdfs dfs -rm -r output.task3

/home/ubuntu/hadoop/bin/mapred streaming -input s3a://datasets.elasticmapreduce/ngrams/books/20090715/spa-all/1gram/data -output output.task3 -mapper "/home/ubuntu/task1 -task 0 -phase map" -reducer "/home/ubuntu/task1 -task 0 -phase reduce" -io typedbytes -inputformat SequenceFileInputFormat

hdfs dfs -cat output.task3/part-00000
```
