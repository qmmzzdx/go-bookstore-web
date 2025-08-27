BIN_DIR=bin
TARGET=$(BIN_DIR)/bookstore-manager
ADMIN_TARGET=$(BIN_DIR)/admin-manager
SRC=cmd/bookstore-manager/bookstore-manager.go
ADMIN_SRC=cmd/admin-manager/admin-manager.go

.PHONY: all bookstore-manager admin-manager clean

all: bookstore-manager admin-manager

bookstore-manager:
	@mkdir -p $(BIN_DIR)
	go build -o $(TARGET) $(SRC)

admin-manager:
	@mkdir -p $(BIN_DIR)
	go build -o $(ADMIN_TARGET) $(ADMIN_SRC)

clean:
	rm -rf $(BIN_DIR)
