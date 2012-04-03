#include <QString>
#include <iostream>
#include <string>
#include "gui.h"
#include "main.h"

using namespace std;

Settings settings;

int main(int argc, char *argv[])
{
    cout << endl<<"Welcome to cfmedias."<<endl;

    settings=Settings();

    //QApplication a(argc, argv);
    if(! parseArgs(settings, argc, argv)) {
        return 0;
    }

    if(! settings.nogui) {
        startGUI();
    } else {
        startCLI();
    }

    //return a.exec();
    //TODO don't exit now
    return 0;
}

bool parseArgs(Settings settings, int argc, char *argv[]) {
    cout << "argc: "<<argc<<" argv:"<< " " ;
    for(int i=0;i<argc;i++) {
        cout << argv[i]<<" ";
    } cout<<endl;


    for(int i=1;i<argc;i++) {
        QString arg=QString(argv[i]);

        if(arg.compare("--nogui")==0 || arg.compare("-n")==0) {
            settings.nogui=true;
        } else if (arg.compare("--help")==0 || arg.compare("-h")==0) {
            printHelp();
            return false;
        } else {
            cout << "Unknown parameters.";
            printHelp();
            return false;
        }
    }
    return true;
}

void printHelp() {
    cout <<"cmd parameters:"<<endl;
    cout << "--help\t -h\t print help and exit"<<endl;
    cout << "--nogui\t -n\t run in CLI only mode"<<endl;
}

void startGUI() {
    GUI w;
    w.show();
}

void startCLI() {
    //start expecting input on cmd or net


}

void startCore(){

}

void startModule(QString name){

}

void stopGUI(){

}

void stopCore(){

}

void stopModule(QString name){

}
