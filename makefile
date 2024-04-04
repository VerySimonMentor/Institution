TARGET_NAME=institution

institution: clean
	go clean ./...
	go build -o ${TARGET_NAME} ./main.go

clean:
	@if [ -f "${TARGET_NAME}" ]; then \
		rm ${TARGET_NAME}; \
	fi