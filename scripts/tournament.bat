@ECHO off

IF NOT DEFINED CONF ( SET CONF=conf\TestConf.yml )
IF NOT DEFINED DATABASE ( SET DATABASE=db\crobots.db )
SET SQL=%TMP%\%RANDOM%.SQL
SET COMMAND=%1

IF NOT DEFINED COMMAND (
    ECHO "Usage: tournament.bat [test|clean|reset|init|f2f|3vs3|4vs4|all]"
    EXIT /B 1
)

IF DEFINED ROBOT (
    SET OPT=-bench %ROBOT%
    IF DEFINED RANDOM_LIMIT ( SET OPT=%OPT% -random -limit %RANDOM_LIMIT% )
) ELSE (
    SET OPT=
)

IF DEFINED CPU ( SET OPT=%OPT% -cpu %CPU% )

IF "%COMMAND%" == "clean" (
    IF DEFINED ROBOT (
        IF EXIST %ROBOT%.ro ( DEL %ROBOT%.ro )
        FOR %%r IN (%ROBOT%) DO (
            ECHO DELETE FROM results_f2f WHERE robot='%%~nr'; | sqlite3 %DATABASE%
            ECHO DELETE FROM results_3vs3 WHERE robot='%%~nr'; | sqlite3 %DATABASE%
            ECHO DELETE FROM results_4vs4 WHERE robot='%%~nr'; | sqlite3 %DATABASE%
        )
    ) ELSE (
        ECHO DELETE FROM results_f2f; | sqlite3 %DATABASE%
        ECHO DELETE FROM results_3vs3; | sqlite3 %DATABASE%
        ECHO DELETE FROM results_4vs4; | sqlite3 %DATABASE%
        ECHO VACUUM; | sqlite3 %DATABASE%
    )
    EXIT /B
)

IF "%COMMAND%" == "init" (
    IF DEFINED ROBOT (
        crobots -c %ROBOT%.r <NUL >NUL
        FOR %%r IN (%ROBOT%) DO (
            ECHO INSERT INTO results_f2f(robot^) VALUES('%%~nr'^); | sqlite3 %DATABASE%
            ECHO INSERT INTO results_3vs3(robot^) VALUES('%%~nr'^); | sqlite3 %DATABASE%
            ECHO INSERT INTO results_4vs4(robot^) VALUES('%%~nr'^); | sqlite3 %DATABASE%
        ) 
    ) ELSE (
        FOR /F %%y IN ('yq -r .listRobots[] %CONF%') DO (
            ECHO INSERT INTO results_f2f(robot^) VALUES('%%~ny'^); | sqlite3 %DATABASE%
            ECHO INSERT INTO results_3vs3(robot^) VALUES('%%~ny'^); | sqlite3 %DATABASE%
            ECHO INSERT INTO results_4vs4(robot^) VALUES('%%~ny'^); | sqlite3 %DATABASE%
        )
    )
    EXIT /B
)

IF "%COMMAND%" == "reset" (
    ECHO UPDATE results_f2f SET games=0,ties=0,wins=0,points=0; | sqlite3 %DATABASE%
    ECHO UPDATE results_3vs3 SET games=0,ties=0,wins=0,points=0; | sqlite3 %DATABASE%
    ECHO UPDATE results_4vs4 SET games=0,ties=0,wins=0,points=0; | sqlite3 %DATABASE%
    EXIT /B
)

IF DEFINED ROBOT (
    IF NOT EXIST %ROBOT%.ro (
        crobots -c %ROBOT%.r <NUL >NUL
    ) ELSE (
        FOR %%r IN (%ROBOT%) DO SET NAME=%%~nr
        FOR /F %%i IN ('DIR /B /O:D %ROBOT%.r %ROBOT%.ro') DO SET NEWEST=%%i

        IF "%NAME%.r" == "%NEWEST%" ( crobots -c %ROBOT%.r <NUL >NUL )
    )
)

IF "%COMMAND%" == "test" (
    IF NOT EXIST %DATABASE% ECHO %DATABASE% doesn't exist
    gorobots -test -config %CONF% -out NUL -sql %SQL% -type f2f %OPT%
    EXIT /B
)

IF "%COMMAND%" == "f2f" (
    gorobots -config %CONF% -out NUL -sql %SQL% -type f2f %OPT%
    IF EXIST %SQL% TYPE %SQL% | sqlite3 %DATABASE%
    EXIT /B
)

IF "%COMMAND%" == "3vs3" (
    gorobots -config %CONF% -out NUL -sql %SQL% -type 3vs3 %OPT%
    IF EXIST %SQL% TYPE %SQL% | sqlite3 %DATABASE%
    EXIT /B
)

IF "%COMMAND%" == "4vs4" (
    gorobots -config %CONF% -out NUL -sql %SQL% -type 4vs4 %OPT%
    IF EXIST %SQL% TYPE %SQL% | sqlite3 %DATABASE%
    EXIT /B
)

IF "%COMMAND%" == "all" (
    gorobots -config %CONF% -out NUL -sql %SQL% -type f2f %OPT%
    IF EXIST %SQL% TYPE %SQL% | sqlite3 %DATABASE%
    gorobots -config %CONF% -out NUL -sql %SQL% -type 3vs3 %OPT%
    IF EXIST %SQL% TYPE %SQL% | sqlite3 %DATABASE%
    gorobots -config %CONF% -out NUL -sql %SQL% -type 4vs4 %OPT%
    IF EXIST %SQL% TYPE %SQL% | sqlite3 %DATABASE%
    EXIT /B
)

IF EXIST %SQL% DEL %SQL%
