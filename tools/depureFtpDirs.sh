#Delete all directories 2 days olders on /home/ariel/fotos_pasillo/* 
find /home/ec2-user/wftpserver/pasillo/* -type d -ctime +1 -exec rm -rf {} \;

#Configure a cron: https://kvz.io/schedule-tasks-on-linux-using-crontab.html
#Execute every day. Minute 0 of hour 1
#0 1 * * * /home/ec2-user/wftpserver/depureFtpDirs.sh >> /home/ec2-user/wftpserver/script_output.log 2>&1

# To see what crontabs are currently running on your system, you can open a terminal and run:
#$ sudo crontab -l
# To edit the list of cronjobs you can run:
#$ sudo crontab -e
