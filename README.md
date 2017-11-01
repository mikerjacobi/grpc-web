
#web-client

`npm install`
`webpack`

#protoc

`protoc --plugin=protoc-gen-ts=./node_modules/.bin/protoc-gen-ts --js_out=import_style=commonjs,binary:generated --service=true --ts_out=service=true:generated -I ../pb ../pb/app.proto`

#standup server

`sudo grpcwebproxy --server_tls_cert_file=tls.crt --server_tls_key_file=tls.key --backend_addr=localhost:5051 --backend_tls_noverify --backend_tls true`


