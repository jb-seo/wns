#!/bin/bash

RED='\033[0;31m'
GREEN='\033[0;32m'
GREENB='\033[1;32m'
NC='\033[0m'

function goToLastUpdatedVersion() {
    DESTDIR=$1
    for VERSION in `ls ${DESTDIR} -t`
    do
        VERSIONDIR="${DESTDIR}/${VERSION}"
        if [ -d $VERSIONDIR ]; then
            cd $VERSIONDIR
            return 0
        fi
    done
}

function goto() {
    if [ ! $# -eq 1 ]; then
        echo "Please specify the component name"
        return 0;
    fi

    PROJECT=$1
    WORKDIR="${BUILDDIR}/work"
    CUTLEN=$(( ${#WORKDIR} + 2 ))
    CANDIDATES=`find ${WORKDIR} -maxdepth 2 | cut -c ${CUTLEN}- | grep "${PROJECT}"`

    COUNT=$(echo "$CANDIDATES" | wc -l)
    if [ ${COUNT} -eq 1 ]; then
        if [ ${#CANDIDATES} -gt 0 ]; then
            goToLastUpdatedVersion "${WORKDIR}/${CANDIDATES}"
        else
            echo -e "${RED}Cannot find...${NC}"
            return 1
        fi
    else
        x=1
        for DIR in $CANDIDATES
        do
            if [[ "$DIR" == */${PROJECT} ]]; then
                echo -e "${GREENB}${x}${NC}: ${DIR}"
            elif [[ "$DIR" == *${PROJECT} ]]; then
                echo -e "${GREEN}${x}${NC}: ${DIR}"
            else
                echo "${x}: ${DIR}"
            fi
            x=$(( $x + 1 ))
        done
        echo -en "Choose one >> "
        read -r SEL
        x=1
        for DIR in $CANDIDATES
        do
            if [ "$SEL" == "$x" ]; then
                goToLastUpdatedVersion "${WORKDIR}/${DIR}"
                return 0
            fi
            x=$(( $x + 1 ))
        done
        echo -e "${RED}Wrong choice...${NC}"
        return 1
    fi
}

function oe() {
    trap 'continue' SIGINT
    if [ ! $# -eq 2 ]; then
        echo "Usage {# oe command component} e.g. '# oe compile qtbase'"
        return 1;
    fi
    pushd . &> /dev/null
    COMMAND=$1
    PROJECT=$2
    RET=0
    goto $PROJECT
    if [ $? -eq 0 ]; then
        cd temp
        FILE="run.do_${COMMAND}"
        if [ -f $FILE ]; then
            ./${FILE}
            RET=$?
        else
            echo "[ERROR] Can't find '${COMMAND}' command file. Please to bitbake job first."
        fi
    fi
    popd &> /dev/null
    trap - SIGINT
    return ${RET}
}

function bb() {
    trap 'continue' SIGINT
    pushd . &> /dev/null
    cd $WEBOS_BASE_DIR
    bitbake $@
    popd &> /dev/null
    trap - SIGINT
}

function @add() {
   NAME=$1
   if [ "${NAME:0:1}" == "@" ]; then
       echo -e "${RED}Don't use '@' as first character.${NC}"
       return 1
   fi

   EXIST=`alias @${NAME} 2> /dev/null`
   if [ $? -eq 0 ]; then
       echo -e "${RED}Already exist...${NC}"
       return 1
   fi

   # replace to relational dir if it is a child dir
   CMD="alias @${NAME}='cd ${PWD/${WEBOS_BASE_DIR}/\$\{WEBOS_BASE_DIR\}}'"
   eval $CMD
   echo $CMD >> $GO_ALIAS_FILE
   return 0
}

function @del() {
   NAME=$1
   EXIST=`alias @${NAME} 2> /dev/null`
   if [ $? -ne 0 ]; then
       echo -e "${RED}Doesn't exist...${NC}"
       return 1
   fi

   eval "unalias @${NAME}"
   sed -i "/@${NAME}=/d" $GO_ALIAS_FILE
}

if [ "$WEBOS_BASE_DIR" == "" ]; then
    export WEBOS_BASE_DIR=$PWD
    source ~/.bashrc
    source ./oe-init-build-env

    # PS1
    export PS1="${BASH_BRed}RPi3 ${BASH_BBlue}\w${BASH_BGreen}\$(__git_ps1 \" [%s]\")${BASH_Color_Off} \$ "
    export PROMPT_DIRTRIM=2

    # set aliases
    alias @base="cd $WEBOS_BASE_DIR"
    alias @sysroot="cd ${WEBOS_BASE_DIR}/BUILD/sysroots/${MACHINE}"

    # load saved aliases
    export GO_ALIAS_FILE="${WEBOS_BASE_DIR}/.go_aliases"
    if [ -f $GO_ALIAS_FILE ]; then
        source $GO_ALIAS_FILE
    fi

    # add bb completion
    if [ `type -t _bitbake` == "function" ]; then
        complete -F _bitbake bb
    fi
else
    echo "Already in go shell"
fi
