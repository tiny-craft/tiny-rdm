#!/bin/bash
clear
BLACK="\033[0;30m"
DARK_GRAY="\033[1;30m"
BLUE="\033[0;34m"
LIGHT_BLUE="\033[1;34m"
GREEN="\033[0;32m"
LIGHT_GREEN="\033[1;32m"
CYAN="\033[0;36m"
LIGHT_CYAN="\033[1;36m"
RED="\033[0;31m"
LIGHT_RED="\033[1;31m"
PURPLE="\033[0;35m"
LIGHT_PURPLE="\033[1;35m"
BROWN="\033[0;33m"
YELLOW="\033[0;33m"
LIGHT_GRAY="\033[0;37m"
WHITE="\033[1;37m"
NC="\033[0m"

parentPath=$( cd "$(dirname "${BASH_SOURCE[0]}")" ; pwd -P )
cd "$parentPath"
appPath=$( find "$parentPath" -name '*.app' -maxdepth 1)
appName=${appPath##*/}
appBashName=${appName// /\ }
appDIR="/Applications/${appBashName}"
echo -e "This tool fix these situations: \"${appBashName}\" is damaged and can't not be opened."
echo ""
if [ ! -d "$appDIR" ];then
  echo ""
  echo -e "Execution result: ${RED}You haven't installed ${appBashName} yet, please install it first.${NC}"
  else
  echo -e "${YELLOW}Please enter your login password, and then press enter. (The password is invisible during input)${NC}"
  sudo spctl --master-disable
  sudo xattr -rd com.apple.quarantine /Applications/"$appBashName"
  sudo xattr -rc /Applications/"$appBashName"
  sudo codesign --sign - --force --deep /Applications/"$appBashName"
  echo -e "Execution result: ${GREEN}Already fixed! ${NC} ${appBashName} will work correctly.${NC}"
fi
echo -e "You can close this window now"
