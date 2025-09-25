# VozNotif
## App made to track Serbian Rail annoucments and email them to the user

# Usage
The app can either be compiled or ran.
```
go build [-o target] main.go
./target
```
OR
```
go run main.go
```

My usage is with the executable as a cron job. Define yours.
```
0-45/5 16 * * 2,3,4 ~/VozNotif/checkVozNotif >> notif.log 2>&1
```
 
