#ifndef CLI_H
#define CLI_H

#include <QThread>

class CLI : public QThread
{
    Q_OBJECT
public:
    explicit CLI(QObject *parent = 0);
    void run();
    
signals:
    
public slots:
    
};

#endif // CLI_H
