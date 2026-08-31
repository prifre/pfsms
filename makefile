# Variabler
VERSION=0.1.4
BINARY_NAME=pfsms
BUILD_DIR_LINUX=pfsms_linux
BUILD_DIR_WIN=pfsms_windows
BUILD_DIR_MAC=pfsms_darwin
OUTPUT_DIR=public_html

# Standardmål
.PHONY: all
all: clean pack

# 1. Skapa mappar och kompilera lokalt för Linux & Windows
.PHONY: build-local
build-local:
	@echo "Kompilerar lokalt för Linux och Windows..."
	@mkdir -p $(BUILD_DIR_LINUX) $(BUILD_DIR_WIN) $(BUILD_DIR_MAC)
	GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR_LINUX)/$(BINARY_NAME) .
	CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc go build -ldflags="-s -w" -o $(BUILD_DIR_WIN)/$(BINARY_NAME).exe .

# 2. Skicka till GitHub, vänta på Actions-bygget och hämta Mac
.PHONY: fetch-artifacts
fetch-artifacts: build-local
	@echo "Pushar till GitHub och väntar på molnbyggena..."
	@git add .
	@git commit -m "Auto-build release v$(VERSION)" || true
	@git push
	@echo "Väntar på att bygg-workflowet ska registreras och starta..."
	@sleep 5
	@gh run watch $$(gh run list --workflow="build.yml" --limit 1 --json databaseId -q '.[0].databaseId')
	@echo "Hämtar artifacts från GitHub..."
	@gh run download --name pfsms-macOS --dir $(BUILD_DIR_MAC) || echo "Varning: Kunde inte hämta macOS artifact."

# 3. Kopiera readme och packa alla plattformar till public_html
.PHONY: pack
pack: fetch-artifacts
	@echo "Kopierar readme.md till alla byggmappar..."
	@cp readme.md $(BUILD_DIR_LINUX)/ 2>/dev/null || true
	@cp readme.md $(BUILD_DIR_WIN)/ 2>/dev/null || true
	@cp readme.md $(BUILD_DIR_MAC)/ 2>/dev/null || true

	@echo "Packar filer till $(OUTPUT_DIR)..."
	@mkdir -p $(OUTPUT_DIR)

	# Packa Windows (.zip)
	@zip -j $(OUTPUT_DIR)/pfsms_$(VERSION).zip $(BUILD_DIR_WIN)/*

	# Packa Linux (.tar.gz)
	@tar -czvf $(OUTPUT_DIR)/pfsms_$(VERSION).tar.gz -C $(BUILD_DIR_LINUX) .

	# Packa macOS (.zip) om filer finns i mappen
	@if [ -n "$$(ls -A $(BUILD_DIR_MAC) 2>/dev/null)" ]; then \
		zip -j $(OUTPUT_DIR)/pfsms_$(VERSION)_darwin.zip $(BUILD_DIR_MAC)/* ; \
	else \
		echo "Hoppar över macOS zip (inga filer i $(BUILD_DIR_MAC))" ; \
	fi

# Rensa bygg- och paketmappar
.PHONY: clean
clean:
	@echo "Rensar gamla byggfiler och paket..."
	@rm -rf $(BUILD_DIR_LINUX) $(BUILD_DIR_WIN) $(BUILD_DIR_MAC)
	@rm -rf $(OUTPUT_DIR)/pfsms_*.zip $(OUTPUT_DIR)/pfsms_*.tar.gz