# Intelligent Monitoring System

<h1 align="center"><img alt="gopher-camera" src="documentation/gopher.png"/></h1>

Its a Simple free Software to make a monitoring system based on Face Recognition using standars IP cameras. They only need to be able to send images to an FTP Server.

I developed it implementing Hexagonal Architecture (or Ports & Adapters architecture ). You can see on the diagram that is very simple to add anothers adapters to process images. In fact, you can easily change AWS Recognize component with another.


## Architecture
![Architecture](documentation/Architecture.png)

## UML Graphics

[PlantUML](https://github.com/arielbrizi/go-intelligent-monitoring-system/tree/develop/documentation/puml) 


## API Documentation - Swagger
[editor.swagger.io](https://editor.swagger.io/?url=https://raw.githubusercontent.com/arielbrizi/go-intelligent-monitoring-system/develop/docs/swagger.yaml) (It's also automatically published on .../swagger/index.html)
- Documentation about configuration and service annotations: https://github.com/swaggo/gin-swagger 

## Installation

### - Pre-Requisites

* Install FTP Server. Recommended: [Install WingFTP in very simple steps](documentation/wingFTP/README.md)

* Set the Ip of the FTP server in your IP cameras. Some models have movement detection and send the image captured to the destination defined.

* Set Up Environment variables in "docker-compose.yml" file:

```
    ## General Variables
    - AWS_ACCESS_KEY_ID=XXX

    - AWS_SECRET_ACCESS_KEY=XXX

    - FTP_DIRECTORY=/home/ariel/fotos_pasillo/ #Directory where your FTP Server save your images.

    - CAMARA_DOMAIN=camarasilvia # In case of S3 it's used to define the "bucket" name

    - QUEUE_TOPIC=images #Don't need to change It

    - QUEUE_BROKER_LIST=broker:9092 #Don't need to change It

    - SNS_TOPIC=arn:aws:sns:us-east-1:491728392546:monitoringSystem-eMail #Previously configured On AWS

    - AUTHORIZED_FACES_DIRECTORY=/home/ariel/fotos_pasillo/authorized_faces/ #Faces you want to define as authorized in JPG files  

    - TELEGRAM_BOT_TOKEN=xxxxxxxxxxxxxxxxX
    - TELEGRAM_CHANNEL=@notifims #The telegram username/channel will be used to send notifications like unauthorized person detected
    - TELEGRAM_USER=@arielbrizi


    ## AWS variables
    - AWS_S3_BUCKET_POLICY={"Version":"2012-10-17","Statement":[{"Sid":"PublicRead","Effect":"Allow","Principal":"*","Action":["s3:GetObject","s3:GetObjectVersion"],"Resource":"arn:aws:s3:::camarasilvia/*"}]} ## Where camarasilvia must be same as CAMARA_DOMAIN variable
```
### - Run "sudo docker-compose up -d" from your go-intelligent-monitoring-system directory. If you have problems check 'service docker status'. If it's neccesary run 'sudo service docker start'
