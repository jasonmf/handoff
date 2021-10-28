WASM_EXEC_SRC=$(shell tinygo env TINYGOROOT)/targets/wasm_exec.js

WORK=work
ASSETS=${WORK}/assets

BIN=${ASSETS}/handoff.wasm
HTML=${ASSETS}/index.html
SERVER=${WORK}/server
WASM_EXEC=${ASSETS}/wasm_exec.js

default: wasm

wasm: ${BIN} ${HTML} ${WASM_EXEC}
service: wasm ${SERVER}

container: service
	docker build -t handoff:latest .

${ASSETS}:
	mkdir -p ${ASSETS}

${BIN}: main.go ../../*.go ${ASSETS}
	# GOOS=js GOARCH=wasm go build -o $@
	tinygo build --no-debug -o $@ -target wasm .

${WASM_EXEC}: ${WASM_EXEC_SRC} ${ASSETS}
	cp ${WASM_EXEC_SRC} $@

${HTML}: index.html ${ASSETS}
	cp index.html ${ASSETS}/

${SERVER}: ../server/main.go
	CGO_ENABLED=0 go build -tags netgo -o $@ -ldflags="-s -w" $<

clean:
	rm -rf ${WORK}