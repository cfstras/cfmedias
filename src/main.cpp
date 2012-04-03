#include <QtGui/QApplication>
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

    //GUI w;
    //w.show();
    
    //return a.exec();
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
}

void printHelp() {
    cout <<"cmd parameters:"<<endl;

    cout << "--help\t -h\t print help and exit"<<endl;
    cout << "--nogui\t -n\t run in CLI only mode"<<endl;
}

void startGUI() {

}

void startCore(){

}

void startModule(string module){

}

void stopGUI(){

}

void stopCore(){

}

void stopModule(){

}
