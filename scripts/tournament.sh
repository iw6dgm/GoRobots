#!/bin/bash
#
# Tournament Execution script v.1.15 21/10/2024 (C) Maurizio Camangi
#

set -e

CONFIG=$1
MODE=$2
PROCESS=`basename $0`
GOROBOTS=gorobots
unset OPT

if [ $# -ne 2 ]
then
{
  echo "Usage: ${PROCESS} [config] [f2f|3vs3|4vs4|all|test|init|clean]"
  exit 1
}
fi

if [ ! -s ${CONFIG} ]
then
{
 echo ${CONFIG} does not exist. Exit.
 exit 1
}
fi

if [ -z ${DATABASE+x} ] ; then
  DATABASE=db/`basename $CONFIG .yml`.db
fi

if [ ! -s $DATABASE ]
then
{
 echo $DATABASE does not exist
 exit 1
}
fi

if ! command -v yq >/dev/null 2>&1
then
{
 echo yq is not installed
 exit 1
}
fi

if ! command -v sqlite3 >/dev/null 2>&1
then
{
 echo sqlite3 is not installed
 exit 1
}
fi

if [ ! -z ${ROBOT+x} ] ; then
  OPT="-bench $ROBOT"
fi

if [ "$ROBOT" != "" ] && [ "$MODE" != "clean" ] && [ ! -s ${ROBOT}.ro ]
then
{
  echo "Robot ${ROBOT}.ro not found. Compiling..."
  crobots -c ${ROBOT}.r </dev/null >/dev/null 2>&1
}
fi

if [ "$MODE" = "test" ]
then
{
  echo "Processing $CONFIG ..."
  $GOROBOTS $OPT -type 4vs4 -config $CONFIG -test
  echo "Test mode completed. Script exits."
  exit 0
}
fi

if [ "$MODE" = "clean" ]
then
{
  if [ "$ROBOT" != "" ]
  then
  { # only bench robot
    [ -s ${ROBOT}.ro ] && echo "Deleting ${ROBOT}.ro .." && rm -f ${ROBOT}.ro
    NAME=`basename $ROBOT`
    echo "DELETE FROM results_f2f WHERE robot='${NAME}';" | sqlite3 $DATABASE
    echo "DELETE FROM results_3vs3 WHERE robot='${NAME}';" | sqlite3 $DATABASE
    echo "DELETE FROM results_4vs4 WHERE robot='${NAME}';" | sqlite3 $DATABASE
  }
  else
  { # all robots
    cat <<EOF | sqlite3 $DATABASE
DELETE FROM results_f2f;
DELETE FROM results_3vs3;
DELETE FROM results_4vs4;
VACUUM;
EOF
  }
  fi
  exit 0
}
fi

if [ "$MODE" = "init" ]
then
{
  if [ "$ROBOT" != "" ]
  then
  { # only bench robot
    [ ${ROBOT}.r -nt ${ROBOT}.ro ] && echo "Robot ${ROBOT}.r newer than .ro .. Re-compiling" && crobots -c ${ROBOT}.r </dev/null >/dev/null 2>&1
    NAME=`basename $ROBOT`
    echo "INSERT INTO results_f2f(robot) VALUES('${NAME}');" | sqlite3 $DATABASE
    echo "INSERT INTO results_3vs3(robot) VALUES('${NAME}');" | sqlite3 $DATABASE
    echo "INSERT INTO results_4vs4(robot) VALUES('${NAME}');" | sqlite3 $DATABASE
  }
  else
  { # for all robots
    for i in `yq -r '.listRobots[]' $CONFIG` ; do
      NAME=`basename $i`
      echo "INSERT INTO results_f2f(robot) VALUES('${NAME}');" | sqlite3 $DATABASE
      echo "INSERT INTO results_3vs3(robot) VALUES('${NAME}');" | sqlite3 $DATABASE
      echo "INSERT INTO results_4vs4(robot) VALUES('${NAME}');" | sqlite3 $DATABASE
    done
  }
  fi
  exit 0
}
fi

TMP_SQL=/tmp/$$TMP.sql

# F2F
if  [ "$MODE" = "f2f" ] || [ "$MODE" = "all" ]
then
{
  $GOROBOTS $OPT -config $CONFIG -type f2f -sql $TMP_SQL -out /dev/null
  cat $TMP_SQL | sqlite3 $DATABASE
}
fi

# 3vs3
if  [ "$MODE" = "3vs3" ] || [ "$MODE" = "all" ]
then
{
  $GOROBOTS $OPT -config $CONFIG -type 3vs3 -sql $TMP_SQL -out /dev/null
  cat $TMP_SQL | sqlite3 $DATABASE
}
fi

# 4vs4
if  [ "$MODE" = "4vs4" ] || [ "$MODE" = "all" ]
then
{
  $GOROBOTS $OPT -config $CONFIG -type 4vs4 -sql $TMP_SQL -out /dev/null
  cat $TMP_SQL | sqlite3 $DATABASE
}
fi