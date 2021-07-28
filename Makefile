# https://github.com/aws/amazon-ec2-metadata-mock
run:
	 AWS_EC2_METADATA_SERVICE_ENDPOINT=http://localhost:1338/latest/meta-data go run main.go view

build-docker:
	docker build . -t emr-idle-check