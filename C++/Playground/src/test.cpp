#include <iostream>
#include <string>
using namespace std;

int main()
{
    string name;
    int time;
    cin >> name >> time;
    for (int i = 0; i < time; i++)
    {
        cout << "Hello, " << name << "!" << endl;
    }
    return 0;
}