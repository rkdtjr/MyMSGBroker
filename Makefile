# Makefile for my_msg_broker

# 컴파일러 및 옵션
CC = g++
CFLAGS = -Wall -O2 -Iinclude
LDFLAGS = 

# 디렉토리
SRC_DIR = src
OBJ_DIR = obj
BUILD_DIR = build
TARGET = $(BUILD_DIR)/my_msg_broker

# 소스 및 오브젝트 파일
SRCS = $(wildcard $(SRC_DIR)/*.cpp)
OBJS = $(patsubst $(SRC_DIR)/%.cpp,$(OBJ_DIR)/%.o,$(SRCS))

# 기본 타겟
all: $(BUILD_DIR) $(OBJ_DIR) $(TARGET)

# 빌드 디렉토리 생성
$(BUILD_DIR):
	mkdir -p $(BUILD_DIR)

$(OBJ_DIR):
	mkdir -p $(OBJ_DIR)

# 타겟 빌드
$(TARGET): $(OBJS)
	$(CC) $(LDFLAGS) -o $@ $^

# 오브젝트 파일 빌드
$(OBJ_DIR)/%.o: $(SRC_DIR)/%.cpp
	$(CC) $(CFLAGS) -c $< -o $@

# 클린
clean:
	rm -rf $(OBJ_DIR)/* $(BUILD_DIR)/*

