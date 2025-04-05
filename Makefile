TARGET=cloudstatus

UPX := $(shell command -v upx 2>/dev/null)

all: frontend build compress

frontend:
	cd cloudstatus-fe && pnpm install && pnpm build

build:
	go build -ldflags "-s -w" -o $(TARGET) .

compress: $(TARGET)
ifdef UPX
	$(UPX) $(TARGET)
else
	@echo "UPX not found, skipping compression."
endif

clean:
	rm -f $(TARGET)