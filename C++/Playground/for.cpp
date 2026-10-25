#include <iostream>
#include <cmath>

using namespace std;

int main() {
    int n,m,h,ans = 1;
    cin >> n >> m;
    h = pow(n-1, 2) + m;
    for (int i = 1; i < h; i++) {
        ans += 2;
    }
    cout << ans << "\n";
    return 0;
}
