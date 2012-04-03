#include "cli.h"
#include <iostream>
#include <QString>
using namespace std;

CLI::CLI(QObject *parent) :
    QThread(parent)
{

}

void CLI::run() {
    //TODO query methods
    //TODO evaluate commands
    bool again=true;

    string input;
    do {
        getline(cin, input);

    } while (again);

}
