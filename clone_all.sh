#!/bin/bash -e

#clone all the gonet2 repo into current dir.
mkdir gonet2
cd gonet2

git clone git@github.com:gonet2/agent.git
git clone git@github.com:gonet2/game.git
git clone git@github.com:gonet2/snowflake.git
git clone git@github.com:gonet2/chat.git
git clone git@github.com:gonet2/auth.git
git clone git@github.com:gonet2/libs.git
git clone git@github.com:gonet2/rank.git
git clone git@github.com:gonet2/geoip.git
git clone git@github.com:gonet2/arch.git
git clone git@github.com:gonet2/bgsave.git
git clone git@github.com:gonet2/wordfilter.git

