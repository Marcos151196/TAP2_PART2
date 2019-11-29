package main

import (
	"fmt"
	"os"
	"os/exec"

	aws "github.com/aws/aws-sdk-go/aws"
	session "github.com/aws/aws-sdk-go/aws/session"
	s3manager "github.com/aws/aws-sdk-go/service/s3/s3manager"
	log "github.com/sirupsen/logrus"
	viper "github.com/theherk/viper"
)

var sess *session.Session = session.Must(session.NewSessionWithOptions(session.Options{
	SharedConfigState: session.SharedConfigEnable,
}))

func main() {
	var args []string
	var err error
	var output string

	// CONFIG FILE
	cfgFile := "config/config.toml"
	viper.SetConfigFile(cfgFile)
	if err := viper.ReadInConfig(); err != nil {
		log.Errorf("Unable to read config from file %s: %v", cfgFile, err)
		os.Exit(1)
	} else {
		log.Infof("Read configuration from file %s", cfgFile)
	}

	// REMOVE TASK3 OUTPUT FROM LOCAL IF IT ALREADY EXISTS
	args = []string{"-rf", viper.GetString("localoutputdirectory") + "output.task3"}
	output, err = RunCMD("rm", args, true)
	if err != nil {
		log.Warnf("Could not remove task3 output from local: %v", err)
	} else {
		log.Debugf("Result: %s", output)
	}

	// REMOVE TASK3 OUTPUT FROM DFS IF IT ALREADY EXISTS
	args = []string{"dfs", "-rm", "-r", viper.GetString("dfsoutputdirectory") + "output.task3"}
	output, err = RunCMD("/home/ubuntu/hdfs", args, true)
	if err != nil {
		log.Warnf("Could not remove task3 output from dfs: %v", err)
	} else {
		log.Debugf("Result: %s", output)
	}

	// RUN TASK 3
	args = []string{"streaming", "-input", viper.GetString("ngrams3bucket"), "-output", "output.task3", "-mapper", "\"" + viper.GetString("task1binary") + " -task 0 -phase map" + "\"", "-reducer", "\"" + viper.GetString("task1binary") + " -task 0 -phase reduce" + "\"", "-io", "typedbytes", "-inputformat", "SequenceFileInputFormat"}
	output, err = RunCMD("/home/ubuntu/hadoop/bin/mapred", args, true)
	if err != nil {
		log.Errorf("Could not run task3 mapred: %v", err)
		return
	} else {
		log.Debugf("Result: %s", output)
	}

	// COPY TASK3 OUTPUT TO LOCAL
	args = []string{"fs", "-copyToLocal", viper.GetString("dfsoutputdirectory") + "output.task3", viper.GetString("localoutputdirectory")}
	output, err = RunCMD("/home/ubuntu/hadoop/bin/hadoop", args, true)
	if err != nil {
		log.Errorf("Could not copy task3 output to local directory: %v", err)
		return
	} else {
		log.Debugf("Result: %s", output)
	}

	// UPLOAD OUTPUT TO S3 BUCKET
	fileName := fmt.Sprintf("%s/output.task3/part-00000", viper.GetString("localoutputdirectory"))
	f, err := os.Open(fileName)
	if err != nil {
		log.Errorf("Failed to open file %q, %v", fileName, err)
		return
	}
	err = UploadFileToS3(f)
	if err != nil {
		f.Close()
		log.Errorf("Could not upload file to S3: %v", err)
		return
	}
	err = f.Close()
	if err != nil {
		log.Errorf("Could not save and close file: %v", err)
		return
	}

}

// RunCMD is a simple wrapper around terminal commands
func RunCMD(path string, args []string, debug bool) (out string, err error) {
	cmd := exec.Command(path, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	_ = cmd.Run()
	return
}

// UPLOAD TASK3 OUT FILE TO S3
func UploadFileToS3(f *os.File) error {
	uploadPath := "MostCommonUnigramsPerDecade/task3_output.txt"
	uploader := s3manager.NewUploader(sess)
	_, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String("tap2"),
		Key:    aws.String(uploadPath),
		Body:   f,
	})
	if err != nil {
		return fmt.Errorf("Failed to upload file, %v", err)
	}
	log.Infof("File uploaded to, %s\n", uploadPath)
	return nil
}
