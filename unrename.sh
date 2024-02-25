#!/bin/bash

unrename() {
    local directory="$1"
    shopt -s dotglob
    for file in "$directory"/*; do
        if [ -f "$file" ]; then
            echo $file;
            new_file=${file%.*};
            mv "$file" "$new_file"
        elif [ -d "$file" ]; then
            unrename "$file"
        fi
    done
    shopt -u dotglob
}

directory_path="./template/files"
unrename "$directory_path"
