# https://github.com/FlixMa/webauthn-pocketbase

include caddy.MK

print:

## Install tools
dep-tools: caddy-dep

## Run backend
run-back:
	cd backend && go run . serve
	# http://127.0.0.1:8090
	# http://localhost:8090/_/
	# user: gedw99@gmail.com
	# password: password10

	# then setup fields as per readme
	
## Run frontend
run-front:
	cd app && npm install
	cd app && npm run dev
	# http://localhost:5173
## Run Caddy reverse proxy over https
caddy:
	$(MAKE) caddy-server-run
	# https://localhost
	# https://hello.localhost
	# https://pb.localhost
	# https://pb-admin.localhost/_/