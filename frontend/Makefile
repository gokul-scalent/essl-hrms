# Mailora server deployment
deploy-web-mailora:
	rm -rf build
	npm run build
	rm -rf ~/Downloads/build.tar.xz
	tar --exclude=".git" -cJf ~/Downloads/build.tar.xz build/
	scp ~/Downloads/build.tar.xz root@185.230.138.44:/tmp/
	rm -rf ~/Downloads/build.tar.xz
	ssh root@185.230.138.44 sudo rm -rf /var/www/dashboard.mailora.media/html/*
	ssh root@185.230.138.44 sudo mv /tmp/build.tar.xz /var/www/dashboard.mailora.media/html
	ssh root@185.230.138.44 sudo tar -xf /var/www/dashboard.mailora.media/html/build.tar.xz -C /var/www/dashboard.mailora.media/html
	ssh root@185.230.138.44 sudo mv /var/www/dashboard.mailora.media/html/build/* /var/www/dashboard.mailora.media/html
	ssh root@185.230.138.44 sudo rm -rf /var/www/dashboard.mailora.media/html/build.tar.xz
	ssh root@185.230.138.44 sudo rm -rf /var/www/dashboard.mailora.media/html/build
	ssh root@185.230.138.44 sudo chown root:root /var/www/dashboard.mailora.media/html -R