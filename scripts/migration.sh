#!/bin/bash
if [[ -z $1 ]]; then
    echo 'Path is required!'
    exit 1
fi

if [[ -z $2 ]]; then
    echo 'Database is required!'
    exit 1
fi

if [[ -z $3 ]]; then
    echo 'Migration version minimum is required!'
    exit 1
fi

if [[ -z $4 ]]; then
    echo 'Migration version target is required!'
    exit 1
fi

path=$1
database=$2
min_version=$3
target_version=$4
last_file=$(ls -t "$path" | grep '.sql' | sort | tail -n 1)
max_version=$(echo $last_file | grep '^[0-9]' | sed 's/^\([0-9]*\).*/\1/' | awk '{print $1 + 0}')
current_version=$(migrate -path $path -database $database version 2>&1)
if [[ $current_version == "error: no migration" ]]; then
    current_version=0
fi

if  [[ $target_version == -1 ]]; then
    target_version=$max_version
fi

if [[ $target_version == 0  && $current_version == 0 ]]; then
    echo 'There is no migration that has been executed. Please specify the migration version target!'
    exit 1
fi

if  (( $target_version < -1 || $target_version > $max_version)); then
    echo 'Invalid version!'
    exit 1
fi

if  (( $target_version < $min_version && $target_version != -1)); then
    echo 'Unfortunately, this action is forbidden. Migrating a database into a lower-than-last version is dangerous!'
    exit 1
fi

if [[ $current_version == $target_version ]]; then 
    echo 'No migration change!'
    exit 0
fi

if  [[ $target_version == $max_version ]]; then
    migrate -path $path -database $database -verbose up
    echo 'Migrate up is successfully executed!'
    exit 0
fi

if [[ $target_version == 0 ]]; then 
    migrate -path $path -database $database -verbose down
    echo 'Migrate down is successfully executed!'
    exit 0
fi

if [[ $current_version < $target_version ]]; then
    gap=$(expr $target_version - $current_version)
    migrate -path $path -database $database -verbose up $gap
    echo 'Migrate up for specific version is successfully executed!'
    exit 0
fi

if [[ $current_version > $target_version ]]; then 
    gap=$(expr $current_version - $target_version)
    migrate -path $path -database $database -verbose down $gap
    echo 'Migrate down for specific version is successfully executed!'
    exit 0
fi

echo 'Unhandled migration condition!'
exit 1