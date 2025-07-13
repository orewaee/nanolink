without-tls:
	go run cmd/cli/main.go run --port 2000 --host 127.0.0.1 redirect

with-tls:
	go run cmd/cli/main.go run --port 2000 --host 127.0.0.1 --cert-file certs/cert.crt --key-file certs/private.key --tls redirect
