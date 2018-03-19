webOS Navigation Shell
======================

It's a simple navigation tool for webOS (OSE) - <https://github.com/webosose/build-webos>.

## How to run
1. prepare build-webos and run ./mcf command at least once
2. copy to 'wns' to your build-webos folder
3. execute '# bash --rcfile wns'

## Usage
1. goto '{keyword}' : move last version of component folder (e.g. # goto nodejs)
2. oe {command} {component} : direct call to run pre-created bitbake scripts (e.g. # oe compile chromium53)
3. bb ~~~ : just a alias of bitbake but can use anywhere
4. btail {component} ({command}) : tailing last activity of bitbake work (e.g. # btail chromium53 compile)
5. @add {key} : add a change directory shortcut for current PWD
6. @del {key} : delete added shortcut
7. @list : list added shortcuts
8. @base : cd to build-webos root
9. @sysroot : cd to target board's sysroot
