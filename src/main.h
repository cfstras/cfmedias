#ifndef MAIN_H
#define MAIN_H

#include <string>

struct Settings {
    bool nogui;
    bool headless;

};

bool parseArgs(Settings, int, char*[]);

void startGUI();

void startCore();

void startModule(std::string);

void stopGUI();

void stopCore();

void stopModule();

void printHelp();



#endif // MAIN_H
