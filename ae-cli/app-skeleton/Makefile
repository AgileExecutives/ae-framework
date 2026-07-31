## Makefile for bringing up and managing the project
## Usage examples:
##   make help
##   make init      # initial build + start server + build api client and app
##   make up        # start server (background) then build api client + app
##   make dev       # start server (background) and run frontend dev (foreground)

PROJECT ?= $(notdir $(CURDIR))
BACKEND_DIR ?= backend
API_CLIENT_DIR ?= api-client-$(PROJECT)
APP_DIR ?= $(PROJECT)-app

SERVER_BIN := $(BACKEND_DIR)/server
LOG_DIR := logs
PID_FILE := .server.pid

SERVER_PORT ?= 8080
WAIT_RETRIES ?= 30
WAIT_SLEEP ?= 1

HEALTH_CHECK_URL := http://localhost:$(SERVER_PORT)/api/v1/health

.PHONY: help init up dev build build-backend run-server stop-server restart-server build-api-client build-app wait-server logs clean

help:
	@echo "Makefile targets:"
	@echo "  init            - initial bring-up: build backend, start server, build api client and app"
	@echo "  up              - start server (if not running) and build api client + app build"
	@echo "  dev             - start server (background) and run frontend dev (foreground)"
	@echo "  build           - build backend + api client + app (production builds where applicable)"
	@echo "  build-backend   - run 'go mod tidy' and build the backend server"
	@echo "  run-server      - start the backend server in background (logs/$(notdir $(SERVER_BIN))).pid"
	@echo "  stop-server     - stop the background server started by run-server"
	@echo "  build-api-client- build API client (npm install && npm run build)"
	@echo "  build-app       - build the app (npm install && npm run build)"
	@echo "  dev             - start server (background) and run 'npm run dev' in $(APP_DIR)"
	@echo "  logs            - tail server logs"
	@echo "  clean           - remove build artifacts and logs"

## Initial bring up: build backend, start server, wait, build clients
init: build-backend run-server wait-server build-api-client build-app

## Start server, then build clients. Useful after changes.
up: run-server wait-server build-api-client build-app

## Start server and run frontend dev (foreground)
dev: run-server wait-server
	@echo "Starting frontend dev in $(APP_DIR) (foreground)"
	@cd $(APP_DIR) && npm install && npm run dev

build: build-backend build-api-client build-app

## Backend build
build-backend:
	@echo "==> Backend: go mod tidy"
	@cd $(BACKEND_DIR) && go mod tidy
	@echo "==> Backend: go build"
	@cd $(BACKEND_DIR) && go build -o ../$(notdir $(SERVER_BIN)) ./...

## Start server in background (writes PID to $(PID_FILE))
run-server: build-backend
	@mkdir -p $(LOG_DIR)
	@if [ -f $(PID_FILE) ] && kill -0 `cat $(PID_FILE)` 2>/dev/null; then \
		echo "Server already running (pid=`cat $(PID_FILE)`)"; \
		exit 0; \
	fi
	@echo "Starting server..."
	@nohup ./$(notdir $(SERVER_BIN)) > $(LOG_DIR)/server.log 2>&1 & echo $$! > $(PID_FILE)
	@sleep 0.5
	@echo "Server started, pid=`cat $(PID_FILE)`"

stop-server:
	@if [ -f $(PID_FILE) ]; then \
		pid=`cat $(PID_FILE)`; \
		if kill -0 $$pid 2>/dev/null; then \
			kill $$pid && echo "Stopped server ($$pid)" || echo "Failed to stop $$pid"; \
		fi; \
		rm -f $(PID_FILE); \
	else \
		echo "No pidfile ($(PID_FILE)) found"; \
	fi

restart-server: stop-server run-server

## Wait until the server responds on the given port (basic health check)
wait-server:
	@i=0; \
	while [ $$i -lt $(WAIT_RETRIES) ]; do \
		if curl -sSf "$(HEALTH_CHECK_URL)" >/dev/null 2>&1; then \
			echo "Server is responding"; exit 0; \
		fi; \
		sleep $(WAIT_SLEEP); i=$$((i+1)); \
	done; \
	echo "Timed out waiting for server after $(WAIT_RETRIES) attempts"; exit 2

logs:
	@mkdir -p $(LOG_DIR)
	@tail -n +1 -f $(LOG_DIR)/server.log

## Build api client
build-api-client:
	@if [ -d $(API_CLIENT_DIR) ]; then \
		cd $(API_CLIENT_DIR) && npm install && npm run build; \
	else \
		echo "API client dir '$(API_CLIENT_DIR)' not found, skipping"; \
	fi

## Build app (production build)
build-app:
	@if [ -d $(APP_DIR) ]; then \
		cd $(APP_DIR) && npm install && npm run build; \
	else \
		echo "App dir '$(APP_DIR)' not found, skipping"; \
	fi

clean:
	@echo "Cleaning logs and server binary"
	@rm -f $(PID_FILE)
	@rm -rf $(LOG_DIR)
	@rm -f $(notdir $(SERVER_BIN))
