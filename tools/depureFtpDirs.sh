dateMenos1=$(TZ=America/Argentina/Buenos_Aires date -d '1 day ago' +'%Y%m%d')
find /home/ec2-user/wftpserver/pasillo/$dateMenos1* -type d  -exec rm -rf {} \;

msg="script executed: folder $dateMenos1 "
echo "$msg"

#Configure a cron: https://kvz.io/schedule-tasks-on-linux-using-crontab.html
#Execute every day. Minute 0 of hour 1
#1 1 * * * /home/ec2-user/wftpserver/depureFtpDirs.sh >> /home/ec2-user/wftpserver/script_output.log 2>&1

# To see what crontabs are currently running on your system, you can open a terminal and run:
#$ sudo crontab -l
# To edit the list of cronjobs you can run:
#$ sudo crontab -e



