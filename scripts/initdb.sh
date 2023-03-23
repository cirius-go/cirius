#!/bin/bash

# Use find to list the YAML files in the specified directory,
# and remove the extension suffix from the filenames
mapfile -t dbs < <(find ./scripts/assets/docker-composes -type f -name '*.yaml' -printf '%f\n' | sed 's/\.[^.]*$//')

echo "Please select a database:"

# Loop through the array of filenames and print each one
for file in "${dbs[@]}"
do
  echo "$file"
done

# Prompt the user for input and read the selected filename
read -p "Enter the name of the database: " file

# Loop until a valid filename is entered
while [[ -z "$file" ]] || ! [ -e "./scripts/assets/docker-composes/$file.yaml" ]
do
  if [[ -z "$file" ]]; then
    echo "You must enter a filename."
  else
    echo "The file '$file.yaml' does not exist."
  fi
  read -p "Enter the name of the database: " file
done

# Copy the selected file to the output file
cp "./scripts/assets/docker-composes/$file.yaml" "./docker-compose.yaml"
echo "Selected database: $file. New 'docker-compose.yaml' file is created"
