build:
	env GOOS=linux GOARCH=arm GOARM=7 go build -o app
	sshpass -p "alexx4355" scp app pi@10.0.1.9:/home/pi/raspi/app
	sshpass -p "alexx4355" ssh -t pi@10.0.1.9 "./raspi/app"