# stresstest

Stress tool to test graphite install (initially).
Usage of https://github.com/fatih/pool

## Setup
> mkdir -p ~/stress/{bin,src}
Set GOPATH to GOPATH to your stresstest directory.
> export GOPATH=~/stress
Next, get pool src
> go get gopkg.in/fatih/pool.v2
Get the stresstest src
> go get github.com/mlambrichs/stresstest
cd to stresstest dir
> cd $GOPATH/src/github.com/mlambrichs/stresstest
...and install
> go install
This will result in a binary 'stresstest' in $GOPATH/bin
So, do this:
> export PATH=$PATH:$GOPATH/bin
so you can call 'stresstest'.

## Usage
stresstest -h


## Contributors

https://github.com/mlambrichs/stresstest/graphs/contributors

