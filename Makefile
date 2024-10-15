.phony: build run push
build:
	@go build bin/app
run:build
	@./bin/app

push:
	@echo "git inialised..."
	@git init
	@echo "staging all files..."
	@git add .
	@echo "commiting all files..."
	@git commit -s -m "$(msg)"
	@echo "pushing all files to git repository..."
	@git push origin main
