BIN_DIR=bin
TARGET=$(BIN_DIR)/bookstore-manager
ADMIN_TARGET=$(BIN_DIR)/admin-manager
SRC=cmd/bookstore-manager.go

.PHONY: all bookstore-manager clean

all: bookstore-manager

bookstore-manager:
	@mkdir -p $(BIN_DIR)
	go build -o $(TARGET) $(SRC)

clean:
	rm -rf $(BIN_DIR)
