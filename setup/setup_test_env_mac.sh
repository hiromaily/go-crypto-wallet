#!/bin/sh

#access github.com/this page and copy and paste on local

sudo xcodebuild -license
xcode-select --install

ruby -e "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install)"
brew update
brew tap caskroom/cask

brew install git
brew install putty
brew install vim
brew install zsh
brew install zsh-completions
brew install go
brew install node.js
brew install python
brew install pyenv
brew install mysql
brew install redis
brew install mercurial
brew install tree
brew install wget
brew install nmap
brew install readline
brew install tmux
brew install jq

export HOMEBREW_CASK_OPTS="--appdir=/Applications"
brew cask install iterm2

brew cask install google-chrome
brew cask install slack
brew cask install sequel-pro
brew cask install docker