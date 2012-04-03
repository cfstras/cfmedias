#ifndef MAIN_H
#define MAIN_H

#include <string>

struct Settings {
    bool nogui;
    bool headless;

};

bool parseArgs(Settings, int, char*[]);

void startGUI();

void startCLI();

void startCore();

void startModule(QString name);

void stopGUI();

void stopCLI();

void stopCore();

void stopModule(QString name);

void printHelp();



#endif // MAIN_H
