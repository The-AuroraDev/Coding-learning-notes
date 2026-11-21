CXX      ?= g++
CXXFLAGS ?= -std=c++14 -Wall -Wextra -O2
OUT      := output

# 任意目标 hello，自动找 hello.cpp
%: %.cpp
	@mkdir -p $(OUT)
	$(CXX) $(CXXFLAGS) -o $(OUT)/$@ $<
	./$(OUT)/$@