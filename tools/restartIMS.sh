#Stop and delete all docker images
sudo docker stop $(sudo docker ps -a -q)
sudo docker rm $(sudo docker ps -a -q)

#Set code path
cd /home/pi/go/src/github.com/arielbrizi/go-intelligent-monitoring-system

#run docker
sudo docker-compose up -d ##The first time it's necessary "sudo docker-compose up --build -d"


