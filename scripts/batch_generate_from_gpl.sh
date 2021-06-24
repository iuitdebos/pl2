#!/bin/sh

gplDir=$1
outDir=$2

for gplPath in $(ls ${gplDir}/*.gpl); do

	gplBasename=$(basename $gplPath)
	outPath=${outDir}/${gplBasename%.gpl}.pl2

	echo "${gplBasename} => ${outPath}"

	pl2-from-gpl -gpl $gplPath -pl2 $outPath &
done
