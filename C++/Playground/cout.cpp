// 行数：总个数
// 1:1
// 2:4
// 3:9
// 4:14
// n:pow(n,2)

// 第n个：与n的关系
// 1:1
// 2:3
// 3:5
// 4:7
// 5:9
// n:2 * n - 1

#include <iostream>
#include <cmath>

using namespace std;

int main() {
    int n,m;
    cin >> n >> m;
    int f = pow(n - 1, 2) + m;
    cout << 2 * f - 1;
    return 0;
}
