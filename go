#!/bin/bash

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

    RED='\033[0;31m'
    GREEN='\033[1;32m'
    NC='\033[0m'

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
    trap '' SIGINT
    return ${RET}
}

function bb() {
    trap 'continue' SIGINT
    pushd . &> /dev/null
    cd $WEBOS_BASE_DIR
    bitbake $@
    popd &> /dev/null
    trap '' SIGINT
}

if [ "$WEBOS_BASE_DIR" == "" ]; then
    export WEBOS_BASE_DIR=$PWD
    source ~/.bashrc
    source ./oe-init-build-env

    # PS1
    export PS1="${BASH_BRed}RPi3 ${BASH_BBlue}\w${BASH_BGreen}\$(__git_ps1 \" [%s]\")${BASH_Color_Off} \$ "
    export PROMPT_DIRTRIM=3

    # set aliases
    alias @base="cd $WEBOS_BASE_DIR"
    alias @work="cd ${WEBOS_BASE_DIR}/BUILD/work"
    alias @sysroot="cd ${WEBOS_BASE_DIR}/BUILD/sysroots"
else
    echo "Already in go shell"
fi
