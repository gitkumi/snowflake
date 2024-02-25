#!/bin/bash

rename() {
    local directory="$1"
    shopt -s dotglob
    for file in "$directory"/*; do
        if [ -f "$file" ]; then
            mv "$file" "$file.templ"
        elif [ -d "$file" ]; then  
            rename "$file"
        fi
    done
    shopt -u dotglob
}

directory_path="./template/files"
rename "$directory_path"
