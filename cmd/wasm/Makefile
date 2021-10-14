WASM_EXEC_SRC=$(shell go env GOROOT)/misc/wasm/wasm_exec.js

OUTDIR=assets

BIN=${OUTDIR}/handoff.wasm
HTML=${OUTDIR}/index.html
WASM_EXEC=${OUTDIR}/wasm_exec.js


wasm: ${BIN} ${HTML} ${WASM_EXEC}

${OUTDIR}:
	mkdir ${OUTDIR}

${BIN}: main.go ../../*.go ${OUTDIR}
	GOOS=js GOARCH=wasm go build -o $@

${WASM_EXEC}: ${WASM_EXEC_SRC} ${OUTDIR}
	cp ${WASM_EXEC_SRC} $@

${HTML}: index.html ${OUTDIR}
	cp index.html ${OUTDIR}/

clean:
	rm -rf ${OUTDIR}